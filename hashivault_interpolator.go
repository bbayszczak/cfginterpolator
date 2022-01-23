package cfginterpolator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

func (i *Interpolators) HashivaultInterpolator(interpolatorConf string, value string, reloadChan chan string) string {
	client, err := getVaultClient()
	if err != nil {
		return ""
	}
	switch interpolatorConf {
	case "":
		return hashivaultKV(interpolatorConf, value, client)
	case "kvv1":
		return hashivaultKV(interpolatorConf, value, client)
	case "kvv1_json":
		return hashivaultKVJSON(interpolatorConf, value, client)
	case "kvv2":
		return hashivaultKV(interpolatorConf, value, client)
	default:
		return ""
	}
}

func hashivaultKV(interpolatorConf string, value string, client *vault.Client) string {
	splitted := strings.Split(value, ":")
	if len(splitted) != 2 {
		return ""
	}
	secret, err := client.Logical().Read(splitted[0])
	if err != nil {
		return ""
	}
	if secret == nil {
		return ""
	}
	// Default value is KVV1
	if interpolatorConf == "kvv1" || interpolatorConf == "" {
		return secret.Data[splitted[1]].(string)
	}
	if interpolatorConf == "kvv2" {
		return secret.Data["data"].(map[string]interface{})[splitted[1]].(string)
	}
	return ""
}

func hashivaultKVJSON(interpolatorConf string, value string, client *vault.Client) string {
	secret, err := client.Logical().Read(value)
	if err != nil {
		return ""
	}
	if secret == nil {
		return ""
	}
	b, err := json.Marshal(secret.Data)
	if err != nil {
		return ""
	}
	return string(b)
}

func getVaultClient() (*vault.Client, error) {
	config := vault.DefaultConfig()
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}
	if len(os.Getenv("VAULT_TOKEN")) == 0 {
		if token, err := getVaultTokenFromFile(); err == nil {
			client.SetToken(token)
		}
	}
	return client, nil
}

func getVaultTokenFromFile() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("$HOME var empty")
	}
	content, err := os.ReadFile(fmt.Sprintf("%s/.vault-token", home))
	if err != nil {
		return "", err
	}
	if len(content) == 0 {
		return "", fmt.Errorf("vault-token file empty")
	}
	return string(content), nil
}
