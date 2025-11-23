package testrunner

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"testrunner": RunCode,
		"add":        RunCode,
		"multiply":   RunCode,
		"version":    RunCode,

		"testrunner.false": RunCode,
		"testrunner.wasm":  RunCodeWithAliases,
	})
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
	})
}
