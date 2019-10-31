package templates

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func (t Template) InitMenu(botPath, applicationName, menuName string, lineNumber int) (err error) {
	fmt.Println("Adding menu to application file")
	applicationFile := fmt.Sprintf(applicationFilePath, botPath) + "application.go"
	var content []byte
	content, err = ioutil.ReadFile(applicationFile)
	if err != nil {
		return
	}
	var str string
	str, err = TemplateData{}.GetMenuText(applicationName, menuName)
	if err != nil {
		return
	}
	if lineNumber == 0 {
		content = append(content, []byte("\n\n"+str)...)
	} else {
		body := string(content)
		split := strings.Split(body, "\n")
		if lineNumber >= len(split) {
			content = append(content, []byte("\n\n"+str)...)
		} else {
			c := strings.Join(split[:lineNumber], "\n") + "\n" + str + "\n" + strings.Join(split[lineNumber:], "\n")
			content = []byte(c)
		}
	}
	err = ioutil.WriteFile(applicationFile, content, os.ModePerm)
	if err != nil {
		return
	}

	fmt.Println("Adding menu text to language files")
	languageInterfaceFile := fmt.Sprintf(languagesFilePath, botPath) + "interface.go"
	content, err = ioutil.ReadFile(languageInterfaceFile)
	if err != nil {
		return
	}
	body := string(content)
	pos := strings.LastIndex(body, "}")
	if pos == -1 {
		err = errors.New("invalid language interface file")
		return
	}
	str, err = TemplateData{}.GetMenuTextInterface(menuName)
	if err != nil {
		return
	}
	body = body[:pos] + str + "\n" + body[pos:]
	err = ioutil.WriteFile(languageInterfaceFile, []byte(body), os.ModePerm)
	if err != nil {
		return
	}
	var languageFiles []os.FileInfo
	languageFiles, err = ioutil.ReadDir(fmt.Sprintf(languagesFilePath, botPath))
	if err != nil {
		return
	}
	for _, f := range languageFiles {
		if f.Name() == "interface.go" || !strings.Contains(f.Name(), ".go") {
			continue
		}
		fmt.Println("Adding to " + f.Name() + " file")
		fPath := fmt.Sprintf(languagesFilePath, botPath) + f.Name()
		content, err = ioutil.ReadFile(fPath)
		if err != nil {
			return
		}
		body = string(content)
		r := regexp.MustCompile(`type\s+(.+)\s+struct\s+{`)
		matches := r.FindAllStringSubmatch(body, -1)
		if len(matches) != 1 || len(matches[0]) != 2 {
			err = errors.New("invalid language file")
			return
		}
		structName := matches[0][1]
		strcutSign := strings.ToLower(structName[:1])
		str, err = TemplateData{}.GetMenuTextLanguage(structName, strcutSign, menuName)
		if err != nil {
			return
		}
		body = body + "\n\n" + str
		err = ioutil.WriteFile(fPath, []byte(body), os.ModePerm)
		if err != nil {
			return
		}
	}
	return
}
