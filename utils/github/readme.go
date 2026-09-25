// Package github
package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ghRepoOwner = "siwakasen"
	ghRepoName  = "siwakasen"
)

var ghToken = os.Getenv("GH_TOKEN")

var baseHeaders = map[string]string{
	"Accept":       "application/vnd.github.v3+json",
	"Content-Type": "application/json",
	"User-Agent":   "siwakasen-gh-readme",
}

type githubContentResponse struct {
	Content string `json:"content"`
	SHA     string `json:"sha"`
}

var url = fmt.Sprintf(
	"https://api.github.com/repos/%s/%s/contents/README.md",
	ghRepoOwner,
	ghRepoName,
)
var client = &http.Client{Timeout: 15 * time.Second}

func spanNotFound(emojiType string) error {
	return fmt.Errorf("span not found for emoji %q", emojiType)
}

func GetReadme(emojiType string) ([]byte, error) {
	if strings.TrimSpace(ghToken) == "" {
		return nil, fmt.Errorf("GH_TOKEN is not set")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range baseHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("Authorization", "Bearer "+ghToken)

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	log.Printf("GET README response time: %.2fs", duration.Seconds())
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close GET README body response: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("github GET README failed: %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var ghResp githubContentResponse
	if err := json.Unmarshal(body, &ghResp); err != nil {
		return nil, err
	}
	if ghResp.Content == "" || ghResp.SHA == "" {
		return nil, fmt.Errorf("github response missing README content or sha")
	}

	decoded, err := base64.StdEncoding.DecodeString(
		strings.ReplaceAll(ghResp.Content, "\n", ""),
	)
	if err != nil {
		return nil, err
	}
	readme := string(decoded)
	newReadme, err := incrementEmojiCount(readme, emojiType)
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(newReadme))

	payload, err := json.Marshal(map[string]string{
		"message": fmt.Sprintf("chore: Add %s count", emojiType),
		"content": encoded,
		"sha":     ghResp.SHA,
	})
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func incrementEmojiCount(readme, emojiType string) (string, error) {
	spanRegex := regexp.MustCompile(
		fmt.Sprintf(`<span[^>]*id=["']count-%s["'][^>]*>(\d+)</span>`, regexp.QuoteMeta(emojiType)),
	)
	match := spanRegex.FindStringSubmatchIndex(readme)

	if len(match) < 4 {
		return "", spanNotFound(emojiType)
	}

	countStart, countEnd := match[2], match[3]
	prev, err := strconv.Atoi(readme[countStart:countEnd])
	if err != nil {
		return "", fmt.Errorf("invalid count value for emoji %q: %w", emojiType, err)
	}

	return readme[:countStart] +
		fmt.Sprintf("%d", prev+1) +
		readme[countEnd:], nil
}

func UpdateReadme(payload []byte) error {
	// PUT README
	putReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	for k, v := range baseHeaders {
		putReq.Header.Set(k, v)
	}
	putReq.Header.Set("Authorization", "Bearer "+ghToken)

	start := time.Now()
	putResp, err := client.Do(putReq)
	duration := time.Since(start)
	log.Printf("PUT README response time: %.2fs", duration.Seconds())

	if err != nil {
		return err
	}
	defer func() {
		if err := putResp.Body.Close(); err != nil {
			log.Printf("failed to close body PUT README response: %v", err)
		}
	}()
	if putResp.StatusCode >= 300 {
		errBody, err := io.ReadAll(putResp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("github PUT README failed: %d: %s", putResp.StatusCode, strings.TrimSpace(string(errBody)))
	}

	return nil
}
