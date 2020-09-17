package codegen

import (
	"errors"
	"fmt"
	"github.com/MattIzSpooky/tf2.rest/responses"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

type Scraper struct {
	url      string
	class    string
	client   *http.Client
	document *goquery.Document
}

const BaseUrl = "https://wiki.teamfortress.com"
const url = BaseUrl + "/wiki/%s_responses"
const timeout = 30 * time.Second

func (s *Scraper) GetURL() string {
	return s.url
}

func NewScraper(class string) *Scraper {
	scraper := &Scraper{
		class: class,
		client: &http.Client{
			Timeout: timeout,
		},
	}

	return scraper
}

func (s *Scraper) FetchDocument() error {
	responsePageUrl := fmt.Sprintf(url, s.class)
	s.url = responsePageUrl

	rsp, err := s.client.Get(responsePageUrl)

	if err != nil {
		return err
	}

	document, err := goquery.NewDocumentFromReader(rsp.Body)

	if err != nil {
		return err
	}

	s.document = document

	return rsp.Body.Close()
}

func (s Scraper) Run() ([]responses.Response, error) {
	var responseSlice []responses.Response
	var _type string
	var subType string

	tables := s.document.Find(".headertemplate")

	if tables.Length() == 0 {
		return nil, errors.New(fmt.Sprintf("could not scrape, no tables found for class: %s", s.class))
	}

	tables.Each(func(i int, table *goquery.Selection) {
		previousNode := table.Prev()
		prvNodeName := goquery.NodeName(previousNode)

		// The checks below are horrible. I am well aware.
		// Gotta do this nasty bit to get data reliably...

		// A title
		if prvNodeName == "h2" && goquery.NodeName(previousNode.Children().First()) == "span" {
			tempType := previousNode.Text()

			// Reset subType when moving to another type.
			if tempType != _type {
				subType = ""
			}

			_type = tempType
		} else if prvNodeName == "h3" { // A subtitle

			// There is usually a h2 before a h3
			beforeThatElement := previousNode.Prev()
			if goquery.NodeName(beforeThatElement) == "h2" {
				_type = beforeThatElement.Text()
			}

			subType = strings.TrimSpace(previousNode.Text())
		} else if prvNodeName == "div" { // Text below a subtitle is just a div??
			// Move up 2 elements to get the right text.
			beforeThatElement := previousNode.Prev().Prev()

			if goquery.NodeName(beforeThatElement) == "h2" {
				_type = beforeThatElement.Text()
			}
			subType = strings.TrimSpace(previousNode.Prev().Text())
		}

		// Remove "responses" from texts such as "Objective-related responses".
		_type = strings.TrimSpace(strings.ReplaceAll(_type, "responses", ""))
		subType = strings.TrimSpace(strings.ReplaceAll(subType, "responses", ""))

		context := table.Find("td > b").Text()

		var condition string

		table.Children().Find("li").Each(func(i int, listElement *goquery.Selection) {
			conditionElement := listElement.Find("ul") // The condition is sometimes a parent list of sorts.
			otherConditionElement := listElement.Parent().Prev() // The condition can also be just a random paragraph it seems. Nice consistency you got there wiki.

			if goquery.NodeName(otherConditionElement) == "p" {
				cond := strings.TrimSpace(otherConditionElement.Text())

				// Ignore notes.
				if !strings.Contains(cond, "Note:") {
					condition = strings.TrimSuffix(cond, `'`) // Some thingies end with a single quote for some reason.
				}
			}

			if conditionElement.Length() == 1 {
				condition = strings.TrimSpace(listElement.Contents().Not("ul").Text())
				return
			}

			if goquery.NodeName(conditionElement.Prev()) == "p" {
				condition = strings.TrimSpace(conditionElement.Prev().Contents().Not("ul").Text())
				return
			}

			response := strings.ReplaceAll(listElement.Text(), `"`, ``)

			audioURI, _ := listElement.Children().First().Attr("href")

			responseSlice = append(responseSlice, responses.Response{
				Class:     s.class,
				Response:  strings.TrimSpace(response),
				AudioFile: BaseUrl + audioURI,
				Type:      _type,
				SubType:   subType,
				Context:   context,
				Condition: condition,
			})
		})
	})

	if len(responseSlice) == 0 {
		return nil, errors.New(fmt.Sprintf("no results were found for %s", s.class))
	}

	return responseSlice, nil
}
