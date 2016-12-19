package service


import (
	"time"
	"strconv"
	"strings"
	"bytes"

	"github.com/joliva-ob/pod-doublecheck/config"
	"github.com/joliva-ob/pod-doublecheck/kubernetes"
	"github.com/joliva-ob/pod-doublecheck/eureka"
	"github.com/joliva-ob/jacaranda/jacaranda-client"
	"github.com/hudl/fargo"
	"github.com/joliva-ob/pod-doublecheck/handler"
)


// Goroutine to schedule and keep looking for differences between kubernetes pods
// and eureka registered apps
func doubleCheckProcessor( checkIntervalTimeSec int, RefreshTimeChan chan int ) {

	timer := time.NewTimer(time.Duration(checkIntervalTimeSec * 1000) * time.Millisecond)
	config.Log.Infof("DoubleChecker started every %v seconds.", checkIntervalTimeSec)

	for {
		select {
		case <- timer.C:
			processPods( checkIntervalTimeSec, timer )
		case checkIntervalTimeSec = <- RefreshTimeChan:
			reScheduleCheckProcessor(checkIntervalTimeSec, timer)
		}
	}

}


func reScheduleCheckProcessor( newRefreshTime int, timer *time.Timer ) {

	if newRefreshTime <= 0 {
		timer.Stop()
		config.Log.Noticef("Double-Check timer is now stopped (%v).", newRefreshTime)
	} else {
		timer.Reset(time.Duration(newRefreshTime * 1000) * time.Millisecond)
		config.Log.Noticef("New refresh time is: %v", newRefreshTime)
	}
}



func processPods( refreshTime int, timer *time.Timer ) {

	podsMap := kubernetes.GetPodsMap()	// k: pod name  v: bool found in eureka list
	appsList := eureka.GetAppsList()	// k: app name	v: Eureka application

	// Compare pods and apps to be able to report results
	go compareToReport( podsMap, appsList )

	timer.Reset(time.Duration(refreshTime * 1000) * time.Millisecond)
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
	i := 0
	var appsNotFoundBuffer bytes.Buffer
	var appsNotFoundList []string
	for p, b := range pods {
		if !b {
			appsNotFoundList = append(appsNotFoundList, p)
			i++
			appsNotFoundBuffer.WriteString(strconv.Itoa(i)+". "+p+"\n\r")
			config.Log.Warningf("Pod not found in Eureka app list: %v", p)
		}
	}
	handler.AddMetric("Pods not found", int64(i), 0, appsNotFoundList)
	config.Log.Noticef("%v pods not found in Eureka apps list.", i)

	if i > 0 {
		message := "Alert: ["+config.Configuration["ENV"].(string)+"] There are "+strconv.Itoa(i)+" pods not registered into Eureka!\n\r"+appsNotFoundBuffer.String()
		chatId := config.Configuration["TELEGRAM_CHAT_ID"].(string)
		url := "http://10.1.2.173:30002/jacaranda/1.0/sendMessage"
		res, err := jacaranda_client.SendTelegramMessage(url, message, chatId)
		if err != nil {
			config.Log.Errorf("ERROR sending message <%v> to <%v>",message,chatId)
		} else {
			config.Log.Infof("Alert message %v successfuly sent to %v with response:%v",message,chatId,res.Status)
		}
	}

}


func applySearchTransformations( podName string ) string {

	if strings.HasPrefix(podName, "API") && !strings.HasPrefix(podName, "API-") {
		podName = strings.Replace(podName, "API", "API-", -1)
	} else if strings.HasPrefix(podName, "DALCOUCHBASE") {
		podName = strings.Replace(podName, "DALCOUCHBASE", "DAL-COUCH", -1)
	} else if strings.HasPrefix(podName, "DAL") {
		podName = strings.Replace(podName, "DAL", "DAL-", -1)
	} else if strings.HasPrefix(podName, "INT-PARTNER-") {
		podName = strings.Replace(podName, "INT-PARTNER-", "PARTNER-CONNECTOR-SERVICE", -1)
	} else if strings.HasPrefix(podName, "INT-TICKETING-") {
		podName = strings.Replace(podName, "INT-TICKETING-", "TICKETING-CONNECTOR-SERVICE", -1)
	} else if strings.HasPrefix(podName, "INT-VENUECONFIG") {
		podName = strings.Replace(podName, "INT-VENUECONFIG-", "VENUECONFIG-SERVICE", -1)
	} else if strings.HasSuffix(podName, "-SVC") {
		podName = strings.Replace(podName, "-SVC", "-SERVICE", -1)
	} else if strings.HasPrefix(podName, "INT-VCONFIG-CONV") {
		podName = strings.Replace(podName, "INT-VCONFIG-CONV", "VENUECONFIG-CONVERTER", -1)
	} else if strings.HasPrefix(podName, "INT-VCONFIGCON-") {
		podName = strings.Replace(podName, "INT-VCONFIGCON-", "VENUECONFIG-CONNECTOR", -1)
	} else if strings.HasPrefix(podName, "INT-") {
		podName = strings.Replace(podName, "INT-", "", -1)
	} else if strings.HasPrefix(podName, "MSCLIENTS") {
		podName = strings.Replace(podName, "MSCLIENTS", "CLIENTS", -1)
	} else if strings.HasPrefix(podName, "MS") {
		podName = strings.Replace(podName, "MS", "MS-", -1)
	} else if strings.HasPrefix(podName, "PAYMENTS") {
		podName = strings.Replace(podName, "PAYMENTS", "PAYMENT", -1)
	} else if strings.HasPrefix(podName, "GODOUBLECHECK") {
		podName = strings.Replace(podName, "GO", "POD-", -1)
	}

	return podName
}


func searchPodNameInEurekaAppList( inPodTransformedName string, inAppList map[string]*fargo.Application ) bool {

	isFound := false
	if inPodTransformedName == "CONFIG" || inPodTransformedName == "SPLUNKFORWARDER" || strings.HasSuffix(inPodTransformedName,"TEST") || inPodTransformedName == "PROMETHEUS" {
		isFound = true // It just fake the result in order to avoid its verification, those never been registered into Eureka.
		//config_pre.Log.Warningf("Pod FOUND in Eureka special apps list: %v", inPodTransformedName)
	} else {
		for _, app := range inAppList {
			if strings.Contains(app.Name, inPodTransformedName) {
				//config_pre.Log.Warningf("Pod %v FOUND into Eureka app %v", inPodTransformedName, app.Name)
				isFound = true
				break
			}
		}
	}

	return isFound
}