package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"git.thomasvoss.com/euro-cash.eu/mintage"
)

const cols = unsafe.Sizeof(mintage.Row{}.Cols) /
	unsafe.Sizeof(mintage.Row{}.Cols[0])

var rv int

func main() {
	if len(os.Args) == 1 {
		if err := mfmt("-", os.Stdin, os.Stdout); err != nil {
			warn(err)
		}
	}
	for _, arg := range os.Args[1:] {
		f, err := os.OpenFile(arg, os.O_RDWR, 0)
		if err != nil {
			warn(err)
			continue
		}
		defer f.Close()
		if err := mfmt(arg, f, f); err != nil {
			warn(err)
		}
	}
	os.Exit(rv)
}

func mfmt(path string, in io.Reader, out io.Writer) error {
	data, err := mintage.Parse(in, path)
	if err != nil {
		return err
	}

	f, outfile := out.(*os.File)
	if outfile {
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(out, `BEGIN %d`, data.StartYear)
	if len(data.Circ) != 0 {
		fmt.Fprintln(out, "\n\nBEGIN CIRC")
		formatSection(out, data.Circ)
	}
	if len(data.BU) != 0 {
		fmt.Fprintln(out, "\n\nBEGIN BU")
		formatSection(out, data.BU)
	}
	if len(data.Proof) != 0 {
		fmt.Fprintln(out, "\n\nBEGIN PROOF")
		formatSection(out, data.Proof)
	}
	fmt.Fprintln(out, "")

	if outfile {
		if off, err := f.Seek(0, io.SeekCurrent); err != nil {
			return err
		} else if err = f.Truncate(off); err != nil {
			return err
		}
	}
	return nil
}

func formatSection(out io.Writer, rows []mintage.Row) {
	var (
		year       string
		longestMM  int
		longestNum [cols]int
	)

	for _, row := range rows {
		s, mm, ok := strings.Cut(row.Label, "\u00A0")
		if ok {
			n := len(mm)
			if s != year {
				n++
			}
			longestMM = max(longestMM, n)
		}
		year = s
		for i, col := range row.Cols {
			n := intlen(col)
			longestNum[i] = max(longestNum[i], n+(n-1)/3)
		}
	}

	extraSpace := 0
	if longestMM != 0 {
		extraSpace = 2
	}

	year = ""

	for i, row := range rows {
		var label string

		s, mm, ok := strings.Cut(row.Label, "\u00A0")
		if ok {
			if s != year {
				label = mm + "*: "
			} else {
				label = mm + ": "
			}
		}
		year = s

		fmt.Fprintf(out, "%-*s", longestMM+extraSpace, label)
		for j, n := range row.Cols {
			if j != 0 {
				fmt.Fprintf(out, " ")
			}
			fmt.Fprintf(out, "%*s", longestNum[j], formatInt(n))
		}

		if i != len(rows)-1 {
			fmt.Fprintln(out, "")
		}
	}
}

func formatInt(n int) string {
	if n <= mintage.Invalid {
		panic(fmt.Sprintf("invalid input %d", n))
	} else if n == mintage.Unknown {
		return "?"
	}

	digits := intlen(n)
	dots := (digits - 1) / 3
	out := make([]byte, digits+dots)

	for i, j := len(out)-1, 0; ; i-- {
		out[i] = byte(n%10) + 48
		n /= 10
		if n == 0 {
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
	case v <= mintage.Invalid:
		panic("mintage count is negative and not -1")
	case v == 0, v == mintage.Unknown:
		return 1
	default:
		n := 0
		for x := v; x != 0; x /= 10 {
			n++
		}
		return n
	}
}

func warn(e error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), e)
	rv = 1
}
