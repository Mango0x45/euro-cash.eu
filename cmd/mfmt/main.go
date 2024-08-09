package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"git.thomasvoss.com/euro-cash.eu/mintages"
)

const cols = unsafe.Sizeof(mintages.Row{}.Cols) /
	unsafe.Sizeof(mintages.Row{}.Cols[0])

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
	data, err := mintages.Parse(in, path)
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

func formatSection(out io.Writer, rows []mintages.Row) {
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
			longestNum[i] = max(longestNum[i], intlen(col))
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
	if n < -1 {
		panic("mintage count is negative and not -1")
	} else if n == -1 {
		return "?"
	}

	in := strconv.Itoa(n)
	dots := (len(in) - 1) / 3
	out := make([]byte, len(in)+dots)

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = '.'
		}
	}
}

func intlen(v int) int {
	switch {
	case v < -1:
		panic("mintage count is negative and not -1")
	case v == 0, v == -1:
		return 1
	default:
		n := 0
		for x := v; x != 0; x /= 10 {
			n++
		}
		return n + (n-1)/3 /* Thousands separators */
	}
}

func warn(e error) {
	fmt.Fprintf(os.Stderr, "%s: %+v\n", filepath.Base(os.Args[0]), e)
	rv = 1
}
