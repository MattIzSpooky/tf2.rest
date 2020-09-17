// +build ignore

package main

import (
	"fmt"
	"github.com/MattIzSpooky/tf2.rest/responses"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const url = "https://wiki.teamfortress.com/wiki/%s_responses"
const baseUrl = "https://wiki.teamfortress.com"

var responseTemplate = template.Must(
	template.New("responseTemplate").
		Funcs(template.FuncMap{
			"ToUpper": strings.ToUpper,
		}).
		Parse(
			`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
// using data from
// {{ .URL }}
package responses

var {{ .Class }}Responses = []Response{
{{- with .Responses }}
	{{ range . }}
	{
		Class: 		{{ .Class | ToUpper }},
		Response:  "{{ .Response }}",
		AudioFile: "{{ .AudioFile }}",
		Type: 	   "{{ .Type }}",
		SubType:   "{{ .SubType }}",
		Context:   "{{ .Context }}",
		Condition: "{{ .Condition }}",
	},
	{{ end }}
{{- end }}
}
`))

var classes = [...]string{
	responses.SCOUT,
	responses.SOLDIER,
	responses.PYRO,
	responses.DEMOMAN,
	responses.HEAVY,
	responses.ENGINEER,
	responses.MEDIC,
	responses.SNIPER,
	responses.SPY,
}

func main() {
	for _, class := range classes {
		os.Remove(fmt.Sprintf("responses/%s.go", class)) // TODO: remove after code is done

		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		responsePageUrl := fmt.Sprintf(url, class)
		rsp, err := client.Get(responsePageUrl)
		if err != nil {
			panic(err)
		}

		defer rsp.Body.Close()

		document, err := goquery.NewDocumentFromReader(rsp.Body)
		if err != nil {
			panic(err)
		}
		var responseSlice []responses.Response
		var _type string
		var subType string

		document.Find(".headertemplate").Each(func(i int, table *goquery.Selection) {
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
				if (goquery.NodeName(beforeThatElement) == "h2") {
					_type = beforeThatElement.Text()
				}

				subType = strings.TrimSpace(previousNode.Text())
			} else if (prvNodeName == "div") {
				beforeThatElement := previousNode.Prev().Prev()
				if (goquery.NodeName(beforeThatElement) == "h2") {
					_type = beforeThatElement.Text()
				}
				subType = strings.TrimSpace(previousNode.Prev().Text())
			}

			_type = strings.TrimSpace(strings.ReplaceAll(_type, "responses", ""))

			fmt.Println(_type)
			context := table.Find("td > b").Text()

			var condition string

			table.Children().Find("li").Each(func(i int, listElement *goquery.Selection) {
				hasCondition := listElement.Find("ul")

				// Condition is found
				if (hasCondition.Length() == 1) {
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
					Class:     class,
					Response:  response,
					AudioFile: baseUrl + audioURI,
					Type:      _type,
					SubType:   subType,
					Context:   context,
					Condition: condition,
				})
			})

			f, err := os.Create(fmt.Sprintf("responses/%s.go", class))

			if err != nil {
				panic(err)
			}
			defer f.Close()

			responseTemplate.Execute(f, struct {
				Timestamp time.Time
				URL       string
				Responses []responses.Response
				Class     string
			}{
				Timestamp: time.Now(),
				URL:       responsePageUrl,
				Responses: responseSlice,
				Class:     class,
			})
		})

		fmt.Println(fmt.Sprintf("Generated responses for: %s", class))
	}
}