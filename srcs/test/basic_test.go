package testsample

import (
	"testing"
	"fmt"
)

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	got := Add(1, 2)
	if got != 3 {
		t.Errorf("expect 3, but %d", got)
	}
}

func Calc(a, b int, operator string) (int, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("0 division is undefined.")
		}
		return a / b, nil
	}
	return 0, fmt.Errorf("unexpected operator: %v", operator)
}

func TestCalc(t *testing.T) {
	type args struct {
		a int
		b int
		operator string
	}

	type testCase struct {
		name string
		args args
		want int
		wantErr bool
	}

	tests := [] testCase{
		{
			name: "plus",
			args: args{
				a: 10,
				b: 2,
				operator: "+",
			},
			want: 12,
			wantErr: false,
		},
		{
			name: "abnormal",
			args: args{
				a: 10,
				b: 2,
				operator: "?",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		got, err := Calc(tt.args.a, tt.args.b, tt.args.operator)
		if err != nil && tt.wantErr == false {
			t.Errorf("wantErr is wrong!!")
			return
		}
		if got != tt.want {
			t.Errorf("Calc() = %v want %v", got, tt.want)
		}
	}
}
