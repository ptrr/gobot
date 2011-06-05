proc sumto max {
    set sum 0
    for { set i 0 } { $i < $max } { incr i } {
        incr sum $i
    }
    return $sum
}

puts [sumto 4]