package console

import(
	"fmt"
)

/**
var(
	console struct{
		Log func(string, ...interface{})
	}
)

func init() {
	console.Log = func(format string, args... interface{}) {
		fmt.Printf(format+"\n", args...)
	}
}
 **/

func Log(format string, args... interface{}) {
		fmt.Printf(format+"\n", args...)
}	
