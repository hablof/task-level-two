package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findUnpackSequenses(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{
			s:    `🤔5`,
			want: `🤔🤔🤔🤔🤔`,
		},
		{
			s:    `a4bc2d5e`,
			want: `aaaabccddddde`,
		},
		{
			s:    `abcd`,
			want: `abcd`,
		},
		{
			s:    `45`,
			want: ``,
		},
		{
			s:    `qwe\4\5`,
			want: `qwe45`,
		},
		{
			s:    `qwe\45`,
			want: `qwe44444`,
		},
		{
			s:    `qwe\\5`,
			want: `qwe\\\\\`,
		},
		{
			s:    `абырвалг\`,
			want: `абырвалг\`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := findUnpackSequenses(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
