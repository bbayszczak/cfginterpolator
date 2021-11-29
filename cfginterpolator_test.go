package cfginterpolator_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/bbayszczak/cfginterpolator"
	"gopkg.in/yaml.v3"
)

var data = `
---
key1: "{{env::ENV_VAR1}}"
key2:
  subkey1: "{{env::ENV_VAR2}}"
  subkey2: "{{env::wrongly_formatted::value}}"
key3: value
key4:
  - listkey1: listvalue1
  - listkey2: "{{env::ENV_VAR3}}"
  - listkey3:
      listsubkey1: listsubvalue1
      listsubkey2: "{{env::ENV_VAR4}}"
`

var expected = `
---
key1: "env_var_val1"
key2:
  subkey1: "env_var_val2"
  subkey2: "{{env::wrongly_formatted::value}}"
key3: value
key4:
  - listkey1: listvalue1
  - listkey2: "env_var_val3"
  - listkey3:
      listsubkey1: listsubvalue1
      listsubkey2: "env_var_val4"
`

// TestInterpolate with yaml
func TestInterpolate(t *testing.T) {
	var conf map[string]interface{}
	var expectedConf map[string]interface{}
	os.Setenv("ENV_VAR1", "env_var_val1")
	os.Setenv("ENV_VAR2", "env_var_val2")
	os.Setenv("ENV_VAR3", "env_var_val3")
	os.Setenv("ENV_VAR4", "env_var_val4")
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		t.Fatalf("cannot unmarshall data: %v", err)
	}
	if err := yaml.Unmarshal([]byte(expected), &expectedConf); err != nil {
		t.Fatalf("cannot unmarshall target: %v", err)
	}

	cfginterpolator.Interpolate(conf)

	if !reflect.DeepEqual(conf, expectedConf) {
		t.Fatalf("expecting %s, has %s", expectedConf, conf)
	}
}
