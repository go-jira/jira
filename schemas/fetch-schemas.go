package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func mayPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	resp, err := http.Get("https://docs.atlassian.com/software/jira/docs/api/REST/7.12.0/")
	mayPanic(err)
	defer resp.Body.Close()

	capture := false
	nextCodeBlock := false
	stream := html.NewTokenizer(resp.Body)
	var buffer string
	for tt := stream.Next(); tt != html.ErrorToken; tt = stream.Next() {
		t := stream.Token()
		if t.Data == "h6" && tt == html.StartTagToken {
			capture = true
		}
		if t.Data == "h6" && tt == html.EndTagToken {
			capture = false
			if strings.Contains(strings.ToLower(buffer), "schema") {
				nextCodeBlock = true
			}
			buffer = ""
		}
		if nextCodeBlock && t.Data == "code" && tt == html.StartTagToken {
			capture = true
		}
		if nextCodeBlock && t.Data == "code" && tt == html.EndTagToken {
			capture = false
			schema := map[string]interface{}{}
			err := json.Unmarshal([]byte(buffer), &schema)
			mayPanic(err)
			title, ok := schema["title"].(string)
			if ok {
				title = strings.ReplaceAll(title, " ", "")
				fileName := fmt.Sprintf("%s.json", title)
				fmt.Printf("Writing %s\n", fileName)
				err := ioutil.WriteFile(fileName, []byte(buffer), 0644)
				mayPanic(err)
			}
			buffer = ""
			nextCodeBlock = false
		}
		if capture && tt == html.TextToken {
			buffer += t.Data
		}
	}
}
