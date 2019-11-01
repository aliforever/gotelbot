package templates

import (
	"bytes"
	"text/template"

	"github.com/gobuffalo/flect"
)

type TemplateData struct {
	ApplicationName  string
	IsMultiLingual   bool
	Languages        []string
	Language         string
	Text             string
	LanguageSign     string
	DefaultLanguage  string
	BotUsername      string
	BotToken         string
	BotGoPathSrc     string
	WebhookUrl       string
	ServerPort       string
	AdminUserId      int
	SuperCommandSign string
	LanguageName     string
	MenuName         string
}

func (td TemplateData) main() string {
	return `package main

import (
	"fmt"
	"os"
	"net/http"
	"errors"
	tgbotapi "github.com/aliforever/telegram-bot-api"
	goerrors "github.com/go-errors/errors"
	"{{.BotGoPathSrc}}/application"
	"{{.BotGoPathSrc}}/configs"
	"{{.BotGoPathSrc}}/models"
	"{{.BotGoPathSrc}}/api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(configs.BotToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	if setWebHook() {
		res, err := bot.SetWebhook(tgbotapi.NewWebhook(configs.WebhookUrl))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	}
	err = models.InitMongoDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	var updates tgbotapi.UpdatesChannel
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	if !getUpdates() && configs.WebhookUrl != "" {
		updates = bot.ListenForWebhook("/" + bot.Token)
	} else {
		bot.RemoveWebhook()
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, err = bot.GetUpdatesChan(u)
	}
	go http.ListenAndServe("0.0.0.0:"+configs.ServerPort, nil)
	api.Client = bot

	app := application.{{.ApplicationName}}{}
	for update := range updates {	
		go PanicHandler(&update, app.ProcessUpdate)
	}
}

func PanicHandler(update *tgbotapi.Update, fun func(update *tgbotapi.Update) error) {
	var err error
	defer func() {
		r := recover()
		if r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("unknown error")
			}
			fmt.Println(goerrors.Wrap(err, 2).ErrorStack())
		}
	}()
	fun(update)
}

func setWebHook() bool {
	return os.Getenv("SET_WEBHOOK") == "YES"
}

func getUpdates() bool {
	return os.Getenv("GET_UPDATES") == "YES"
}
`
}

func (td TemplateData) configConstants() string {
	return `package configs

const (
	BotUsername       = "{{.BotUsername}}"
	BotToken          = "{{.BotToken}}"
)

var (
	WebhookUrl        = "{{if .WebhookUrl}}{{.WebhookUrl}}{{end}}"
	ServerPort        = "{{if .ServerPort}}{{.ServerPort}}{{else}}6000{{end}}"
	AdminUserId       = "{{if .AdminUserId}}{{.AdminUserId}}{{else}}81997375{{end}}"
	SupperCommandSign = "{{if .SuperCommandSign}}{{.SuperCommandSign}}{{else}}/_{{end}}"
)
`

}

func (td TemplateData) apiFile() string {
	return `package api

import tgbotapi "github.com/aliforever/telegram-bot-api"

var Client *tgbotapi.BotAPI`
}

func (td TemplateData) methods() string {
	return `package application

import (
	tgbotapi "github.com/aliforever/telegram-bot-api"
	"{{.BotGoPathSrc}}/languages"
)

func GetUpdateSender(update *tgbotapi.Update) (sender *tgbotapi.User) {
	if update.Message != nil {
		sender = update.Message.From
	} else if update.ChannelPost != nil {
		sender = update.ChannelPost.From
	} else if update.CallbackQuery != nil {
		sender = update.CallbackQuery.From
	}
	return 
}
func BackKeyboard(lang languages.Language) *tgbotapi.ReplyKeyboardMarkup {
	var rows [][]string
	firstRow := []string{lang.BackBtn()}
	rows = append(rows, firstRow)
	return MakeReplyKeyboardFromArray(rows)
}

func MakeReplyKeyboardFromArray(rows ...[][]string) *tgbotapi.ReplyKeyboardMarkup {
	var keyboardButtons [][]tgbotapi.KeyboardButton
	var buttonRow []tgbotapi.KeyboardButton
	for _, row := range rows {
		for _, buttons := range row {
			for _, button := range buttons {
				buttonRow = append(buttonRow, tgbotapi.KeyboardButton{Text: button})
			}
			keyboardButtons = append(keyboardButtons, buttonRow)
			buttonRow = []tgbotapi.KeyboardButton{}
		}
	}
	keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: keyboardButtons, ResizeKeyboard: true}
	return keyboard
}

func MakeInlineReplyKeyboardFromArray(rows ...[][]map[string]string) *tgbotapi.InlineKeyboardMarkup {
	var InlineKeyboardButton [][]tgbotapi.InlineKeyboardButton
	var keyboardRow []tgbotapi.InlineKeyboardButton
	for _, row := range rows {
		for _, buttons := range row {
			for _, button := range buttons {
				buttonObj := tgbotapi.InlineKeyboardButton{}
				if _, ok := button["text"]; ok {
					buttonObj.Text = button["text"]
				}
				if _, ok := button["url"]; ok {
					str := button["url"]
					buttonObj.URL = &str
				}
				if value, ok := button["data"]; ok {
					buttonObj.CallbackData = &value
				}
				keyboardRow = append(keyboardRow, buttonObj)
			}
			InlineKeyboardButton = append(InlineKeyboardButton, keyboardRow)
			keyboardRow = []tgbotapi.InlineKeyboardButton{}
		}
	}
	keyboard := &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardButton}
	return keyboard
}`

}

