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
	class 	 string
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

	s.document.Find(".headertemplate").Each(func(i int, table *goquery.Selection) {
		previousNode := table.Prev()
		prvNodeName := goquery.NodeName(previousNode)

		if prvNodeName == "h2" && goquery.NodeName(previousNode.Children().First()) == "span" {
			tempType := previousNode.Text()

			if tempType != _type {
				subType = ""
			}

			_type = tempType
		} else if prvNodeName == "h3" {
			beforeThatElement := previousNode.Prev()
			if goquery.NodeName(beforeThatElement) == "h2" {
				_type = beforeThatElement.Text()
			}

			subType = strings.TrimSpace(previousNode.Text())
		} else if prvNodeName == "div" {
			beforeThatElement := previousNode.Prev().Prev()

			if goquery.NodeName(beforeThatElement) == "h2" {
				_type = beforeThatElement.Text()
			}
			subType = strings.TrimSpace(previousNode.Prev().Text())
		}

		_type = strings.TrimSpace(strings.ReplaceAll(_type, "responses", ""))

		context := table.Find("td > b").Text()

		var condition string

		table.Children().Find("li").Each(func(i int, listElement *goquery.Selection) {
			conditionElement := listElement.Find("ul")

			if conditionElement.Length() == 1 {
				tempCondition := strings.TrimSpace(listElement.Contents().Not("ul").Text())

				if tempCondition == "Rare" {
					context = context + " (" + tempCondition + ")"
				} else {
					condition = tempCondition
				}
				return
			}

			response := strings.ReplaceAll(listElement.Text(), `"`, ``)

			audioURI, _ := listElement.Children().First().Attr("href")

			responseSlice = append(responseSlice, responses.Response{
				Class:     s.class,
				Response:  response,
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
