package main

import (
	"log"
	"http"
	"goirc"
	"fmt"
	"strings"
	"io/ioutil"
)

var bots map[string]*goirc.IRC = make(map[string]*goirc.IRC)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func handler(w http.ResponseWriter, req *http.Request) {
	file, err := ioutil.ReadFile("html/index.html")
	text := ""
	if err == nil {
		content := string(file)
		for name, b := range bots {
			text += "<div class='bot' id='" + name + "'>" + name
			if b.Connected {
				text += " <span style='color: green'>Connected</span><button id='but_" + name + "' class='disconnect'>Disconnect</button>"
			} else {
				text += " <span style='color: red'>Not connected</span><button id='but_" + name + "' class='connect'>Connect</button>"
			}
			text += "</div>"
		}
		content = strings.Replace(content, "{{BOTS}}", text, -1)
		fmt.Fprintf(w, content)
	}
}

func SourceHandler(w http.ResponseWriter, r *http.Request) {
	//log.Printf("IN SOURCE:%s\n", r.URL.Path[1:])
    http.ServeFile(w, r, r.URL.Path[1:])
	
}

func (i *Manager) StartManager() {
	http.HandleFunc("/html/", SourceHandler)
	http.HandleFunc("/new", NewBot)
	http.HandleFunc("/create", CreateBot)
	http.HandleFunc("/init", InitializeBot)		
	http.HandleFunc("/", handler)
	log.Printf("Manager started")
	err := http.ListenAndServe(":1337", nil);
	if err != nil {
		log.Fatal(err)
	}
}

func TestDing(w http.ResponseWriter, req *http.Request) {
	
}

func NewBot(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "html/newbot.html")
}

func CreateBot(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	name := req.Form["name"][0]
	server := req.Form["server"][0]
	channel := req.Form["channel[]"]
	bot := goirc.NewIRC(server, "6667", name)
	for _, c := range channel {
		bot.AddChannel(c, "", false)
	}
	bots[name] = bot
	http.Redirect(w, req, "/", http.StatusFound)
}

func InitializeBot(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	if(bots[name] != nil){
		bot := bots[name]
		if success, err := bot.Connect(); !success {
			println(err)
			return
		}
		bot.ReceiveFunc = ReceiveIRC
		go bot.Receive()
	}
}