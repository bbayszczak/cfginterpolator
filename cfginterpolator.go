package cfginterpolator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	leftSeparator         = "{{"
	rightSeparator        = "}}"
	interpolatorSeparator = "::"
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
			conf[key] = useInterpolator(splitted[0], splitted[1])
		case map[string]interface{}:
			Interpolate(valueWithType)
		}
	}
	return nil
}

func useInterpolator(interpolatorName string, value string) string {
	var interpolators Interpolators
	interpotalorFuncName := fmt.Sprintf("%sInterpolator", strings.Title(interpolatorName))
	ret := reflect.ValueOf(&interpolators).MethodByName(interpotalorFuncName).Call([]reflect.Value{reflect.ValueOf(value)})
	retVal := fmt.Sprintf("%v", reflect.ValueOf(ret[0]))
	return retVal
}
