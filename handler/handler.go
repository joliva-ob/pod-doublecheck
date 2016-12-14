package handler


import (

	"net/http"
	"strconv"
	"encoding/json"

)



type MetricType struct {
	Name string
	Value int64
	Threshold int64
	Alert bool
}


var Metrics = make(map[string]*MetricType)
var RefreshTimeChan (chan int)



func StartWebServer(port int) {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/monitoring", MonitoringHandler)
	http.HandleFunc("/refreshtime", RefreshTimeHandler)
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}



func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("200"))
}


func MonitoringHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(Metrics)
	w.Header().Add("Content-Type","application/json;charset=UTF-8")
	w.Write([]byte(string(b)))
}


func AddMetric(name string, value int64, threshold int64) *MetricType {

	alert := false
	if value >= threshold && value > 0 {
		alert = true
	}
	metric := &MetricType{Name: name, Value: value, Threshold: threshold, Alert: alert}
	Metrics[name] = metric

	return metric
}


// PUT http://localhost:8080/refreshtime?20
// Receive: 20
func RefreshTimeHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
	// return the current refresh time
	case "PUT":
		//config_pre.Log.Debugf("PUT param received: %v", r.URL.RawQuery)
		newRefreshTime, _ := strconv.Atoi(r.URL.RawQuery)
		RefreshTimeChan <- newRefreshTime
	}

}