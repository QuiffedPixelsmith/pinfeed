package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	origin  = "https://www.pinterest.com/"
	repoURL = "https://github.com/attilaolah/"
)

var (
	thumb       = regexp.MustCompile("\\b(https?://[0-9a-z-]+.pinimg.com/)192(x/[/0-9a-f]+.jpg)\\b")
	replacement = []byte("${1}736${2}")
)

func main() {
	http.HandleFunc("/", pinFeed)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func pinFeed(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, repoURL, http.StatusMovedPermanently)
		return
	}
	res, err := http.Get(feedURL(username(r.URL.Path)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	for key, vals := range res.Header {
		for _, val := range vals {
			w.Header().Add(key, val)
		}
	}
	w.WriteHeader(res.StatusCode)
	buf, err := replaceThumbs(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(buf)
}

func feedURL(username string) string {
	return origin + "/" + username + "/feed.rss"
}

func username(path string) string {
	return strings.SplitN(path, "/", 3)[1]
}

func replaceThumbs(r io.Reader) (buf []byte, err error) {
	if buf, err = ioutil.ReadAll(r); err == nil {
		buf = thumb.ReplaceAll(buf, replacement)
	}
	return
}
