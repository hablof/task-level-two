package mysort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bigText = `Я не могу дышать, мне не видно неба
Я не могу понять был ты или не был
Ветром по волосам, солнце в ладони - твоя…
Красные облака вечером дарят в спину
Я с тобой так легка, я с тобою красива
Бешено так в груди бьётся сердце твоё

Отпускаю и в небо
Улетает жёлтыми листьями
Наше прошлое лето
С телефонными, глупыми письмами
Отпускаю и слёзы
Высыхают на ресницах
Но как же синие звёзды
Нам с тобой могли присниться

Рано ещё ни быть, поздно уже поверить
Я не могла любить, я не могла измерить
Месяцы за окном, солнце в закатах с тобой
И опускаюсь вниз, и поднимаюсь в небо
Я не могу понять был ты или не был
В сотнях ночных дорог ты остаёшься со мной`
	bigTextSorted = `

Бешено так в груди бьётся сердце твоё
В сотнях ночных дорог ты остаёшься со мной
Ветром по волосам, солнце в ладони - твоя…
Высыхают на ресницах
И опускаюсь вниз, и поднимаюсь в небо
Красные облака вечером дарят в спину
Месяцы за окном, солнце в закатах с тобой
Нам с тобой могли присниться
Наше прошлое лето
Но как же синие звёзды
Отпускаю и в небо
Отпускаю и слёзы
Рано ещё ни быть, поздно уже поверить
С телефонными, глупыми письмами
Улетает жёлтыми листьями
Я не могла любить, я не могла измерить
Я не могу дышать, мне не видно неба
Я не могу понять был ты или не был
Я не могу понять был ты или не был
Я с тобой так легка, я с тобою красива`
)

func Test_sortLinesInString(t *testing.T) {
	tests := []struct {
		s       string
		opt     SortOptions
		name    string
		want    string
		wantErr error
	}{
		{
			name: "maksim",
			s:    bigText,
			opt: SortOptions{
				SortType: Alphabetical,
			},
			want:    bigTextSorted,
			wantErr: nil,
		},
		{
			name: "dontIgnoreRearSpaces",
			s:    "Вова\nВова \nПетя \nПетя\nАлиса\nАлиса ",
			opt: SortOptions{
				SortType: Alphabetical,
			},
			want:    "Алиса\nАлиса \nВова\nВова \nПетя\nПетя ",
			wantErr: nil,
		},
		{
			name: "ignoreRearSpaces",
			s:    "Алиса\nАлиса \nАлиса   \nАлиса \nВова\nВова \nВова \nВова  \nПетя \nПетя",
			opt: SortOptions{
				SortType:             Alphabetical,
				IgnoreTrailingSpaces: true,
				CheckSorted:          true,
			},
			want:    isSortedMsg,
			wantErr: nil,
		},
		{
			name: "byColumn1WithDelim\"|\"",
			s:    "Илон|Маск\nДжефф|Безос\nБернар|Арно\nБилл|Гейтс\nЛарри|Эллисон\nЛарри|Пейдж",
			opt: SortOptions{
				SortType:     Alphabetical,
				ByColumn:     true,
				ColumnNumber: 1,
				Delim:        "|",
			},
			want:    "Бернар|Арно\nДжефф|Безос\nБилл|Гейтс\nИлон|Маск\nЛарри|Пейдж\nЛарри|Эллисон",
			wantErr: nil,
		},
		{
			name: "checkSortedbyColumn0",
			s:    "Илон Маск\nДжефф Безос\nБернар Арно\nБилл Гейтс\nЛарри Эллисон\nЛарри Пейдж",
			opt: SortOptions{
				SortType:     Alphabetical,
				ReverseOrder: false,
				Unique:       false,
				CheckSorted:  true,
				ByColumn:     true,
				ColumnNumber: 0,
				Delim:        " ",
			},
			want: "",
			wantErr: ErrNotSorted{
				lineNotInOrder: 1,
			},
		},
		{
			name: "byColumn1",
			s:    "Илон Маск\nДжефф Безос\nБернар Арно\nБилл Гейтс\nЛарри Эллисон\nЛарри Пейдж",
			opt: SortOptions{
				SortType:     Alphabetical,
				ReverseOrder: false,
				Unique:       false,
				ByColumn:     true,
				ColumnNumber: 1,
				Delim:        " ",
			},
			want:    "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			wantErr: nil,
		},
		{
			name: "byNumbersWithSuffix",
			s:    "1k\n15G\n15u\n1M\n1m\n100",
			opt: SortOptions{
				SortType: HumanNumberic,
			},
			want:    "15u\n1m\n100\n1k\n1M\n15G",
			wantErr: nil,
		},
		{
			name: "byNumbers",
			s:    "115\n27\n282\n9\n1000\n5",
			opt: SortOptions{
				SortType: Numeric,
			},
			want:    "5\n9\n27\n115\n282\n1000",
			wantErr: nil,
		},
		{
			name: "byMonth",
			s:    "JAN\nJUL\nAUG\nDEC\nAPR\nMAR",
			opt: SortOptions{
				SortType: ByMonth,
			},
			want:    "JAN\nMAR\nAPR\nJUL\nAUG\nDEC",
			wantErr: nil,
		},
		{
			name: "byStrings",
			s:    "каждый\nохотник\nжелает\nзнать\nгде\nсидит\nфазан",
			opt: SortOptions{
				SortType: Alphabetical,
			},
			want:    "где\nжелает\nзнать\nкаждый\nохотник\nсидит\nфазан",
			wantErr: nil,
		},
		{
			name: "unique",
			s:    "желает\nзнать\nкаждый\nохотник\nжелает\nзнать\nгде\nсидит\nфазан",
			opt: SortOptions{
				SortType: Alphabetical,
				Unique:   true,
			},
			want:    "где\nжелает\nзнать\nкаждый\nохотник\nсидит\nфазан",
			wantErr: nil,
		},
		{
			name: "byStringsreverseOrder",
			s:    "каждый\nохотник\nжелает\nзнать\nгде\nсидит\nфазан",
			opt: SortOptions{
				SortType:     Alphabetical,
				ReverseOrder: true,
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
