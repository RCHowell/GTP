package gtp

import "testing"

func Test_columnToLetter(t *testing.T) {
	tests := []struct {
		name string
		c    int32
		want string
	}{
		{
			name: "low out-of-range",
			c:    -1,
			want: "",
		},
		{
			name: "high out-of-range",
			c:    20,
			want: "",
		},
		{
			name: "start",
			c:    0,
			want: "A",
		},
		{
			name: "H cutoff",
			c:    7,
			want: "H",
		},
		{
			name: "J cutoff",
			c:    8,
			want: "J",
		},
		{
			name: "end",
			c:    18,
			want: "T",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := columnToLetter(tt.c); got != tt.want {
				t.Errorf("columnToLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_letterToColumn(t *testing.T) {
	tests := []struct {
		name string
		l rune
		want int32
	}{
		{
			name: "low out-of-range",
			l:    ' ',
			want: -1,
		},
		{
			name: "high out-of-range",
			l:    'a',
			want: -1,
		},
		{
			name: "start",
			l:    'A',
			want: 0,
		},
		{
			name: "H cutoff",
			l:    'H',
			want: 7,
		},
		{
			name: "J cutoff",
			l:    'J',
			want: 8,
		},
		{
			name: "end",
			l:    'T',
			want: 18,
		},
		{
			name: "I",
			l:    'I',
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := letterToColumn(tt.l); got != tt.want {
				t.Errorf("letterToColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}