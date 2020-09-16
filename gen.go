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
		Class: {{ .Class | ToUpper }},
		Response: "{{ .Response }}",
		AudioFile: "{{ .AudioFile }}",
		Context: "{{ .Context }}",
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
		_type := document.Find("#Kill-related_responses").Text()

		table := document.Find(".headertemplate").First()
		context := table.Find("#mw-content-text > table:nth-child(8) > tbody > tr:nth-child(1) > td > b").Text()

		var responseSlice []responses.Response
		var condition string

		table.Children().Find("li").Each(func(i int, selection *goquery.Selection) {
			hasCondition := selection.Find("ul")

			// Condition is found
			if (hasCondition.Length() == 1) {
				condition = strings.TrimSpace(selection.Contents().Not("ul").Text())
				return
			}

			response := strings.ReplaceAll(selection.Text(), `"`, ``)

			audioURI, _ := selection.Children().First().Attr("href")

			responseSlice = append(responseSlice, responses.Response{
				Class:     class,
				Response:  response,
				AudioFile: baseUrl + audioURI,
				Type:      _type,
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

		fmt.Println(fmt.Sprintf("Generated responses for: %s", class))
	}
}
