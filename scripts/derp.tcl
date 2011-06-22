bind pub - yo greet
bind join - #PU_HORSES welcome

proc sumto max {
    set sum 0
    for { set i 0 } { $i < $max } { incr i } {
        incr sum $i
    }
    return $sum
}

proc welcome { nick uhost hand chan } {
	putserv "PRIVMSG $chan :welcome $nick!"
}

proc greet { nick uhost hand chan msg } {
	putserv "PRIVMSG $chan :yo $nick"
}

puts [sumto 4]