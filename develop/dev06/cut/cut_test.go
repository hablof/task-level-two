package cut

import "testing"

func TestCut(t *testing.T) {
	tests := []struct {
		name string
		s    string
		opts CutOptions
		want string
	}{
		{
			name: "",
			s:    "first words of file\nnextwords|kek\nand another in .txt",
			opts: CutOptions{
				OnlyDelimited: false,
				Delimiter:     " ",
				Fields:        []int{2},
			},
			want: "of\nnextwords|kek\nin",
		},
		{
			name: "",
			s:    "first words of file\nnextwords|kek\nand another in .txt",
			opts: CutOptions{
				OnlyDelimited: true,
				Delimiter:     " ",
				Fields:        []int{2},
			},
			want: "of\nin",
		},
		{
			name: "",
			s:    "first words of file\nnextwords kek\nand another in .txt",
			opts: CutOptions{
				OnlyDelimited: false,
				Delimiter:     " ",
				Fields:        []int{2},
			},
			want: "of\n\nin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cut(tt.s, tt.opts); got != tt.want {
				t.Errorf("Cut() = %v, want %v", got, tt.want)
			}
		})
	}
}
