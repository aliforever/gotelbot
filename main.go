package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aliforever/gotelbot/functions"

	"github.com/aliforever/gotelbot/templates"
)

func main() {
	var init, botToken, path, langs, menu string
	var botUsername, botPath *string
	var languages []string
	flag.StringVar(&init, "init", "", "--init=bot_token[,bot_username]")
	flag.StringVar(&langs, "langs", "", "--langs=english,farsi")
	flag.StringVar(&path, "path", "", "--path=/home/go/src/")
	flag.StringVar(&menu, "menu", "", "--menu=Main[:20]")
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
		return
	}
	if menu != "" {
		var err error
		if path != "" {
			botPath = &path
		}
		lineNumber := 0
		split := strings.Split(menu, ":")
		if len(split) == 2 {
			lineNumber, err = strconv.Atoi(split[1])
			if err != nil {
				fmt.Println("invalid line number, should be integer")
				return
			}
		}
		err = InitMenu(split[0], lineNumber, botPath)
		if err != nil {
			fmt.Println(err)
		}
		return
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

func InitMenu(menuName string, lineNumber int, botPath *string) (err error) {
	var path string
	if botPath == nil {
		path, err = functions.CurrentPath()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if _, err = os.Stat(*botPath); err != nil {
			return
		}
		path = *botPath
	}
	applicationName := ""
	applicationPath := fmt.Sprintf("%s/application/application.go", path)
	var app []byte
	app, err = ioutil.ReadFile(applicationPath)
	if err != nil {
		err = errors.New("application file not found in path, " + err.Error())
		return
	}
	r := regexp.MustCompile(`type\s+(.+)\s+struct\s+{`)
	applicationStructs := r.FindAllStringSubmatch(string(app), -1)
	if len(applicationStructs) != 1 || len(applicationStructs[0]) != 2 {
		err = errors.New("more than one structs in application file, please pass --app=name")
		return
	}
	applicationName = applicationStructs[0][1]
	err = templates.Template{}.InitMenu(path, applicationName, menuName, lineNumber)
	if err == nil {
		err = functions.FmtPath(path)
		if err == nil {
			err = functions.ImportsPath(path)
		}
	}
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
