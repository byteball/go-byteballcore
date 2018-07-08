package core

import(
	"fmt"
)

type(
)

func Throw(format string, args... interface{}) {
	panic(fmt.Sprintf(format, args...))
}
