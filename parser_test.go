package gtp

import (
	"reflect"
	"testing"
)

func Test_ParseColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Color
	}{
		{
			"Basic White",
			"white",
			Color_WHITE,
		},
		{
			"Basic Black",
			"black",
			Color_BLACK,
		},
		{
			"Mix-case White",
			"wHiTe",
			Color_WHITE,
		},
		{
			"Mix-case Black",
			"BlAcK",
			Color_BLACK,
		},
		{
			"Single White",
			"w",
			Color_WHITE,
		},
		{
			"Single Black",
			"b",
			Color_BLACK,
		},
		{
			"Else Empty",
			"anything here",
			Color_EMPTY,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseColor(tt.input); got != tt.want {
				t.Errorf("ParseColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitOnDoubleLF(t *testing.T) {
	type args struct {
		data  []byte
		atEOF bool
	}
	tests := []struct {
		name        string
		args        args
		wantAdvance int
		wantToken   []byte
		wantErr     bool
	}{
		{
			name: "Basic ends in two new lines",
			args: args{
				data:  []byte{'A', 'B', 'C', '\n', '\n'},
				atEOF: false,
			},
			wantAdvance: 5,
			wantToken:   []byte{'A', 'B', 'C'},
			wantErr:     false,
		},
		{
			name: "Extra input",
			args: args{
				data:  []byte{'C', 'A', 'T', '\n', '\n', 'X', 'Y', 'Z'},
				atEOF: false,
			},
			wantAdvance: 5,
			wantToken:   []byte{'C', 'A', 'T'},
			wantErr:     false,
		},
		{
			name: "Input multiple standalone newlines",
			args: args{
				data:  []byte{'A', '\n', 'B', '\n', 'C', 'D', '\n', '\n'},
				atEOF: false,
			},
			wantAdvance: 8,
			wantToken:   []byte{'A', '\n', 'B', '\n', 'C', 'D'},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdvance, gotToken, err := SplitOnDoubleLF(tt.args.data, tt.args.atEOF)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitOnDoubleLF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAdvance != tt.wantAdvance {
				t.Errorf("SplitOnDoubleLF() gotAdvance = %v, want %v", gotAdvance, tt.wantAdvance)
			}
			if !reflect.DeepEqual(gotToken, tt.wantToken) {
				t.Errorf("SplitOnDoubleLF() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func TestParseGenMoveResponse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *GenMoveResponse
		wantErr bool
	}{
		{
			name:  "Standard move, single digit row",
			input: "=1 A7",
			want: &GenMoveResponse{
				Id: 1,
				Move: &Move{
					Type: Move_PLACE,
					Vertex: &Vertex{
						Row:    6,
						Column: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "Standard move, double digit row",
			input: "=99 J11",
			want: &GenMoveResponse{
				Id: 99,
				Move: &Move{
					Type: Move_PLACE,
					Vertex: &Vertex{
						Row:    10,
						Column: 9,
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "Resign",
			input: "=1 rEsIgN",
			want: &GenMoveResponse{
				Id:   1,
				Move: &Move{Type: Move_RESIGN},
			},
			wantErr: false,
		},
		{
			name:  "Pass",
			input: "=1 paSS",
			want: &GenMoveResponse{
				Id:   1,
				Move: &Move{Type: Move_PASS},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGenMoveResponse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGenMoveResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGenMoveResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO negative testing
func TestParseVertex(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    *Vertex
		wantErr bool
	}{
		{
			name: "ll",
			in: "A1",
			want: &Vertex{Row: 0, Column: 0},
			wantErr: false,
		},
		{
			name: "ul",
			in: "A19",
			want: &Vertex{Row: 18, Column: 0},
			wantErr: false,
		},
		{
			name: "lr",
			in: "S1",
			want: &Vertex{Row: 0, Column: 18},
			wantErr: false,
		},
		{
			name: "ur",
			in: "S19",
			want: &Vertex{Row: 18, Column: 18},
			wantErr: false,
		},
		{
			name: "lowercase-ll",
			in: "a1",
			want: &Vertex{Row: 0, Column: 0},
			wantErr: false,
		},
		{
			name: "lowercase-ul",
			in: "a19",
			want: &Vertex{Row: 18, Column: 0},
			wantErr: false,
		},
		{
			name: "lowercase-lr",
			in: "s1",
			want: &Vertex{Row: 0, Column: 18},
			wantErr: false,
		},
		{
			name: "lowercase-ur",
			in: "s19",
			want: &Vertex{Row: 18, Column: 18},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVertex(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVertex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseVertex() got = %v, want %v", got, tt.want)
			}
		})
	}
}
