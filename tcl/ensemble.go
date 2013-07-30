package gotcl

import (
	"fmt"
	"sort"
	"strings"
)

func formatNames(sv []string) string {
	if len(sv) == 0 {
		return ""
	}
	if len(sv) == 1 {
		return sv[0]
	}
	sort.Strings(sv)
	return strings.Join(sv[0:len(sv)-1], ", ") + ", or " + sv[len(sv)-1]
}

type ensembleSpec map[string]interface{}

func (es ensembleSpec) makeCmd() TclCmd {
	cmds := make(map[string]TclCmd, len(es))
	for k, v := range es {
		cmds[k] = MakeCmd(v)
	}
	return func(i *Interp, args []*TclObj) TclStatus {
		if len(args) == 0 {
			return i.FailStr("wrong # args")
		}
		return doEnsemble(cmds, args[0].AsString(), i, args[1:])
	}
}

func doEnsemble(e map[string]TclCmd, cmd string, i *Interp, args []*TclObj) TclStatus {
	c, ok := e[cmd]
	if ok {
		return c(i, args)
	}
	sv := make([]string, len(e))
	ind := 0
	for k := range e {
		sv[ind] = k
		ind++
	}
	return i.FailStr(
		fmt.Sprintf("unknown or ambiguous subcommand \"%s\". Must be %s.", cmd, formatNames(sv)))
}
