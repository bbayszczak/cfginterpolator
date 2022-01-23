package cfginterpolator_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bbayszczak/cfginterpolator"
	"gopkg.in/yaml.v3"
)

func TestEnvInterpolator(t *testing.T) {
	var i cfginterpolator.Interpolators
	name := "ENV_VAR1"
	target := "env_var_VALUE1"
	os.Setenv(name, target)
	interpolated := i.EnvInterpolator("", name, nil)
	if interpolated != target {
		t.Fatalf("env var '%s' is '%s' and should be '%s'", name, interpolated, target)
	}
}

func ExampleInterpolator_Env() {
	var conf map[string]interface{}
	os.Setenv("ENV_VAR_1", "env_var_VAL_1")
	os.Setenv("ENV_VAR_2", "env_var_VAL_2")
	os.Setenv("ENV_VAR_3", "env_var_VAL_3")
	data := `
---
key1: "{{env::ENV_VAR_1}}"
key2:
  subkey1: "{{env::ENV_VAR_1}}"
key4:
  - listkey2: "{{env::ENV_VAR_2}}"
    listkey3:
      listsubkey2: "{{env::ENV_VAR_3}}"
`
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		panic(err)
	}
	cfginterpolator.Interpolate(conf, nil)
	fmt.Println(conf)
	// Output: map[key1:env_var_VAL_1 key2:map[subkey1:env_var_VAL_1] key4:[map[listkey2:env_var_VAL_2 listkey3:map[listsubkey2:env_var_VAL_3]]]]

}
