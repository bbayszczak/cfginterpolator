package cfginterpolator_test

import (
	"fmt"
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

func TestMain(m *testing.M) {
	initVault()
	exitCode := m.Run()
	os.Exit(exitCode)
}

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

func ExampleInterpolateFromYAMLFile() {
	type Config struct {
		Key1 string
		Key2 map[string]string
		Key3 []map[string]string
	}
	os.Setenv("ENV_VAR_1", "secret_value_1_kv_V1")
	os.Setenv("ENV_VAR_2", "secret_value_2_kv_V1")
	os.Setenv("ENV_VAR_3", "secret_value_3_kv_V1")
	var conf Config
	if err := cfginterpolator.InterpolateFromYAMLFile("/cfginterpolator/example_files/config.yml", &conf); err != nil {
		panic(err)
	}
	fmt.Println(conf)
	// Output: {secret_value_1_kv_V1 map[subkey1:secret_value_2_kv_V1] [map[listkey1:secret_value_3_kv_V1]]}
}
