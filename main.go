package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ashwin95r/jsonfilter/filter"
)

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	mp := map[string]string{
		"error": err.Error(),
	}
	js, _ := json.Marshal(mp)
	fmt.Fprint(w, string(js))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		// Only accept POST.
		writeError(w, errors.New("Invalid method"))
		return
	}

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, err)
		return
	}

	js, err := filter.Parse(req)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func Serve() {
	http.HandleFunc("/", queryHandler)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Println("Listening on port ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	Serve()
}
