package env

import (
	"os"
	"testing"
)

func TestEnvString(t *testing.T) {

	t.Log("TEST_STR missing")
	{
		if String("TEST_STR", "x") != "x" {
			t.Errorf("String did not provide default value")
		}
	}

	t.Log("TEST_STR present")
	{
		os.Setenv("TEST_STR", "y")
		defer os.Unsetenv("TEST_STR")

		if String("TEST_STR", "x") != "y" {
			t.Errorf("String did not obtain env value")
		}
	}
}
