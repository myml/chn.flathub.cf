package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func server() {
	server, _ := url.Parse("https://dl.flathub.org")
	proxy := httputil.NewSingleHostReverseProxy(server)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = server.Host
		r.URL.Scheme = server.Scheme
		r.Host = server.Host
		log.Println("proxy", r.URL.String())

		proxy.ServeHTTP(w, r)
	})
	http.HandleFunc("/repo/summary.sig", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	http.HandleFunc("/repo/summary", func(w http.ResponseWriter, r *http.Request) {
		log.Println("summary")
		r.URL.Host = server.Host
		r.URL.Scheme = server.Scheme
		resp, err := http.Get(r.URL.String())
		if err != nil {
			log.Println("error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if resp.StatusCode >= 400 {
			w.WriteHeader(resp.StatusCode)
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data = bytes.Replace(data, []byte("https://dl.flathub.org/repo/"), []byte("http://localhost:18080/repo/"), 1)
		w.Write(data)
	})
	log.Println(http.ListenAndServe(":18080", nil))
}

func summaryReplace(repo [30]byte) {
}
func main() {
	server()
	return
	b, err := ioutil.ReadFile("./summary")
	if err != nil {
		log.Panic(err)
	}
	b = bytes.Replace(b, []byte("https://dl.flathub.org/repo/"), []byte("http://localhost:18080/repo/"), 1)
	ioutil.WriteFile("/home/myml/.local/share/flatpak/repo/summary", b, 0655)
}

// https://dl.flathub.org/repo/
// http://localhost:18080/repo/
