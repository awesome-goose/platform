package sources

import (
	"os"
	"strings"

	"github.com/awesome-goose/contracts"
)

type osEnvSource struct{}

func NewOsEnvSource() *osEnvSource {
	return &osEnvSource{}
}

func (v *osEnvSource) Load(env contracts.Env) {
	all := os.Environ()
	for _, kv := range all {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			env.Set(parts[0], parts[1])
		}
	}
}
