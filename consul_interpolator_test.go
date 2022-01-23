package cfginterpolator_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bbayszczak/cfginterpolator"
	consul "github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v3"
)

var consul_http_addr string = "http://consul:8500"
var consul_http_token string = ""

func initConsul() {
	if os.Getenv("CONSUL_HTTP_ADDR") == "" {
		os.Setenv("CONSUL_HTTP_ADDR", consul_http_addr)
	} else {
		consul_http_addr = os.Getenv("CONSUL_HTTP_ADDR")
	}
	if os.Getenv("CONSUL_HTTP_TOKEN") == "" {
		os.Setenv("CONSUL_HTTP_TOKEN", consul_http_token)
	} else {
		consul_http_token = os.Getenv("CONSUL_HTTP_TOKEN")
	}
	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		fmt.Printf("cannot instanciate Consul client: '%s'\n", err)
	}

	p := &consul.KVPair{Key: "path/to/consul_key", Value: []byte("consul_value")}
	_, err = consulClient.KV().Put(p, nil)
	if err != nil {
		fmt.Printf("cannot write in Consul KV: '%s'\n", err)
	}
}

func TestConsulInterpolator(t *testing.T) {
	var i cfginterpolator.Interpolators
	interpolated := i.ConsulInterpolator("", "path/to/consul_key")
	if interpolated != "consul_value" {
		t.Fatalf("value read from Consul is '%s' instead of 'consul_value'", interpolated)
	}
}

func ExampleInterpolator_Consul() {
	var conf map[string]interface{}
	data := `
---
key1: "{{consul::path/to/consul_key}}"
`
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		panic(err)
	}
	cfginterpolator.Interpolate(conf, nil)
	fmt.Println(conf)
	// Output: map[key1:consul_value]
}
