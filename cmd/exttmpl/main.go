package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template/parse"

	"golang.org/x/text/language"
	"golang.org/x/text/message/pipeline"
	"golang.org/x/tools/go/packages"
)

const (
	pkgbase  = "git.thomasvoss.com/euro-cash.eu"
	srclang  = "en"
	srcdir   = "./src"
	transdir = srcdir + "/rosetta"
	outfile  = "catalog.gen.go"
	transfn  = "T"
)

func main() {
	/* cd to the project root directory */
	try(os.Chdir(filepath.Dir(os.Args[0])))

	pkgnames := packageList(".")

	var paths []string
	pkgs := try2(packages.Load(&packages.Config{
		Mode: packages.NeedFiles | packages.NeedEmbedFiles,
	}, pkgnames...))

	for _, pkg := range pkgs {
		if len(pkg.Errors) != 0 {
			for _, err := range pkg.Errors {
				warn(err.Msg)
			}
			os.Exit(1)
		}
		for _, f := range pkg.EmbedFiles {
			if filepath.Ext(f) == ".tmpl" {
				paths = append(paths, f)
			}
		}
	}

	msgs := make([]pipeline.Message, 0, 1024)
	for _, path := range paths {
		f := try2(os.ReadFile(path))
		trees := make(map[string]*parse.Tree)
		t := parse.New("name")
		t.Mode |= parse.SkipFuncCheck
		try2(t.Parse(string(f), "", "", trees))
		for _, t := range trees {
			process(&msgs, t.Root)
		}
	}

	pconf := &pipeline.Config{
		Supported:      languages(),
		SourceLanguage: language.Make(srclang),
		Packages:       pkgnames,
		Dir:            transdir,
		GenFile:        outfile,
		GenPackage:     srcdir,
	}

	state := try2(pipeline.Extract(pconf))
	state.Extracted.Messages = append(state.Extracted.Messages, msgs...)

	try(state.Import())
	try(state.Merge())
	try(state.Export())
	try(state.Generate())
}

func process(tmplMsgs *[]pipeline.Message, node parse.Node) {
	switch node.Type() {
	case parse.NodeList:
		if ln, ok := node.(*parse.ListNode); ok {
			for _, n := range ln.Nodes {
				process(tmplMsgs, n)
			}
		}
	case parse.NodeIf:
		if in, ok := node.(*parse.IfNode); ok {
			process(tmplMsgs, in.List)
			if in.ElseList != nil {
				process(tmplMsgs, in.ElseList)
			}
		}
	case parse.NodeWith:
		if wn, ok := node.(*parse.WithNode); ok {
			process(tmplMsgs, wn.List)
			if wn.ElseList != nil {
				process(tmplMsgs, wn.ElseList)
			}
		}
	case parse.NodeRange:
		if rn, ok := node.(*parse.RangeNode); ok {
			process(tmplMsgs, rn.List)
			if rn.ElseList != nil {
				process(tmplMsgs, rn.ElseList)
			}
		}
	case parse.NodeAction:
		an, ok := node.(*parse.ActionNode)
		if !ok {
			break
		}

		for _, cmd := range an.Pipe.Cmds {
			if !hasIndent(cmd, transfn) {
				continue
			}
			for _, arg := range cmd.Args {
				if arg.Type() != parse.NodeString {
					continue
				}
				if sn, ok := arg.(*parse.StringNode); ok {
					txt := collapse(sn.Text)
					*tmplMsgs = append(*tmplMsgs, pipeline.Message{
						ID:      pipeline.IDList{txt},
						Key:     txt,
						Message: pipeline.Text{Msg: txt},
					})
					break
				}
			}
		}
	}
}

func hasIndent(cmd *parse.CommandNode, s string) bool {
	if len(cmd.Args) == 0 {
		return false
	}
	arg := cmd.Args[0]
	var idents []string
	switch arg.Type() {
	case parse.NodeField:
		idents = arg.(*parse.FieldNode).Ident
	case parse.NodeVariable:
		idents = arg.(*parse.VariableNode).Ident
	}
	return slices.Contains(idents, s)
}

func packageList(path string) []string {
	ents := try2(os.ReadDir(path))
	xs := make([]string, 0, len(ents))
	foundOne := false
	for _, ent := range ents {
		switch {
		case filepath.Ext(ent.Name()) == ".go":
			if !foundOne {
				xs = append(xs, pkgbase+"/"+path)
				foundOne = true
			}
		case !ent.IsDir(), ent.Name() == "cmd", ent.Name() == "vendor":
			continue
		default:
			xs = append(xs, packageList(path+"/"+ent.Name())...)
		}
	}
	return xs
}

func languages() []language.Tag {
	ents := try2(os.ReadDir(transdir))
	tags := make([]language.Tag, len(ents))
	for i, e := range ents {
		tags[i] = language.MustParse(e.Name())
	}
	return tags
}

func collapse(s string) string {
	var (
		sb   strings.Builder
		prev bool
	)
	const spc = " \t\n"

	for _, ch := range strings.Trim(s, spc) {
		if strings.ContainsRune(spc, ch) {
			if prev {
				continue
			}
			ch = ' '
			prev = true
		} else {
			prev = false
		}
		sb.WriteRune(ch)
	}

	return sb.String()
}

func try(err error) {
	if err != nil {
		die(err)
	}
}

func try2[T any](val T, err error) T {
	if err != nil {
		die(err)
	}
	return val
}

func warn(err any) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
}

func die(err any) {
	warn(err)
	os.Exit(1)
}
