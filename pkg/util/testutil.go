package util

import (
	"runtime"
	"strconv"
	"strings"
)

func GetGoRoutineId() int {
	buf := make([]byte, 30)
	runtime.Stack(buf, false)
	res, _ := strconv.ParseInt(strings.Split(string(buf), " ")[1], 10, 32)
	return int(res)
}
