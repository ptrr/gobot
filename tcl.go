package main 

import (
	"gotcl"
	"fmt"
	"os"
	"strings"
)

const (
	CMD_PUB = "pub"
)

var g_scripts *ScriptManager = &ScriptManager{scripts : make([]*Script, 0)}

type ScriptManager struct {
	scripts []*Script
}

func (man *ScriptManager) AddScript(_script *Script) {
	man.scripts = append(man.scripts, _script)
}

func (man *ScriptManager) OnPub(_nickname, _host, _handle, _channel, _text string) {
	command := _text
	space := strings.Index(_text, " ")
	if space != -1 {
		command = command[0:space]
	}
	command = strings.ToLower(command)
	
	if len(_nickname) == 0 {
		_nickname = "-"
	}
	if len(_host) == 0 {
		_host = "-"
	}
	if len(_handle) == 0 {
		_handle = "-"
	}
	if len(_channel) == 0 {
		_channel = "-"
	}
	if len(_text) == 0 {
		_text = "-"
	}
	_text = "\"" + _text + "\""
	
	for _, script := range man.scripts {
		if sp, contains := script.binds[CMD_PUB]; contains {
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

type Script struct {
	interpreter *gotcl.Interp
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

func LoadScript(_filename string) *Script{
	file, err := os.Open(_filename)
	if err != nil {
		fmt.Printf("Error loading script %s: %s\n", _filename, err.String())
		
	}
	defer file.Close()
	
	i := gotcl.NewInterp()
	
	script := &Script{filename : _filename, interpreter : i, binds : make(map[string]*ScriptProc)}
	
	LoadScriptCommands(script)
	
	_, err = i.Run(file)
	if err != nil {
		fmt.Printf("Error running script %s: %s\n", _filename, err.String())
	}
	
	g_scripts.AddScript(script)
	
	return script
}

func LoadScriptCommands(_script *Script) {
	_script.interpreter.SetCmd("bind", func(_i *gotcl.Interp, _args []*gotcl.TclObj) gotcl.TclStatus {
		_script.Cmd_Bind(_i, _args)
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
/*
func tclPuts(i *Interp, args []*TclObj) TclStatus {
	newline := true
	var s string
	file := i.chans["stdout"].(io.Writer)
	if len(args) == 1 {
		s = args[0].AsString()
	} else if len(args) == 2 || len(args) == 3 {
		if args[0].AsString() == "-nonewline" {
			newline = false
			args = args[1:]
		}
		if len(args) > 1 {
			outfile, ok := i.chans[args[0].AsString()]
			if !ok {
				return i.FailStr("wrong args")
			}
			file, ok = outfile.(io.Writer)
			if !ok {
				return i.FailStr("channel wasn't opened for writing")
			}
			args = args[1:]
		}
		s = args[0].AsString()
	}
	if newline {
		fmt.Fprintln(file, s)
	} else {
		fmt.Fprint(file, s)
	}
	return i.Return(kNil)
}


MSG

bind msg <flags> <command> <proc>
procname <nick> <user@host> <handle> <text>

Description: used for /msg commands. The first word of the user's msg is the command, and everything else becomes the text argument.

Module: server*/