package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

var teaSites = []string{
	"http://www.adagio.com",
	"https://verdanttea.com",
}

var teaTypes = []string{
	`pu'er`,
	`puer`,
	`pu 'er`,
	`pu er`,
	`pu-er`,
	`chai`,
	`matcha`,
	`rooibos`,
	`oolong`,
	`black`,
	`white`,
	`green`,
	`herbal`,
	`yellow`,
	`fermented`,
}

var relativePathPattern = regexp.MustCompile("^/")
var teaCategoryPattern = strings.Join(teaTypes, " tea ") + " tea" + strings.Join(teaTypes, " teas ") + " teas"
var originPattern = regexp.MustCompile("(https://www..+.com|http://www..+.com)")
var urlDelimeterReplacer = strings.NewReplacer("_", " ", "-", " ", ".", " ", "/", " ", "html", "")
var multilinePattern = regexp.MustCompile("\n")
var tabPattern = regexp.MustCompile("\t")

func ExtractHyperlinkContent(link string) string {
	origin := GetOrigin(link)
	withoutOrigin := strings.Replace(link, origin, "", -1)
	return urlDelimeterReplacer.Replace(withoutOrigin)
}

func RemoveUrlDelmeters(url string) string {
	return urlDelimeterReplacer.Replace(url)
}

