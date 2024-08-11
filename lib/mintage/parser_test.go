package mintage

import (
	"bytes"
	"errors"
	"testing"
)

func TestParserComplete(t *testing.T) {
	data, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN 2020
		BEGIN CIRC
		1.000 1001 1002 1003 1004 1005 1006 1007
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106 2107
		BEGIN PROOF
		1.200 1201 1202 1203 1204 1205 1206 1207
		 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	if err != nil {
		t.Fatalf(`Expected err=nil; got "%s"`, err)
	}
	if data.StartYear != 2020 {
		t.Fatalf("Expected data.StartYear=2020; got %d",
			data.StartYear)
	}

	for i, row := range data.Circ {
		for j, col := range row.Cols {
			var n int
			if i == 1 && j == 1 {
				n = -1
			} else {
				n = 1000*i + j + 1000
			}
			if col != n {
				t.Fatalf("Expected data.Circ[i][j]=%d; got %d", n, col)
			}
		}
	}

	for i, row := range data.BU {
		for j, col := range row.Cols {
			var n int
			if i == 1 && j == 1 {
				n = -1
			} else {
				n = 1000*i + j + 1100
			}
			if col != n {
				t.Fatalf("Expected data.BU[i][j]=%d; got %d", n, col)
			}
		}
	}

	for i, row := range data.Proof {
		for j, col := range row.Cols {
			var n int
			if i == 1 && j == 1 {
				n = -1
			} else {
				n = 1000*i + j + 1200
			}
			if col != n {
				t.Fatalf("Expected data.Proof[i][j]=%d; got %d", n, col)
			}
		}
	}

	if len(data.Circ) != 2 {
		t.Fatalf("Expected len(data.Circ)=2; got %d", len(data.Circ))
	}
	if len(data.BU) != 2 {
		t.Fatalf("Expected len(data.BU)=2; got %d", len(data.BU))
	}
	if len(data.Proof) != 2 {
		t.Fatalf("Expected len(data.Proof)=2; got %d", len(data.Proof))
	}
}

func TestParserNoProof(t *testing.T) {
	data, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN 2020
		BEGIN CIRC
		1.000 1001 1002 1003 1004 1005 1006 1007
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106 2107
	`)), "-")

	if err != nil {
		t.Fatalf(`Expected err=nil; got "%s"`, err)
	}

	if len(data.Proof) != 0 {
		t.Fatalf("Expected len(data.Proof)=0; got %d", len(data.Proof))
	}
}

func TestParserMintmarks(t *testing.T) {
	data, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN 2020
		BEGIN CIRC
		      1.000 1001 1002 1003 1004 1005 1006 1007
		KNM*:  2000    ? 2002 2003 2004 2005 2006 2007
		MdP:   3000    ? 3002 3003 3004 3005 3006 3007
	`)), "-")

	if err != nil {
		t.Fatalf(`Expected err=nil; got "%s"`, err)
	}

	for i, row := range data.Circ {
		for j, col := range row.Cols {
			var n int
			if i > 0 && j == 1 {
				n = -1
			} else {
				n = 1000*i + j + 1000
			}
			if col != n {
				t.Fatalf("Expected data.Circ[i][j]=%d; got %d", n, col)
			}
		}
	}

	for i, y := range [...]int{2020, 2021, 2021} {
		if data.Circ[i].Year != y {
			t.Fatalf("Expected data.Circ[%d].Year=%d; got %d",
				i, y, data.Circ[i].Year)
		}
	}
	for i, s := range [...]string{"", "KNM", "MdP"} {
		if data.Circ[i].Mintmark != s {
			t.Fatalf(`Expected data.Circ[%d].Mintmark="%s"; got "%s"`,
				i, s, data.Circ[i].Mintmark)
		}
	}
}

func TestParserNoYear(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN CIRC
		1.000 1001 1002 1003 1004 1005 1006 1007
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106 2107
		BEGIN PROOF
		1.200 1201 1202 1203 1204 1205 1206 1207
		 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserNoType(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN 2020
		1.000 1001 1002 1003 1004 1005 1006 1007
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106 2107
		BEGIN PROOF
		1.200 1201 1202 1203 1204 1205 1206 1207
		 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserNoYearOrType(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		1.000 1001 1002 1003 1004 1005 1006 1007
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106 2107
		BEGIN PROOF
		1.200 1201 1202 1203 1204 1205 1206 1207
		 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserBadToken(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN 2020
		BEGIN CIRC
		1.000 1001 1002 1003 1004 1005 1006 1007
		I’m bad
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106 2107
		BEGIN PROOF
		1.200 1201 1202 1203 1204 1205 1206 1207
		 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserShortRow(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN 2020
		BEGIN CIRC
		1.000 1001 1002 1003 1004 1005 1006 1007
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106
		BEGIN PROOF
		1.200 1201 1202 1203 1204 1205 1206 1207
		 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserLongRow(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN 2020
		BEGIN CIRC
		1.000 1001 1002 1003 1004 1005 1006 1007
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106 2107 2108
		BEGIN PROOF
		1.200 1201 1202 1203 1204 1205 1206 1207
		 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserBadCoinType(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		BEGIN 2020
		BEGIN CIRCULATED
		1.000 1001 1002 1003 1004 1005 1006 1007
		 2000    ? 2002 2003 2004 2005 2006 2007
		BEGIN BU
		1.100 1101 1102 1103 1104 1105 1106 1107
		 2100    ? 2102 2103 2104 2105 2106 2107
		BEGIN PROOF
		1.200 1201 1202 1203 1204 1205 1206 1207
		 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}
