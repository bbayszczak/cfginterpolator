package cfginterpolator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	leftSeparator                      = "{{"
	rightSeparator                     = "}}"
	interpolatorSeparator              = "::"
	interpolatorConfigurationSeparator = ":"
)

type Interpolators struct{}

func Interpolate(conf map[string]interface{}) error {
	for key, val := range conf {
		switch valueWithType := val.(type) {
		case string:
			re := regexp.MustCompile(fmt.Sprintf("%s.*%s.*%s", leftSeparator, interpolatorSeparator, rightSeparator))
			match := re.Find([]byte(valueWithType))
			if match == nil {
				continue
			}
			trimmedValue := strings.TrimRight(strings.TrimLeft(string(match), leftSeparator), rightSeparator)
			splitted := strings.Split(trimmedValue, interpolatorSeparator)
			if len(splitted) != 2 {
				continue
			}
			conf[key] = strings.ToLower(useInterpolator(splitted[0], splitted[1]))
		case map[string]interface{}:
			Interpolate(valueWithType)
		case []interface{}:
			for index, elem := range valueWithType {
				interpolateInterface(elem, index, valueWithType)
			}
		}
	}
	return nil
}

func interpolateInterface(i interface{}, index int, list []interface{}) {
	switch valueWithType := i.(type) {
	case string:
		re := regexp.MustCompile(fmt.Sprintf("%s.*%s.*%s", leftSeparator, interpolatorSeparator, rightSeparator))
		match := re.Find([]byte(valueWithType))
		if match == nil {
			return
		}
		trimmedValue := strings.TrimRight(strings.TrimLeft(string(match), leftSeparator), rightSeparator)
		splitted := strings.Split(trimmedValue, interpolatorSeparator)
		if len(splitted) != 2 {
			return
		}
		list[index] = useInterpolator(strings.ToLower(splitted[0]), splitted[1])
	case map[string]interface{}:
		Interpolate(valueWithType)
	case []interface{}:
		for _, elem := range valueWithType {
			interpolateInterface(elem, index, valueWithType)
		}
	}
}

func useInterpolator(interpolatorName string, value string) string {
	interpolatorConf := ""
	if strings.Contains(interpolatorName, interpolatorConfigurationSeparator) {
		splitted := strings.Split(interpolatorName, interpolatorConfigurationSeparator)
		interpolatorName = splitted[0]
		interpolatorConf = splitted[1]
	}
	var interpolators Interpolators
	interpotalorFuncName := fmt.Sprintf("%sInterpolator", strings.Title(interpolatorName))
	ret := reflect.ValueOf(&interpolators).MethodByName(interpotalorFuncName).Call([]reflect.Value{reflect.ValueOf(interpolatorConf), reflect.ValueOf(value)})
	retVal := fmt.Sprintf("%v", reflect.ValueOf(ret[0]))
	return retVal
}
