package main

import (
	"strings"
	"regexp"
)

var irc *IRC
const command string = "^![a-z0-9]"

func main() {
	irc = NewIRC("irc.rizon.net", "6667", "GoBOT")
	if success, err := irc.Connect(); !success {
		println(err)
		return
	}
	irc.SendJoin("#pokemon-universe", "")
	irc.receiveFunc = ReceiveIRC
	irc.receive()
}

func ReceiveIRC(_command string, _arguments []string, _message, _nickname string) {
	switch _command {
		case "PING":
			time := _message
			ProcessPING(time)
		case "PRIVMSG":
			channel := _arguments[0]
			ProcessPRIVMSG(channel, _message, _nickname)
	}
}

func ProcessPING(_time string) {
	irc.Send("PONG "+_time+"\r\n")
}

func ProcessPRIVMSG(_channel, _message, _nickname string) {
	matched, error := regexp.MatchString(command, _message)
	if error == nil {
		if matched {
			ProcessCommand(_message, _channel)
		} else {
			if _message == "herp" {
				if _channel == irc.nickname {
					_channel = _nickname
				}
				irc.SendPriv(_channel,  "derp")
			}	
		}
	} else {
		println("ProcessPRIVMSG error: "+error.String())
	}
}

func ProcessCommand(_message string, _channel string){
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