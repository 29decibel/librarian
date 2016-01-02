package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// AlfredItem is a simple item
type AlfredItem struct {
	XMLName  xml.Name `xml:"item"`
	Title    string   `xml:"title,omitempty"`
	Subtitle string   `xml:"subtitle,omitempty"`
	Arg      string   `xml:"arg,omitempty"`
	Icon     string   `xml:"icon,omitempty"`
}

// AlfredItems is
type AlfredItems struct {
	XMLName xml.Name     `xml:"items"` // root of this element
	Items   []AlfredItem `xml:"items"`
}

func trimTrim(input string) string {
	return strings.Replace(strings.TrimSpace(input), "\n", "", -1)
}

// SearchItems is
func SearchItems(query string) {
	doc, err := goquery.NewDocument(fmt.Sprintf("https://smcl.bibliocommons.com/search?locale=en-US&t=smart&formats=BK&q=%s", url.QueryEscape(query)))
	if err != nil {
		log.Fatal(err)
	}

	alfredItems := AlfredItems{}

	r, err := regexp.Compile(`\d`)

	doc.Find(".cp_bib_list .listItem").Each(func(i int, s *goquery.Selection) {
		link := s.Find("span.title a")
		href, _ := link.Attr("href")
		title := link.Text()
		subtitle := s.Find(".subTitle").Text()
		yearInfo := s.Find(".hidden-md .format").Text()
		year := strings.Join(r.FindAllString(yearInfo, -1), "")

		if len(subtitle) == 0 {
			subtitle = title
		}

		if len(title) > 0 {
			item := AlfredItem{Title: fmt.Sprintf("[%s]-%s", string(year), title), Arg: fmt.Sprintf("https://smcl.bibliocommons.com%s", href), Subtitle: trimTrim(subtitle), Icon: "books.png"}
			alfredItems.Items = append(alfredItems.Items, item)
		}
	})

	buf, _ := xml.Marshal(alfredItems)
	fmt.Printf("<xml>%s</xml>", string(buf))
}

func main() {
	args := os.Args[1:len(os.Args)]
	SearchItems(strings.Join(args, " "))
}
