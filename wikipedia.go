package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"slidemonk/slidenlp/smlib"
)

var (
	wikiToHtmlUrl = "http://localhost"
)

type page struct {
	Title string `xml:"title"`
	Text  string `xml:"revision>text"`
}

func usage() {
	fmt.Println("wikipedia <wikipedia dump file name>")
}

func main() {
	url := flag.String("url", "http://localhost", "wikitext to HTML service URL")
	port := flag.Int("port", 3000, "wikitext to HTML service port")
	ignoreErrors := flag.Bool("ignore-errors", true,
		"ignore and continue on errors")
	flag.Parse()

	// parse rest of the flags
	tail := flag.Args()
	var filename string
	if len(tail) != 0 {
		filename = tail[0]
	}
	if len(filename) == 0 {
		usage()
		return
	}

	wikiToHtmlUrl := *url + ":" + strconv.Itoa(*port) + "/parse"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	r, err := smlib.Decompressor(f)
	if err != nil {
		log.Fatalln(err)
	}

	decoder := xml.NewDecoder(r)
	done := false
	for !done {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "page" {
				var p page
				decoder.DecodeElement(&p, &se)
				// parse the page and store in DB
				fmt.Println("Parsing:", p.Title)
				text, err := parse(p.Text, wikiToHtmlUrl)
				if err != nil {
					log.Println(err)
					dumpPage(p)
					if !*ignoreErrors {
						done = true
					}
				} else {
					// store in db?
					fmt.Println(string(text))
				}
			}
		}
	}
}

// parse document in wikitext format to HTML
//
// @doc: docunent in wikitext format
// @url: URL of wikitext to HTML conversion service
//
func parse(doc string, url string) ([]byte, error) {
	buf := bytes.NewBufferString(doc)
	resp, err := http.Post(url, "text/plain", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func dumpPage(p page) {
	log.Println("Error parsing page:", p.Title)
	/*log.Println("Dumping page:")
	  log.Println("TITLE:\n", p.Title)
	  log.Println("TEXT:\n", p.Text)*/

}
