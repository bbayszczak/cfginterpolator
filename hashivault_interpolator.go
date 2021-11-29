package cfginterpolator

import (
	"strings"

	vault "github.com/hashicorp/vault/api"
)

func (i *Interpolators) HashiVaultInterpolator(value string) string {
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
	return secret.Data[splitted[1]].(string)
}
