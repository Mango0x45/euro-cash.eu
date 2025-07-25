package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template/parse"
)

type config struct {
	arg     int
	plural  int
	context int
	domain  int
}

type translation struct {
	msgid       string
	msgidPlural string
	msgctx      string
	domain      string
}

type transinfo struct {
	comment string
	locs    []loc
}

type loc struct {
	file string
	line int
}

func (l loc) String() string {
	return fmt.Sprintf("%s:%d", l.file, l.line)
}

var (
	outfile      *os.File
	currentFile  []byte
	currentPath  string
	lastComment  string
	translations map[translation]transinfo = make(map[translation]transinfo)
	configs                                = map[string]config{
		"Get":    {1, -1, -1, -1},
		"GetC":   {1, -1, 2, -1},
		"GetD":   {2, -1, -1, 1},
		"GetDC":  {2, -1, 3, 1},
		"GetN":   {1, 2, -1, -1},
		"GetNC":  {1, 2, 4, -1},
		"GetND":  {2, 3, -1, 1},
		"GetNDC": {2, 3, 5, 1},
	}
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] template...\n",
		filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	try(os.Chdir(filepath.Dir(os.Args[0])))

	outpath := flag.String("out", "-", "output file")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *outpath == "-" {
		outfile = os.Stdout
	} else {
		outfile = try2(os.Create(*outpath))
	}

	for _, f := range flag.Args() {
		process(f)
	}

	for tl, ti := range translations {
		if ti.comment != "" {
			fmt.Fprintln(outfile, "#.", ti.comment)
		}

		slices.SortFunc(ti.locs, func(a, b loc) int {
			if x := strings.Compare(a.file, b.file); x != 0 {
				return x
			}
			return a.line - b.line
		})
		for _, x := range ti.locs {
			fmt.Fprintln(outfile, "#:", x)
		}

		if tl.msgctx != "" {
			writeField(outfile, "msgctx", tl.msgctx)
		}
		writeField(outfile, "msgid", tl.msgid)
		if tl.msgidPlural != "" {
			writeField(outfile, "msgid_plural", tl.msgidPlural)
			writeField(outfile, "msgstr[0]", "")
			writeField(outfile, "msgstr[1]", "")
		} else {
			writeField(outfile, "msgstr", "")
		}
		fmt.Fprint(outfile, "\n")
	}
}

func process(path string) {
	currentPath = path
	currentFile = try2(os.ReadFile(path))
	trees := make(map[string]*parse.Tree)
	t := parse.New(path)
	t.Mode |= parse.ParseComments | parse.SkipFuncCheck
	try2(t.Parse(string(currentFile), "", "", trees))
	for _, t := range trees {
		processNode(t.Root)
	}
}

func processNode(node parse.Node) {
	switch n := node.(type) {
	case *parse.ListNode:
		for _, m := range n.Nodes {
			processNode(m)
		}
	case *parse.IfNode:
		processBranch(n.BranchNode)
	case *parse.RangeNode:
		processBranch(n.BranchNode)
	case *parse.WithNode:
		processBranch(n.BranchNode)
	case *parse.ActionNode:
		for _, cmd := range n.Pipe.Cmds {
			if len(cmd.Args) == 0 {
				continue
			}

			f, ok := cmd.Args[0].(*parse.FieldNode)
			if !ok || len(f.Ident) == 0 {
				continue
			}

			cfg, ok := configs[f.Ident[len(f.Ident)-1]]
			if !ok {
				continue
			}

			var (
				tl     translation
				linenr int
			)

			if sn, ok := cmd.Args[cfg.arg].(*parse.StringNode); ok {
				tl.msgid = sn.Text
				linenr = getlinenr(sn.Pos.Position())
			} else {
				continue
			}
			if cfg.plural != -1 {
				if sn, ok := cmd.Args[cfg.plural].(*parse.StringNode); ok {
					tl.msgidPlural = sn.Text
				}
			}
			if cfg.context != -1 {
				if sn, ok := cmd.Args[cfg.context].(*parse.StringNode); ok {
					tl.msgctx = sn.Text
				}
			}

			ti := translations[tl]
			if lastComment != "" {
				ti.comment = lastComment
				lastComment = ""
			}
			/* FIXME: Add filename and line number */
			ti.locs = append(ti.locs, loc{currentPath, linenr})
			translations[tl] = ti
		}
	case *parse.CommentNode:
		if strings.HasPrefix(n.Text, "/* TRANSLATORS:") {
			lastComment = strings.TrimSpace(n.Text[2 : len(n.Text)-2])
		}
	}
}

func processBranch(n parse.BranchNode) {
	processNode(n.List)
	if n.ElseList != nil {
		processNode(n.ElseList)
	}
}

func writeField(w io.Writer, pfx, s string) {
	fmt.Fprintf(w, "%s ", pfx)
	if strings.ContainsRune(s, '\n') {
		fmt.Fprintln(w, "\"\"")
		lines := strings.SplitAfter(s, "\n")
		n := len(lines)
		if n > 1 && lines[n-1] == "" {
			lines = lines[:n-1]
		}
		for _, ss := range lines {
			writeLine(w, ss)
		}
		fmt.Fprintln(w, "\"\"")
	} else {
		writeLine(w, s)
	}
}

func writeLine(w io.Writer, s string) {
	fmt.Fprint(w, "\"")
	for _, c := range s {
		switch c {
		case '\\', '"':
			fmt.Fprintf(w, "\\%c", c)
		case '\n':
			fmt.Fprint(w, "\\n")
		default:
			fmt.Fprintf(w, "%c", c)
		}
	}
	fmt.Fprintln(w, "\"")
}

func getlinenr(p parse.Pos) int {
	return bytes.Count(currentFile[:p], []byte{'\n'}) + 1
}

func try(err error) {
	if err != nil {
		die(err)
	}
}

func try2[T any](val T, err error) T {
	try(err)
	return val
}

func warn(err any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
}

func die(err any) {
	warn(err)
	os.Exit(1)
}
