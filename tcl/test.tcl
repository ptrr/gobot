source testlib.tcl

proc has_command cmd_name {
    set matching [llength [info commands $cmd_name]]
    return [expr "$matching > 0"]
}

if {![has_command +]} {
    proc + {a b} {
        return [expr {$a + $b}]
    }
}
if {![has_command -]} {
    proc - {a b} {
        return [expr {$a - $b}]
    }
}

if {[has_command go]} {
    source chanstest.tcl
}


set x 12
proc foo {x} {
    return [expr {1 + 1 + 1 + $x}]
}

assert [foo 3] == 6 "globals unused"
assert $x == 12 "locals unchanged"

test {same line} {
    set x 0
    incr x;incr x;  incr x   ; incr x
    assert $x == 4
}

test {has_command} {
    assert_noerr {
        foreach x [info commands] {
            if {![has_command $x]} {
                error "mismatch"
            }
        }
    }
}

test {set test} {
    assert [set x 4] == 4
    assert [set x] == 4
}

test {if test} {
    set x ""
    if { 2 < 10 } then {
        set x yes
    } else {
        set x no
    }
    assert $x == yes

    set y ""
    if { 10 < 2 } {
        set y no
    } else {
        set y yes
    }
    assert $y == yes
}

set fizzle 10
test global_test {
    assert $::fizzle == 10

    set ::fizzle 55
    assert $::fizzle == 55
    assert ${::fizzle} == 55
}

test {incr test} {
    set x 10
    incr x
    incr x
    incr x -1
    assert $x == 11
}

test {incr return} {
    set x 5
    assert [incr x] == 6
}

test {return in string} {
    proc somereturn {} {
        set x "foo [return ok]"
        return failed
    }
    assert [somereturn] == ok 
}

test for_test {
    set res 0
    for { set x 0 } { $x < 10 } { incr x } {
        set res [+ $res $x]
    }
    assert $res == 45
}

test {while test} {
    set x 0
    assert_noerr {
        while {$x > 2} {
            error "failed."
        }
    }
    while {$x < 10} {
        incr x
    }
    assert $x == 10
    while {$x >= 0} {
        break
        incr x -1
    }
    assert $x == 10
    while {$x > 0} {
        incr x -1
        continue
        error "failed"
    }
    assert $x == 0
}

test break_test {
    set val ""
    for { set x 0 } { $x < 5 } { incr x } {
        set val yes
        break
        set val fail
    }
    assert $val == yes
}

test {continue test} {
    set val ""
    for { set x 0 } { $x < 5 } { incr x } {
        set val yes
        continue
        set val fail
    }
    assert $val == yes
}

test {command with spaces} {
    set x [+ 1 2 ]
    assert $x == 3
}

test {empty subcommand} {
    assert "foo[ ]bar" == "foobar"
}

test {list procs} {
    set x [list 1 2 3 4 5]
    assert [llength $x] == 5
    assert [lindex $x 0] == 1
    assert [lindex $x 3] == 4
    assert [llength [list]] == 0
}

test {list parse} {
    set x [list {one two} three]
    assert [llength $x] == 2
    set y "$x four"
    assert $y == "{one two} three four"
    assert [llength $y] == 3
    assert [llength { xxx{}xxx }] == 1
}

test {list parse with commands} {
    set x { [ + 4 5 ] }
    assert [llength $x] == 5
}

test {nested list parse} {
    set x { 1 2 { x 3 } }
    assert [llength $x] == 3
    assert [lindex $x 2] == { x 3 }
    assert [llength [lindex $x 2]] == 2
}

