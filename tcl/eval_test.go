package gotcl

import (
	"testing"
	"os"
	"strings"
	"io"
)

func TestFull(t *testing.T) {
	file, err := os.Open("test.tcl", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, e := NewInterp().Run(file)
	if e != nil {
		t.Fatal(e)
	}
}

func RunString(it *Interp, s string) {
	var r io.Reader = strings.NewReader(s)
	_, e := it.Run(r)
	if e != nil {
		panic(e.String())
	}
}

func runCmd(setup, cmd string, b *testing.B) {
	b.StopTimer()
	it := NewInterp()
	RunString(it, setup)
	v := FromStr(cmd)
	it.EvalObj(v)
	if it.err != nil {
		panic(it.err.String())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		it.EvalObj(v)
	}
}

func Benchmark_Plus(b *testing.B) {
	runCmd("", "+ 1 8", b)
}

func Benchmark_Plus4(b *testing.B) {
	runCmd("", "+ 1 [+ 1 [+ 1 [+ 1 8]]]", b)
}

func Benchmark_ExprPlus4(b *testing.B) {
	runCmd("", "expr { 1 + 1 + 1 + 1 + 8 }", b)
}

func Benchmark_IncrX4(b *testing.B) {
	runCmd("set x 0", "incr x; incr x; incr x; incr x", b)
}

func Benchmark_Fib(b *testing.B) {
	fib := `
proc fib {n} {
    if { $n < 2 } {
        return 1
    } else {
        return [+ [fib [- $n 1]] [fib [- $n 2]]]
    }
}
`
	runCmd(fib, "fib 17", b)
}

func Benchmark_Fib2(b *testing.B) {
	fib2 := `
proc fib2 {n} {
    set a 1
    set b 1
    for { set nn $n } { 0 < $nn } { incr nn -1 } {
        set tmp [+ $a $b]
        set a $b
        set b $tmp
    }
    return $a
}
`
	runCmd(fib2, "fib2 70", b)
}


func BenchmarkSumTo(b *testing.B) {
	sumto := `
proc sum_to {n} {
    set x 0
    for { set i 0 } { $i < $n } { incr i } {
        set x [+ $x $i]
    }
}
`
	runCmd(sumto, "sum_to 20000", b)
}

func Benchmark_SumIota(b *testing.B) {
	code := `
proc iota {n} {
    set result [list]
    for {set i 1} { $i <= $n } { incr i } {
       lappend result $i
    }
    return $result
}

proc sum {lst} {
    set result 0
    foreach x $lst {
        incr result $x
    }
    return $result
}`
	runCmd(code, "sum [iota 10000]", b)
}
