package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		// Only accept POST.
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
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Println("Listening on port ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	Serve()
}