func (td TemplateData) menu() string {
	return `func (bot *{{.ApplicationName}}) {{.MenuName}}Menu() {
	bot.CurrentUser.SaveField("menu", "MainMenu")

	if !bot.IsSwitched {
		if bot.Update.Message.Text != "" {
			items := map[string]func(){
				bot.Language.BackBtn(): bot.MainMenu,
			}
			if _, ok := items[bot.Update.Message.Text]; ok {
				bot.IsSwitched = true
				items[bot.Update.Message.Text]()
				return
			}
		}
	}
	cfg := tgbotapi.NewMessage(int64(bot.CurrentUser.Id), bot.Language.{{.MenuName}}Menu())
	cfg.ReplyMarkup = bot.{{.MenuName}}MenuKeyboard()
	api.Client.Send(cfg)
}

func (bot *{{.ApplicationName}}) {{.MenuName}}MenuKeyboard() (keyboard *tgbotapi.ReplyKeyboardMarkup) {
	var rows [][]string
	rows = append(rows, []string{bot.Language.BackBtn()})
	keyboard = MakeReplyKeyboardFromArray(rows)
	return
}`
}

func (td TemplateData) menuTextInterface() string {
	return `{{.MenuName}}Menu() string`
}

func (td TemplateData) menuTextLanguage() string {
	return `func ({{.LanguageSign}} {{.Language}}) {{.MenuName}}Menu() string {
	return "Welcome to {{.MenuName}} Menu"
}`
}

func (td TemplateData) textInterface() string {
	return `{{.Text}}() string`
}

func (td TemplateData) textLanguage() string {
	return `func ({{.LanguageSign}} {{.Language}}) {{.Text}}() string {
	return "Text Goes Here"
}`
}

