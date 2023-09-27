package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cut(t *testing.T) {

	data := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		"1\t2\t3\t4\t5\t6\t7",
		"Thisisonecolumn",
		"one\ttwo\tthree\tfour\tfive\tsix\tseven\teight\tnine",
		"a\tb\tc\td\te",
	)

	dataWithDelimiter := fmt.Sprintf("%s\n%s\n%s\n%s",
		"1 2 3 4 5 6 7",
		"Thisisonecolumn",
		"one two three four five six seven eight nine",
		"a b c d e",
	)

	tests := []struct {
		name string
		in   string
		f    *flags
		want string
	}{
		{
			name: "fields",
			in:   data,
			f: &flags{
				fields:    [][2]int{{1, 3}, {6, 6}, {8, 8}},
				delimiter: "\t",
			},
			want: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"1\t2\t3\t6",
				"Thisisonecolumn",
				"one\ttwo\tthree\tsix\teight",
				"a\tb\tc",
			),
		},

		{
			name: "fields with delimeter",
			in:   dataWithDelimiter,
			f: &flags{
				fields:    [][2]int{{1, 3}, {6, 6}, {8, 8}},
				delimiter: " ",
			},
			want: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				"1 2 3 6",
				"Thisisonecolumn",
				"one two three six eight",
				"a b c",
			),
		},

		{
			name: "fields with separated",
			in:   data,
			f: &flags{
				fields:    [][2]int{{1, 3}, {6, 6}, {8, 8}},
				delimiter: "\t",
				separated: true,
			},
			want: fmt.Sprintf("%s\n%s\n%s\n",
				"1\t2\t3\t6",
				"one\ttwo\tthree\tsix\teight",
				"a\tb\tc",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			r := strings.NewReader(tt.in)
			cut(r, w, tt.f)
			assert.Equal(t, tt.want, w.String())
		})
	}
}

func Test_parseFields(t *testing.T) {
	tests := []struct {
		name string
		in   string
		f    *flags
		want [][2]int
		err  error
	}{
		{
			name: "fields with comma",
			in:   "1,3,6",
			f:    &flags{},
			want: [][2]int{{1, 1}, {3, 3}, {6, 6}},
		},
		{
			name: "fields with comma 2",
			in:   "1,6,3",
			f:    &flags{},
			want: [][2]int{{1, 1}, {3, 3}, {6, 6}},
		},
		{
			name: "fields with intterval",
			in:   "1-3,6",
			f:    &flags{},
			want: [][2]int{{1, 3}, {6, 6}},
		},
		{
			name: "fields with intterval 2",
			in:   "-3,6",
			f:    &flags{},
			want: [][2]int{{1, 3}, {6, 6}},
		},
		{
			name: "fields with intterval 3",
			in:   "-",
			f:    &flags{},
			want: [][2]int{{1, math.MaxInt}},
		},
		{
			name: "fields with intterval 4",
			in:   "3-",
			f:    &flags{},
			want: [][2]int{{3, math.MaxInt}},
		},
		{
			name: "fields with intterval 5",
			in:   "8-,3",
			f:    &flags{},
			want: [][2]int{{3, 3}, {8, math.MaxInt}},
		},
		{
			name: "fields with intterval 6",
			in:   "8-,3-5,1",
			f:    &flags{},
			want: [][2]int{{1, 1}, {3, 5}, {8, math.MaxInt}},
		},
		{
			name: "fields with intterval 7",
			in:   "8-6,3-0,1",
			f:    &flags{},
			want: [][2]int{{1, 1}, {3, 3}, {8, 8}},
		},
		{
			name: "fields with intterval 8",
			in:   "0-6",
			f:    &flags{},
			want: [][2]int{{1, 6}},
		},
		{
			name: "fields with intterval 9",
			in:   "0-6,7,5-7",
			f:    &flags{},
			want: [][2]int{{1, 7}},
		},
		{
			name: "error 1",
			in:   "8-,3-5,a,1",
			f:    &flags{},
			err:  ErrNotNumber,
		},
		{
			name: "error 2",
			in:   "5-a,",
			f:    &flags{},
			err:  ErrNotNumber,
		},
		{
			name: "error 3",
			in:   "b-1",
			f:    &flags{},
			err:  ErrNotNumber,
		},
		{
			name: "error 4",
			in:   "-1-3,6",
			f:    &flags{},
			err:  ErrInvalidInterval,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseFields(tt.f)(tt.in)
			assert.ErrorIs(t, err, tt.err)
			if err == nil {
				assert.Equal(t, tt.want, tt.f.fields)
			}
		})
	}
}
