package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	ghRepoOwner = "siwakasen"
	ghRepoName  = "siwakasen"
)

var ghToken = os.Getenv("GH_TOKEN")

var baseHeaders = map[string]string{
	"Accept":        "application/vnd.github.v3+json",
	"Content-Type":  "application/json",
	"User-Agent":    "siwakasen-gh-readme",
	"Authorization": "",
}

type githubContentResponse struct {
	Content string `json:"content"`
	SHA     string `json:"sha"`
}

func UpdateReadme(emojiType string) error {
	baseHeaders["Authorization"] = "token " + ghToken

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/README.md",
		ghRepoOwner,
		ghRepoName,
	)

	client := &http.Client{}

	// GET README
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range baseHeaders {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var ghResp githubContentResponse
	if err := json.Unmarshal(body, &ghResp); err != nil {
		return err
	}

	decoded, _ := base64.StdEncoding.DecodeString(
		strings.ReplaceAll(ghResp.Content, "\n", ""),
	)
	readme := string(decoded)

	target := `<span id="count-` + emojiType + `">`
	start := strings.Index(readme, target)
	if start == -1 {
		return fmt.Errorf("span not found")
	}

	start += len(target)
	end := strings.Index(readme[start:], "</span>") + start

	var prev int
	fmt.Sscanf(readme[start:end], "%d", &prev)

	newReadme := readme[:start] +
		fmt.Sprintf("%d", prev+1) +
		readme[end:]

	encoded := base64.StdEncoding.EncodeToString([]byte(newReadme))

	payload, _ := json.Marshal(map[string]string{
		"message": fmt.Sprintf("chore: Add %s count", emojiType),
		"content": encoded,
		"sha":     ghResp.SHA,
	})

	// PUT README
	putReq, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	for k, v := range baseHeaders {
		putReq.Header.Set(k, v)
	}

	putResp, err := client.Do(putReq)
	if err != nil {
		return err
	}
	defer putResp.Body.Close()

	if putResp.StatusCode >= 300 {
		errBody, _ := io.ReadAll(putResp.Body)
		return fmt.Errorf("github error: %s", errBody)
	}

	return nil
}
