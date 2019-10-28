package templates

import (
	"fmt"
	"io/ioutil"
	"os"
)

const apiFilePath string = "%s/api/"

func (t Template) makeApiFile(botPath string) (err error) {
	folder := fmt.Sprintf(apiFilePath, botPath)
	err = t.makeDirectory(folder)
	if err != nil {
		return
	}
	path := folder + "api.go"
	err = ioutil.WriteFile(path, []byte(TemplateData{}.apiFile()), os.ModePerm)
	return
}
