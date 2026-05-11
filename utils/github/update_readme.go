// Package github
package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func UpdateReadme(emojiType string) error {
	if strings.TrimSpace(ghToken) == "" {
		return fmt.Errorf("GH_TOKEN is not set")
	}

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/README.md",
		ghRepoOwner,
		ghRepoName,
	)

	client := &http.Client{}

	// GET README
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	for k, v := range baseHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("Authorization", "Bearer "+ghToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("github GET README failed: %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var ghResp githubContentResponse
	if err := json.Unmarshal(body, &ghResp); err != nil {
		return err
	}
	if ghResp.Content == "" || ghResp.SHA == "" {
		return fmt.Errorf("github response missing README content or sha")
	}

	decoded, err := base64.StdEncoding.DecodeString(
		strings.ReplaceAll(ghResp.Content, "\n", ""),
	)
	if err != nil {
		return err
	}
	readme := string(decoded)

	spanRegex := regexp.MustCompile(
		fmt.Sprintf(`<span[^>]*id=["']count-%s["'][^>]*>(\d+)</span>`, regexp.QuoteMeta(emojiType)),
	)
	match := spanRegex.FindStringSubmatchIndex(readme)
	if match == nil || len(match) < 4 {
		return fmt.Errorf("span not found for emoji %q", emojiType)
	}

	countStart, countEnd := match[2], match[3]
	prev, err := strconv.Atoi(readme[countStart:countEnd])
	if err != nil {
		return fmt.Errorf("invalid count value for emoji %q: %w", emojiType, err)
	}

	newReadme := readme[:countStart] +
		fmt.Sprintf("%d", prev+1) +
		readme[countEnd:]

	encoded := base64.StdEncoding.EncodeToString([]byte(newReadme))

	payload, err := json.Marshal(map[string]string{
		"message": fmt.Sprintf("chore: Add %s count", emojiType),
		"content": encoded,
		"sha":     ghResp.SHA,
	})
	if err != nil {
		return err
	}

	// PUT README
	putReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	for k, v := range baseHeaders {
		putReq.Header.Set(k, v)
	}
	putReq.Header.Set("Authorization", "Bearer "+ghToken)

	putResp, err := client.Do(putReq)
	if err != nil {
		return err
	}
	defer putResp.Body.Close()

	if putResp.StatusCode >= 300 {
		errBody, err := io.ReadAll(putResp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("github PUT README failed: %d: %s", putResp.StatusCode, strings.TrimSpace(string(errBody)))
	}

	return nil
}
