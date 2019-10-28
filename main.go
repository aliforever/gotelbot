package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"net/http"
	"strings"

	"github.com/aliforever/gotelbot/functions"

	"github.com/aliforever/gotelbot/templates"
)

func main() {
	var init, botToken, path, langs string
	var botUsername, botPath *string
	var languages []string
	flag.StringVar(&init, "init", "", "--init=bot_token[,bot_username]")
	flag.StringVar(&langs, "langs", "", "--langs=english,farsi")
	flag.StringVar(&path, "path", "", "--path=/home/go/src/")
	flag.Parse()
	if init != "" {
		split := strings.Split(init, ",")
		if len(split) == 2 {
			botUsername = &split[1]
		}
		botToken = split[0]
		if langs != "" {
			languages = strings.Split(langs, ",")
		}
		if path != "" {
			botPath = &path
		}
		botPath, botUsername, err := InitTelegramBot(botToken, botUsername, botPath, languages)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(fmt.Sprintf("%s structure is created at %s", botUsername, botPath))
	}
}

type UsernameResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		Id        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	} `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func GetTelegramUsernameWithToken(token string) (usernameResp *UsernameResponse, err error) {
	address := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token)
	var resp *http.Response
	resp, err = http.Get(address)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	j := json.NewDecoder(resp.Body)
	err = j.Decode(&usernameResp)
	return
}

func InitTelegramBot(token string, username, path *string, languages []string) (botPath, botUsername string, err error) {
	var usernameResp *UsernameResponse
	if username == nil {
		usernameResp, err = GetTelegramUsernameWithToken(token)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error getting username for the given token %s: %s", token, err.Error()))
			return
		}
		if !usernameResp.Ok {
			err = errors.New(fmt.Sprintf("Error getting username for the given token %s: %d %s", token, usernameResp.ErrorCode, usernameResp.Description))
			return
		}
		botUsername = usernameResp.Result.Username
	} else {
		botUsername = *username
	}
	fmt.Println(fmt.Sprintf("Initializing %s", botUsername))
	if path == nil {
		botPath = build.Default.GOPATH + "/src/" + botUsername
	} else {
		botPath = *path + "/" + botUsername
	}
	err = templates.Template{}.Init(botUsername, token, botPath, languages)
	if err == nil {
		err = functions.FmtPath(botPath)
		/*if err == nil {
			err = functions.ImportsPath(botPath)
		}*/
	}
	return
}
