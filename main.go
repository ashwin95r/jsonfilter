package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	ERR_INVALID = errors.New("Could not decode request")
)

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
			resp := Response{
				Title:    f.Title,
				ImageURL: f.ImageData.ShowImage,
				Slug:     f.Slug,
			}
			respList = append(respList, resp)
		}

		//	fmt.Println(f)
	}
	respMap["response"] = respList
	return json.Marshal(respMap)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid method")
		return
	}

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "%v", err)
		return
	}
	js, err := parse(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		mp := map[string]string{
			"error": err.Error(),
		}
		js, _ := json.Marshal(mp)
		fmt.Fprint(w, string(js))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func Serve() {
	http.HandleFunc("/", queryHandler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func main() {
	Serve()
}
