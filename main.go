package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var teaSites = [...]string{
	"https://verdanttea.com/",
}

var teaTypes = [...]string{
	"oolong",
	"black",
	"white",
	"grean",
	"herbal",
	"yellow",
	"fermented",
}

var tags = [...]string{
	"span",
	"div",
	"p",
	"text",
	"li",
	"button",
	"a",
	"select",
	"option",
	"h1",
	"h2",
	"h3",
	"h4",
	"article",
}

/*
	visit site
	look for teaTypes
	if hyperlink, crawl if not seen
	if not hyperlink, save for language processing
*/

func ExampleScrape() {
	doc, err := goquery.NewDocument(teaSites[0])
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	found := doc.Find("a").FilterFunction(func(i int, node *goquery.Selection) bool {
		text := node.Text()

		matched, err := regexp.MatchString("Oolong", text)
		if err != nil {
			log.Fatal(err)
		}

		return matched
	})

	found.Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		fmt.Printf("href: %s, text: %s ", href, s.Text())
	})

}

func main() {
	ExampleScrape()
}
