package mintages

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
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

type coinset [8]int

type Data struct {
	StartYear       int
	Circ, BU, Proof []coinset
}

func ForCountry(code string) (Data, error) {
	path := filepath.Join("data", "mintages", code)
	f, err := os.Open(path)
	if err != nil {
		return Data{}, err
	}
	defer f.Close()
	return parse(f, path)
}

func parse(reader io.Reader, file string) (Data, error) {
	var (
		data  Data       // Our data struct
		slice *[]coinset // Where to append mintages
	)

	scanner := bufio.NewScanner(reader)
	for linenr := 1; scanner.Scan(); linenr++ {
		line := scanner.Text()
		tokens := strings.FieldsFunc(strings.TrimSpace(line), unicode.IsSpace)

		switch {
		case len(tokens) == 0:
			continue
		case tokens[0] == "BEGIN":
			if len(tokens)-1 != 1 {
				return Data{}, SyntaxError{
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
					return Data{}, SyntaxError{
						expected: "‘CIRC’, ‘BU’, ‘PROOF’, or a year",
						got:      arg,
						file:     file,
						linenr:   linenr,
					}
				}
				data.StartYear, _ = strconv.Atoi(arg)
			}
		case isNumeric(tokens[0], true), tokens[0] == "?":
			switch {
			case slice == nil:
				return Data{}, SyntaxError{
					expected: "coin type declaration",
					got:      tokens[0],
					file:     file,
					linenr:   linenr,
				}
			case data.StartYear == 0:
				return Data{}, SyntaxError{
					expected: "start year declaration",
					got:      tokens[0],
					file:     file,
					linenr:   linenr,
				}
			}

			numcoins := len(coinset{})
			tokcnt := len(tokens)

			if tokcnt != numcoins {
				word := "entries"
				if tokcnt == 1 {
					word = "entry"
				}
				return Data{}, SyntaxError{
					expected: fmt.Sprintf("%d mintage entries", numcoins),
					got:      fmt.Sprintf("%d %s", tokcnt, word),
					file:     file,
					linenr:   linenr,
				}
			}

			var row coinset
			for i, tok := range tokens {
				if tok == "?" {
					row[i] = -1
				} else {
					row[i] = atoiWithDots(tok)
				}
			}
			*slice = append(*slice, row)
		default:
			return Data{}, SyntaxError{
				expected: "‘BEGIN’ directive or mintage row",
				got:      fmt.Sprintf("invalid token ‘%s’", tokens[0]),
				file:     file,
				linenr:   linenr,
			}
		}
	}

	/* Pad rows of ‘unknown’ mintages at the end of each set of mintages
	   for each year that we haven’t filled in info for. This avoids
	   things accidentally breaking if the new year comes and we forget
	   to add extra rows. */
	for _, ms := range [...]*[]coinset{&data.Circ, &data.BU, &data.Proof} {
		finalYear := len(*ms) + data.StartYear - 1
		missing := time.Now().Year() - finalYear
		for i := 0; i < missing; i++ {
			*ms = append(*ms, coinset{-1, -1, -1, -1, -1, -1, -1, -1})
		}
	}

	return data, nil
}

func isNumeric(s string, dot bool) bool {
	for _, ch := range s {
		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			if !dot {
				return false
			}
		default:
			return false
		}
	}
	return true
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
