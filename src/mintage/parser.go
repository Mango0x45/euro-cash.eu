package mintage

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type SyntaxError struct {
	expected, got string
	file          string
	linenr        int
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s:%d: syntax error: expected %s but got %s",
		e.file, e.linenr, e.expected, e.got)
}

type Data struct {
	Standard      []SRow
	Commemorative []CRow
}

type SRow struct {
	Year     int
	Mintmark string
	Mintages [typeCount][denoms]int
}

type CRow struct {
	Year     int
	Mintmark string
	Name     string
	Mintage  [typeCount]int
}

const (
	TypeCirc = iota
	TypeNIFC
	TypeProof
	typeCount
)

const (
	Unknown = -iota - 1
	Invalid
)

const (
	denoms = 8
	ws     = " \t"
)

func Parse(r io.Reader, file string) (Data, error) {
	yearsSince := time.Now().Year() - 1999 + 1
	data := Data{
		Standard:      make([]SRow, 0, yearsSince),
		Commemorative: make([]CRow, 0, yearsSince),
	}

	scanner := bufio.NewScanner(r)
	for linenr := 1; scanner.Scan(); linenr++ {
		line := strings.Trim(scanner.Text(), ws)
		if isBlankOrComment(line) {
			continue
		}

		if len(line) < 4 || !isNumeric(line[:4], false) {
			return Data{}, SyntaxError{
				expected: "4-digit year",
				got:      line,
				linenr:   linenr,
				file:     file,
			}
		}

		var (
			commem   bool
			mintmark string
		)
		year, _ := strconv.Atoi(line[:4])
		line = line[4:]

		if len(line) != 0 {
			if strings.ContainsRune(ws, rune(line[0])) {
				commem = true
				goto out
			}
			if line[0] != '-' {
				return Data{}, SyntaxError{
					expected: "end-of-line or mintmark",
					got:      line,
					linenr:   linenr,
					file:     file,
				}
			}

			if line = line[1:]; len(line) == 0 {
				return Data{}, SyntaxError{
					expected: "mintmark",
					got:      "end-of-line",
					linenr:   linenr,
					file:     file,
				}
			}

			switch i := strings.IndexAny(line, ws); i {
			case 0:
				return Data{}, SyntaxError{
					expected: "mintmark",
					got:      "whitespace",
					linenr:   linenr,
					file:     file,
				}
			case -1:
				mintmark = line
			default:
				mintmark, line = line[:i], line[i:]
				commem = true
			}
		}
	out:

		if !commem {
			row := SRow{
				Year:     year,
				Mintmark: mintmark,
			}
			for i := range row.Mintages {
				line = ""
				for isBlankOrComment(line) {
					if !scanner.Scan() {
						return Data{}, SyntaxError{
							expected: "mintage row",
							got:      "end-of-file",
							linenr:   linenr,
							file:     file,
						}
					}
					line = strings.Trim(scanner.Text(), ws)
					linenr++
				}

				tokens := strings.FieldsFunc(line, func(r rune) bool {
					return strings.ContainsRune(ws, r)
				})
				if tokcnt := len(tokens); tokcnt != denoms {
					word := "entries"
					if tokcnt == 1 {
						word = "entry"
					}
					return Data{}, SyntaxError{
						expected: fmt.Sprintf("%d mintage entries", denoms),
						got:      fmt.Sprintf("%d %s", tokcnt, word),
						linenr:   linenr,
						file:     file,
					}
				}

				for j, tok := range tokens {
					if tok != "?" && !isNumeric(tok, true) {
						return Data{}, SyntaxError{
							expected: "numeric mintage figure or ‘?’",
							got:      tok,
							linenr:   linenr,
							file:     file,
						}
					}

					if tok == "?" {
						row.Mintages[i][j] = Unknown
					} else {
						row.Mintages[i][j] = atoiWithDots(tok)
					}
				}
			}

			data.Standard = append(data.Standard, row)
		} else {
			row := CRow{
				Year:     year,
				Mintmark: mintmark,
			}
			line = strings.TrimLeft(line, ws)
			if line[0] != '"' {
				return Data{}, SyntaxError{
					expected: "string",
					got:      line,
					linenr:   linenr,
					file:     file,
				}
			}

			line = line[1:]
			switch i := strings.IndexByte(line, '"'); i {
			case -1:
				return Data{}, SyntaxError{
					expected: "closing quote",
					got:      "end-of-line",
					linenr:   linenr,
					file:     file,
				}
			case 0:
				return Data{}, SyntaxError{
					expected: "commemorated event",
					got:      "empty string",
					linenr:   linenr,
					file:     file,
				}
			default:
				row.Name, line = line[:i], line[i+1:]
			}

			if len(line) != 0 {
				return Data{}, SyntaxError{
					expected: "end-of-line",
					got:      line,
					linenr:   linenr,
					file:     file,
				}
			}

			for isBlankOrComment(line) {
				if !scanner.Scan() {
					return Data{}, SyntaxError{
						expected: "mintage row",
						got:      "end-of-file",
						linenr:   linenr,
						file:     file,
					}
				}
				line = strings.Trim(scanner.Text(), ws)
				linenr++
			}

			tokens := strings.FieldsFunc(line, func(r rune) bool {
				return strings.ContainsRune(ws, r)
			})
			if tokcnt := len(tokens); tokcnt != typeCount {
				word := "entries"
				if tokcnt == 1 {
					word = "entry"
				}
				return Data{}, SyntaxError{
					expected: fmt.Sprintf("%d mintage entries", typeCount),
					got:      fmt.Sprintf("%d %s", tokcnt, word),
					linenr:   linenr,
					file:     file,
				}
			}

			for i, tok := range tokens {
				if tok == "?" {
					row.Mintage[i] = Unknown
				} else {
					row.Mintage[i] = atoiWithDots(tok)
				}
			}

			data.Commemorative = append(data.Commemorative, row)
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

func isBlankOrComment(s string) bool {
	return len(s) == 0 || s[0] == '#'
}
