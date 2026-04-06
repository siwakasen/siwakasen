// Package handlers
package handlers

import (
	"log"
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

	log.Println(req, "request emoji: ", emojiType)

	if emojiType == "" {
		log.Println(req, "emoji type not found")
		http.Redirect(w, req, redirectURL, http.StatusFound)
		return
	}

	if !allowedEmojiTypes[emojiType] {
		log.Println(req, "emoji type not allowed")

		http.Redirect(w, req, redirectURL, http.StatusFound)
		return
	}

	err := github.UpdateReadme(emojiType)
	if err != nil {
		log.Println(req, err)
		http.Redirect(w, req, redirectURL, http.StatusFound)
		return
	}

	http.Redirect(w, req, redirectURL, http.StatusFound)
}
