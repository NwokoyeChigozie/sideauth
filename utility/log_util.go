package utility

import (
	"fmt"
)

func LogAndPrint(data interface{}, args ...interface{}) {
	logger := NewLogger()
	fmt.Println(data, args)
	logger.Info(data, args)
}
