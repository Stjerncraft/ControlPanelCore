package web

import (
	"html/template"
	"io"
)

type clientTemplate struct {
	name string
	templateMain *template.Template
}

type templateData struct {
	CoreModule template.HTML
}

func NewTemplate(name string) *clientTemplate {
	newTemplate := new(clientTemplate)
	newTemplate.name = name

	return newTemplate
}

func (tem *clientTemplate) GetName() string {
	return tem.name
}

func (tem *clientTemplate) WriteMain(wr io.Writer) error {
	data := &templateData{CoreModule: "<script data-main='module/core/sc-cpm-core' src='module/core/lib/require.js'></script>"}

	err := tem.templateMain.ExecuteTemplate(wr, "main.html", data)
	if err != nil {
		return err;
	}

	return nil
}

func (tem *clientTemplate) WriteLogin() error {
	return nil
}
