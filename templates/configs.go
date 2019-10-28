package templates

import (
	"fmt"
	"io/ioutil"
	"os"
)

const configsFilePath string = "%s/configs/"

func (t Template) makeConfigsConstantFile(botUsername, botToken, botPath, webhookUrl, serverPort, supperCommandSign string, adminUserId int) (err error) {
	folder := fmt.Sprintf(configsFilePath, botPath)
	var result string
	result, err = TemplateData{}.GetConfigConstantsFile(botUsername, botToken, webhookUrl, serverPort, supperCommandSign, adminUserId)
	if err != nil {
		return
	}
	err = t.makeDirectory(folder)
	if err != nil {
		return
	}
	path := folder + "constants.go"
	err = ioutil.WriteFile(path, []byte(result), os.ModePerm)
	return
}
