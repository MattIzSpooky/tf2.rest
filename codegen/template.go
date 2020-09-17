package codegen

import (
	"github.com/MattIzSpooky/tf2.rest/responses"
	"strings"
	"text/template"
	"time"
)

type ResponseTemplate struct {
	Timestamp time.Time
	URL       string
	Class     string
	Responses []responses.Response
}

const tmplContent = `// Code generated by go generate; DO NOT EDIT.
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
		Condition: ` + "`" + "{{ .Condition }}" + "`"+`,
	},
	{{ end }}
{{- end }}
}
`

func NewResponseTemplate() *template.Template  {
	tmpl := template.New("responseTemplate")

	tmpl.Funcs(template.FuncMap{
		"ToUpper": strings.ToUpper,
	})

	return template.Must(tmpl.Parse(tmplContent))

}