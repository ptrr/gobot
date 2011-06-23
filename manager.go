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
			text += "<div class='bot' id='" + name + "'><div class='inner'>" + name
			if b.Connected {
				text += " <span style='color: #6add4b;'>Connected</span>"
				text += "</div><div class='buttons'><button id='but_" + name + "' style='display:none' class='connect'>Connect</button><button  id='but_" + name + "' class='disconnect'>Disconnect</button></div></div>"
			} else {
				text += " <span style='color: #dd4b4b;'>Not connected</span>"
				text += "</div><div class='buttons'><button id='but_" + name + "' class='connect'>Connect</button><button style='display:none' id='but_" + name + "' class='disconnect'>Disconnect</button></div></div>"
			}
			
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
	http.HandleFunc("/kill", KillBot)		
	http.HandleFunc("/", handler)
	log.Printf("Manager started")
	err := http.ListenAndServe(":8080", nil);
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
	LoadScript("scripts/derp.tcl", bot)
	for _, c := range channel {
		bot.AddChannel(c, "", false)
	}
	bots[name] = bot
	http.ServeFile(w, req, "html/complete.html")
	//http.Redirect(w, req, "/", http.StatusFound)
}

func InitializeBot(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	if(bots[name] != nil){
		bot := bots[name]
		if !bot.Connected {
			if success, err := bot.Connect(); !success {
				println(err)
				return
			}
			bot.ReceiveFunc = ReceiveIRC
			go bot.Receive()
		}
		
	}
}

func KillBot(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	if(bots[name] != nil){
		bot := bots[name]
		bot.Disconnect()
	}
}