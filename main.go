package main

import (
	"strings"
	"regexp"
	//"http"
	//"gotcl"
	"goirc"
)

var irc *goirc.IRC
var manager *Manager

const command string = "^![a-z0-9]"



func main() {
	println("Connecting to irc")
	manager = NewManager()
	manager.StartManager()
}

/*func TCLTest() {
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
}*/

func ReceiveIRC(_command string, _arguments []string, _message, _nickname string, _irc *goirc.IRC) {
	println(_command)
	println(_message)
	switch _command {
		case "PRIVMSG":
			channel := _arguments[0]
			if channel == _irc.Nickname {
				channel = _nickname
			}
			ProcessPRIVMSG(channel, _message, _nickname, _irc)
		case "433":
			Process443( _irc)
	}
}

func Process443(_irc *goirc.IRC) {
	_irc.Nickname = _irc.Nickname + "_"
	_irc.SendNick(_irc.Nickname)
	_irc.SendJoin("#PU_HORSES", "")
}

func ProcessPRIVMSG(_channel, _message, _nickname string, _irc *goirc.IRC) {
	matched, error := regexp.MatchString(command, _message)
	if error == nil {
		if matched {
			ProcessCommand(_message, _channel, _nickname, _irc)
		} else {
			if _message == "herp" {
				_irc.SendPriv(_channel,  "derp")
			}	
		}
	} else {
		println("ProcessPRIVMSG error: "+error.String())
	}
}

func ProcessCommand(_message string, _channel string, _nickname string, _irc *goirc.IRC){
	clean := strings.Replace(_message, "!", "", -1)
	stem := strings.Split(clean, " ", -1);
	switch stem[0] {
		case "kill":
			if(len(stem) < 2){
				msg := "PRIVMSG "+_channel+" :No one to kill.\r\n"
				_irc.Send(msg)
			} else {
				msg := "PRIVMSG "+_channel+" :\001ACTION kills "+ stem[1] +"! \001\r\n"
				_irc.Send(msg)
			}
		default : 
			println("Unknown command")		
	}
}