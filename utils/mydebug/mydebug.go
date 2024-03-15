package mydebug

import (
	"fmt"
	"sync"
)

var DEBUG bool = true

// MuPrint
// debug打印锁
// 防止打印信息格式凌乱
var MuPrint sync.Mutex

// DPrintf
// debug函数
func DPrintf(args ...interface{}) {
	if DEBUG {
		MuPrint.Lock()
		defer MuPrint.Unlock()
		fmt.Printf("[Debug] ")
		for _, arg := range args {
			fmt.Printf("%v", arg)
		}
		fmt.Println()
	}
}
