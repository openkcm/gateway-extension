//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	configFile, err := os.ReadFile("../config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	config := string(configFile)

	// create the test cases
	tests := []struct {
		name    string
		config  string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "zero values",
			wantErr: assert.Error,
		}, {
			name:    "invalid config, unknown keys",
			config:  `foo: bar`,
			wantErr: assert.Error,
		}, {
			name:    "valid config",
			config:  config,
			wantErr: assert.NoError,
		},
	}

	// run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// write config.yaml
			if tc.config != "" {
				file := "./config.yaml"
				err := os.WriteFile(file, []byte(tc.config), 0640)
				if err != nil {
					t.Errorf("could not write file: %v, got: %s", file, err)
				}
				defer os.Remove(file)
			}

			// create the command with a timeout context
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			cmd := exec.CommandContext(ctx, "./"+binary)

			type res struct {
				out []byte
				err error
			}
			resCh := make(chan res)
			var b bytes.Buffer
			cmd.Stdout = &b
			cmd.Stderr = &b
			require.NoError(t, cmd.Start())
			go func() {
				err := cmd.Wait()
				resCh <- res{b.Bytes(), err}
			}()

			var result res
			select {
			case <-time.After(5 * time.Second):
				cmd.Process.Signal(os.Interrupt)
				result = <-resCh
			case result = <-resCh:
			}

			if !tc.wantErr(t, result.err, fmt.Sprintf("%s\nerr = %v\nout = %s", cmd.String(), result.err, result.out)) {
				return
			}
		})
	}
}
