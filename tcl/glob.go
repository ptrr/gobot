package gotcl

import "utf8"

func uncons(s string) (int, string) {
	head, sz := utf8.DecodeRuneInString(s)
	if head == utf8.RuneError {
		return head, ""
	}
	return head, s[sz:]
}

func matchcharset(pat, strin string) (bool, string, string) {
	if strin == "" {
		return false, "", ""
	}
	sh, str := uncons(strin)
	ph, rest := uncons(pat)
	got_match := false
	for ph != ']' && ph != utf8.RuneError {
		if !got_match {
			if sh == ph {
				got_match = true
			} else if rh, rt := uncons(rest); rh == '-' {
				var ph2 int
				ph2, rest = uncons(rt)
				if ph2 == utf8.RuneError {
					return false, "", ""
				}
				got_match = sh <= ph2 && sh >= ph
			}
		}
		ph, rest = uncons(rest)
	}
	return got_match, rest, str
}

func GlobMatch(pat, str string) bool {
	for pat != "" {
		ph, rest := uncons(pat)
		switch ph {
		case '?':
			if str == "" {
				return false
			}
			_, str = uncons(str)
		case '[':
			is_match := false
			is_match, rest, str = matchcharset(rest, str)
			if !is_match {
				return false
			}
		case '*':
			if rest == "" {
				return true
			}
			for ; str != ""; _, str = uncons(str) {
				if GlobMatch(rest, str) {
					return true
				}
			}
			return false
		default:
			if str == "" {
				return false
			}
			if ph == '\\' {
				if rest == "" {
					return false
				}
				ph, rest = uncons(rest)
			}
			var sh int
			sh, str = uncons(str)
			if sh != ph {
				return false
			}
		}
		pat = rest
	}
	return str == ""
}
