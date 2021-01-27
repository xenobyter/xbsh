package line

import (
	"errors"
	"testing"

	"github.com/eiannone/keyboard"
)

func Test_mockInput_next(t *testing.T) {
	tests := []struct {
		name      string
		wantRunes []rune
		wantKeys  []keyboard.Key
		wantErrs  []error
	}{
		{"Simple Runes", []rune{'a', 'b'}, []keyboard.Key{0, 0}, []error{nil, nil}},
		{"Runes and key", []rune{'a', 'b', 0}, []keyboard.Key{0, 0, 1}, []error{nil, nil, nil}},
		{"Runes, key, error", []rune{'a', 'b', 0, 0}, []keyboard.Key{0, 0, 1, 0}, []error{nil, nil, nil, errors.New("keyboard")}},
	}

	for _, tt := range tests {
		mockIn := newMockInput()
		mockIn.set(tt.wantRunes, tt.wantKeys, tt.wantErrs)
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.wantErrs {
				gotRune, gotKey, gotErr := mockIn.next()
				if gotRune != tt.wantRunes[i] {
					t.Errorf("mockKeys.next() rune = %c, wantRune %c", gotRune, tt.wantRunes[i])
				}
				if gotKey != tt.wantKeys[i] {
					t.Errorf("mockKeys.next() key = %c, wantKey %c", gotKey, tt.wantKeys[i])
				}
				if gotErr != tt.wantErrs[i] {
					t.Errorf("mockKeys.next() err = %c, wantErr %c", gotErr, tt.wantErrs[i])
				}

			}
		})
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		mock     mockInput
		wantLine string
	}{
		{"Simple chars", "", mockInput{[]rune{'a', 'b', 0}, []keyboard.Key{0, 0, keyboard.KeyEnter}, []error{nil, nil, nil}, 0}, "ab"},
		{"Simple chars with blank", "", mockInput{[]rune{'a', 0, 'b', 0}, []keyboard.Key{0, keyboard.KeySpace, 0, keyboard.KeyEnter}, []error{nil, nil, nil, nil}, 0}, "a b"},
		{"ArrowLeft", "", mockInput{
			[]rune{'a', 'b', 0, 'c', 0},
			[]keyboard.Key{0, 0, keyboard.KeyArrowLeft, 0, keyboard.KeyEnter},
			[]error{nil, nil, nil, nil, nil}, 0},
			"acb"},
		{"Boundary ArrowLeft", "", mockInput{
			[]rune{'a', 0, 0, 'c', 0},
			[]keyboard.Key{0, keyboard.KeyArrowLeft, keyboard.KeyArrowLeft, 0, keyboard.KeyEnter},
			[]error{nil, nil, nil, nil, nil}, 0},
			"ca"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLine := Read(tt.prompt, tt.mock); gotLine != tt.wantLine {
				t.Errorf("Read() = %v, want %v", gotLine, tt.wantLine)
			}
		})
	}
}

func Test_insert(t *testing.T) {
	type args struct {
		line string
		pos  int
		c    rune
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{"Append", args{"lin", 3, 'e'}, "line"},
		{"Insert", args{"lie", 2, 'n'}, "line"},
		{"Append umlaut", args{"lin", 3, 'ä'}, "linä"},
		{"Insert umlaut", args{"lie", 2, 'ä'}, "liäe"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := insert(tt.args.line, tt.args.pos, tt.args.c); gotRes != tt.wantRes {
				t.Errorf("insert() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_output(t *testing.T) {
	type args struct {
		prompt, line string
		lx, ox, ld   int
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantCx  int
	}{
		{"Simple line", args{"1234", "5678", 4, 0, 100}, "5678", 8},
		{"Cursor in line", args{"1234", "5678", 3, 0, 100}, "5678", 7},
		{"Cursor at 0", args{"1234", "5678", 0, 0, 100}, "5678", 4},
		{"Cursor at 0, short line", args{"1234", "5678", 0, 0, 3}, "567", 4},
		{"Cursor at 1, short line", args{"1234", "5678", 1, 0, 3}, "567", 5},
		{"Cursor at 1, ox 1", args{"1234", "234567890", 1, 1, 3}, "345", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotCx := output(tt.args.prompt, tt.args.line, tt.args.lx, tt.args.ox, tt.args.ld)
			if gotRes != tt.wantRes {
				t.Errorf("output() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
			if gotCx != tt.wantCx {
				t.Errorf("output() gotCx = %v, want %v", gotCx, tt.wantCx)
			}
		})
	}
}

func Test_move(t *testing.T) {
	type args struct {
		dx, lx, ox, ll, ld int
	}
	tests := []struct {
		name   string
		args   args
		wantnx int
		wantno int
	}{
		{"Simple right", args{1, 0, 0, 10, 10}, 1, 0},
		{"Respect line length right", args{11, 0, 0, 10, 10}, 10, 0},
		{"Respect line length left", args{-11, 2, 0, 10, 10}, 0, 0},
		{"move ox right", args{2, 9, 0, 100, 10}, 11, 1},
		{"move ox left", args{-2, 9, 9, 100, 10}, 7, 7},
		{"Don't move ox", args{-2, 13, 10, 100, 10}, 11, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotnx, gotno := move(tt.args.dx, tt.args.lx, tt.args.ox, tt.args.ll, tt.args.ld)
			if gotnx != tt.wantnx {
				t.Errorf("move() gotnx = %v, want %v", gotnx, tt.wantnx)
			}
			if gotno != tt.wantno {
				t.Errorf("move() gotno = %v, want %v", gotno, tt.wantno)
			}
		})
	}
}

func Test_delete(t *testing.T) {
	type args struct {
		line string
		pos  int
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{"Delete last char", args{"abc", 2}, "ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := delete(tt.args.line, tt.args.pos); gotRes != tt.wantRes {
				t.Errorf("delete() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_jump(t *testing.T) {
	type args struct {
		line string
		dx   int
		lx   int
		ox   int
		ll   int
		ld   int
	}
	tests := []struct {
		name   string
		args   args
		wantNx int
		wantNo int
	}{
		{"Jump word right", args{"one two", 1, 0, 0, 7, 100}, 4, 0},
		{"Jump word left", args{"one two", -1, 7, 0, 7, 100}, 4, 0},
		{"Jump right to end", args{"onetwo", 1, 0, 0, 6, 100}, 6, 0},
		{"Jump left to start", args{"onetwo", -1, 3, 0, 6, 100}, 0, 0},
		{"Jump left from space", args{"one two", -1, 3, 0, 6, 100}, 0, 0},
		{"Jump left from word start", args{"one two", -1, 4, 0, 6, 100}, 0, 0},
		{"Jump word over ld", args{"one two", 1, 0, 0, 7, 3}, 4, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNx, gotNo := jump(tt.args.line, tt.args.dx, tt.args.lx, tt.args.ox, tt.args.ll, tt.args.ld)
			if gotNx != tt.wantNx {
				t.Errorf("jump() gotNx = %v, want %v", gotNx, tt.wantNx)
			}
			if gotNo != tt.wantNo {
				t.Errorf("jump() gotNo = %v, want %v", gotNo, tt.wantNo)
			}
		})
	}
}
