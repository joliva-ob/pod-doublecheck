package service


import (
	 "time"
	"strings"

	"github.com/joliva-ob/pod-doublecheck/config"
	"github.com/joliva-ob/pod-doublecheck/kubernetes"
	"github.com/joliva-ob/pod-doublecheck/eureka"
	"github.com/hudl/fargo"
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

	podsMap := kubernetes.GetPodsMap()	// k: pod name  v: bool found in eureka list
	appsList := eureka.GetAppsList()	// k: app name	v: Eureka application

	// Compare pods and apps to be able to report results
	go compareToReport( podsMap, appsList )
}




func compareToReport( pods map[string]bool, apps  map[string]*fargo.Application ){

	// Retrieve search results
	for podName, _ := range pods {

		transformedPodName := applySearchTransformations( strings.ToUpper(podName) )
		found := searchPodNameInEurekaAppList( transformedPodName, apps )
		if !found {
			pods[podName] = false
		}
	}

	// Comunicate search results
	// TODO: send Telegram message as an alert in case of any missing pod registered
	// TODO: expose results to /metrics endpoint!
	i := 0
	for p, b := range pods {
		if !b {
			config.Log.Warningf("Pod not found in Eureka apps list: %v", p)
			i++
		}
	}
	config.Log.Noticef("%v pods not found in Eureka apps list.", i)

}


func applySearchTransformations( podName string ) string {

	if strings.HasPrefix(podName, "API") {
		podName = strings.Replace(podName, "API", "API-", -1)
	} else if strings.HasPrefix(podName, "DAL") {
		podName = strings.Replace(podName, "DAL", "DAL-", -1)
	} else if strings.HasPrefix(podName, "INT-PARTNER-") {
		podName = strings.Replace(podName, "INT-PARTNER-", "PARTNER-CONNECTOR-SERVICE", -1)
	} else if strings.HasPrefix(podName, "INT-TICKETING-") {
		podName = strings.Replace(podName, "INT-TICKETING-", "TICKETING-CONNECTOR-SERVICE", -1)
	} else if strings.HasPrefix(podName, "INT-") {
		podName = strings.Replace(podName, "INT-", "", -1)
	} else if strings.HasSuffix(podName, "-SVC") {
		podName = strings.Replace(podName, "-SVC", "-SERVICE", -1)
	} else if strings.HasPrefix(podName, "INT-VCONFIG-CONV") {
		podName = strings.Replace(podName, "INT-VCONFIG-CONV", "VENUECONFIG-CONVERTER", -1)
	} else if strings.HasPrefix(podName, "INT-VCONFIGCON-") {
		podName = strings.Replace(podName, "INT-VCONFIGCON-", "VENUECONFIG-CONNECTOR", -1)
	} else if strings.HasPrefix(podName, "MSCLIENTS") {
		podName = strings.Replace(podName, "MSCLIENTS", "CLIENTS", -1)
	} else if strings.HasPrefix(podName, "MS") {
		podName = strings.Replace(podName, "MS", "MS-", -1)
	} else if strings.HasPrefix(podName, "PAYMENTS") {
		podName = strings.Replace(podName, "PAYMENTS", "PAYMENT", -1)
	} else if strings.HasPrefix(podName, "DALCOUCH") {
		podName = strings.Replace(podName, "DALCOUCHBASE", "DAL-COUCH", -1)
	}

	return podName
}


func searchPodNameInEurekaAppList( inPodTransformedName string, inAppList map[string]*fargo.Application ) bool {

	isFound := false
	if inPodTransformedName != "CONFIG" || inPodTransformedName != "SPLUNKFORWARDER" {
		isFound = true // It just fake the result in order to avoid it within the final review list
	} else {
		for _, app := range inAppList {
			if strings.Contains(app.Name, inPodTransformedName) {
				isFound = true
				break
			}
		}
	}

	return isFound
}