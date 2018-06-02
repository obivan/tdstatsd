package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
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
		showVersion bool
	)
	flag.StringVar(
		&tdURL,
		"url",
		"http://localhost:8000/_perf.txt",
		"Traffic Director stats page URL",
	)
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
	log.Println("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
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

func parse(data []byte) ([]TDPool, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		// Scan until we see header
		if scanner.Text() == "Origin server statistics (for http):" {
			break
		}
	}
	// then consume next 4 lines
	for i := 0; i < 4; i++ {
		scanner.Scan()
	}
	// If we see the next header, then we do not have data
	// it's not an error, just return an empty slice
	if scanner.Text() == "Origin server statistics (for tcp):" {
		return nil, nil
	}
	// now we see 1st line of pools. store them
	var poolLines []string
	if p := scanner.Text(); p != "" {
		poolLines = append(poolLines, p)
	}
	for scanner.Scan() {
		// and scan until whitespace line
		if scanner.Text() != "" {
			poolLines = append(poolLines, scanner.Text())
		} else {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning error: %s", err)
	}
	return parseLines(poolLines)
}

func parseLines(lines []string) ([]TDPool, error) {
	var pools []TDPool
	for _, line := range lines {
		scanner := bufio.NewScanner(strings.NewReader(line))
		scanner.Split(bufio.ScanWords)
		// advance and return scanned token
		scanWord := func() string {
			scanner.Scan()
			return scanner.Text()
		}
		pool := TDPool{
			Name:   scanWord(),
			URL:    scanWord(),
			Status: scanWord(),
		}
		pools = append(pools, pool)
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("line scanning error: %s", err)
		}
	}
	return pools, nil
}
