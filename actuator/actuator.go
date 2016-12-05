package actuator

import (
	"encoding/json"
	"github.com/joliva-ob/pod-doublecheck/config"
)



func ConfigProps() string {
	configJson := config.Configuration
	b, err := json.Marshal(configJson)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func Info() string {
	infoJson := createInfo()
	b, err := json.Marshal(infoJson)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func createInfo() *InfoJson {
	infoJson := new(InfoJson)

	return infoJson
}

func Health() string {
	healthJson := new(HealthJson)
	b, err := json.Marshal(healthJson)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
