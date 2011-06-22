package main 

import (
	"gotcl"
	"goirc"
	"fmt"
	"os"
	"strings"
)

const (
	CMD_BIND = "bind"
	CMD_PUTSERV = "putserv"
)

const (
	BIND_PUB = "pub"
	BIND_JOIN = "join"
)

var g_scripts *ScriptManager = &ScriptManager{scripts : make([]*Script, 0)}

type ScriptManager struct {
	scripts []*Script
}

func (man *ScriptManager) AddScript(_script *Script) {
	man.scripts = append(man.scripts, _script)
}

func (man *ScriptManager) OnPub(_nickname, _host, _handle, _channel, _text string) {
	command := man.firstWord(_text)
	command = strings.ToLower(command)
	man.fillEmptyParams(&_nickname, &_host, &_handle, &_channel, &_text)
	_text = man.quote(_text)
	
	for _, script := range man.scripts {
		if sp, contains := script.binds[BIND_PUB]; contains {
			if sp.param == command {
				call := fmt.Sprintf("%s %s %s %s %s %s", sp.proc, _nickname, _host, _handle, _channel, _text)
				_, err := script.interpreter.EvalString(call)
				if err != nil {
					fmt.Printf("Error evaluating OnPub procedure '%s' in %s: %s", sp.proc, script.filename, err.String())
				}
			}
		}
	}
}

func (man *ScriptManager) OnJoin(_nickname, _host, _handle, _channel string) {
	man.fillEmptyParams(&_nickname, &_host, &_handle, &_channel)

	for _, script := range man.scripts {
		if sp, contains := script.binds[BIND_JOIN]; contains {
			if sp.param == _channel {
				call := fmt.Sprintf("%s %s %s %s %s", sp.proc, _nickname, _host, _handle, _channel)
				_, err := script.interpreter.EvalString(call)
				if err != nil {
					fmt.Printf("Error evaluating OnPub procedure '%s' in %s: %s", sp.proc, script.filename, err.String())
				}
			}
		}
	}
}

func (man *ScriptManager) quote(_text string) string {
	return "\"" + _text + "\""
}

func (man *ScriptManager) firstWord(_text string) (word string) {
	word = _text
	space := strings.Index(_text, " ")
	if space != -1 {
		word = word[0:space]
	}
	return
}

func (man *ScriptManager) fillEmptyParams(_params ... *string) {
	for _, param := range _params {
		if len(*param) == 0 {
			*param = "-"
		}
	}
}

type Script struct {
	interpreter *gotcl.Interp
	irc *goirc.IRC
	filename string
	binds map[string]*ScriptProc
}

type ScriptProc struct {
	command string
	flags string
	param string
	proc string
}

func NewScriptProc(_command, _flags, _param, _proc string) *ScriptProc {
	return &ScriptProc{command : _command, flags : _flags, param : _param, proc : _proc}
}

func LoadScript(_filename string, _irc *goirc.IRC) *Script{
	file, err := os.Open(_filename)
	if err != nil {
		fmt.Printf("Error loading script %s: %s\n", _filename, err.String())
		
	}
	defer file.Close()
	
	i := gotcl.NewInterp()
	
	script := &Script{filename : _filename, interpreter : i, binds : make(map[string]*ScriptProc), irc : _irc}
	
	LoadScriptCommands(script)
	
	_, err = i.Run(file)
	if err != nil {
		fmt.Printf("Error running script %s: %s\n", _filename, err.String())
	}
	
	g_scripts.AddScript(script)
	
	return script
}

func LoadScriptCommands(_script *Script) {
	_script.interpreter.SetCmd(CMD_BIND, func(_i *gotcl.Interp, _args []*gotcl.TclObj) gotcl.TclStatus {
		_script.Cmd_Bind(_i, _args)
		return 0
	})
	
	_script.interpreter.SetCmd(CMD_PUTSERV, func(_i *gotcl.Interp, _args []*gotcl.TclObj) gotcl.TclStatus {
		_script.Cmd_Putserv(_i, _args)
		return 0
	})
}

func (s *Script) Cmd_Bind(_i *gotcl.Interp, _args []*gotcl.TclObj) gotcl.TclStatus {
	cmd := _args[0].AsString()
	if len(_args) < 4 {
		fmt.Printf("Not enough arguments for bind command '%s' in %s", cmd, s.filename)
	}
	
	flags  	:= _args[1].AsString()
	param := _args[2].AsString()
	proc 	:= _args[3].AsString()

	s.binds[cmd] = NewScriptProc(cmd, flags, param, proc)

	return 0
}

func (s *Script) Cmd_Putserv(_i *gotcl.Interp, _args []*gotcl.TclObj) gotcl.TclStatus {
	if len(_args) == 1 {
		msg := _args[0].AsString()
		s.irc.Send(msg+"\r\n")
	}
	return 0
}
