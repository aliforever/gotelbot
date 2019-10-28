package templates

import (
	gomondel "github.com/aliforever/gomondel/funcs"
)

const modelsFilePath string = "%s/models/"

func (t Template) modelsGlobalsFile() string {
	return `package models

import "github.com/aliforever/telegram-bot-api"

var API *tgbotapi.BotAPI`
}

func (t Template) initDatabase(botPath, botUsername string) (err error) {
	_, err = gomondel.InitDatabase(botPath, botUsername)
	return
}

func (t Template) makeUsersModel(botPath string) (err error) {
	modelName := "User"
	modelIdType := "int"
	fields := map[string]string{
		"FirstName": "string",
		"LastName":  "string",
		"Username":  "string",
		"IsAdmin":   "bool",
		"Menu":      "string",
		"Language":  "string",
	}
	modelFields := gomondel.MakeModelFieldsFromMap(fields)
	_, err = gomondel.CreateModel(botPath, modelName, &modelIdType, nil, nil, modelFields)
	return
}
