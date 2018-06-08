package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

var (
	tdstatsdVersion = "dev"
	tdstatsdCommit  = "none"
	tdstatsdDate    = "unknown"
)

func main() {
	var (
		tdURL       string
		bindAddr    string
		showVersion bool
	)
	flag.StringVar(&tdURL, "url", "http://localhost:8000/_perf.txt",
		"Traffic Director stats page URL")
	flag.StringVar(&bindAddr, "bind", "0.0.0.0:8081",
		"Listen on address:port")
	flag.BoolVar(&showVersion, "version", false, "Show version and exit")
	flag.Parse()

	if showVersion {
		fmt.Printf("%v, commit %v, built at %v\n",
			tdstatsdVersion, tdstatsdCommit, tdstatsdDate)
		os.Exit(0)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	router := http.NewServeMux()
	router.Handle("/pools", pools(tdURL, client))
	router.HandleFunc("/", index)
	handleStatic(router)

	log.Println("Listening on", bindAddr)
	log.Fatal(http.ListenAndServe(bindAddr, router))
}

func pools(url string, c *http.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pools" {
			http.Error(w, http.StatusText(http.StatusNotFound),
				http.StatusNotFound)
			return
		}
		data, err := getTdData(url, c)
		internalServerError := func(err error) {
			log.Println(err)
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
		}
		if err != nil {
			internalServerError(err)
			return
		}
		pools, err := parse(data)
		if err != nil {
			internalServerError(err)
			return
		}
		sort.Sort(byStatus(pools))
		j, err := json.Marshal(pools)
		if err != nil {
			internalServerError(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(j); err != nil {
			internalServerError(err)
			return
		}
	})
}

func getTdData(url string, c *http.Client) ([]byte, error) {
	resp, err := c.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error getting url '%s': %s", url, err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing body: %s", cerr)
		}
	}()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %s", err)
	}
	return buf, nil
}

func handleStatic(router *http.ServeMux) {
	// serve vue
	vueData, err := unpack(vueCompressedData)
	if err != nil {
		log.Fatal("decompressing vue error:", err)
	}
	router.Handle("/static/vue.js", vue(vueData))

	// serve bootstrap
	bootstrapData, err := unpack(bootstrapCompressedData)
	if err != nil {
		log.Fatal("decompressing bootstrap error:", err)
	}
	router.Handle("/static/bootstrap.css", bootstrap(bootstrapData))

	// serve bootstrap map
	bootstrapMapData, err := unpack(bootstrapMapCompressedData)
	if err != nil {
		log.Fatal("decompressing bootstrap map error:", err)
	}
	router.Handle("/static/bootstrap.css.map",
		bootstrapMap(bootstrapMapData))
}
