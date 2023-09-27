package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pattern(t *testing.T) {

	text := `1 edc
2 asd
3 dsa
4 Qq
5 zxcv
5 cxz
4 asd
3 sss
2 dsa
1 qq`

	tests := []struct {
		name string
		f    *flags
		want string
	}{
		{
			name: "A",
			f:    &flags{after: 1, pattern: "asd"},
			want: "2 asd\n3 dsa\n4 asd\n3 sss\n",
		},
		{
			name: "B",
			f:    &flags{before: 1, pattern: "asd"},
			want: "1 edc\n2 asd\n5 cxz\n4 asd\n",
		},
		{
			name: "C",
			f:    &flags{context: 1, pattern: "asd"},
			want: "1 edc\n2 asd\n3 dsa\n5 cxz\n4 asd\n3 sss\n",
		},
		{
			name: "count",
			f:    &flags{count: true, pattern: "asd"},
			want: "2",
		},
		{
			name: "ignoreCase",
			f:    &flags{ignoreCase: true, context: 1, pattern: "qQ"},
			want: "3 dsa\n4 Qq\n5 zxcv\n2 dsa\n1 qq\n",
		},
		{
			name: "invert",
			f:    &flags{invert: true, count: true, pattern: "qq"},
			want: "9",
		},
		{
			name: "fixed",
			f:    &flags{fixed: true, pattern: "asd"},
			want: "",
		},
		{
			name: "lineNum",
			f:    &flags{lineNum: true, pattern: "asd"},
			want: "2:2 asd\n7:4 asd\n",
		},
		{
			name: "lineNum",
			f:    &flags{lineNum: true, context: 1, pattern: "asd"},
			want: "1:1 edc\n2:2 asd\n3:3 dsa\n6:5 cxz\n7:4 asd\n8:3 sss\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			r := bytes.NewBufferString(text)
			grep(r, w, tt.f)
			assert.Equal(t, tt.want, w.String())
		})
	}
}
