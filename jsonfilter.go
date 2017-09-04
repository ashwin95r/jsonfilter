package main

import (
	"encoding/json"
	"errors"
)

var (
	ERR_INVALID = errors.New("Could not decode request")
)

// Request is the struct for Decoding the request.
type Request struct {
	Payload []Message `json:"payload"`
}

type Image struct {
	ShowImage string `json:"showImage"`
}

type Message struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	ImageData    *Image `json:"image"`
	Drm          bool   `json:"drm"`
	EpisodeCount int    `json:"episodeCount"`
}

// Response is the struct for encoding the result.
type Response struct {
	ImageURL string `json:"image"`
	Slug     string `json:"slug"`
	Title    string `json:"title"`
}

func parse(payload []byte) ([]byte, error) {
	var req Request
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, ERR_INVALID
	}

	respMap := make(map[string]interface{})
	var respList []Response
	for _, f := range req.Payload {
		if f.Drm && f.EpisodeCount > 0 && f.ImageData != nil && f.ImageData.ShowImage != "" {
			// Add those objects that satisfy the filter conditions.
			resp := Response{
				Title:    f.Title,
				ImageURL: f.ImageData.ShowImage,
				Slug:     f.Slug,
			}
			respList = append(respList, resp)
		}
	}
	respMap["response"] = respList
	return json.Marshal(respMap)
}
