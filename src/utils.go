package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getVersion(document Document) []byte {
	// version: timestamp (milliseconds, 13 digits) + deleted
	var delStatus byte
	delStatus = 48
	if document.Deleted {
		delStatus = 49
	}
	versionGen := append([]byte(document.Timestamp), delStatus)

	return versionGen
}

func request(url string, method string, payload string, contentType string) (int, []byte) {

	client := &http.Client{
		Jar: jar,
	}

	var req *http.Request
	var err error

	if payload == "" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(payload))
	}

	if err != nil {
		fmt.Println(err)
	}

	if contentType == "x-www-form-urlencoded" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else if contentType == "application/json" {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return res.StatusCode, body
}
