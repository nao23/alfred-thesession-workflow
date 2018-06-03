package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"github.com/deanishe/awgo"
	"github.com/PuerkitoBio/goquery"
)

var	(
	wf         *aw.Workflow
	helpUrl    = "http://www.deanishe.net"
	maxResults = 10
)

const BASE_URL   = "http://thesession.org"

func GetHTML(target string, keyword string) *goquery.Document {
	// Request the HTML page.
	values := url.Values{}
	values.Add("q", keyword)
	res, err := http.Get(BASE_URL + "/" + target + "/search?" + values.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func search(target string, keyword string) {
	doc := GetHTML(target, keyword)
	doc.Find("li.manifest-item").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s %s", i+1, titleAnchor.Text(), titleAnchor.Next().Text())
		href, _ := titleAnchor.Attr("href")
		href = BASE_URL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func init() {
	wf = aw.New(aw.HelpURL(helpUrl), aw.MaxResults(maxResults))
}

func run() {
	allTargets := []string{"tunes", "recordings"}

	// Parse cmd-line arguments
	target  := flag.String("target", "tunes", "search target")
	keyword := flag.String("keyword", "Kesh", "keyword")
	flag.Parse()

	search(*target, *keyword)

	wf.WarnEmpty("No matching target", "Try: " + strings.Join(allTargets, ","))
	wf.SendFeedback()
}

func main()  {
	aw.Run(run)
}
