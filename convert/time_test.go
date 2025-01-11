package convert

import (
	"fmt"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestCurrentMonthStartDay(t *testing.T) {
	fmt.Println(CurrentMonthStartDay(carbon.Shanghai))
}
