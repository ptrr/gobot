package goirc

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type IRC struct {
	socket      *net.TCPConn
	rw          *bufio.ReadWriter
	Connected   bool
	ReceiveFunc func(_command string, _arguments []string, _message, _nickname string, _bot *IRC)

	Address, Nickname string
	channels          map[string]string
	joinTries         int
}

func NewIRC(_ip, _port, _nickname string) *IRC {
	return &IRC{Address: _ip + ":" + _port, Nickname: _nickname, channels: make(map[string]string)}
}

func (i *IRC) Disconnect() {
	i.SendQuit()
	i.socket.Close()
	i.Connected = false
}

func (i *IRC) Connect() (bool, string) {
	var error error
	sock, error := net.Dial("tcp", i.Address)
	if error != nil {
		return false, error.Error()
	}
	i.socket = sock.(*net.TCPConn)
	i.rw = bufio.NewReadWriter(bufio.NewReader(i.socket), bufio.NewWriter(i.socket))

	i.SendNick(i.Nickname)
	i.sendUser(i.Nickname, i.socket.LocalAddr().String(), i.socket.RemoteAddr().String(), i.Nickname)
	i.JoinChannels()
	i.Connected = true
	return true, ""
}

func (i *IRC) Receive() {
	for i.Connected {
		text, err := i.rw.ReadString('\n')
		if err != nil {
			fmt.Printf("Receive error: %s\n", err.Error())
			break
		}

		text = strings.Replace(text, "\r\n", "", -1)
		i.parseMsg(text)
	}
	i.Connected = false
}

func (i *IRC) parseMsg(_text string) {
	var msg, nick string
	msg = _text
	if _text[0] == ':' {
		hostend := strings.Index(_text, " ")
		if hostend != -1 {
			host := _text[1:hostend]
			msg = _text[hostend:len(_text)]
			nick = ""

			nickend := strings.Index(host, "!")
			if nickend != -1 {
				nick = host[0:nickend]
			}
		}
	}
	parts := strings.SplitN(msg, " :", 2)
	if len(parts) >= 1 {
		if len(parts) == 1 {
			parts[0] = strings.Trim(parts[0], " ")
			parts = strings.SplitN(parts[0], " ", 2)
		}
		cmd := strings.Trim(parts[0], " ")
		args := make([]string, 0)
		cmdparts := strings.Split(cmd, " ")
		if len(cmdparts) > 1 {
			cmd = cmdparts[0]
			args = append(args, cmdparts[1:]...)
		}
		msg := parts[1]
		i.preReceiveFunc(cmd, args, msg, nick)
	}
}

func (i *IRC) preReceiveFunc(_command string, _arguments []string, _message, _nickname string) {
	switch _command {
	case "PING":
		time := _message
		i.Send("PONG" + time + "\n\r")

	case "433": //Nickname taken
		i.Nickname = fmt.Sprintf("%s%d", i.Nickname, i.joinTries)
		i.SendNick(i.Nickname)
		i.joinTries++
		i.JoinChannels()

	default:
		if i.ReceiveFunc != nil {
			i.ReceiveFunc(_command, _arguments, _message, _nickname, i)
		}
	}
}

func (i *IRC) sendUser(_nickname, _localhost, _remotehost, _realname string) {
	msg := "USER " + _nickname + " " + _localhost + " " + _remotehost + " :" + _realname + "\r\n"
	i.Send(msg)
}

func (i *IRC) Send(_text string) {
	if _, err := i.rw.WriteString(_text); err != nil {
		fmt.Printf("Send error: %s\n", err.Error())
	}
	i.rw.Flush()
}

func (i *IRC) SendNick(_nickname string) {
	msg := "NICK " + _nickname + "\r\n"
	i.Send(msg)
}

func (i *IRC) SendJoin(_channel, _password string) {
	msg := "JOIN " + _channel + " " + _password + "\r\n"
	i.Send(msg)
}

func (i *IRC) SendPriv(_channel, _message string) {
	msg := "PRIVMSG " + _channel + " :" + _message + "\r\n"
	i.Send(msg)
}

func (i *IRC) SendQuit() {
	msg := "QUIT Bot shut down. \r\n"
	i.Send(msg)
}

func (i *IRC) AddChannel(_channel, _password string, _connect bool) {
	if _, contains := i.channels[_channel]; contains {
		return //already in list
	}

	i.channels[_channel] = _password

	if _connect {
		i.SendJoin(_channel, _password)
	}
}

func (i *IRC) JoinChannels() {
	for channel, password := range i.channels {
		i.SendJoin(channel, password)
	}
}
