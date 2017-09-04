package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	ERR_INVALID = errors.New("Could not decode request")
)

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

func findFeatureArray(dec *json.Decoder) error {
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return ERR_INVALID
		}
		if s, ok := t.(string); ok && s == "payload" && dec.More() {
			//	fmt.Println(s)
			// we found the payload element
			d, err := dec.Token()
			if err != nil {
				return ERR_INVALID
			}
			if delim, ok := d.(json.Delim); ok {
				if delim.String() == "[" {
					// we have our start of the array
					break
				} else {
					// A different kind of delimiter
					return fmt.Errorf("Expected features to be an array.")
				}
			}
		}
	}

	if !dec.More() {
		return ERR_INVALID // fmt.Errorf("Cannot find any features.")
	}
	return nil
}

func parse(payload io.Reader) ([]byte, error) {
	dec := json.NewDecoder(payload)
	err := findFeatureArray(dec)
	if err != nil {
		return nil, ERR_INVALID
	}

	respMap := make(map[string]interface{})
	var respList []Response

	for dec.More() {
		var f Message
		err = dec.Decode(&f)
		if err != nil {
			return nil, ERR_INVALID
		}

		if f.Drm && f.EpisodeCount > 0 && f.ImageData != nil && f.ImageData.ShowImage != "" {
			resp := Response{
				Title:    strings.Trim(strings.Split(f.Title, "(")[0], " \t\n"),
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

	js, err := parse(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func Serve() {
	http.HandleFunc("/", queryHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	Serve()
}
