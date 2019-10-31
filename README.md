# **Golang Telegram Bot Framework**

#### **Features:**
+ Multilingual
+ MongoDB Database Storage
+ Menu Based
+ Super Easy to Develop


#### **Installation**
1. `go get -u github.com/aliforever/gotelbot`
2. `go install`

#### **Usage**
1. Message [BotFather](https://t.me/botfather) and create a new bot.
2. Copy the token.
3. Run `gotelbot --init=token` _(Replace token)_.
4. Run `go run main.go`.
5. Say Hello to your bot on Telegram. 

#### **Commands**
1. ##### `gotelbot --init=token[,username] [--path=path/to/go/src] [--langs=english,etc]` 
        - (if you don't pass ',username', gotelbot will grab the username)
        - (if you don't pass --path, gotelbot will use GOPATH environment variable)
        - (if you don't pass --langs, gotelbot will use English as default language)
    example: `gotelbot --init=bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`
2. ##### `gotelbot --menu=name[:line] [--path=path/to/project]`
        - (if you don't pass ':line', gotelbot will append menu to end of application.go file)
        - (if you don't pass --path, gotelbot will read current terminal directory, 
        where gotelbot command was executed. Make sure to cd bot's path.)
    example: `gotelbot --menu=Welcome`
#### **Flags**
+ `--langs=english,persian,italian`
+ `--path=/home/go/src/`
