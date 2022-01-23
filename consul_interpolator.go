package cfginterpolator

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
	consulwatch "github.com/hashicorp/consul/api/watch"
)

func (i *Interpolators) ConsulInterpolator(interpolatorConf string, value string, reloadChan chan string) string {
	defaultConsulConfig := consul.DefaultConfig()
	if reloadChan != nil {
		go consulWatchKey(defaultConsulConfig, value, reloadChan)
	}
	return consulReadKey(defaultConsulConfig, value)
}

func consulReadKey(consulConfig *consul.Config, keyPath string) string {
	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return ""
	}
	pair, _, err := consulClient.KV().Get(keyPath, nil)
	if err != nil {
		return ""
	}
	if pair == nil {
		return ""
	}
	return string(pair.Value)
}

func consulWatchKey(consulConfig *consul.Config, keyPath string, reloadChan chan string) {
	params := make(map[string]interface{})
	params["type"] = "key"
	params["token"] = consulConfig.Token
	params["key"] = keyPath
	// Create the watch plan
	wp, err := consulwatch.Parse(params)
	if err != nil {
		return
	}
	// Create and test that the API is accessible before starting a blocking
	// loop for the watch.
	//
	// Consul does not have a /ping endpoint, so the /status/leader endpoint
	// will be used as a substitute since it does not require an ACL token to
	// query, and will always return a response to the client, unless there is a
	// network communication error.
	// consulClient, err := consul.NewClient(consulConfig)
	// if err != nil {
	// 	return
	// }
	// _, err = consulClient.Status().Leader()
	// if err != nil {
	// 	return
	// }
	wp.Handler = func(idx uint64, data interface{}) {
		reloadChan <- fmt.Sprintf("consul:%s", keyPath)
	}
	if err := wp.Run(consulConfig.Address); err != nil {
		return
	}
}