func Match(toFind string, in string) bool {
	return regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`, toFind)).MatchString(in)
}

func HasOverlap(a string, b string) bool {
	aPattern := regexp.MustCompile(a)
	bPattern := regexp.MustCompile(b)

	return aPattern.FindString(b) != "" || bPattern.FindString(a) != ""
}

func GetOrigin(url string) string {
	return originPattern.FindString(url)
}

func NormalizeLink(link string, originLink string) string {
	if relativePathPattern.MatchString(link) {
		// normalize the relative path
		extractedOrigin := GetOrigin(originLink)
		return fmt.Sprintf("%s%s", extractedOrigin, link)
	}

	return link
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

func GetText(node *goquery.Selection) string {
	return TrimContent(node.Text())
}

func TrimContent(content string) string {
	text := multilinePattern.ReplaceAllString(content, "")
	text = tabPattern.ReplaceAllString(text, "")
	return strings.Trim(text, " ")
}

/*
	Creates regex that will match content that contains
	phrases that uniquely identify the tea from which the regex
	is derived
*/
func CreateTeaNamePattern(teaName string) *regexp.Regexp {
	// break name into sequential elements
	// use anything that isn't solely a tea category in an "or" pattern

	teaElements := strings.Split(teaName, " ")
	subTeas := make([]string, 0)

	for i := 0; i < len(teaElements)-1; i++ {
		for j := i + 1; j < len(teaElements)+1; j++ {
			nextSubTea := ""
			for k := i; k < j; k++ {
				nextSubTea += fmt.Sprintf(" %s", teaElements[k])
			}

			if !Match(nextSubTea, teaCategoryPattern) {
				subTeas = append(subTeas, strings.Trim(nextSubTea, " "))
			}

		}
	}

	patternString := fmt.Sprintf(`(?i)%s`, strings.Join(subTeas, "|"))
	return regexp.MustCompile(patternString)
}

/*
	Valid document content is a terminal or of height 1
	This should allow for DOM elements that serve as text containers
	rather than structural elements
	Shouldn't contain JS
	Content must contain substrings that may refer to found teaName
*/
func GetFormattedDocContent(doc *goquery.Document, teaName string) string {
	fmt.Println("formatting content")
	doc.Find("script").Remove()
	doc.Find("link").Remove()
	doc.Find("style").Remove()
	teaNamePattern := CreateTeaNamePattern(teaName)

	content := ""
	selection := []*goquery.Selection{doc.Find("body")}

	for len(selection) > 0 {

		nextSelection := selection[0]

		if len(selection) == 1 {
			selection = make([]*goquery.Selection, 0)
		} else {
			selection = selection[1:]
		}

		// if selection height 0 and/or only has children of height zero
		// and match teaName substring, keep
		nextSelection.Each(func(i int, node *goquery.Selection) {
			children := node.Children()
			totalZeroDepthChildren := children.FilterFunction(func(i int, node *goquery.Selection) bool {
				return node.Children().Length() == 0
			}).Length()
			totalChildren := children.Length()

			if totalChildren == 0 || totalChildren == totalZeroDepthChildren {
				text := GetText(node)

				if teaNamePattern.MatchString(text) {
					// no children/shallow child tree
					// at least contains substrings that likely refer to tea name
					content += fmt.Sprintf(" %s", text)
				}

			} else {
				// put on selection stack for processing
				selection = append(selection, children)
			}

		})
	}

	return TrimContent(content)
}

func (t *MaybeTea) ConfirmConvertTeaType() (*Tea, bool) {
	doc, err := t.GetDocument()
	if err == nil {
		headers := doc.Find("h1").FilterFunction(func(i int, node *goquery.Selection) bool {
			content := GetText(node)
			return HasOverlap(content, t.name)
		})

		if headers.Length() == 1 {
			header := GetText(headers.First())
			data := GetFormattedDocContent(doc, header)

			// in the case that the previously found name has extra stuff on the end
			// and assuming that the header will only contain the name
			return t.Convert(header, data), true
		}
	}

	return &Tea{MaybeTea{"", ""}, ""}, false
}

func (t *MaybeTea) GetDocument() (*goquery.Document, error) {
	return goquery.NewDocument(t.link)
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
	db          *DataBase
	links       []string
	seen        map[string]bool
	possibleTea []*MaybeTea
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

		if !t.Visited(next) {
			t.UpdateVisited(next)
			return next
		}
		return t.GetNextLink()
	}
	return ""
}

func (t *Crawler) ScrapeSites() *Crawler {
	nextLink := t.GetNextLink()

	if nextLink != "" {
		doc, err := goquery.NewDocument(nextLink)
		if err != nil {
			fmt.Println("There was an error while getting the document for this link: %s", nextLink)
			return t.ScrapeSites()
		}

		t.ScrapePage(doc, nextLink)
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

func (t *Crawler) ScrapePage(doc *goquery.Document, baseLink string) *Crawler {
	for _, teaType := range teaTypes {
		doc.Find("a").Each(func(i int, node *goquery.Selection) {

			href, exists := node.Attr("href")
			href = NormalizeLink(href, baseLink)

			if exists && !t.Visited(href) {
				content := GetText(node)

				if content == "" {
					// when hyperlink has no content, work with the link
					content = ExtractHyperlinkContent(href)
				}

				if Match(teaType, content) {
					// if basic tea type is found in the hyperlink's relevant content
					// check to see if could be specific tea type

					if Match(content, teaCategoryPattern) {
						// let main crawler handle getting through tea categories and going
						// between sites
						t.links = append(t.links, href)

					} else {
						// let MaybeTea handle more specific tea finding
						t.AddMaybeTea(href, content)
						t.ProcessMaybes()
					}
				}
			}

		})
	}

	return t

}

func (t *Crawler) AddTea(tea *Tea) {
	t.db.AddTea(tea.name, tea.link, tea.data)
}

func (t *Crawler) Visited(link string) bool {
	return t.seen[link]
}

func (t *Crawler) UpdateVisited(link string) {
	t.seen[link] = true
}

/*
	ProcessMaybes handles examining elements that may be specific tea types
	Updates seen in crawler so that the crawler doesn't explore it
*/
func (t *Crawler) ProcessMaybes() {
	total := len(t.possibleTea)

	if total > 0 {
		next := t.possibleTea[0]
		tea, converted := next.ConfirmConvertTeaType()

		if converted {
			t.AddTea(tea)
		}

		t.UpdateVisited(tea.link)

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
func main() {
	tg := Crawler{
		CreateDataBase("root", "root", "127.0.0.1", "3307", "mysql"),
		teaSites,
		make(map[string]bool),
		make([]*MaybeTea, 0),
	}

	tg.ScrapeSites()
}
