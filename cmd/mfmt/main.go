/* Simple formatter for mintage files.  This is not perfect and doesn’t
   check for syntactic correctness, it’s just to get stuff aligned
   nicely.  Maybe in the future I will construct a military-grade
   formatter. */

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ws         = " \t"
	longestNum = len("1.000.000.000")
)

var (
	rv int

	reMintageYear = regexp.MustCompile(`^\d{4}(-[^ \t]+)?`)
	reMintageRow  = regexp.MustCompile(`^(([0-9.]+|\?)[ \t]+)*([0-9.]+|\?)$`)
)

func main() {
	if len(os.Args) == 1 {
		mfmt("-", os.Stdin, os.Stdout)
	}
	for _, arg := range os.Args[1:] {
		f, err := os.OpenFile(arg, os.O_RDWR, 0)
		if err != nil {
			warn(err)
			continue
		}
		defer f.Close()
		mfmt(arg, f, f)
	}
	os.Exit(rv)
}

func mfmt(file string, r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for linenr := 1; scanner.Scan(); linenr++ {
		line := strings.Trim(scanner.Text(), ws)

		switch {
		case len(line) == 0, line[0] == '#':
			fmt.Fprintln(w, line)
		case reMintageYear.MatchString(line):
			fmtMintageYear(line, w)
		case reMintageRow.MatchString(line):
			fmtMintageRow(line, w)
		default:
			warn(fmt.Sprintf("%s:%d: potential syntax error", file, linenr))
			fmt.Fprintln(w, line)
		}
	}
}

func fmtMintageYear(line string, w io.Writer) {
	switch i := strings.IndexAny(line, ws); i {
	case -1:
		fmt.Fprintln(w, line)
	default:
		fmt.Fprintf(w, "%s %s\n", line[:i], strings.TrimLeft(line[i:], ws))
	}
}

func fmtMintageRow(line string, w io.Writer) {
	tokens := strings.FieldsFunc(line, func(r rune) bool {
		return strings.ContainsRune(ws, r)
	})

	for i, tok := range tokens {
		s := formatMintage(tok)

		if i == 0 {
			fmt.Fprintf(w, "\t%*s", longestNum, s)
		} else {
			fmt.Fprintf(w, "%*s", longestNum+1, s)
		}
	}

	fmt.Fprintln(w)
}

func formatMintage(s string) string {
	if s == "?" {
		return s
	}

	n := atoiWithDots(s)
	digits := intlen(n)
	dots := (digits - 1) / 3
	out := make([]byte, digits+dots)

	for i, j := len(out)-1, 0; ; i-- {
		out[i] = byte(n%10) + 48
		if n /= 10; n == 0 {
			return string(out)
		}
		if j++; j == 3 {
			i, j = i-1, 0
			out[i] = '.'
		}
	}
}

func intlen(v int) int {
	switch {
	case v == 0:
		return 1
	default:
		n := 0
		for x := v; x != 0; x /= 10 {
			n++
		}
		return n
	}
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

func warn(e any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), e)
	rv = 1
}
