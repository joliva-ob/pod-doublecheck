package service


import (

	"github.com/joliva-ob/pod-doublecheck/handler"

	"os"
	"strconv"
)


var (
	RefreshTimeChan = make(chan int)
)



func StartService() {


	refreshTime, _ := strconv.Atoi(os.Getenv("REFRESH_TIME_SECONDS")) //config.Configuration["REFRESH_TIME_SECONDS"].(int)
	handler.RefreshTimeChan = RefreshTimeChan

	go doubleCheckProcessor( refreshTime, RefreshTimeChan )


}




