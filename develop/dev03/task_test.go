package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData_Sort(t *testing.T) {

	col2 := int(2)
	col3 := int(3)

	tests := []struct {
		name string
		d    *Data
		cmp  func(a, b string) int
		want []string
		col  *int
	}{
		{
			name: "numbers",
			d:    &Data{text: []string{"1", "2", "3", "41", "5", "11"}},
			cmp:  cmpNumbersAndMonths(nil, false, " ", false),
			want: []string{"1", "2", "3", "5", "11", "41"},
		},
		{
			name: "numbers reverse",
			d:    &Data{text: []string{"1", "2", "3", "41", "5", "11"}},
			cmp:  cmpNumbersAndMonths(nil, false, " ", true),
			want: []string{"41", "11", "5", "3", "2", "1"},
		},
		{
			name: "months",
			d: &Data{text: []string{
				"March", "April", "May", "June", "September", "October",
				"July", "August", "January", "February", "November", "December",
			}},
			cmp: cmpNumbersAndMonths(nil, true, " ", false),
			want: []string{
				"January", "February", "March", "April", "May", "June",
				"July", "August", "September", "October", "November", "December",
			},
		},
		{
			name: "months reverse",
			d: &Data{text: []string{
				"March", "April", "May", "June", "September", "October",
				"July", "August", "January", "February", "November", "December",
			}},
			cmp: cmpNumbersAndMonths(nil, true, " ", true),
			want: []string{
				"December", "November", "October", "September", "August", "July",
				"June", "May", "April", "March", "February", "January",
			},
		},
		{
			name: "months reverse col2",
			d: &Data{text: []string{
				"2023-March-18", "2023-February-01", "2023-January-14",
				"2023-September-26", "2023-August-09", "2023-July-02",
				"2023-June-07", "2023-May-20", "2023-April-05",
				"2023-December-13", "2023-November-18", "2023-October-11",
			}},
			cmp: cmpNumbersAndMonths(&col2, true, "-", true),
			want: []string{
				"2023-December-13", "2023-November-18", "2023-October-11",
				"2023-September-26", "2023-August-09", "2023-July-02",
				"2023-June-07", "2023-May-20", "2023-April-05",
				"2023-March-18", "2023-February-01", "2023-January-14",
			},
		},
		{
			name: "string",
			d: &Data{text: []string{
				"2023-March-18", "2023-February-01", "2023-January-14",
				"2023-September-26", "2023-August-09", "2023-July-02",
				"2023-June-07", "2023-May-20", "2023-April-05",
				"2023-December-13", "2023-November-18", "2023-October-11",
			}},
			cmp: cmpDefault(nil, " ", false),
			want: []string{
				"2023-April-05", "2023-August-09", "2023-December-13",
				"2023-February-01", "2023-January-14", "2023-July-02",
				"2023-June-07", "2023-March-18", "2023-May-20",
				"2023-November-18", "2023-October-11", "2023-September-26",
			},
		},
		{
			name: "string col3",
			d: &Data{text: []string{
				"2023-March-18", "2023-February-01", "2023-January-14",
				"2023-September-26", "2023-August-09", "2023-July-02",
				"2023-June-07", "2023-May-20", "2023-April-05",
				"2023-December-13", "2023-November-19", "2023-October-11",
			}},
			cmp: cmpDefault(&col3, "-", false),
			want: []string{
				"2023-February-01", "2023-July-02", "2023-April-05",
				"2023-June-07", "2023-August-09", "2023-October-11",
				"2023-December-13", "2023-January-14", "2023-March-18",
				"2023-November-19", "2023-May-20", "2023-September-26",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Sort(tt.cmp)
			assert.Equal(t, tt.want, tt.d.text)
		})
	}
}

func TestData_IsSort(t *testing.T) {

	col3 := int(3)

	tests := []struct {
		name string
		d    *Data
		cmp  func(a, b string) int
		exp  bool
	}{
		{
			name: "string col3 false",
			d: &Data{text: []string{
				"2023-March-18", "2023-February-01", "2023-January-14",
				"2023-September-26", "2023-August-09", "2023-July-02",
				"2023-June-07", "2023-May-20", "2023-April-05",
				"2023-December-13", "2023-November-19", "2023-October-11",
			}},
			cmp: cmpDefault(&col3, "-", false),
			exp: false,
		},

		{
			name: "string col3 true",
			d: &Data{text: []string{
				"2023-February-01", "2023-July-02", "2023-April-05",
				"2023-June-07", "2023-August-09", "2023-October-11",
				"2023-December-13", "2023-January-14", "2023-March-18",
				"2023-November-19", "2023-May-20", "2023-September-26",
			}},
			cmp: cmpDefault(&col3, "-", false),
			exp: true,
		},

		{
			name: "string ",
			d: &Data{text: []string{
				"2023-April-05", "2023-August-09", "2023-December-13",
				"2023-February-01", "2023-January-14", "2023-July-02",
				"2023-June-07", "2023-March-18", "2023-May-20",
				"2023-November-18", "2023-October-11", "2023-September-26",
			}},
			cmp: cmpDefault(nil, " ", false),
			exp: true,
		},
		{
			name: "string ",
			d: &Data{text: []string{
				"2023-September-26", "2023-October-11", "2023-November-18",
				"2023-May-20", "2023-March-18", "2023-June-07",
				"2023-July-02", "2023-January-14", "2023-February-01",
				"2023-December-13", "2023-August-09", "2023-April-05",
			}},
			cmp: cmpDefault(nil, " ", true),
			exp: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.exp, tt.d.IsSorted(tt.cmp))

		})
	}
}

func TestData_DeleteRepeatedLines(t *testing.T) {

	col2 := int(2)

	type args struct {
		col *int
		sep string
	}
	tests := []struct {
		name string
		d    *Data
		args args
		want []string
	}{
		{
			name: "test1",
			d: &Data{text: []string{
				"1", "1", "1", "1", "1", "1",
			}},
			want: []string{
				"1",
			},
		},
		{
			name: "test2",
			d: &Data{text: []string{
				"1 asd", "1 asd", "1 bcf", "1 bcf", "1 cc", "1 cda",
			}},
			args: args{
				col: &col2,
				sep: " ",
			},
			want: []string{
				"1 asd", "1 bcf", "1 cc", "1 cda",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.DeleteRepeatedLines(tt.args.col, tt.args.sep)
			assert.Equal(t, tt.want, tt.d.text)
		})
	}
}
