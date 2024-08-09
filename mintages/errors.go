package mintages

import "fmt"

type location struct {
	file   string
	linenr int
}

func (loc location) String() string {
	return fmt.Sprintf("%s: %d", loc.file, loc.linenr)
}

type BadTokenError struct {
	location
	token string
}

func (e BadTokenError) Error() string {
	return fmt.Sprintf("%s: unknown token ‘%s’", e.location, e.linenr, e.token)
}

type ArgCountMismatchError struct {
	location
	token         string
	expected, got int
}

func (e ArgCountMismatchError) Error() string {
	var suffix string
	if e.expected != 1 {
		suffix = "s"
	}
	return fmt.Sprintf("%s: ‘%s’ token expects %d argument%s but got %d",
		e.location, e.token, e.expected, suffix, e.got)
}

type SyntaxError struct {
	location
	expected, got string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s: syntax error: expected %s but got %s",
		e.location, e.expected, e.got)
}
