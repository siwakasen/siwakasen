package handlers

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")

	emojiType := req.URL.Query().Get("type")

	if emojiType == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "type is required",
		})
		return
	}

	if !allowedEmojiTypes[emojiType] {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "emoji not allowed",
		})
		return
	}

	err := github.UpdateReadme(emojiType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "😱",
			"error":   err.Error(),
		})
		return
	}

	w.Header().Add("Location", "https://github.com/siwakasen")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": emojiType,
	})

}
