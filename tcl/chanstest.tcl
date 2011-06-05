proc gento {max chan} {
    for {set i 0} { $i <= $max } { incr i } {
        sendchan $chan $i
    }
    closechan $chan
}

proc gen max {
    set ch [newchan]
    go gento $max $ch
    return $ch
} 

proc zip_with {fn ch1 ch2} {
    set out [newchan]
    set code { { op ch1 ch2 res } {
        forchan v $ch1 {
            set v2 [<- $ch2]
            sendchan $res [$op $v $v2]
        }
        closechan $res
    }}
    go [list apply $code $fn $ch1 $ch2 $out]
    return $out
}

proc sumchan ch {
    set res 0
    forchan v $ch {
        incr res $v
    }
    return $res
}

test {sumchan} {
    expect [sumchan [gen 100]] == 5050
}

test {zip} {
    expect [sumchan [zip_with + [gen 10] [gen 10]]] == 110
}

if false {
proc fibx {n} {
    if { $n < 2 } {
        return 1
    } else {
        return [+ [fibx [- $n 1]] [fibx [- $n 2]]]
    }
}

proc pfib {n} {
    set ch1 [newchan]
    set ch2 [newchan]
    set code { { n out } {
        sendchan $out [fibx $n]
    }}
    go [list apply $code [- $n 1] $ch1]
    go [list apply $code [- $n 2] $ch2]
    return [+ [<- $ch1] [<- $ch2]]
}

puts [time { fibx 22 }]
puts [time { pfib 22 }]
}
