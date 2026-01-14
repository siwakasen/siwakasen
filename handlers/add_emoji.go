package handlers

import (
	"net/http"

	"github.com/siwakasen/siwakasen/utils/github"
)

var allowedEmojiTypes = map[string]bool{
	"👊":  true,
	"😎":  true,
	"❤️": true,
	"👋":  true,
	"👍":  true,
	"😁":  true,
	"😅":  true,
	"😜":  true,
	"🤩":  true,
	"🤯":  true,
}

func AddMoji(w http.ResponseWriter, req *http.Request) {
	redirectURL := "https://github.com/siwakasen"

	emojiType := req.URL.Query().Get("type")

	if emojiType == "" {
		http.Redirect(w, req, redirectURL, http.StatusFound)
		return
	}

	if !allowedEmojiTypes[emojiType] {
		http.Redirect(w, req, redirectURL, http.StatusFound)
		return
	}

	err := github.UpdateReadme(emojiType)
	if err != nil {
		http.Redirect(w, req, redirectURL, http.StatusFound)
		return
	}

	http.Redirect(w, req, redirectURL, http.StatusFound)
}
