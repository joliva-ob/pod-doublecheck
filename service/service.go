package service


import (

	"github.com/joliva-ob/pod-doublecheck/handler"
	"github.com/oneboxtm/onebox-go-message-processor/config"
)


var (
	RefreshTimeChan = make(chan int)
)



func StartService() {


	refreshTime, _ := config.Configuration["REFRESH_TIME_SECONDS"].(int) //strconv.Atoi(os.Getenv("REFRESH_TIME_SECONDS"))
	handler.RefreshTimeChan = RefreshTimeChan

	go doubleCheckProcessor( refreshTime, RefreshTimeChan )


}




