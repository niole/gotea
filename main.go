package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var teaSites = [1]string{
	"https://verdanttea.com/",
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
}

func ExampleScrape() {
	doc, err := goquery.NewDocument(teaSites[0])
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	found := doc.Find("span").FilterFunction(func(i int, node *goquery.Selection) bool {
		text := node.Text()

		matched, err := regexp.MatchString("oolong", text)
		if err != nil {
			log.Fatal(err)
		}

		return matched
	})

	fmt.Printf("%s", found)

	found.Each(func(i int, s *goquery.Selection) {
		fmt.Printf(s.Text())
	})

}

func main() {
	ExampleScrape()
}
