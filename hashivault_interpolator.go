package cfginterpolator

import (
	"fmt"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

func (i *Interpolators) HashivaultInterpolator(interpolatorConf string, value string) string {
	splitted := strings.Split(value, ":")
	if len(splitted) != 2 {
		return ""
	}
	config := vault.DefaultConfig()
	client, err := vault.NewClient(config)
	if err != nil {
		return ""
	}
	if len(os.Getenv("VAULT_TOKEN")) == 0 {
		if token, err := getVaultTokenFromFile(); err == nil {
			client.SetToken(token)
		}
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
