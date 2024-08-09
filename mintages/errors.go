package mintages

import "fmt"

type location struct {
	file   string
	linenr int
}

func (loc location) String() string {
	return fmt.Sprintf("%s: %d", loc.file, loc.linenr)
}

type SyntaxError struct {
	location
	expected, got string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s: syntax error: expected %s but got %s",
		e.location, e.expected, e.got)
}
