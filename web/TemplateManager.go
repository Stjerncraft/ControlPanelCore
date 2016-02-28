package web

import (
	"io/ioutil"
	"path/filepath"
	"os"
"strings"
	"html/template"
	"fmt"
)

type TemplateManager struct {
	templateList map[string]*clientTemplate
	currentTemplate *clientTemplate
}

func NewTemplateManager() *TemplateManager {
	return new(TemplateManager)
}

func (manager *TemplateManager) ReloadTemplates() error {
	//Get list of templates, and load the current one.
	//Template directory = name
	//Current template if chosen in config
	//Default to "default", or the first in list if it does not exist

	manager.templateList = make(map[string]*clientTemplate)

	//Load list of templates, folder name = template name
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			newTemplate := NewTemplate(file.Name())
			manager.templateList[newTemplate.GetName()] = newTemplate
		}
	}

	//Parse the template files
	for name, tem := range manager.templateList {
		files := make([]string, 0)
		err := filepath.Walk("./templates/" + name, func(path string, f os.FileInfo, err error) error {
			if strings.HasSuffix(f.Name(), ".html") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return err
		}

		tem.templateMain, err = template.ParseFiles(files...)
		if err != nil {
			return err
		}

		//Verify we got the required pages
		//TODO: Rewrite into a 'check' function
		if tem.templateMain.Lookup("main.html") == nil {
			fmt.Println("Template is missing main.html, Unloading ", tem.GetName())
			delete(manager.templateList, tem.GetName())
			continue
		}
		if tem.templateMain.Lookup("login.html") == nil {
			fmt.Println("Template is missing login.html, Unloading ", tem.GetName())
			delete(manager.templateList, tem.GetName())
			continue
		}
	}

	//TODO: Read current template from config
	for _, tmpl := range manager.templateList {
		manager.currentTemplate = tmpl
		break;
	}

	return nil
}

func (manager *TemplateManager) GetCurrentTemplate() *clientTemplate {
	return manager.currentTemplate
}
