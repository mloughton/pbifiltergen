package dax

import "testing"

func TestCreateColumn(t *testing.T) {
	type want struct {
		col Column
		err error
	}
	type test struct {
		input string
		want  want
	}

	tests := []test{
		{
			input: "TEST[TEST],STRING",
			want: want{
				col: Column{
					fullName:   "TEST[TEST]",
					tableName:  "TEST",
					columnName: "TEST",
					columnType: "STRING",
				},
				err: nil,
			},
		},
		{
			input: "'TEST SPACE'[TEST SPACE],STRING",
			want: want{
				col: Column{
					fullName:   "'TEST SPACE'[TEST SPACE]",
					tableName:  "TEST SPACE",
					columnName: "TEST SPACE",
					columnType: "STRING",
				},
				err: nil,
			},
		},
		{
			input: "TEST[TEST],STRIN",
			want: want{
				col: Column{},
				err: errInvalidType,
			},
		},
	}
	for i, test := range tests {
		col, err := createColumn(test.input)
		if col != test.want.col || err != test.want.err {
			t.Fatalf("test %v | expected: col %v, err %v; got: col %v, err: %v", i+1, test.want.col, test.want.err, col, err)
		}
	}
}

func TestGenerateDax(t *testing.T) {
	type test struct {
		input []Column
		want  string
	}

	tests := []test{
		{
			input: []Column{
				{
					fullName:   "TEST[TEST]",
					tableName:  "TEST",
					columnName: "TEST",
					columnType: "STRING",
				},
			},
			want: `VAR V0 =
    IF(
        ISFILTERED( TEST[TEST] ) = TRUE(),
        VAR vals = SELECTCOLUMNS(ADDCOLUMNS(VALUES( TEST[TEST] ), "clean",
            SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(SUBSTITUTE(TEST[TEST],
                "%", "%25"),
                " ", "%20"),
                "'", "%27%27"),
                "+", "%2B"),
                "/", "%2F"),
                "?", "%3F"),
                "#", "%23"),
                "&", "%26")
            ), "clean", [clean])
        VAR cnt = COUNTROWS(vals)
        VAR valsStr = IF(cnt = 1, "eq%20%27" & vals & "%27", "in%20(%27" & CONCATENATEX(vals, [clean], "%27,%20%27") & "%27)")
        VAR fullStr = "TEST/TEST%20" & valsStr
        VAR fullStrBlank = IF(valsStr = BLANK(), BLANK(), fullStr)
    	RETURN
        	fullStrBlank,   
        	BLANK()
    )

VAR link = CONCATENATEX(FILTER({V0}, [Value] <> BLANK()), [Value], "%20and%20")
VAR fullLink = IF(link = BLANK(), BLANK(), [Dashboard Link] & "?rs:embed=true&filter=" & link)
RETURN
	fullLink`,
		},
	}

	for i, test := range tests {
		result, err := GenerateDax(&test.input)
		if err != nil || result != test.want {
			t.Fatalf("test %v | expected:\n%v\n got\n%v\nerr %v", i+1, test.want, result, err)
		}
	}
}
