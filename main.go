package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var teaSites = []string{
	"https://verdanttea.com/",
}

var teaTypes = []string{
	`(?i)\bpu'er\b`,
	`(?i)\bpuer\b`,
	`(?i)\bpu 'er\b`,
	`(?i)\bpu er\b`,
	`(?i)\bpu-er\b`,
	`(?i)\bchai\b`,
	`(?i)\bmatcha\b`,
	`(?i)\brooibos\b`,
	`(?i)\boolong\b`,
	`(?i)\bblack\b`,
	`(?i)\bwhite\b`,
	`(?i)\bgrean\b`,
	`(?i)\bherbal\b`,
	`(?i)\byellow\b`,
	`(?i)\bfermented\b`,
}

func Match(pattern string, in string) bool {
	return regexp.MustCompile(pattern).MatchString(in)
}

func MatchStart(substring string, in string) bool {
	return regexp.MustCompile(`(?i)^`+substring).MatchString(in) ||
		regexp.MustCompile(`(?i)^`+in).MatchString(substring)
}

var tags = []string{
	"a",
}

type Tea struct {
	MaybeTea
	data string
}

type MaybeTea struct {
	name string
	link string
}

func (t *MaybeTea) Convert(name string, data string) *Tea {
	return &Tea{
		MaybeTea{
			name,
			t.link,
		},
		data,
	}
}

func (t *MaybeTea) ConfirmConvertTeaType() (*Tea, bool) {
	doc := t.GetDocument()
	headers := doc.Find("h1").FilterFunction(func(i int, node *goquery.Selection) bool {
		title := node.Text()
		return MatchStart(title, t.name)
	})

	if headers.Length() == 1 {
		header := headers.First().Text()
		data := doc.Text()

		// in the case that the previously found name has extra stuff on the end
		// and assuming that the header will only contain the name
		return t.Convert(header, data), true
	}

	return &Tea{MaybeTea{"", ""}, ""}, false
}

func (t *MaybeTea) GetDocument() *goquery.Document {
	doc, err := goquery.NewDocument(t.link)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

/*
	The crawler must be able to keep track of information about the user flow
	spawn multiple crawlers?
	how to find matching data between pages?
	if a link's labeled as a certain thing which is then found in a header on the page it leads to, that's a tea
	but it's not always that easy, sometimes the descriptor is not a link
	but the link is only near the descriptor
	how to determine that a link and a descriptor go together?
	either descriptor is child to the link or the link and the descriptor have a common parent
	the href and the descriptor will overlap in some way
	confirm that tea type is real by getting next page in flow
	confirm that it's in the page's data
*/
type Crawler struct {
	links       []string
	seen        map[string]bool
	seenTeas    map[string]bool
	possibleTea []*MaybeTea
	tea         []*Tea
	data        []string
}

func (t *Crawler) GetNextLink() string {
	totalLinks := len(t.links)

	if totalLinks > 0 {
		next := t.links[0:1][0]

		if totalLinks > 1 {
			t.links = t.links[1:]
		} else {
			t.links = make([]string, 0)
		}

		if !t.seen[next] {
			t.seen[next] = true
			return next
		}
		return t.GetNextLink()
	}
	return ""
}

func (t *Crawler) ScrapeSites() *Crawler {
	nextLink := t.GetNextLink()
	fmt.Println("nextLink", nextLink)

	if nextLink != "" {
		doc, err := goquery.NewDocument(nextLink)
		if err != nil {
			fmt.Println("There was an error while getting the document for this link: %s", nextLink)
			return t.ScrapeSites()
		}

		t.ScrapePage(doc)
		return t.ScrapeSites()
	}

	fmt.Println("done")
	return t
}

func (t *Crawler) AddMaybeTea(link string, name string) *MaybeTea {
	tea := &MaybeTea{
		name,
		link,
	}

	t.possibleTea = append(t.possibleTea, tea)

	return tea
}

func (t *Crawler) ScrapePage(doc *goquery.Document) *Crawler {
	for _, teaType := range teaTypes {
		for _, tag := range tags {
			found := doc.Find(tag).FilterFunction(func(i int, node *goquery.Selection) bool {
				text := node.Text()
				return Match(teaType, text)
			})

			found.Each(func(i int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if exists {
					t.links = append(t.links, href)
				}

				text := s.Text()

				t.AddMaybeTea(href, text)
				t.ProcessMaybes()
				t.data = append(t.data, text)
			})

		}
	}

	return t

}

func (t *Crawler) ProcessMaybes() {
	fmt.Println("processmaybes")
	total := len(t.possibleTea)

	if total > 0 {
		next := t.possibleTea[0:1][0]
		tea, converted := next.ConfirmConvertTeaType()

		if converted {
			fmt.Println("tea name: %s", tea.name)
			t.tea = append(t.tea, tea)
		}

		if total > 1 {
			t.possibleTea = t.possibleTea[1:]
			t.ProcessMaybes()
		} else {
			t.possibleTea = make([]*MaybeTea, 0)
		}

	}
}

/*
	visit site
	look for teaTypes
	if hyperlink, crawl if not seen
	if not hyperlink, save for language processing
*/

func ScrapeSite() {
	tg := Crawler{
		teaSites,
		make(map[string]bool),
		make(map[string]bool),
		make([]*MaybeTea, 0),
		make([]*Tea, 0),
		make([]string, 0),
	}

	tg.ScrapeSites()
}

func main() {
	ScrapeSite()
}
