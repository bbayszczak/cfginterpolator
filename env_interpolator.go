package cfginterpolator

import "os"

func (i *Interpolators) EnvInterpolator(interpolatorConf string, value string, reloadChan chan string) string {
	return os.Getenv(value)
}
