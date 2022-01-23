package cfginterpolator_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

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

	cfginterpolator.Interpolate(conf, nil)

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
	if err := cfginterpolator.InterpolateFromYAMLFile("example_files/config.yml", &conf); err != nil {
		panic(err)
	}
	fmt.Println(conf)
	// Output: {secret_value_1_kv_V1 map[subkey1:secret_value_2_kv_V1] [map[listkey1:secret_value_3_kv_V1]]}
}

func ExampleInterpolateAndWatchYAMLFile() {
	type Config struct {
		Key1 string
		Key2 map[string]string
		Key3 []map[string]string
	}
	os.Setenv("ENV_VAR_1", "var1_0")
	os.Setenv("ENV_VAR_2", "var2")
	os.Setenv("ENV_VAR_3", "var3")
	var conf Config

	cfginterpolator.ReloadInterval = time.Second
	go func() {
		time.Sleep(500 * time.Millisecond)
		if err := cfginterpolator.InterpolateAndWatchYAMLFile("example_files/config.yml", &conf); err != nil {
			fmt.Println(err)
		}
	}()
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Second)
		os.Setenv("ENV_VAR_1", fmt.Sprintf("var1_%d", i))
		fmt.Printf("%d: %s\n", i, conf)
	}

	// Output:
	//1: {var1_0 map[subkey1:var2] [map[listkey1:var3]]}
	//2: {var1_1 map[subkey1:var2] [map[listkey1:var3]]}
	//3: {var1_2 map[subkey1:var2] [map[listkey1:var3]]}
}

func ExampleInterpolateAndWatchYAMLFile_Consul() {
	type Config struct {
		Key1 string
	}
	var conf Config

	cfginterpolator.ReloadInterval = 1 * time.Minute
	go func() {
		if err := cfginterpolator.InterpolateAndWatchYAMLFile("example_files/config_consul.yml", &conf); err != nil {
			fmt.Println(err)
		}
	}()
	for i := 1; i <= 3; i++ {
		writeConsulKey("path/to/consul_key_watch", fmt.Sprintf("consul_watch_val%d", i))
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%d: %s\n", i, conf)
	}

	// Output:
	//1: {consul_watch_val1}
	//2: {consul_watch_val2}
	//3: {consul_watch_val3}
}
