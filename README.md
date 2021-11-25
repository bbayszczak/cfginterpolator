# cfginterpolator

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
    "github.com/mitchellh/mapstructure"
    "github.com/bbayszczak/cfginterpolator"
)

type Config struct {
    Username string
    Password string
}

var rawConfig = `
---
username: "John-David"
password: "{{env::PASSWORD}}"
`

func main() {
    var config map[string]interface{}
    var configStruct Config
    if err := yaml.Unmarshal([]byte(rawConfig), &config); err != nil {
		panic("cannot unmarshall data: %v", err)
	}
    cfginterpolator.Interpolate(config)
    mapstructure.Decode(config, &configStruct)
}
```

## Available external datasources

### environment variables

`{{env::ENV_VAR1}}` will be replaced by the value of the environment variable `ENV_VAR1`

## External datasources to be implemented

- [x] environment variables

- [ ] file

- [ ] Hashicorp Vault

- [ ] Hashicorp Consul

## Improvements

- [ ] allow to interpolate several times in the same value

- [ ] not panic when interpolator name does not exists

- [ ] add error returns
