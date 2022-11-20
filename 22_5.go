package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	resp, err := http.Get("https://petition.parliament.uk/petitions")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	src := string(body)

	re := regexp.MustCompile("\\<h2\\>.*\\</h2\\>")
	rHTML := regexp.MustCompile("<[^>]*>")
	titles := re.FindAllString(src, -1)

	for _, title := range titles {
		cleanTitle := rHTML.ReplaceAllString(title, "")
		fmt.Println(cleanTitle)
	}
}
