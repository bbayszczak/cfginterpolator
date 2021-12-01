package cfginterpolator

import (
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
