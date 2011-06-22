package main

import (
	"log"
	"http"
	"goirc"
)

var bots map[string]*goirc.IRC = make(map[string]*goirc.IRC)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	log.Printf("Client connected")
}

func (i *Manager) StartManager() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/new", NewBot)
	http.HandleFunc("/create", CreateBot)
	http.HandleFunc("/init", InitializeBot)
	log.Printf("Manager started")
	err := http.ListenAndServe(":1337", nil);
	if err != nil {
		log.Fatal(err)
	}
}

func NewBot(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

func CreateBot(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	server := req.FormValue("server")
	//channel := req.FormValue("channel")
	
	bot := goirc.NewIRC(server, "6667", name)
	bots[name] = bot
	http.Redirect(w, req, "/", -1)
}

func InitializeBot(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	if(bots[name] != nil){
		bot := bots[name]
		if success, err := bot.Connect(); !success {
			println(err)
			return
		}
		bot.SendJoin("#PU_HORSES", "")
		bot.ReceiveFunc = ReceiveIRC
		go bot.Receive()
	}
}

const page = `
	<form action="/create" method="post">
Botname: <input type="text" name="name" /><br />
Server: <input type="text" name="server" /><br />
Channel: <input type="text" name="channel" /><br />
<input type="submit" value="Submit" /></form>
`