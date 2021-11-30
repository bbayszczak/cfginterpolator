package cfginterpolator

import "os"

func (i *Interpolators) EnvInterpolator(interpolatorConf string, value string) string {
	return os.Getenv(value)
}
