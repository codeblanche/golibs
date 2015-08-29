package config

import (
	"fmt"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/codeblanche/golibs/logr"
	"github.com/codeblanche/golibs/slice"
)

// Load TOML config files from given root dir into conf interface{}. Load loads prod.toml, stage.toml, test.toml,
// and dev.toml in this order so configurations can be inherited and overriden. Load stops after loading the file
// matching the given env value: prod|stage|test|dev.
func Load(root, env string, conf interface{}) {
	file := fmt.Sprintf("%s.toml", env)
	avail := []string{"prod.toml", "stage.toml", "test.toml", "dev.toml"}
	// default to prod
	if !slice.Strings(avail).Contains(file) {
		file = avail[0]
	}
	logr.Infof("Loading config for env %s", env)
	for _, f := range avail {
		f = path.Join(root, f)
		toml.DecodeFile(f, conf)
		// break after
		if f == file {
			break
		}
	}
	logr.Debugf("%+v", conf)
}
