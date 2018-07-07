package mutex

import(
	"sync"
	"fmt"
	"strings"
)

type(
	UnlockCbT	func ()
)

func Lock_sync(tags []string) UnlockCbT {
	mtx := sync.Mutex{}
	mtx.Lock()

	return func () {
		mtx.Unlock()
		fmt.Printf("mutex.Unlock [%s]\n", strings.Join(tags, ", "))
	}
}
