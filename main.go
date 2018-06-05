package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/deanishe/awgo"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	wf         *aw.Workflow
	helpURL    = "https://github.com/nao23/alfred-thesession-workflow"
	maxResults = 10
)

const baseURL = "http://thesession.org"

func GetHTML(target string, keyword string) *goquery.Document {
	// Request the HTML page.
	values := url.Values{}
	values.Add("q", keyword)
	res, err := http.Get(baseURL + "/" + target + "/search?" + values.Encode())
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

func SearchTunes(target string, keyword string) {
	doc := GetHTML(target, keyword)
	doc.Find("li.manifest-item").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s %s", i+1, titleAnchor.Text(), titleAnchor.Next().Text())
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func SearchRecordings(target string, keyword string) {
	doc := GetHTML(target, keyword)
	doc.Find("li.manifest-item").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s by %s", i+1, titleAnchor.Text(), titleAnchor.Next().Text())
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func SearchSessions(target string, keyword string) {
	doc := GetHTML(target, keyword)
	doc.Find("li.manifest-item:not(:has(del))").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s", i+1, strings.Trim(titleAnchor.Parent().Text(), "\n"))
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func SearchDiscussions(target string, keyword string) {
	doc := GetHTML(target, keyword)
	doc.Find("li.manifest-item").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s %s", i+1, titleAnchor.Text(), s.Find("span.manifest-item-extra").Text())
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func init() {
	wf = aw.New(aw.HelpURL(helpURL), aw.MaxResults(maxResults))
}

func run() {
	// Parse cmd-line arguments
	target := flag.String("target", "tunes", "search target")
	keyword := flag.String("keyword", "Kesh", "keyword")
	flag.Parse()

	switch *target {
	case "tunes":
		SearchTunes(*target, *keyword)
	case "recordings":
		SearchRecordings(*target, *keyword)
	case "sessions":
		SearchSessions(*target, *keyword)
	case "discussions":
		SearchDiscussions(*target, *keyword)
	default:
	}

	wf.WarnEmpty("No matching "+(*target), "Try a different keyword?")
	wf.SendFeedback()
}

func main() {
	aw.Run(run)
}