func (td TemplateData) application() string {
	return `package application

import (
	tgbotapi "github.com/aliforever/telegram-bot-api"
	"{{.BotGoPathSrc}}/models"
	"{{.BotGoPathSrc}}/languages"
	"{{.BotGoPathSrc}}/configs"
	"{{.BotGoPathSrc}}/api"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"errors"
	"reflect"
	"strings"
)

type {{.ApplicationName}} struct {
	Update      *tgbotapi.Update
	CurrentUser *models.User
	IsSwitched  bool
	Language    languages.Language
}

func (bot *{{.ApplicationName}}) ProcessSuperCommands(command string) {
if command == "users:count" {
		count, err := models.User{}.Count()
		if err != nil {
			cfg := tgbotapi.NewMessage(int64(bot.CurrentUser.Id), err.Error())
			api.Client.Send(cfg)
			return
		}
		cfg := tgbotapi.NewMessage(int64(bot.CurrentUser.Id), fmt.Sprintf("User Count is: %d", count))
		api.Client.Send(cfg)
	}
	return
}

func (bot *{{.ApplicationName}}) SwitchMenu(name string) {
	st := reflect.TypeOf(bot)
	_, ok := st.MethodByName(name)
	if !ok {
		cfg := tgbotapi.NewMessage(int64(bot.CurrentUser.Id), "Invalid Menu " + name)
		api.Client.Send(cfg)
		bot.CurrentUser.SaveField("menu", "MainMenu")
		return
	}
	bot.IsSwitched = true
	reflect.ValueOf(bot).MethodByName(name).Call([]reflect.Value{})
}

func (bot *{{.ApplicationName}}) registerOrUpdateUser() (err error) {
	if bot.Update == nil {
		err = errors.New("empty update")
		return
	}
	from := GetUpdateSender(bot.Update)
	if from == nil {
		err = errors.New("nil_from")
		return
	}
	bot.CurrentUser, err = models.User{}.FindById(from.ID)
	if err == nil && bot.CurrentUser != nil {
		bsonMap := bson.M{}
		bsonMap["first_name"] = from.FirstName
		if from.LastName != "" {
			bsonMap["last_name"] = from.LastName
		}
		if from.UserName != "" {
			bsonMap["username"] = from.UserName
		}
		err = bot.CurrentUser.SaveCustom(bsonMap)
	} else {
		bot.CurrentUser = &models.User{}
		bot.CurrentUser.Id = from.ID
		bot.CurrentUser.FirstName = from.FirstName
		bot.CurrentUser.LastName = from.LastName
		bot.CurrentUser.Username = from.UserName
		err = bot.CurrentUser.Create()
	}
	return
}

{{- if .IsMultiLingual}}
func (bot *{{.ApplicationName}}) ChooseLanguageMenu() {
	bot.CurrentUser.SaveField("menu", "ChooseLanguageMenu")
	
	if !bot.IsSwitched {
		if bot.Update.Message.Text != "" {
			text := bot.Update.Message.Text
			{{range $element := .Languages}}select{{$element}} := (languages.{{$element}}{}).LanguageFlag() + " " + (languages.{{$element}}{}).LanguageName()
			if text == select{{$element}} {
				bot.Language = &languages.{{$element}}{}
				bot.CurrentUser.Language = (languages.{{$element}}{}).LanguageCode()
				bot.CurrentUser.Save()
				bot.SwitchMenu("MainMenu")
				return
			}
			{{end}}
		}
	}
	text := {{- range $element := .Languages}} (languages.{{$element}}{}).SelectLanguageMenu() + " " + (languages.{{$element}}{}).LanguageFlag() + {{- end}} ""
	cfg := tgbotapi.NewMessage(int64(bot.CurrentUser.Id), text)
	cfg.ReplyMarkup = bot.ChooseLanguageMenuKeyboard()
	api.Client.Send(cfg)
}

func (bot *{{.ApplicationName}}) ChooseLanguageMenuKeyboard() (keyboard *tgbotapi.ReplyKeyboardMarkup) {
	var rows [][]string
	var langs []string

	{{range $element := .Languages}}langs = append(langs, (languages.{{$element}}{}).LanguageFlag() + " " + (languages.{{$element}}{}).LanguageName())
	{{end}}
	rows = append(rows, langs)
	keyboard = MakeReplyKeyboardFromArray(rows)
	return
}
{{- end}}

func (bot *{{.ApplicationName}}) ProcessGroup() {

}

func (bot *{{.ApplicationName}}) ProcessCallback() {
	callback := bot.Update.CallbackQuery
	split := strings.Split(callback.Data, ":")
	itemsMap := map[string]func(data []string){}
	if _, ok := itemsMap[split[0]]; ok {
		itemsMap[split[0]](split[1:])
		return
	}
	cfg := tgbotapi.NewMessage(int64(bot.CurrentUser.Id), "Callback Query Handler not Found: " + callback.Data)
	api.Client.Send(cfg)
}

func (bot *{{.ApplicationName}}) ProcessMenu() {
	menu := bot.CurrentUser.Menu
	st := reflect.TypeOf(bot)
	_, ok := st.MethodByName(menu)
	if ok {
		bot.IsSwitched = false
		reflect.ValueOf(bot).MethodByName(menu).Call([]reflect.Value{})
		return
	}
	bot.SwitchMenu("MainMenu")
}

func (bot *{{.ApplicationName}}) ProcessUpdate(update *tgbotapi.Update) (err error) {
	bot.Update = update
	message := bot.Update.Message
	callback := bot.Update.CallbackQuery
	inlineQuery := bot.Update.InlineQuery

	if message == nil && callback == nil && inlineQuery == nil {
		return
	}

	err = bot.registerOrUpdateUser()
	if err != nil {
		return
	}

	{{- if .IsMultiLingual}}
	if bot.CurrentUser.Language == "" {
	{{- range $element := .Languages}} 
		if bot.CurrentUser.Language == (languages.{{$element}}{}).LanguageCode() {
			bot.Language = &languages.{{$element}}{}
		}
	{{- end}}
	} else {
		if bot.CurrentUser.Menu != "ChooseLanguageMenu" {
			bot.IsSwitched = true
			bot.ChooseLanguageMenu()
			return
		}
	}
	{{else}}
	if bot.CurrentUser.Language == "" {
		bot.CurrentUser.Language = (languages.{{.DefaultLanguage}}{}).LanguageCode()
		bot.CurrentUser.Save()
	}
	bot.Language = languages.{{.DefaultLanguage}}{}
	{{- end}}
	if message != nil {
		if message.Text != "" {
			if strings.Contains(message.Text, configs.SupperCommandSign) {
				split := strings.Split(message.Text, configs.SupperCommandSign)
				bot.ProcessSuperCommands(split[1])
				return
			}
		}
		if strings.Contains(message.Chat.Type, "group") {
			bot.ProcessGroup()
		} else if message.Chat.Type == "private" {
			bot.ProcessMenu()
		} else {
			return
		}
	} else if callback != nil {
		bot.ProcessCallback()
	}
	return
}

func (bot *{{.ApplicationName}}) MainMenu() {
	bot.CurrentUser.SaveField("menu", "MainMenu")

	if !bot.IsSwitched {
		if bot.Update.Message.Text != "" {
			items := map[string]func(){
				{{- if .IsMultiLingual}}
				bot.Language.SelectLanguageBtn(): bot.ChooseLanguageMenu,
				{{- end}}
			}
			if _, ok := items[bot.Update.Message.Text]; ok {
				bot.IsSwitched = true
				items[bot.Update.Message.Text]()
				return
			}
		}
	}

	cfg := tgbotapi.NewMessage(int64(bot.CurrentUser.Id), bot.Language.MainMenu())
	{{- if .IsMultiLingual}}
	cfg.ReplyMarkup = bot.MainMenuKeyboard()
	{{- end}}
	api.Client.Send(cfg)
}

{{if .IsMultiLingual}}
func (bot *{{.ApplicationName}}) MainMenuKeyboard() (keyboard *tgbotapi.ReplyKeyboardMarkup) {
	var rows [][]string
	rows = append(rows, []string{bot.Language.SelectLanguageBtn()})
	keyboard = MakeReplyKeyboardFromArray(rows)
	return
}
{{end}}
`

}

