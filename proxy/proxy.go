package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"regexp"
)

// Server used for proxy
func Server(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	re := regexp.MustCompile(`(\:\/)([^\/])`)
	url := re.ReplaceAllString(r.URL.Path[1:], "$1/$2")

	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	for key, val := range resp.Header {
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}
