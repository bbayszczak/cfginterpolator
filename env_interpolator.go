package cfginterpolator

import "os"

func (i *Interpolators) EnvInterpolator(value string) string {
	return os.Getenv(value)
}
