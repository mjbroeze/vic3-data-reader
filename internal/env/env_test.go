package env

import (
	"os"
	"path/filepath"
	"testing"

	"vic3-data-reader/internal/testframework/tempenv"
)

func TestVic3Dir_default(t *testing.T) {
	hm, err := os.UserHomeDir()
	if err != nil {
		t.Fatal("could not retrieve user home in test setup: ", err)
	}
	defaultPath := filepath.Join(hm, ".local/share/Steam/steamapps/common/Victoria 3")
	defaultVal, err := Vic3Dir.GetValue()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if defaultPath != defaultVal {
		t.Errorf("defaultPath: %s, defaultVal: %s", defaultPath, defaultVal)
	}
}

func TestVic3Dir_custom(t *testing.T) {
	customPath := "/my/custom/path"
	test := func() {
		envVal, err := Vic3Dir.GetValue()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if customPath != envVal {
			t.Errorf("customPath: %s, envVal: %s", customPath, envVal)
		}
	}

	err := tempenv.Mock(t, string(Vic3Dir), customPath, test)
	if err != nil {
		t.Fatal("Mock error: ", err)
	}
}
