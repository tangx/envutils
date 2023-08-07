package envutils

import "testing"

func Test_Import(t *testing.T) {

	config := &struct{}{}
	MustImport("app", config)
}
