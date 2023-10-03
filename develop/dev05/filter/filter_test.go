package filter

import "testing"

func TestFilter(t *testing.T) {
	tests := []struct {
		name        string
		inputString string
		opts        FilterOptions
		want        string
	}{
		{
			name:        "invert all",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern: "noMatch",
				Invert:  true,
			},
			want: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
		},
		{
			name:        "Билл Гейтс + next",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern:    "Билл Гейтс",
				LinesAfter: 1,
			},
			want: "\033[32mБилл Гейтс\033[0m\nИлон Маск",
		},
		{
			name:        "Билл Гейтс line",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern: "Билл Гейтс",
				Fixed:   true,
			},
			want: "\033[32mБилл Гейтс\033[0m",
		},
		{
			name:        "Билл Гейтс + prev",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern:     "Билл Гейтс",
				LinesBefore: 1,
			},
			want: "Джефф Безос\n\033[32mБилл Гейтс\033[0m",
		},
		{
			name:        "count р",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern: "р",
				Count:   true,
			},
			want: "matched 3 lines",
		},
		{
			name:        "count lines == р",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern: "р",
				Fixed:   true,
				Count:   true,
			},
			want: "matched 0 lines",
		},
		{
			name:        "find р",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern:  "р",
				WithNums: true,
			},
			want: "1:Бе\033[32mр\033[0mнар Арно\n5:Ла\033[32mр\033[0mри Пейдж\n6:Ла\033[32mр\033[0mри Эллисон",
		},
		{
			name:        "find space",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern: " ",
			},
			want: "Бернар\033[32m \033[0mАрно\nДжефф\033[32m \033[0mБезос\nБилл\033[32m \033[0mГейтс\nИлон\033[32m \033[0mМаск\nЛарри\033[32m \033[0mПейдж\nЛарри\033[32m \033[0mЭллисон",
		},
		{
			name:        "find д",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern: "д",
			},
			want: "Ларри Пей\033[32mд\033[0mж",
		},
		{
			name:        "find Д",
			inputString: "Бернар Арно\nДжефф Безос\nБилл Гейтс\nИлон Маск\nЛарри Пейдж\nЛарри Эллисон",
			opts: FilterOptions{
				Pattern:    "д",
				IgnoreCase: true,
			},
			want: "\033[32mД\033[0mжефф Безос\nЛарри Пей\033[32mд\033[0mж",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.inputString, tt.opts)
			if got != tt.want {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
