package mintage

import (
	"bytes"
	"errors"
	"testing"
)

func TestParserComplete(t *testing.T) {
	data, err := Parse(bytes.NewBuffer([]byte(`
		2020
			 1000 1001  1002 1003 1004 1005 1006 1007
			 1100 1101  1102 1103 1104 1105 1106 1107
			 1200 1201  1202 1203 1204 1205 1206 1207
		2021-KNM
			2.000    ?  2002 2003 2004 2005 2006 2007
			2.100    ?  2102 2103 2104 2105 2106 2107
			2.200    ?  2202 2203 2204 2205 2206 2207
		2021-MdP
			 3000 3001  3002 3003 3004 3005 3006 3007
			 3100 3101  3102 3103 3104 3105 3106 3107
			 3200 3201  3202 3203 3204 3205 3206 3207
		2022
			4000 4001 4.002 4003 4004 4005 4006 4007
			4100 4101 4.102 4103 4104 4105 4106 4107
			4200 4201 4.202 4203 4204 4205 4206 4207

		2009 "10th Anniversary of Economic and Monetary Union"
			1000 2000 3000
		2022-⋆ "35th Anniversary of the Erasmus Programme"
			1001    ? 3001
	`)), "-")

	if err != nil {
		t.Fatalf(`Expected err=nil; got "%s"`, err)
	}

	for i, row := range data.Standard {
		for k := TypeCirc; k <= TypeProof; k++ {
			for j, col := range row.Mintages[k] {
				n := 1000*(i+1) + 100*k + j
				if i == 1 && j == 1 {
					n = Unknown
				}
				if col != n {
					t.Fatalf("Expected data.Standard[%d].Mintages[%d][%d]=%d; got %d",
						i, k, j, col, n)
				}
			}
		}
	}

	for i, row := range data.Commemorative {
		for k := TypeCirc; k <= TypeProof; k++ {
			n := 1000*(k+1) + i
			if i == 1 && k == 1 {
				n = Unknown
			}
			if row.Mintage[k] != n {
				t.Fatalf("Expected row.Mintage[%d]=%d; got %d",
					k, n, row.Mintage[k])
			}
		}
	}

	if len(data.Standard) != 4 {
		t.Fatalf("Expected len(data.Standard)=2; got %d", len(data.Standard))
	}
	if len(data.Commemorative) != 2 {
		t.Fatalf("Expected len(data.Commemorative)=2; got %d", len(data.Commemorative))
	}

	for i, x := range [...]struct {
		year           int
		mintmark, name string
	}{
		{2009, "", "10th Anniversary of Economic and Monetary Union"},
		{2022, "⋆", "35th Anniversary of the Erasmus Programme"},
	} {
		if data.Commemorative[i].Year != x.year {
			t.Fatalf("Expected data.Commemorative[%d].Year=%d; got %d",
				i, x.year, data.Commemorative[i].Year)
		}
		if data.Commemorative[i].Mintmark != x.mintmark {
			t.Fatalf(`Expected data.Commemorative[%d].Mintmark="%s"; got "%s"`,
				i, x.mintmark, data.Commemorative[i].Mintmark)
		}
		if data.Commemorative[i].Name != x.name {
			t.Fatalf(`Expected data.Commemorative[%d].Name="%s"; got "%s"`,
				i, x.name, data.Commemorative[i].Name)
		}
	}
}

func TestParserMintmarks(t *testing.T) {
	data, err := Parse(bytes.NewBuffer([]byte(`
		2020
			 1000 1001  1002 1003 1004 1005 1006 1007
			 1100 1101  1102 1103 1104 1105 1106 1107
			 1200 1201  1202 1203 1204 1205 1206 1207
		2021-KNM
			2.000    ?  2002 2003 2004 2005 2006 2007
			2.100    ?  2102 2103 2104 2105 2106 2107
			2.200    ?  2202 2203 2204 2205 2206 2207
		2021-MdP
			 3000 3001  3002 3003 3004 3005 3006 3007
			 3100 3101  3102 3103 3104 3105 3106 3107
			 3200 3201  3202 3203 3204 3205 3206 3207
		2022
			 4000 4001 4.002 4003 4004 4005 4006 4007
			 4100 4101 4.102 4103 4104 4105 4106 4107
			 4200 4201 4.202 4203 4204 4205 4206 4207
	`)), "-")

	if err != nil {
		t.Fatalf(`Expected err=nil; got "%s"`, err)
	}

	for i, row := range data.Standard {
		for j, col := range row.Mintages[TypeCirc] {
			n := 1000*(i+1) + j
			if i == 1 && j == 1 {
				n = Unknown
			}
			if col != n {
				t.Fatalf("Expected data.Standard[%d].Mintages[TypeCirc][%d]=%d; got %d",
					i, j, col, n)
			}
		}
	}

	for i, y := range [...]int{2020, 2021, 2021, 2022} {
		if data.Standard[i].Year != y {
			t.Fatalf("Expected data.Standard[%d].Year=%d; got %d",
				i, y, data.Standard[i].Year)
		}
	}

	for i, m := range [...]string{"", "KNM", "MdP", ""} {
		if data.Standard[i].Mintmark != m {
			t.Fatalf(`Expected data.Standard[%d].Mintmark="%s"; got "%s"`,
				i, m, data.Standard[i].Mintmark)
		}
	}
}

func TestParserNoYear(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
			1.000 1001 1002 1003 1004 1005 1006 1007
			1.100 1101 1102 1103 1104 1105 1106 1107
			1.200 1201 1202 1203 1204 1205 1206 1207
		2021
			 2000    ? 2002 2003 2004 2005 2006 2007
			 2100    ? 2102 2103 2104 2105 2106 2107
			 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserBadToken(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		2020
			1.000 1001 1002 1003 1004 1005 1006 1007
			1.100 1101 1102 1103 1104 1105 1106 1107
			1.200 1201 1202 1203 1204 1205 1206 1207
		2021 Naughty!
			 2000    ? 2002 2003 2004 2005 2006 2007
			 2100    ? 2102 2103 2104 2105 2106 2107
			 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserShortRow(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		2020
			1.000 1001 1002 1003 1004 1005 1006 1007
			1.100 1101 1102 1103 1104 1105 1106 1107
			1.200 1201 1202 1203 1204 1205 1206 1207
		2021
			 2000    ? 2002 2003 2004 2005 2006 2007
			 2100    ? 2102 2103 2104 2105 2106
			 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserLongRow(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		2020
			1.000 1001 1002 1003 1004 1005 1006 1007
			1.100 1101 1102 1103 1104 1105 1106 1107
			1.200 1201 1202 1203 1204 1205 1206 1207
		2021
			 2000    ? 2002 2003 2004 2005 2006 2007
			 2100    ? 2102 2103 2104 2105 2106 2107 2108
			 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}

func TestParserMissingRow(t *testing.T) {
	_, err := Parse(bytes.NewBuffer([]byte(`
		2020
			1.000 1001 1002 1003 1004 1005 1006 1007
			1.100 1101 1102 1103 1104 1105 1106 1107
			1.200 1201 1202 1203 1204 1205 1206 1207
		2021
			 2000    ? 2002 2003 2004 2005 2006 2007
			 2200    ? 2202 2203 2204 2205 2206 2207
	`)), "-")

	var sErr SyntaxError
	if !errors.As(err, &sErr) {
		t.Fatalf("Expected err=SyntaxError; got %s", err)
	}
}
