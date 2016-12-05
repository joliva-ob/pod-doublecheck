package service


import (
	 "time"

	"github.com/joliva-ob/pod-doublecheck/config"
)



// Goroutine to schedule and keep looking for differences between kubernetes pods
// and eureka registered apps
func doubleCheckProcessor( checkIntervalTimeSec int, statusChan chan string ) {

	ticker := time.Tick(time.Duration(checkIntervalTimeSec * 1000) * time.Millisecond)
	config.Log.Infof("DoubleChecker started every %v seconds.", checkIntervalTimeSec)

	for {
		select {
		case <- ticker:
			processPods(checkIntervalTimeSec)
		case <- statusChan:
			config.Log.Debugf("BookingChecker is up and running.")
		}
	}

}



func processPods( checkIntervalTimeSec int ) {

	config.Log.Infof("Processing pods...")
	podsMap := make(map[string]string) // k: pod name v: eureka id
	appsMap := make(map[string]string) // k: eureka id v: pod name

	// TODO: get two lists with pods and apps

	// Compare pods and apps to be able to report results
	go compareToReport( podsMap, appsMap )
}




func compareToReport( pods map[string]string, apps map[string]string ){

	config.Log.Infof("Compare and report...")

	// TODO: Compare both items list
}