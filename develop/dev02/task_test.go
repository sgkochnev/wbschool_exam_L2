package main

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    string
		wantErr error
	}{
		{
			name: "correct string 1",
			str:  `a4bc2d5e`,
			want: `aaaabccddddde`,
		},
		{
			name: "correct string 2",
			str:  `abcd`,
			want: `abcd`,
		},
		{
			name:    "invalid string 1",
			str:     `45`,
			wantErr: ErrInvalidtString,
		},
		{
			name:    "invalid string 2",
			str:     `4`,
			wantErr: ErrInvalidtString,
		},
		{
			name:    "invalid string 3",
			str:     `\`,
			wantErr: ErrInvalidtString,
		},
		{
			name:    "invalid string 4",
			str:     `as\fasd`,
			wantErr: ErrInvalidtString,
		},
		{
			name:    "invalid string 5",
			str:     `qwe45`,
			wantErr: ErrInvalidtString,
		},
		{
			name: "empty string",
			str:  ``,
		},
		{
			name: "string with escape characters 1",
			str:  `qwe\4\5`,
			want: `qwe45`,
		},
		{
			name: "string with escape characters 2",
			str:  `qwe\45`,
			want: `qwe44444`,
		},
		{
			name: "string with escape characters 3",
			str:  `qwe\\5`,
			want: `qwe\\\\\`,
		},
		{
			name: "string with escape characters 4",
			str:  `\5`,
			want: `5`,
		},
		{
			name: "string with escape characters 5",
			str:  `\45`,
			want: `44444`,
		},
		{
			name: "space string",
			str:  ` `,
			want: ` `,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Unpack(tt.str); got != tt.want || err != tt.wantErr {
				t.Errorf("str = %v, want %v", got, tt.want)
				t.Errorf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
