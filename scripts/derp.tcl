bind pub - yo greet
bind join - #PU_HORSES welcome
bind part - #PU_HORSES goodbye 

proc sumto max {
    set sum 0
    for { set i 0 } { $i < $max } { incr i } {
        incr sum $i
    }
    return $sum
}

proc welcome { nick uhost hand chan } {
	putserv "PRIVMSG $chan :Welcome $nick!"
}

proc goodbye { nick uhost hand chan } {
	putserv "PRIVMSG $chan :Noooo, $nick! :<"
}

proc greet { nick uhost hand chan msg } {
	putserv "PRIVMSG $chan :Yo $nick"
}

puts [sumto 4]