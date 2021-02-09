package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getVersion(document Document) []byte {
	// version: block height + - + timestamp (milliseconds, 13 digits) + deleted
	var delStatus byte
	delStatus = 48
	if document.Deleted {
		delStatus = 49
	}
	versionGen := append([]byte("0-"+document.Timestamp), delStatus)

	return versionGen
}

func getReplID(sNodeID string, nodeID string, databaseName string) string {
	stringToHash := sNodeID + nodeID + databaseName
	hasher := sha1.New()
	hasher.Write([]byte(stringToHash))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
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
		return 500, nil

	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return res.StatusCode, body

}
