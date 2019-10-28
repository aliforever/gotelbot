package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const mainFilePath string = "%s/"

func (t Template) initMainFile(botPath, applicationName string) (err error) {
	folder := fmt.Sprintf(mainFilePath, botPath)
	var str string
	goPathSrc := botPath[strings.Index(botPath, "src")+len("src")+1:]
	str, err = TemplateData{}.GetMainFile(goPathSrc, applicationName)
	if err != nil {
		return
	}
	path := folder + "main.go"
	err = ioutil.WriteFile(path, []byte(str), os.ModePerm)
	if err != nil {
		return
	}
	return
}
