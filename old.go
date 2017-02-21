package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}

func proxy(w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	fmt.Printf("%+v\n", r)
	fmt.Printf("BODY: %s\n", r.Body)

	url := r.URL.String()
	req, err := http.NewRequest(r.Method, url, r.Body)

	if err != nil {
		fmt.Errorf("Error loading: %s", url)
		io.WriteString(w, "")
	}
	response, err := client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Println("Request Error %s", url)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Errorf("Error reading: %s", url)
	}
	wHeader := w.Header()
	copyHeader(response.Header, &wHeader)
	wHeader.Add("Requested-Host", req.Host)
	w.WriteHeader(response.StatusCode)
	w.Write(body)
}

func main() {
	http.HandleFunc("/", proxy)
	http.ListenAndServe(":8888", nil)
}
