package eureka


import (
	"strconv"
	"log"
	"strings"
	"time"
	"net/http"

	"github.com/hudl/fargo"
	"github.com/satori/go.uuid"
	"github.com/joliva-ob/pod-doublecheck/config"
	"github.com/joliva-ob/pod-doublecheck/handler"
)



var _INSTANCEID string = uuid.NewV4().String()
var _HEARTBEAT_MAX_CONSECUTIVE_ERRORS int = 5
var _HEARTBEAT_SLEEPTIMEBETWEENHEARTBEATINSECONDS time.Duration = 10
var _SECUREPORT int = 8443
var _DATACENTER_NAME string = "MyOwn"


func Register() {

	eurekaUrl := cleanEurekaUrlIfNeeded(config.Configuration["eureka.client.serviceUrl.defaultZone"].(string))
	eurekaConn := fargo.NewConn(eurekaUrl)
	instance := new(fargo.Instance)
	instance.App = config.Configuration["spring.application.name"].(string)
	instance.DataCenterInfo.Name = _DATACENTER_NAME
	instance.HealthCheckUrl = "http://" + config.Configuration["eureka.instance.ip-address"].(string) + ":" + strconv.Itoa(config.Configuration["server.port"].(int)) + "/health"
	instance.HomePageUrl = "http://" + config.Configuration["eureka.instance.ip-address"].(string) + ":" + strconv.Itoa(config.Configuration["server.port"].(int)) + "/"
	instance.StatusPageUrl = "http://" + config.Configuration["eureka.instance.ip-address"].(string) + ":" + strconv.Itoa(config.Configuration["server.port"].(int)) + "/info"
	instance.IPAddr = config.Configuration["eureka.instance.ip-address"].(string)
	instance.HostName = config.Configuration["hostname"].(string)
	instance.SecurePort = _SECUREPORT
	instance.SecureVipAddress = config.Configuration["spring.application.name"].(string)
	instance.VipAddress = config.Configuration["spring.application.name"].(string)
	instance.Status = fargo.StatusType("UP")
	instance.SetMetadataString("instanceId", _INSTANCEID)

	err := eurekaConn.RegisterInstance(instance)

	if err != nil {
		panic("cannot register in eureka")
	}

	startHeartbeat(eurekaUrl, config.Configuration["spring.application.name"].(string), config.Configuration["hostname"].(string), _INSTANCEID)

}

func cleanEurekaUrlIfNeeded(eurekaUrl string) string {
	newEurekaUrl := strings.Split(eurekaUrl, ",")[0]
	if newEurekaUrl[len(newEurekaUrl)-1:] == "/" {
		newEurekaUrl = newEurekaUrl[:len(newEurekaUrl)-1]
	}
	return newEurekaUrl
}

func Deregister() {

	eurekaUrl := cleanEurekaUrlIfNeeded(config.Configuration["eureka.client.serviceUrl.defaultZone"].(string)) + "/apps/" + config.Configuration["spring.application.name"].(string) + "/" + config.Configuration["hostname"].(string) +  ":" + _INSTANCEID
	req, _ := http.NewRequest(http.MethodDelete, eurekaUrl, nil)
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusOK {
		log.Println("Deregistered correctly")
	} else {
		log.Println("Error while deregistering")
	}
}

func startHeartbeat(eurekaUrl string, appName string, hostname string, instance string) {
	consecutiveErrors := 0
	for {
		url := eurekaUrl + "/apps/" + appName + "/" + hostname + ":" + instance

		req, _ := http.NewRequest("PUT", url, nil)
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != http.StatusOK {
			consecutiveErrors++
			if consecutiveErrors >= _HEARTBEAT_MAX_CONSECUTIVE_ERRORS {
				Deregister()
				Register()
			}
		}

		res.Body.Close()
		time.Sleep(_HEARTBEAT_SLEEPTIMEBETWEENHEARTBEATINSECONDS * time.Second)
	}
}


// GetApps returns a map of all Applications
func GetAppsList() map[string]*fargo.Application {

	eurekaUrl := cleanEurekaUrlIfNeeded(config.Configuration["eureka.client.serviceUrl.defaultZone"].(string))
	e := fargo.NewConn(eurekaUrl)
	appsMap, _ := e.GetApps()

	//for i, a := range appsMap {
	//	config_pre.Log.Debugf("%v. Eureka app name: %v", i, a.Name)
	//}
	config.Log.Infof(strconv.Itoa(len(appsMap))+" eureka apps found.")
	handler.AddMetric("Eureka apps", int64(len(appsMap)), 300, nil) // Max apps allowed

	return appsMap
}
