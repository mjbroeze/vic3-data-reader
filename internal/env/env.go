package env

import (
	"fmt"
	"os"
	"path/filepath"
)

type Key string

func (e Key) GetValue() (string, error) {
	var err error
	v := os.Getenv(string(e))
	if len(v) == 0 {
		fn, ok := defaultEnvMap[e]
		if ok {
			v, err = fn()
		} else {
			err = fmt.Errorf("environment variable %q not found, and default not provided", e)
		}
	}
	return v, err
}

const (
	Vic3Dir Key = "VIC3_DIR"
)

var defaultEnvMap = map[Key]func() (string, error){
	Vic3Dir: func() (string, error) {
		defaultRel := ".local/share/Steam/steamapps/common/Victoria 3"
		usrHome, err := os.UserHomeDir()
		return filepath.Join(usrHome, defaultRel), err
	},
}
