package cfginterpolator_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bbayszczak/cfginterpolator"
	vault "github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v3"
)

func initVault() {
	if os.Getenv("VAULT_ADDR") == "" {
		os.Setenv("VAULT_ADDR", "http://vault:8200")
	}
	if os.Getenv("VAULT_TOKEN") == "" {
		os.Setenv("VAULT_TOKEN", "myroot")
	}
	err := os.WriteFile(fmt.Sprintf("%s/.vault-token", os.Getenv("HOME")), []byte("myroot"), 0644)
	if err != nil {
		fmt.Printf("cannot create Vault token file: '%s'\n", err)
	}
	config := vault.DefaultConfig()
	client, err := vault.NewClient(config)
	if err != nil {
		fmt.Printf("cannot instanciate Vault client: '%s'\n", err)
	}
	if err := client.Sys().Mount("secretv1", &vault.MountInput{Type: "kv"}); err != nil {
		fmt.Println(err)
	}
	if err := client.Sys().Mount("secretv2", &vault.MountInput{Type: "kv", Options: map[string]string{"version": "2"}}); err != nil {
		fmt.Println(err)
	}

	inputDataV1 := map[string]interface{}{
		"secret_key_v1": "secret_value_kv_v1",
	}
	inputDataV2 := map[string]interface{}{
		"data": map[string]interface{}{
			"secret_key_v2": "secret_value_kv_v2",
		},
	}
	_, err = client.Logical().Write("secretv1/path/to/secret", inputDataV1)
	if err != nil {
		fmt.Println(err)
	}
	_, err = client.Logical().Write("secretv2/data/path/to/secret", inputDataV2)
	if err != nil {
		fmt.Println(err)
	}
}

func TestHashivaultInterpolatorKVV1(t *testing.T) {
	var i cfginterpolator.Interpolators
	interpolated := i.HashivaultInterpolator("kvv1", "secretv1/path/to/secret:secret_key_v1")
	if interpolated != "secret_value_kv_v1" {
		t.Fatalf("value read from vault is '%s' instead of 'secret_value_kv_v1'", interpolated)
	}
}

func TestHashivaultInterpolatorKVV2(t *testing.T) {
	var i cfginterpolator.Interpolators
	interpolated := i.HashivaultInterpolator("kvv2", "secretv2/data/path/to/secret:secret_key_v2")
	if interpolated != "secret_value_kv_v2" {
		t.Fatalf("value read from vault is '%s' instead of 'secret_value_kv_v2'", interpolated)
	}
}

func TestHashivaultInterpolatorKVV1TokenFromFile(t *testing.T) {
	if err := os.Unsetenv("VAULT_TOKEN"); err != nil {
		t.Fatalf("cannot unset env var VAULT_TOKEN: %s", err)
	}
	err := os.WriteFile(fmt.Sprintf("%s/.vault-token", os.Getenv("HOME")), []byte("myroot"), 0644)
	if err != nil {
		fmt.Printf("cannot create Vault token file: '%s'\n", err)
	}
	var i cfginterpolator.Interpolators
	interpolated := i.HashivaultInterpolator("kvv1", "secretv1/path/to/secret:secret_key_v1")
	if interpolated != "secret_value_kv_v1" {
		t.Fatalf("value read from vault is '%s' instead of 'secret_value_kv_v1'", interpolated)
	}
}

func ExampleInterpolator_HashiVault() {
	var conf map[string]interface{}
	data := `
---
key1: "{{hashivault:kvv1::secretv1/path/to/secret:secret_key_v1}}"
key2:
  subkey1: "{{hashivault::secretv1/path/to/secret:secret_key_v1}}"
key4:
  - listkey2: "{{hashivault:kvv2::secretv2/data/path/to/secret:secret_key_v2}}"
    listkey3:
      listsubkey2: "{{hashivault:kvv2::secretv2/data/path/to/secret:secret_key_v2}}"
`
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		panic(err)
	}
	cfginterpolator.Interpolate(conf)
	fmt.Println(conf)
	// Output: map[key1:secret_value_kv_v1 key2:map[subkey1:secret_value_kv_v1] key4:[map[listkey2:secret_value_kv_v2 listkey3:map[listsubkey2:secret_value_kv_v2]]]]
}
