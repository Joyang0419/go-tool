package number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInRange(t *testing.T) {
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
			if got := InRange(tt.args.number, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("InRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "a is greater than b",
			args: args{
				a: 10,
				b: 5,
			},
			want: 10,
		},
		{
			name: "b is greater than a",
			args: args{
				a: 3,
				b: 7,
			},
			want: 7,
		},
		{
			name: "a is equal to b",
			args: args{
				a: 4,
				b: 4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	numbers := []int64{1, 3, 4, 5}
	result := Min(numbers...)
	var expected int64 = 1

	assert.Equal(t, expected, result)
}

func TestAverage(t *testing.T) {
	numbers := []int64{0, 3, 4, 5}
	result := Average(numbers...)
	var expected int64 = 3

	assert.Equal(t, expected, result)
}

func TestSum(t *testing.T) {
	numbers := []int64{1, 3, 4, 5}
	result := Sum(numbers...)
	var expected int64 = 13

	assert.Equal(t, expected, result)
}
