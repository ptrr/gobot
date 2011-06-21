package main

import (
	"strings"
	"regexp"
	//"http"
	"os"
	"fmt"
	
	"gotcl"
	"goirc"
)

var irc *goirc.IRC
const command string = "^![a-z0-9]"

func main() {
	TCLTest()

	println("Connecting to irc")
	irc = goirc.NewIRC("irc.rizon.net", "6667", "StiertjesBot")
	if success, err := irc.Connect(); !success {
		println(err)
		return
	}
	irc.AddChannel("#PU_HORSES", "", true)
	irc.ReceiveFunc = ReceiveIRC
	irc.Receive()
	
	// Just so it doens't shut down
	//http.ListenAndServe(":6543", nil)
}

func TCLTest() {
	file, e := os.Open("scripts/derp.tcl")
	if e != nil {
		panic(e.String())
	}
	defer file.Close()
	i := gotcl.NewInterp()
	_, err := i.Run(file)
	if err != nil {
		fmt.Println("Error: " + err.String())
	}
}

func ReceiveIRC(_command string, _arguments []string, _message, _nickname string) {
	switch _command {
		case "PRIVMSG":
			channel := _arguments[0]
			if channel == irc.Nickname {
				channel = _nickname
			}
			ProcessPRIVMSG(channel, _message, _nickname)
	}
}

func ProcessPRIVMSG(_channel, _message, _nickname string) {
	matched, error := regexp.MatchString(command, _message)
	if error == nil {
		if matched {
			ProcessCommand(_message, _channel, _nickname)
		} else {
			if _message == "herp" {
				irc.SendPriv(_channel,  "derp")
			}	
		}
	} else {
		println("ProcessPRIVMSG error: "+error.String())
	}
}

func ProcessCommand(_message string, _channel string, _nickname string){
	clean := strings.Replace(_message, "!", "", -1)
	stem := strings.Split(clean, " ", -1);
	switch stem[0] {
		case "kill":
			if(len(stem) < 2){
				msg := "PRIVMSG "+_channel+" :No one to kill.\r\n"
				irc.Send(msg)
			} else {
				msg := "PRIVMSG "+_channel+" :\001ACTION kills "+ stem[1] +"! \001\r\n"
				irc.Send(msg)
			}
		default : 
			println("Unknown command")		
	}
}