func (td TemplateData) language() string {
	return `package languages

type {{.LanguageName}} struct {
}

func ({{.LanguageSign}} {{.LanguageName}}) LanguageFlag() string {
	return "🇺🇸"
}

func ({{.LanguageSign}} {{.LanguageName}}) LanguageName() string {
	return "{{.LanguageName}}"
}

func ({{.LanguageSign}} {{.LanguageName}}) LanguageCode() string {
	return "EN"
}

func ({{.LanguageSign}} {{.LanguageName}}) SelectLanguageMenu() string {
	return "Please Select Your Language"
}

func ({{.LanguageSign}} {{.LanguageName}}) SelectLanguageBtn() string {
	return "🇮🇷🇺🇸 Change Language"
}

func ({{.LanguageSign}} {{.LanguageName}}) MainMenu() string {
	return "You're in Main Menu, Please Choose an Option"
}

func ({{.LanguageSign}} {{.LanguageName}}) BackBtn() string {
	return "🔙 Back 🔙"
}`

}

func (td TemplateData) GetConfigConstantsFile(botUsername, botToken, webhookUrl, serverPort, superCommandSign string, adminUserId int) (result string, err error) {
	content := td.configConstants()
	td.BotUsername = botUsername
	td.BotToken = botToken
	td.WebhookUrl = webhookUrl
	td.ServerPort = serverPort
	td.AdminUserId = adminUserId
	td.SuperCommandSign = superCommandSign
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("config_constants").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) GetApplicationFile(botGoPathSrc, applicationName string, languages []string) (result string, err error) {
	content := td.application()
	if len(languages) > 1 {
		td.IsMultiLingual = true
		td.Languages = languages
		for k, v := range td.Languages {
			td.Languages[k] = flect.Capitalize(v)
		}
	} else {
		td.DefaultLanguage = "English"
	}
	td.ApplicationName = applicationName
	td.BotGoPathSrc = botGoPathSrc
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("application").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) GetMainFile(goPathSrc string, applicationName string) (result string, err error) {
	content := td.main()
	td.ApplicationName = applicationName
	td.BotGoPathSrc = goPathSrc
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("main").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) GetMethodsFile(botGoPathSrc string) (result string, err error) {
	content := td.methods()
	td.BotGoPathSrc = botGoPathSrc
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("methods").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) GetMenuText(applicationName, menuName string) (result string, err error) {
	content := td.menu()
	td.ApplicationName = applicationName
	td.MenuName = menuName
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("menu").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) GetMenuTextInterface(menuName string) (result string, err error) {
	content := td.menuTextInterface()
	td.MenuName = menuName
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("menu_text_interface").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) GetTextInterface(text string) (result string, err error) {
	content := td.textInterface()
	td.Text = text
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("text_interface").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) GetTextLanguage(language, sign, text string) (result string, err error) {
	content := td.textLanguage()
	td.Text = text
	td.Language = language
	td.LanguageSign = sign
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("text_language").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) GetMenuTextLanguage(language, sign, menu string) (result string, err error) {
	content := td.menuTextLanguage()
	td.MenuName = menu
	td.Language = language
	td.LanguageSign = sign
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("menu_text_language").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}

func (td TemplateData) FillLanguageFile(languageSign, languageName string) (result string, err error) {
	content := td.language()
	td.LanguageName = languageName
	td.LanguageSign = languageSign
	var (
		tmpl *template.Template
		bf   bytes.Buffer
	)

	tmpl, err = template.New("language").Parse(content)
	if err != nil {
		return
	}
	err = tmpl.Execute(&bf, td)
	if err != nil {
		return
	}
	result = bf.String()
	return
}
