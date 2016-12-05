package service


import (

)


var (
	statusChan = make(chan string)
)



func StartService() {


	go doubleCheckProcessor( 60, statusChan )


}




