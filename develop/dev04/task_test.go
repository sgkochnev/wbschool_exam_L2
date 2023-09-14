package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findAnagrams(t *testing.T) {

	tests := []struct {
		name  string
		words []string
		want  map[string][]string
	}{
		{
			name:  "test1",
			words: []string{"тяпКа", "пятак", "пятка", "Пятка", "слиток", "листок", "Столик", "свисток", "рулетка"},
			want: map[string][]string{
				"тяпка":  {"пятак", "пятка", "тяпка"},
				"слиток": {"листок", "слиток", "столик"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findAnagrams(tt.words)
			assert.Equal(t, tt.want, got)
		})
	}
}
