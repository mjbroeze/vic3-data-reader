package tempenv

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// Mock temporarily sets a single env value for a test.
// After the test is done, the value is reset/unset to the previous state.
// This can print env values to console, so make something better if that is a concern.
func Mock(t *testing.T, envKey string, mockValue string, test func()) error {
	var err error

	sysValue, sysSet := os.LookupEnv(envKey)

	if sysSet {
		t.Log("Temporarily resetting env variable ", envKey, " to ", mockValue)
	}

	err = os.Setenv(envKey, mockValue)
	if err != nil {
		return fmt.Errorf("could not set mock env variable '%s': %s", envKey, err)
	}
	defer func(sysValue string, sysSet bool) {
		if sysSet {
			err = os.Setenv(envKey, sysValue)
			if err != nil {
				log.Fatalf("could not reset mock env variable '%s' to '%s': %v", envKey, sysValue, err)
			}
		} else {
			err = os.Unsetenv(envKey)
			if err != nil {
				log.Fatalf("could not unset mock env variable '%s': %v", envKey, err)
			}
		}
	}(sysValue, sysSet)

	test()
	return err
}
