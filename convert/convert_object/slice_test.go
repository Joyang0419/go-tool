package convert_object

import (
	"fmt"
	"testing"
)

func TestCreateSlice(t *testing.T) {
	s := CreateSlice[float64]("1", "2", "3")
	fmt.Println(s)
}

func TestCreateSlice2(t *testing.T) {
	type model struct {
		Name string
	}

	s := CreateSlice[model](model{Name: "a"}, model{Name: "b"})
	fmt.Println(s)
}

func TestSliceToMap(t *testing.T) {
	s := CreateSlice[float64]("1", "2", "3")

	toMapFn := func(i float64) (string, int) {
		return fmt.Sprintf("%f", i), int(i)
	}
	m := SliceToMap(s, toMapFn)

	fmt.Println(m)
}

func TestSliceToSet(t *testing.T) {
	s := CreateSlice[float64]("1", "2", "3", "1", "2", "3")

	uniq := SliceToSet(s)

	fmt.Println(uniq)
}

func TestSliceContains(t *testing.T) {
	s := CreateSlice[float64]("1", "2", "3", "1", "2", "3")

	fmt.Println(SliceContains(s, float64(1)))
}

func TestSliceToSpecifiedType(t *testing.T) {
	s := CreateSlice[float64](1, 2, 3)

	s2 := SliceToSpecifiedType(s, func(i float64) string {
		return fmt.Sprintf("%f", i)
	})

	fmt.Println(s2)
}
