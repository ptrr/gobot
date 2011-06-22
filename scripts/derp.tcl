bind pub - yo greet

proc sumto max {
    set sum 0
    for { set i 0 } { $i < $max } { incr i } {
        incr sum $i
    }
    return $sum
}

proc greet { nick uhost hand chan msg } {
	putserv "PRIVMSG $chan :yo $nick"
}

puts [sumto 4]