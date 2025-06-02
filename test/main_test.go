//go:build integration

package integration_test

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

const binary = "gateway-extension"

func buildCmd(cmd string) {
	build := exec.Command("go", "build", "-buildvcs=false", "-race", "-cover", "-o", cmd, "../cmd/"+cmd)
	if output, err := build.CombinedOutput(); err != nil {
		log.Fatalf("go build ../cmd/%s: output = %s, err = %s", cmd, output, err)
	}
}

func TestMain(m *testing.M) {
	// Pre-build the app
	buildCmd(binary)
	exitCode := m.Run()
	_ = os.Remove(binary)
	_ = os.Remove("extension.sock")

	os.Exit(exitCode)
}
