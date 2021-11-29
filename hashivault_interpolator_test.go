package cfginterpolator_test

import (
	"os"
	"testing"

	"github.com/bbayszczak/cfginterpolator"
	vault "github.com/hashicorp/vault/api"
)

func initVault(t *testing.T) {
	if os.Getenv("VAULT_ADDR") == "" {
		os.Setenv("VAULT_ADDR", "http://vault:8200")
	}
	if os.Getenv("VAULT_TOKEN") == "" {
		os.Setenv("VAULT_TOKEN", "myroot")
	}
	config := vault.DefaultConfig()
	client, err := vault.NewClient(config)
	if err != nil {
		t.Fatalf("cannot instanciate Vault client: '%s'", err)
	}
	if err := client.Sys().Mount("secretv1", &vault.MountInput{Type: "kv"}); err != nil {
		t.Fatal(err)
	}
	if err := client.Sys().Mount("secretv2", &vault.MountInput{Type: "kv", Options: map[string]string{"version": "2"}}); err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}
	_, err = client.Logical().Write("secretv2/data/path/to/secret", inputDataV2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHashiVaultInterpolatorKVV1(t *testing.T) {
	initVault(t)
	var i cfginterpolator.Interpolators
	interpolated := i.HashiVaultInterpolator("secretv1/path/to/secret:secret_key_v1")
	if interpolated != "secret_value_kv_v1" {
		t.Fatalf("value read from vault is '%s' instead of 'secret_value_kv_v1'", interpolated)
	}
}
