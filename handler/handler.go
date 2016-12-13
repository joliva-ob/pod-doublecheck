package handler


import (

	"net/http"
	"strconv"

)



func StartWebServer(port int) {
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}



func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("200"))
}

// TODO: PUT /start and PUT /stop alerting
// TODO: set check scheduled time window with PUT /checktime
// TODO: expose GET /metrics with

