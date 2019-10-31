package templates

import (
	"fmt"
	"os"

	"github.com/aliforever/gotelbot/functions"

	"github.com/gobuffalo/flect"

	"strings"

	"github.com/go-errors/errors"
)

type Template struct {
}

func (t Template) Init(botUsername, botToken, botPath string, languages []string) (err error) {
	if strings.Index(botPath, "src") == -1 {
		err = errors.New("Wrong Path")
		return
	}
	fmt.Println(fmt.Sprintf("Creating directory %s", botPath))
	err = t.makeDirectory(botPath)
	if err != nil {
		return
	}
	fmt.Println("Creating configurations...")
	err = t.makeConfigsConstantFile(botUsername, botToken, botPath, "", "", "", 0)
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("Initializing database %s", botUsername))
	err = t.initDatabase(botPath, botUsername)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println("Creating user model")
	err = t.makeUsersModel(botPath)
	if err != nil {
		return
	}
	fmt.Println("Creating api file")
	err = t.makeApiFile(botPath)
	if err != nil {
		return
	}
	fmt.Println("Creating language files")
	err = t.initLanguageFiles(botPath, languages)
	if err != nil {
		return
	}
	fmt.Println("Creating application files")
	err = t.initApplicationFiles(botPath, botUsername, languages)
	if err != nil {
		return
	}
	fmt.Println("Creating main file")
	err = t.initMainFile(botPath, flect.Capitalize(botUsername))
	if err != nil {
		return
	}
	err = functions.FmtPath(botPath)
	/*if err != nil {
		return
	}
	err = functions.ImportsPath(botPath)*/
	return
}

func (t Template) makeDirectory(path string) (err error) {
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
		}
	}
	return
}

func (t Template) CurrentPath() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}
	return
}
