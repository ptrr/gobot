bind pub - yo greet

proc sumto max {
    set sum 0
    for { set i 0 } { $i < $max } { incr i } {
        incr sum $i
    }
    return $sum
}

proc greet { nick uhost hand chan msg } {
	puts "Derp called: "
	puts $nick
	puts $msg
}

puts [sumto 4]