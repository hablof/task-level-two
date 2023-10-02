package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findAnagrams(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		s       []string
		want    map[string][]string
		wantErr error
	}{
		{
			name:    "eng-word",
			s:       []string{"столик", "Jeffrey", "листок", "слИток", "тяпка", "пятка"},
			want:    nil,
			wantErr: fmt.Errorf("error on word #%d-\"%s\": %w", 1, "jeffrey", ErrNonRussianLetter),
		},
		{
			name: "valid",
			s:    []string{"столик", "пятак", "листок", "слИток", "тяпка", "пятка"},
			want: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"столик": {"листок", "слиток", "столик"},
			},
		},
		{
			name: "valid2",
			s:    []string{"ракета", "карета", "ракета", "ракета", "ошибка", "ашибко", "пень"},
			want: map[string][]string{
				"ракета": {"карета", "ракета"},
				"ошибка": {"ашибко", "ошибка"},
			},
		},
		{
			name: "без анаграм",
			s:    []string{"Практический", "опыт", "показывает", "что", "постоянное", "информационно", "техническое", "обеспечение", "нашей", "деятельности", "обеспечивает", "широкому", "кругу", "специалистов", "участие", "в", "формировании", "позиций", "занимаемых", "участниками", "в", "отношении", "поставленных", "задач", "Не", "следует", "однако", "забывать", "о", "том", "что", "курс", "на", "социально", "ориентированный", "национальный", "проект", "требует", "от", "нас", "анализа", "позиций", "занимаемых", "участниками", "в", "отношении", "поставленных", "задач", "Разнообразный", "и", "богатый", "опыт", "сложившаяся", "структура", "организации", "позволяет", "выполнить", "важнейшие", "задания", "по", "разработке", "позиций", "занимаемых", "участниками", "в", "отношении", "поставленных", "задач"},
			want: map[string][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findAnagrams(tt.s)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, len(tt.want), len(got))

			for key, wantSlice := range tt.want {

				gotSlice := got[key]
				assert.Equal(t, len(wantSlice), len(gotSlice))

				for i, elem := range wantSlice {
					assert.Equal(t, elem, gotSlice[i])
				}
			}
		})
	}
}
