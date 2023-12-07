package idx

import (
	"fmt"
	"strconv"
	"strings"
)

// GenJid 数字拼接生成ID
func GenJid(in ...interface{}) int64 {
	var j []string
	for _, i := range in {
		j = append(j, fmt.Sprintf("%v", i))
	}

	j1 := strings.Join(j, "")
	i, err := strconv.ParseInt(j1, 10, 64)
	if err != nil {
		return 0
	}

	return i
}
