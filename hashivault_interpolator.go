package cfginterpolator

import (
	"strings"

	vault "github.com/hashicorp/vault/api"
)

func (i *Interpolators) HashiVaultInterpolator(interpolatorConf string, value string) string {
	splitted := strings.Split(value, ":")
	if len(splitted) != 2 {
		return ""
	}
	config := vault.DefaultConfig()
	client, err := vault.NewClient(config)
	if err != nil {
		return ""
	}
	secret, err := client.Logical().Read(splitted[0])
	if err != nil {
		return ""
	}
	if secret == nil {
		return ""
	}
	if interpolatorConf == "KVV2" {
		return secret.Data["data"].(map[string]interface{})[splitted[1]].(string)
	}
	// Default value is KVV1
	return secret.Data[splitted[1]].(string)
}
