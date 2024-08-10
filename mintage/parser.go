package mintage

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	_       = -iota
	Unknown // Unknown mintage
	Invalid // All mintages <= than this are invalid
)

type SyntaxError struct {
	expected, got string
	file          string
	linenr        int
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s:%d: syntax error: expected %s but got %s",
		e.file, e.linenr, e.expected, e.got)
}

type Row struct {
	Year     int
	Mintmark string
	Cols     [8]int
}

type Set struct {
	StartYear       int
	Circ, BU, Proof []Row
}

func (r Row) Label() string {
	if r.Mintmark != "" {
		return fmt.Sprintf("%d %s", r.Year, r.Mintmark)
	}
	return strconv.Itoa(r.Year)
}

func Parse(reader io.Reader, file string) (Set, error) {
	var (
		data  Set    // Our data struct
		slice *[]Row // Where to append mintages
	)

	scanner := bufio.NewScanner(reader)
	for linenr := 1; scanner.Scan(); linenr++ {
		var mintmark struct {
			s    string
			star bool
		}

		line := scanner.Text()
		tokens := strings.FieldsFunc(strings.TrimSpace(line), unicode.IsSpace)

		switch {
		case len(tokens) == 0:
			continue
		case tokens[0] == "BEGIN":
			if len(tokens)-1 != 1 {
				return Set{}, SyntaxError{
					expected: "single argument to ‘BEGIN’",
					got:      fmt.Sprintf("%d arguments", len(tokens)-1),
					file:     file,
					linenr:   linenr,
				}
			}

			arg := tokens[1]

			switch arg {
			case "CIRC":
				slice = &data.Circ
			case "BU":
				slice = &data.BU
			case "PROOF":
				slice = &data.Proof
			default:
				if !isNumeric(arg, false) {
					return Set{}, SyntaxError{
						expected: "‘CIRC’, ‘BU’, ‘PROOF’, or a year",
						got:      arg,
						file:     file,
						linenr:   linenr,
					}
				}
				data.StartYear, _ = strconv.Atoi(arg)
			}
		case isLabel(tokens[0]):
			n := len(tokens[0])
			if n > 2 && tokens[0][n-2] == '*' {
				mintmark.star = true
				mintmark.s = tokens[0][:n-2]
			} else {
				mintmark.s = tokens[0][:n-1]
			}
			tokens = tokens[1:]
			if !isNumeric(tokens[0], true) && tokens[0] != "?" {
				return Set{}, SyntaxError{
					expected: "mintage row after label",
					got:      tokens[0],
					file:     file,
					linenr:   linenr,
				}
			}
			fallthrough
		case isNumeric(tokens[0], true), tokens[0] == "?":
			switch {
			case slice == nil:
				return Set{}, SyntaxError{
					expected: "coin type declaration",
					got:      tokens[0],
					file:     file,
					linenr:   linenr,
				}
			case data.StartYear == 0:
				return Set{}, SyntaxError{
					expected: "start year declaration",
					got:      tokens[0],
					file:     file,
					linenr:   linenr,
				}
			}

			numcoins := len(Row{}.Cols)
			tokcnt := len(tokens)

			if tokcnt != numcoins {
				word := "entries"
				if tokcnt == 1 {
					word = "entry"
				}
				return Set{}, SyntaxError{
					expected: fmt.Sprintf("%d mintage entries", numcoins),
					got:      fmt.Sprintf("%d %s", tokcnt, word),
					file:     file,
					linenr:   linenr,
				}
			}

			row := Row{Mintmark: mintmark.s}
			if len(*slice) == 0 {
				row.Year = data.StartYear
			} else {
				row.Year = (*slice)[len(*slice)-1].Year
				if row.Mintmark == "" || mintmark.star {
					row.Year++
				}
			}

			for i, tok := range tokens {
				if tok == "?" {
					row.Cols[i] = Unknown
				} else {
					row.Cols[i] = atoiWithDots(tok)
				}
			}
			*slice = append(*slice, row)
		default:
			return Set{}, SyntaxError{
				expected: "‘BEGIN’ directive or mintage row",
				got:      fmt.Sprintf("invalid token ‘%s’", tokens[0]),
				file:     file,
				linenr:   linenr,
			}
		}
	}

	return data, nil
}

func isNumeric(s string, dot bool) bool {
	for _, ch := range s {
		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			if ch != '.' || !dot {
				return false
			}
		}
	}
	return true
}

func isLabel(s string) bool {
	n := len(s)
	return (n > 2 && s[n-1] == ':' && s[n-2] == '*') ||
		(n > 1 && s[n-1] == ':')
}

func atoiWithDots(s string) int {
	n := 0
	for _, ch := range s {
		if ch == '.' {
			continue
		}
		n = n*10 + int(ch) - '0'
	}
	return n
}
