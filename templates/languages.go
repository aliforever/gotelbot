package templates

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gobuffalo/flect"
)

const languagesFilePath string = "%s/languages/"

func (t Template) languageInterfaceFile() string {
	return `package languages

type Language interface {
	LanguageName() string
	LanguageFlag() string
	LanguageCode() string
	SelectLanguageMenu() string
	SelectLanguageBtn() string
	MainMenu() string
	BackBtn() string
}`
}

func (t Template) languageFile() string {
	return `package languages

type {{.LanguageName}} struct {
}

func (e {{.LanguageName}}) LanguageFlag() string {
	return "🇺🇸"
}

func (e {{.LanguageName}}) LanguageName() string {
	return "{{.LanguageName}}"
}

func (e {{.LanguageName}}) LanguageCode() string {
	return "EN"
}

func (e {{.LanguageName}}) SelectLanguageMenu() string {
	return "Please Select Your Language"
}

func (e {{.LanguageName}}) SelectLanguageBtn() string {
	return "🇮🇷🇺🇸 Change Language"
}

func (e {{.LanguageName}}) MainMenu() string {
	return "You're in Main Menu, Please Choose an Option"
}

func (e {{.LanguageName}}) BackBtn() string {
	return "🔙 Back 🔙"
}`
}

func (t Template) initLanguageFiles(botPath string, languages []string) (err error) {
	folder := fmt.Sprintf(languagesFilePath, botPath)
	err = t.makeDirectory(folder)
	if err != nil {
		return
	}
	path := folder + "interface.go"
	err = ioutil.WriteFile(path, []byte(t.languageInterfaceFile()), os.ModePerm)
	if err != nil {
		return
	}
	if len(languages) >= 1 {
		for _, l := range languages {
			path = folder + l + ".go"
			lang := l
			var str string
			str, err = TemplateData{}.FillLanguageFile(t.languageFile(), flect.Capitalize(lang))
			if err != nil {
				return
			}
			err = ioutil.WriteFile(path, []byte(str), os.ModePerm)
			if err != nil {
				return
			}
		}
	} else {
		path = folder + "english.go"
		lang := "english"
		var str string
		str, err = TemplateData{}.FillLanguageFile(t.languageFile(), flect.Capitalize(lang))
		if err != nil {
			return
		}
		err = ioutil.WriteFile(path, []byte(str), os.ModePerm)
		if err != nil {
			return
		}
	}

	return
}