test {weird parse} {
    set x(x 44
    assert ${x(x} == 44
    set y{}y 12
    # assert $y{}y == 12
}

test {list with empties} {
    set x [list {} {} {}]
    assert $x == {{} {} {}}
}

test {lappend} {
    lappend x 0
    lappend x 1 2 3 4 5
    assert [llength $x] == 6
}

test {lappend more} {
    set x [list 1 2]
    set y $x
    set z $x
    lappend y 3 4 5
    lappend z 9 9 9
    assert $y != $z
}

test {args} {
    proc count_args {args} {
        return [llength $args]
    }
    assert [count_args 1 2 3] == 3
}

test {bad proc} {
    proc fizzle {x} { " }
    set ec [catch { fizzle 4 } msg]
    assert $ec == 1
    assert [string index $msg end] == {"}
}

test unset_test {
    set foo 11
    unset foo
    assert [catch { puts $foo }] == 1
}

test {catch sets resultvar} {
    catch { expr { 3 + 4 } } woozle
    assert [info exists woozle] == 1
    assert $woozle == 7
}

test test_time {
    set x 0
    time { set x [+ $x 1] } 4
    assert $x == 4 "time runs multiple"
}

test {escape test} {
    assert " \[ whee \] " == { [ whee ] }
    assert " \" " == { " }
        # " to make vim happy
}

test {escaped word} {
    set a \n
    set b "\n"
    assert $a == $b
}

test {uplevel test} {
    proc incrvar {vn} {
        uplevel "incr $vn"
    }
    set x 10
    incrvar x
    assert $x == 11
}

test {upvar} {
    proc add {vn} {
        upvar $vn x
        incr x
    }
    set x 0
    add x
    assert $x == 1
}

test {upvar multi} {
    proc proc2 {} {
        upvar foo zz
        incr zz 2
    }
    proc proc1 {} {
        upvar x foo
        proc2
    }

    set x 0
    proc1
    assert $x == 2
}

test {double upvar} {
    proc proc2 {} {
        upvar 2 foo zz
        incr zz 2
    }
    proc proc1 {} {
        proc2
    }

    set foo 0
    proc1
    assert $foo == 2
}

test {default arg} {
    proc foo { { x 1 }  { y 0 } } {
        return [+ $x $y]
    }
    assert [foo] == 1
    assert [foo 5] == 5
    assert [foo 5 5] == 10
}


test {foreach break} {
    set x {yes no no}
    set y no
    foreach i $x {
        set y $i
        break
    }
    assert $y == yes
}

test {foreach continue} {
    set x {1 2 3}
    set y no
    foreach i $x {
        set y yes
        continue
        set y no
    }
    assert $y == yes
}

test {foreach var remains} {
    foreach x {1 2} {
    }
    assert $x == 2
}

test {string length} {
    assert [string length ""] == 0
    assert [string length "xxx"] == 3
    assert [string bytelength "xxx"] == 3
    assert [string length "世界"] == 2
    assert [string bytelength "世界"] == 6
}

test {string index} {
    assert [string index "" 4] == ""
    assert [string index "abcdefg" 0] == "a"
    assert [string index "abcdefg" 2] == "c"
}

test {string trim} {
    assert [string trim " X "] == "X"
    assert [string trim "  "] == ""
    assert [string trim foo] == foo
}

test {string match} {
    assert [string match a?c abc] == 1
    assert [string match ab* abcde] == 1
}

test {split} {
    assert [split "a  b  c"] == [list a {} b {} c]
}

test {split with chars} {
    assert [split "axxbyc" "xy"] == [list a {} b c]
}

test {split empty} {
    assert [split "abc" ""] == "a b c"
}

test {concat} {
    assert [concat "  a  " "  b  "] == "a b"
    assert [concat " a b {c   " d " e} f"] == "a b {c d e} f"
}

test {info exists} {
    assert [info exists candycane] == 0
    set x 4
    assert [info exists x] == 1
    unset x
    assert [info exists x] == 0
}

test {info cmdcount} {
    set x [info cmdcount]
    set y [info cmdcount]
    assert $x < $y
}

test {lsearch} {
    assert [lsearch {a b c d} b] == 1
    assert [lsearch {a b c d} z] == -1
}

test {rename to delete} {
    proc fizzlebuggy {} {}
    assert [lsearch [info commands] fizzlebuggy] > -1
    rename fizzlebuggy wheatgerm
    assert [lsearch [info commands] fizzlebuggy] == -1
    assert [lsearch [info commands] wheatgerm] > -1
    rename wheatgerm {}
    assert [lsearch [info commands] wheatgerm] == -1
    assert [lsearch [info commands] fizzlebuggy] == -1
}

test {expr string eq ne} {
    assert [expr {"foo" eq "foo"}] == 1
    assert [expr {"foo" ne "foo"}] == 0
    assert [expr {"foo" ne "roo"}] == 1
    assert [expr {"foo" ne "foo"}] == 0
}

test {if (true|false|no)} {
    set x boo
    if true { set x ok }
    if false then { set x fail }
    if no then { set x fail }
    assert $x == ok
}

test {apply} {
    assert [apply {{x} { incr x }} 4] == 5
    assert [apply {{x} { return [- $x 1] }} 4] == 3
    assert [apply {{} { return 99 }}] == 99
}

test {expr} {
    assert [expr 1 + 1] == 2
    assert [expr "1 + 1"] == 2
    assert [expr "1 +" 1] == 2
    assert [expr {1 + 1}] == 2
    set x 10
    assert [expr { $x + $x }] == 20
    assert [expr { 2 * $x }] == 20
    assert [expr { $x - -4 }] == 14
    set y 2
    assert [expr { ($y*$y) + $y }] == 6
    assert [expr {$y*$y}] == 4

    assert [expr { 2 ^ 2 }] == 0
    assert [expr { 1 << 2 }] == 4
    assert [expr { 4 >> 2 }] == 1
    assert [expr { ~0 }] == -1
}

test {expr /} {
    assert [expr { 10 / 2 }] == 5
}

test {expr ternary if} {
    assert [expr { false ? 11 : 44 }] == 44
    assert [expr { (99 == 99) ? "yay" : "boo" }] == "yay"
    assert [ expr { 1 < 2 ? 11 : 44 }] == 11
}

test {foreach trick} {
    foreach { a b } { 1 2 } break
    expect $a == 1
    expect $b == 2
}

test {array syntax} {
    set x(14) yes
    set x(15) yes2
    expect $x(14) == yes
    expect $x(15) == yes2
    expect [info exists x] == 1
    expect [info exists x(14)] == 1 "x(14)"
    expect [info exists x(15)] == 1 "x(15)"
    expect [info exists x(16)] == 0 "x(16)"
}

test {array vars} {
    set x(0) "foo"
    assert_err {
        puts $x
    }
    set y(14) 2
    assert_err {
        incr y
    }
}

test {array size} {
    set x(0) 1
    set x(1) 1
    set x(2) 1
    expect [array size x] == 3
}

test {array get} {
    set x(10) 20
    expect [array get x] == {10 20}
}

test {array set} {
    assert_err {
        array set foobar { 0 zero 1 }
    }
    array set whee { 0 zero 1 one 2 two }
    expect $whee(0) == "zero"
    expect $whee(1) == "one"
    expect $whee(2) == "two"
}

test {array exsits} {
    expect [array exists wharg] == 0
    array set wharg {a b c d}
    expect [array exists wharg] == 1
    set x 11
    expect [array exists x] == 0
}

test { expand syntax } {
    set ll {x yes}
    set x no
    set {*}$ll
    expect $x == yes
    expect [list 0 {*}{1 2} 3 {*}{ 4 5 } 6] == {0 1 2 3 4 5 6}
}

test {expand with bad list} {
    assert_err {
    list {*}{ " }
    # "
    }
}

test {for with bad stuff} {
    assert_err {
        for {set i 0} { $i < 10 } { error "boo" } {
            incr i
        }
    }
}

test {lazy || and && in expr} {
    set x yay
    expr { true || [set x boo] }
    expect $x == yay
    set y yay
    expr { false && [set y boo] }
    expect $y == yay
}

proc nothing args {}

test { list parsing again } {
    set x { nothing "["nothing"]" nothing }
    assert_noerr {
        eval $x
    }
    assert_err {
        llength $x
    }
}


proc fib {n} {
    if { $n < 2 } {
        return 1
    } else {
        return [+ [fib [- $n 1]] [fib [- $n 2]]]
    }
}

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
}

assert [iota 4] == {1 2 3 4}
assert [sum [iota 2]] == 3
assert [sum [iota 2]] == 3

assert [fib 8] == 34
assert [fib2 10] == 89

proc sum_to {n} {
    set x 0
    for { set i 0 } { $i < $n } { incr i } {
        set x [+ $x 1]
    }
}

puts "\nPassed $::passcount/$::assertcount assertions."
if false {
    puts "\n----Benchmarks----"

    bench { fib 17 }
    bench { fib2 70 }
    bench { sum_to 20000 }
    bench { iota 10000 } 
    bench { sum [iota 10000] }
}
