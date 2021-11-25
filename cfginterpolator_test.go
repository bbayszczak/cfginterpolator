package cfginterpolator

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

var data = `
---
key1: "{{env::ENV_VAR1}}"
key2:
  subkey1: "{{env::ENV_VAR2}}"
  subkey2: "{{env::wrongly_formatted::value}}"
key3: value
`

var expected = `
---
key1: "env_var_val1"
key2:
  subkey1: "env_var_val2"
  subkey2: "{{env::wrongly_formatted::value}}"
key3: value
`

// TestInterpolate with yaml
func TestInterpolate(t *testing.T) {
	var conf map[string]interface{}
	var expectedConf map[string]interface{}
	os.Setenv("ENV_VAR1", "env_var_val1")
	os.Setenv("ENV_VAR2", "env_var_val2")
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		t.Fatalf("cannot unmarshall data: %v", err)
	}
	if err := yaml.Unmarshal([]byte(expected), &expectedConf); err != nil {
		t.Fatalf("cannot unmarshall target: %v", err)
	}

	Interpolate(conf)

	if !reflect.DeepEqual(conf, expectedConf) {
		t.Fatalf("expecting %s, has %s", expectedConf, conf)
	}
}
