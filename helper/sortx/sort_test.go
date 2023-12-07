package sortx

import (
	"fmt"
	"strconv"
	"testing"

	"gitlab.galaxy123.cloud/base/public/helper/randx"
)

type P struct {
	Age  float64 `json:"age"`
	Name string  `json:"name"`
}

func TestSortStruct(t *testing.T) {
	var ps []*P

	for i := 0; i < 5; i++ {
		p := &P{
			Age:  float64(randx.GetRandNum(2, 0)),
			Name: strconv.FormatInt(randx.GetRandNum(2, 0), 10),
		}
		ps = append(ps, p)
	}

	sortStruct := SortStructByFloat[*P](ps, "Age", 1)
	fmt.Println(sortStruct)
}
