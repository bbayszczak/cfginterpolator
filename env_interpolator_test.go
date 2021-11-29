package cfginterpolator_test

import (
	"os"
	"testing"

	"github.com/bbayszczak/cfginterpolator"
)

func TestEnvInterpolator(t *testing.T) {
	var i cfginterpolator.Interpolators
	name := "ENV_VAR1"
	target := "env_var_value1"
	os.Setenv(name, target)
	interpolated := i.EnvInterpolator(name)
	if interpolated != target {
		t.Fatalf("env var '%s' is '%s' and should be '%s'", name, interpolated, target)
	}
}
