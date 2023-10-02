package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortLinesInString(t *testing.T) {
	tests := []struct {
		s       string
		opt     sortOptions
		name    string
		want    string
		wantErr error
	}{
		{
			name: "byColumn1",
			s:    "Илон Маск\nДжефф Безос\nБернар Арно\nБилл Гейтс\nЛарри Эллисон\nЛарри Пейдж",
			opt: sortOptions{
				sortType:     byStrings,
				reverseOrder: false,
				unique:       false,
				checkSorted:  true,
				byColumn:     true,
				columnNumber: 0,
				delim:        " ",
			},
			want: "",
			wantErr: ErrNotSorted{
				lineNotInOrder: 1,
			},
		},
		{
			name: "byColumn1",
			s:    "Илон Маск\nДжефф Безос\nБернар Арно\nБилл Гейтс\nЛарри Эллисон\nЛарри Пейдж",
			opt: sortOptions{
				sortType:     byStrings,
				reverseOrder: false,
				unique:       false,
				byColumn:     true,
				columnNumber: 1,
				delim:        " ",
			},
			want:    "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			wantErr: nil,
		},
		{
			name: "byNumbersWithSuffix",
			s:    "1k\n15G\n15u\n1M\n1m\n100",
			opt: sortOptions{
				sortType: byNumbersWithSuffix,
			},
			want:    "15u\n1m\n100\n1k\n1M\n15G",
			wantErr: nil,
		},
		{
			name: "byNumbers",
			s:    "115\n27\n282\n9\n1000\n5",
			opt: sortOptions{
				sortType: byNumbers,
			},
			want:    "5\n9\n27\n115\n282\n1000",
			wantErr: nil,
		},
		{
			name: "byMonth",
			s:    "JAN\nJUL\nAUG\nDEC\nAPR\nMAR",
			opt: sortOptions{
				sortType: byMonth,
			},
			want:    "JAN\nMAR\nAPR\nJUL\nAUG\nDEC",
			wantErr: nil,
		},
		{
			name: "byStrings",
			s:    "каждый\nохотник\nжелает\nзнать\nгде\nсидит\nфазан",
			opt: sortOptions{
				sortType: byStrings,
			},
			want:    "где\nжелает\nзнать\nкаждый\nохотник\nсидит\nфазан",
			wantErr: nil,
		},
		{
			name: "unique",
			s:    "желает\nзнать\nкаждый\nохотник\nжелает\nзнать\nгде\nсидит\nфазан",
			opt: sortOptions{
				sortType: byStrings,
				unique:   true,
			},
			want:    "где\nжелает\nзнать\nкаждый\nохотник\nсидит\nфазан",
			wantErr: nil,
		},
		{
			name: "byStringsreverseOrder",
			s:    "каждый\nохотник\nжелает\nзнать\nгде\nсидит\nфазан",
			opt: sortOptions{
				sortType:     byStrings,
				reverseOrder: true,
			},
			want:    "фазан\nсидит\nохотник\nкаждый\nзнать\nжелает\nгде",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortLinesInString(tt.s, tt.opt)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
