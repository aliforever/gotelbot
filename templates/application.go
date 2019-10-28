package templates

import (
	"fmt"
	"io/ioutil"
	"os"

	"strings"

	"github.com/gobuffalo/flect"
)

const applicationFilePath string = "%s/application/"

func (t Template) initApplicationFiles(botPath, botUsername string, languages []string) (err error) {
	folder := fmt.Sprintf(applicationFilePath, botPath)
	err = t.makeDirectory(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
	var str string
	goPathSrc := botPath[strings.Index(botPath, "src")+len("src")+1:]
	str, err = TemplateData{}.GetApplicationFile(goPathSrc, flect.Capitalize(botUsername), languages)
	if err != nil {
		return
	}
	path := folder + "application.go"
	err = ioutil.WriteFile(path, []byte(str), os.ModePerm)
	if err != nil {
		return
	}
	str, err = TemplateData{}.GetMethodsFile(goPathSrc)
	if err != nil {
		return
	}
	path = folder + "methods.go"
	err = ioutil.WriteFile(path, []byte(str), os.ModePerm)
	if err != nil {
		return
	}
	// path = folder + "methods.go"
	return
}
