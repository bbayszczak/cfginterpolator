# cfginterpolator

![example workflow](https://github.com/bbayszczak/cfginterpolator/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbayszczak/cfginterpolator)](https://goreportcard.com/report/github.com/bbayszczak/cfginterpolator)
[![Go Reference](https://pkg.go.dev/badge/github.com/bbayszczak/cfginterpolator.svg)](https://pkg.go.dev/github.com/bbayszczak/cfginterpolator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

cfginterpolator is an interpolate library in golang allowing to include data from external sources in your configuration

cfginterpolator ingest a `map[string]interface{}` with as many nested `map[string]interface{}` as you want.

This data structure can, as example, be produced by:

- a `yaml` file read by `gopkg.in/yaml.v3`

- a `json` file read by `encoding/json`

After `cfginterpolator` `Interpolate` job, you can use https://github.com/mitchellh/mapstructure to inject the `map[string]interface{}` interpolated by `cfginterpolator` to any kind of `struct` you defined.

## Example

```go
package main

import (
    "github.com/bbayszczak/cfginterpolator"
)

type Config struct {
	Key1 string
	Key2 map[string]string
	Key3 []map[string]string
}
os.Setenv("ENV_VAR_1", "secret_value_1_kv_V1")
os.Setenv("ENV_VAR_2", "secret_value_2_kv_V1")
os.Setenv("ENV_VAR_3", "secret_value_3_kv_V1")
var conf Config
if err := cfginterpolator.InterpolateFromYAMLFile("cfginterpolator/example_files/config.yml", &conf); err != nil {
	panic(err)
}
```

Other examples can be found in the go doc [![Go Reference](https://pkg.go.dev/badge/github.com/bbayszczak/cfginterpolator.svg)](https://pkg.go.dev/github.com/bbayszczak/cfginterpolator)

## Available external datasources

### environment variables

`{{env::ENV_VAR1}}` will be replaced by the value of the environment variable `ENV_VAR1`

### Hashicorp Vault

#### Prequisites

- Environment variable `VAULT_ADDR` should contains the Vault address (e.g.: `https://vault.mydomain.com:8200`)

- A Vault token should exists in environment variable `VAULT_TOKEN` on in file `$HOME/.vault-token`. The enviroment
variable takes predence over the file if both are set.

#### K/V v1

`{{hashivault:kvv1::secret/path/to/secret:key}}` will be replaced by the value of the key `key` of secret `secret/path/to/secret`

K/V v1 is the the default value, the two following expressions act identical: `{{hashivault::secret/path/to/secret:key}}` & `{{hashivault:kvv1::secret/path/to/secret:key}}`

#### K/V v1 JSON

`{{hashivault:kvv1_json::secret/path/to/secret}}` will be replaced by secret `secret/path/to/secret` JSON value

#### K/V v2

`{{hashivault:kvv2::secret/data/path/to/secret:key}}` will be replaced by the value of the key `key` of secret `secret/path/to/secret`

With `K/V v2` you need to add `data` after the secret engine name. `apps/my/secret` will become `apps/data/my/secret` 

## External datasources to be implemented

- [x] environment variables

- [ ] file

- [x] Hashicorp Vault

- [ ] Hashicorp Consul

## Improvements

- [x] interpolate directly from YAML file

- [ ] interpolate directly from JSON file

- [ ] allow to interpolate several times in the same value

- [ ] not panic when interpolator name does not exists

- [ ] add error returns
