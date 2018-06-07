package main

import (
	"flag"
	"fmt"
	"html/template"
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
	indexTpl := template.Must(
		template.New("index").Parse(indexPageTemplate))
	http.Handle("/", index(tdURL, client, indexTpl))
	log.Println("Listening on", bindAddr)
	log.Fatal(http.ListenAndServe(bindAddr, nil))
}

func index(url string, c *http.Client, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound),
				http.StatusNotFound)
			return
		}
		data, err := getTdData(url, c)
		internalServerError := func(err error) {
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
		w.WriteHeader(http.StatusOK)
		if err := tpl.Execute(w, pools); err != nil {
			internalServerError(err)
			return
		}
	})
}

func getTdData(url string, c *http.Client) ([]byte, error) {
	resp, err := c.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error %s during %s url processing",
			err, url)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("dude, wtf.. error closing body: %s", cerr)
		}
	}()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("can't read response body: %s", err)
	}
	return buf, nil
}
