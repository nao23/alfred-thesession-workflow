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

func getHTML(target string, keyword string) *goquery.Document {
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

func searchTunes(target string, keyword string) {
	doc := getHTML(target, keyword)
	doc.Find("li.manifest-item").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s %s", i+1, titleAnchor.Text(), titleAnchor.Next().Text())
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func searchRecordings(target string, keyword string) {
	doc := getHTML(target, keyword)
	doc.Find("li.manifest-item").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s by %s", i+1, titleAnchor.Text(), titleAnchor.Next().Text())
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func searchSessions(target string, keyword string) {
	doc := getHTML(target, keyword)
	doc.Find("li.manifest-item:not(:has(del))").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s", i+1, strings.Trim(titleAnchor.Parent().Text(), "\n"))
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func searchEvents(target string, keyword string) {
	doc := getHTML(target, keyword)
	doc.Find("li.manifest-item").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s (%s)", i+1, titleAnchor.Text(), s.Find("time").Text())
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func searchDiscussions(target string, keyword string) {
	doc := getHTML(target, keyword)
	doc.Find("li.manifest-item").Each(func(i int, s *goquery.Selection) {
		titleAnchor := s.Find("a.manifest-item-title")
		titleStr := fmt.Sprintf("%d. %s %s", i+1, titleAnchor.Text(), s.Find("span.manifest-item-extra").Text())
		href, _ := titleAnchor.Attr("href")
		href = baseURL + href
		wf.NewItem(titleStr).Subtitle(href).Arg(href).Valid(true)
	})
}

func run() {
	// Parse cmd-line arguments
	target := flag.String("target", "tunes", "search target")
	keyword := flag.String("keyword", "Kesh", "keyword")
	flag.Parse()

	switch *target {
	case "tunes":
		searchTunes(*target, *keyword)
	case "recordings":
		searchRecordings(*target, *keyword)
	case "sessions":
		searchSessions(*target, *keyword)
	case "events":
		searchEvents(*target, *keyword)
	case "discussions":
		searchDiscussions(*target, *keyword)
	default:
	}

	wf.WarnEmpty("No matching "+(*target), "Try a different keyword?")
	wf.SendFeedback()
}

func init() {
	wf = aw.New(aw.HelpURL(helpURL), aw.MaxResults(maxResults))
}

func main() {
	aw.Run(run)
}
