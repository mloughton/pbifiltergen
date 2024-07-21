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
