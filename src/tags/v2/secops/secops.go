package ydk

import (
	"fmt"
	logger "github.com/ywteam/ydk-go/logger"
)

func init() {
	logger.Log(logger.Info, "Ydk initialized")
}
func Ydk() {
	scanners := Scanners()
	// fmt.Println(scanners[len(scanners)-1].Id)
	logger.Log(logger.Info, fmt.Sprintf("Running %d scanners", len(scanners)))
}