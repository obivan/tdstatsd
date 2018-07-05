package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

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
	// If we see the empty line, then we do not have data
	// it's not an error, just return an empty slice
	if scanner.Text() == "" {
		return []TDPool{}, nil
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
		return []TDPool{}, fmt.Errorf("scanning error: %s", err)
	}
	return parseLines(poolLines)
}

func parseLines(lines []string) ([]TDPool, error) {
	pools := make([]TDPool, 0)
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
