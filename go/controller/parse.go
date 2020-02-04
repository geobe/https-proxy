// Parse all templates once and make them globally available
package controller

import (
	"html/template"
	"os"
)

const viewPath = "/view/*.go.html"

var views *template.Template

// map transports values from go code to templates
type Viewmodel map[string]interface{}

func Templates(projectBase string) {
	//pwd, _ := os.Getwd()
	//pwd += "/"
	//var path string
	//info, err := os.Stat(pwd + "view")
	//if os.IsNotExist(err) || !info.IsDir() {
	//	if strings.HasSuffix(pwd, "main/") {
	//		path = pwd + "../view/*.go.html"
	//	} else {
	//		path = pwd + projectBase + viewPath
	//	}
	//} else {
	//	path = pwd + "view/*.go.html"
	//}
	path := ResourceBase(projectBase) + viewPath
	// create empty template first to eventually add functions
	t := template.New("proxybase.html")
	// then parse all templates in 'path' directory
	template.Must(t.ParseGlob(path))
	// and hold them in a local variable
	views = t
}

func Views() *template.Template {
	return views
}

/**
return base directory of static resources like styles, images or templates
projectBase	base path of the project
*/
func ResourceBase(projectBase string) string {
	pwd, _ := os.Getwd()
	var path string
	info, err := os.Stat(pwd + "/web")
	// check if web directory is in current working directory
	if os.IsNotExist(err) || !info.IsDir() {
		// if not, we are running in the development environment
		path = pwd + "/" + projectBase + "/web"
	} else {
		// if yes, we are running in a distribution environment
		path = pwd + "/web"
	}
	return path
}
