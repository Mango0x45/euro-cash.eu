package mintages

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type coinset [8]int

type Data struct {
	StartYear       int
	Circ, Bu, Proof []coinset
}

func ForCountry(code string) (Data, error) {
	path := filepath.Join("data", "mintages", code)
	f, err := os.Open(path)
	if err != nil {
		return Data{}, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var (
		data  Data       // Our data struct
		slice *[]coinset // Where to append mintages
	)

	for linenr := 1; scanner.Scan(); linenr++ {
		line := scanner.Text()
		tokens := strings.FieldsFunc(strings.TrimSpace(line), unicode.IsSpace)

		switch {
		case len(tokens) == 0:
			continue
		case tokens[0] == "BEGIN":
			if len(tokens) != 2 {
				return Data{}, ArgCountMismatchError{
					token:    tokens[0],
					expected: 1,
					got:      len(tokens) - 1,
					location: location{path, linenr},
				}
			}

			arg := tokens[1]

			switch arg {
			case "CIRC":
				slice = &data.Circ
			case "BU":
				slice = &data.Bu
			case "PROOF":
				slice = &data.Proof
			default:
				if !isNumeric(arg, false) {
					return Data{}, SyntaxError{
						expected: "‘CIRC’, ‘BU’, ‘PROOF’, or a year",
						got:      arg,
						location: location{path, linenr},
					}
				}
				data.StartYear, _ = strconv.Atoi(arg)
			}
		case isNumeric(tokens[0], true):
			numcoins := len(coinset{})
			tokcnt := len(tokens)

			if tokcnt != numcoins {
				return Data{}, SyntaxError{
					expected: fmt.Sprintf("%d mintage entries", numcoins),
					got:      fmt.Sprintf("%d entries", tokcnt),
					location: location{path, linenr},
				}
			}

			var row coinset
			for i, tok := range tokens {
				row[i], _ = strconv.Atoi(strings.ReplaceAll(tok, ".", ""))
			}
			*slice = append(*slice, row)
		default:
			return Data{}, BadTokenError{
				token:    tokens[0],
				location: location{path, linenr},
			}
		}
	}

	/* Pad rows of ‘unknown’ mintages at the end of each set of mintages
	   for each year that we haven’t filled in info for. This avoids
	   things accidentally breaking if the new year comes and we forget
	   to add extra rows. */
	for _, ms := range [...]*[]coinset{&data.Circ, &data.Bu, &data.Proof} {
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
