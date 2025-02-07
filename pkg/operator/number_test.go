package operator

import (
	"testing"
)

func TestNumberInRange(t *testing.T) {
	type args struct {
		number int
		min    int
		max    int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "number is in range",
			args: args{
				number: 5,
				min:    1,
				max:    10,
			},
			want: true,
		},
		{
			name: "number is out of range",
			args: args{
				number: 15,
				min:    1,
				max:    10,
			},
			want: false,
		},
		{
			name: "number is equal to min",
			args: args{
				number: 1,
				min:    1,
				max:    10,
			},
			want: true,
		},
		{
			name: "number is equal to max",
			args: args{
				number: 10,
				min:    1,
				max:    10,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberInRange(tt.args.number, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("NumberInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
