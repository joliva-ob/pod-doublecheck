package config

import (
	"os"
	"net/http"
	"io/ioutil"
	"reflect"
	"strings"
	"strconv"
	"fmt"
	"bufio"

	"github.com/Jeffail/gabs"
	"github.com/op/go-logging"

)


var _DEFAULT_PORT int = 8080
var Configuration map[string]interface{} = make(map[string]interface{})
var Log *logging.Logger


func LoadConfig() {

	preload()
	printBootstrap()

	// Set logger
	format := logging.MustStringFormatter( os.Getenv("LOG_FORMAT") )
	logbackend1 := logging.NewLogBackend(os.Stdout, "", 0)
	logbackend1Formatted := logging.NewBackendFormatter(logbackend1, format)
	f, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		defer f.Close()
	}
	logbackend2 := logging.NewLogBackend(f, "", 0)
	logbackend2Formatted := logging.NewBackendFormatter(logbackend2, format)
	logging.SetBackend(logbackend1Formatted, logbackend2Formatted)
	Log = logging.MustGetLogger( os.Getenv("spring_application_name") )

	Log.Info("Loading config_pre...")
	newConfig := loadBasicsFromEnvironmentVars()
	getConfigFromSpringCloudConfigServer(newConfig["spring.cloud.config_pro.uri"].(string), newConfig)
	Configuration = newConfig
	Log.Info("Config loaded sucessfuly")

}




func loadBasicsFromEnvironmentVars() map[string]interface{} {

	var newConfig map[string]interface{} = make(map[string]interface{})
	newConfig["spring.profiles.active"] = os.Getenv("spring_profiles_active")
	newConfig["spring.cloud.config_pre.uri"] = os.Getenv("spring_cloud_config_uri")
	newConfig["spring.cloud.config_pre.label"] = os.Getenv("spring_cloud_config_label")
	newConfig["server.port"] = os.Getenv("server_port")
	newConfig["eureka.instance.ip-address"] = os.Getenv("eureka_instance_ip_address")
	newConfig["spring.application.name"] = os.Getenv("spring_application_name")
	newConfig["hostname"], _ = os.Hostname()
	newConfig["LOG_FORMAT"] = os.Getenv("LOG_FORMAT")
	newConfig["LOG_FILE"] = os.Getenv("LOG_FILE")
	newConfig["EUREKA_APP_NAME"] = os.Getenv("EUREKA_APP_NAME")
	newConfig["EUREKA_PUBLIC_HOST"] = os.Getenv("EUREKA_PUBLIC_HOST")
	newConfig["REFRESH_TIME_SECONDS"] = os.Getenv("REFRESH_TIME_SECONDS")
	newConfig["ENV"] = os.Getenv("ENV")
	newConfig["TELEGRAM_CHAT_ID"] = os.Getenv("TELEGRAM_CHAT_ID")

	port, err := strconv.Atoi(newConfig["server.port"].(string))
	if err != nil {
		newConfig["server.port"] = _DEFAULT_PORT
	} else {
		newConfig["server.port"] = port
	}

	if newConfig["spring.profiles.active"] == "" || newConfig["spring.cloud.config_pre.uri"] == "" || newConfig["spring.cloud.config_pre.label"] == "" || newConfig["server.port"] == "" || newConfig["eureka.instance.ip-address"] == 0 || newConfig["spring.application.name"] == "" {
		panic("spring_profiles_active , spring_cloud_config_uri , spring_cloud_config_label , server_port , eureka_instance_ip_address, spring_application_name environment vars are mandatories")
	}

	return newConfig
}

func getConfigFromSpringCloudConfigServer(uriEndpoint string, newConfig map[string]interface{}) {
	finalEndpoint := uriEndpoint + "/" + newConfig["spring.application.name"].(string) + "/" + newConfig["spring.profiles.active"].(string) + "/" + newConfig["spring.cloud.config_pre.label"].(string)
	Log.Info("Getting config_pre from " + finalEndpoint)
	rs, err := getJsonFromSpringCloudConfigServer(finalEndpoint)
	if err != nil {
		panic("can't load configuration from " + finalEndpoint)
	}
	rewriteConfig(rs, newConfig)
}

func rewriteConfig(container *gabs.Container, newConfig map[string]interface{}) {
	newConfig["label"], _ = container.Path("label").Data().(string)
	newConfig["name"], _ = container.Path("name").Data().(string)
	source := container.Path("propertySources").Path("source")
	propertySources, _ := source.Children()

	iterateOverEachKeyAndReplaceVars(propertySources, newConfig)
	replaceVars(newConfig)

}

func replaceVars(newConfig map[string]interface{}) {
	for field, value := range newConfig {
		if isString(value) {
			if strings.Contains(value.(string), "${") {
				modifiedValue := value.(string)
				splitted := strings.Split(value.(string), "${")
				for i:=0;i<len(splitted);i++ {
					fieldToFind := strings.Split(splitted[i], "}")[0]
					if newConfig[fieldToFind] != nil {
						modifiedValue = strings.Replace(modifiedValue, "${"+fieldToFind+"}", newConfig[fieldToFind].(string), 10)
						newConfig[field] = modifiedValue
					}
				}
			}
		}
	}
}

func iterateOverEachKeyAndReplaceVars(containers []*gabs.Container, newConfig map[string]interface{}) {
	for _, child := range containers {
		keyvalueconfigurationmap, _ := child.ChildrenMap()
		for configurationField, configurationValue := range keyvalueconfigurationmap {
			modifiedConfigurationValue := configurationValue.Data()
			if isString(modifiedConfigurationValue) {
				if configurationValueThanMustBeReplaced(modifiedConfigurationValue) {
					modifiedConfigurationValue = replaceConfigurationValueAndReturnTheNewValue(modifiedConfigurationValue, newConfig)
				}
			}
			addNewKeyValueToConfigurationIfNotExists(configurationField, modifiedConfigurationValue, newConfig)
		}

	}
}

func replaceConfigurationValueAndReturnTheNewValue(modifiedConfigurationValue interface{}, newConfig map[string]interface{}) interface{} {
	splitted := strings.Split(modifiedConfigurationValue.(string), "${")
	for i:=1; i<len(splitted);i++  {
		fieldName := strings.Split(splitted[i], "}")[0]
		if newConfig[fieldName] != nil {
			fieldValue := newConfig[fieldName].(string)
			modifiedConfigurationValue = strings.Replace(modifiedConfigurationValue.(string), "${"+fieldName+"}", fieldValue, 1)
		}
	}
	return modifiedConfigurationValue
}

func configurationValueThanMustBeReplaced(data interface{}) bool {
	return strings.Contains(data.(string), "${")
}
func isString(data interface{}) bool {
	return reflect.TypeOf(data).String() == "string"
}

func addNewKeyValueToConfigurationIfNotExists(key string, value interface{}, newConfig map[string]interface{}) {
	if newConfig[key] == nil {
		newConfig[key] = value
	}
}

func getJsonFromSpringCloudConfigServer(url string) (*gabs.Container, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	jsonParsed, err := gabs.ParseJSON(b)

	if err != nil {
		return nil, err
	}

	return jsonParsed, nil
}

func AddKeyValueToConfig(key string, value interface{}) {
	Configuration[key] = value
}


func preload() {

	os.Setenv("spring_application_name", os.Getenv("EUREKA_APP_NAME"))
	os.Setenv("eureka_instance_ip_address", os.Getenv("EUREKA_PUBLIC_HOST"))

}


func printBootstrap() {

	lines, _ := readLines("resources/boot_logo.txt")
	for _, line := range lines {
		fmt.Println(line)
	}

}


// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}