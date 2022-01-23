package cfginterpolator

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

var (
	LeftSeparator                      = "{{"
	RightSeparator                     = "}}"
	InterpolatorSeparator              = "::"
	InterpolatorConfigurationSeparator = ":"
	// When using watch, specify interval between reloads
	ReloadInterval = 30 * time.Second
)

type Interpolators struct{}

// Interpolate interpolate values in a map[string]interface{}
func Interpolate(conf map[string]interface{}, reloadChan chan string) error {
	for key, val := range conf {
		switch valueWithType := val.(type) {
		case string:
			re := regexp.MustCompile(fmt.Sprintf("%s.*%s.*%s", LeftSeparator, InterpolatorSeparator, RightSeparator))
			match := re.Find([]byte(valueWithType))
			if match == nil {
				continue
			}
			trimmedValue := strings.TrimRight(strings.TrimLeft(string(match), LeftSeparator), RightSeparator)
			splitted := strings.Split(trimmedValue, InterpolatorSeparator)
			if len(splitted) != 2 {
				continue
			}
			conf[key] = useInterpolator(strings.ToLower(splitted[0]), splitted[1])
		case map[string]interface{}:
			Interpolate(valueWithType, reloadChan)
		case []interface{}:
			for index, elem := range valueWithType {
				interpolateInterface(elem, index, valueWithType, reloadChan)
			}
		}
	}
	return nil
}

func interpolateInterface(i interface{}, index int, list []interface{}, reloadChan chan string) {
	switch valueWithType := i.(type) {
	case string:
		re := regexp.MustCompile(fmt.Sprintf("%s.*%s.*%s", LeftSeparator, InterpolatorSeparator, RightSeparator))
		match := re.Find([]byte(valueWithType))
		if match == nil {
			return
		}
		trimmedValue := strings.TrimRight(strings.TrimLeft(string(match), LeftSeparator), RightSeparator)
		splitted := strings.Split(trimmedValue, InterpolatorSeparator)
		if len(splitted) != 2 {
			return
		}
		list[index] = useInterpolator(strings.ToLower(splitted[0]), splitted[1])
	case map[string]interface{}:
		Interpolate(valueWithType, reloadChan)
	case []interface{}:
		for _, elem := range valueWithType {
			interpolateInterface(elem, index, valueWithType, reloadChan)
		}
	}
}

func useInterpolator(interpolatorName string, value string) string {
	interpolatorConf := ""
	if strings.Contains(interpolatorName, InterpolatorConfigurationSeparator) {
		splitted := strings.Split(interpolatorName, InterpolatorConfigurationSeparator)
		interpolatorName = splitted[0]
		interpolatorConf = splitted[1]
	}
	var interpolators Interpolators
	interpotalorFuncName := fmt.Sprintf("%sInterpolator", strings.Title(interpolatorName))
	ret := reflect.ValueOf(&interpolators).MethodByName(interpotalorFuncName).Call([]reflect.Value{reflect.ValueOf(interpolatorConf), reflect.ValueOf(value)})
	retVal := fmt.Sprintf("%v", reflect.ValueOf(ret[0]))
	return retVal
}

// InterpolateFromYAMLFile interpolate YAML file content and write
// interpolated content to out interface{}
func InterpolateFromYAMLFile(yamlFileName string, out interface{}) error {
	var conf map[string]interface{}
	data, err := os.ReadFile(yamlFileName)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		return err
	}
	if err := interpolateAndDecode(conf, out, nil); err != nil {
		return err
	}
	return nil
}

// InterpolateFromYAMLFile interpolate YAML file content and write
// interpolated content to out interface{}. If the data is updated
// in any external sources, it will be updated int he out interface
// InterpolateFromYAMLFile does not end unless an error occured which
// will be returned
func InterpolateAndWatchYAMLFile(yamlFileName string, out interface{}) error {
	var conf map[string]interface{}
	reloadChan := make(chan string, 10)
	data, err := os.ReadFile(yamlFileName)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		return err
	}
	if err := interpolateAndDecode(conf, out, reloadChan); err != nil {
		return err
	}
	for {
		select {
		case <-reloadChan:
			if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
				return err
			}
			if err := interpolateAndDecode(conf, out, nil); err != nil {
				return err
			}
		case <-time.After(ReloadInterval):
			if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
				return err
			}
			if err := interpolateAndDecode(conf, out, nil); err != nil {
				return err
			}
		}
	}
}

func interpolateAndDecode(conf map[string]interface{}, out interface{}, reload chan string) error {
	if err := Interpolate(conf, reload); err != nil {
		return err
	}
	if err := mapstructure.Decode(conf, out); err != nil {
		return err
	}
	return nil
}
