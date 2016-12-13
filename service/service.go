package service


import (
	"os"
	"strconv"
	"github.com/joliva-ob/pod-doublecheck/handler"
)


var (
	RefreshTimeChan = make(chan int)
)



func StartService() {


	refreshTime, _ := strconv.Atoi(os.Getenv("REFRESH_TIME_SECONDS"))
	handler.RefreshTimeChan = RefreshTimeChan

	go doubleCheckProcessor( refreshTime, RefreshTimeChan )


}




