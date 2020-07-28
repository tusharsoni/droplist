package content

import (
	"html/template"
	"strings"

	"github.com/tusharsoni/copper/cerror"
)

func generatePreviewHTML(uuid, templateHTML string, params map[string]interface{}) (string, error) {
	var html strings.Builder

	tmpl, err := template.New(uuid).Funcs(template.FuncMap{
		"trackURL": func(href string) string {
			return href
		},
	}).Parse(templateHTML)
	if err != nil {
		return "", cerror.New(err, "failed to parse html body", nil)
	}

	err = tmpl.Execute(&html, params)
	if err != nil {
		return "", cerror.New(err, "failed to execute email template", nil)
	}

	return html.String(), nil
}
