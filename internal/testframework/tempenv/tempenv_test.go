package tempenv

import (
	"log"
	"os"
	"testing"
)

// TestMock_valueSet - tempenv.Mock sets an environment variable for the duration of a test
func TestMock_valueSet(t *testing.T) {
	const EnvVar = "VIC3-DATA-READER_test.TestMock.EnvVar"
	const EnvVal = "test-val"
	test := func() {
		expected := EnvVal
		actual := os.Getenv(EnvVar)
		if actual != expected {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	}

	err := Mock(t, EnvVar, EnvVal, test)
	if err != nil {
		t.Fatal("Mock error: ", err)
	}
}

// TestMock_valueReset - tempenv.Mock resets an environment variable if it was set prior to the test
func TestMock_valueReset(t *testing.T) {
	const EnvVar = "VIC3-DATA-READER_test.TestMock.EnvVar"
	const MockInitVal = "init-val"
	test := func() {
		const MockTestVal = "test-val"

		// overwrite MockInitVal with MockTestVal
		mockTest := func() {}
		err := Mock(t, EnvVar, MockTestVal, mockTest)
		if err != nil {
			t.Fatal("Nested Mock error: ", err)
		}

		// check that the value is reset to MockInitVal
		expected := MockInitVal
		actual := os.Getenv(EnvVar)
		if actual != expected {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
	}

	err := Mock(t, EnvVar, MockInitVal, test)
	if err != nil {
		t.Fatal("Mock error: ", err)
	}
}

// TestMock_valueUnset - tempenv.Mock unsets an environment variable if it was unset prior to the test
func TestMock_valueUnset(t *testing.T) {
	// test data
	const EnvVar = "VIC3-DATA-READER_test.TestMock.EnvVar"
	const EnvVal = "test-val"

	// unset env prior to test
	sysValue, sysSet := os.LookupEnv(EnvVar)
	err := os.Unsetenv(EnvVar)
	if err != nil {
		log.Fatalf("could not unset mock env variable '%s': %v", EnvVar, err)
	}
	// reset env after test
	defer func(sysValue string, sysSet bool) {
		if sysSet {
			err := os.Setenv(EnvVar, sysValue)
			if err != nil {
				log.Fatalf("could not reset mock env variable '%s': %v", EnvVar, err)
			}
		} else {
			err := os.Unsetenv(EnvVar)
			if err != nil {
				log.Fatalf("could not unset mock env variable '%s': %v", EnvVar, err)
			}
		}
	}(sysValue, sysSet)

	// test
	err = Mock(t, EnvVar, EnvVal, func() {})
	if err != nil {
		t.Fatal("Mock error: ", err)
	}

	_, testSet := os.LookupEnv(EnvVar)
	if testSet {
		t.Errorf("mock env variable '%s' should have been unset after test", EnvVar)
	}
}
