package cfginterpolator

import (
	consul "github.com/hashicorp/consul/api"
)

func (i *Interpolators) ConsulInterpolator(interpolatorConf string, value string) string {
	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return ""
	}
	pair, _, err := consulClient.KV().Get(value, nil)
	if err != nil {
		return ""
	}
	return string(pair.Value)
}
