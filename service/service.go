package service


import (

)


var (
	statusChan = make(chan string)
)



func StartService() {


	go doubleCheckProcessor( 10, statusChan )


}




