package assert

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, expected any, actual any, msgAndArgs ...any) {
	if eq := reflect.DeepEqual(expected, actual); !eq {
		if len(msgAndArgs) > 0 {
			t.Fatalf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		}
		t.Fatalf("expected `%v` but got `%v`", expected, actual)
	}
}

func Error(t *testing.T, err error, msgAndArgs ...any) {
	if err == nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		}
		t.Fatalf("expected error but got nil")
	}
}

func NoError(t *testing.T, err error, msgAndArgs ...any) {
	if err != nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		}
		t.Fatalf("expected nil but got `%v`", err)
	}
}
