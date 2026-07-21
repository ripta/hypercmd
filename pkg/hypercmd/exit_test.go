package hypercmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitCode(t *testing.T) {
	plain := errors.New("boom")

	tests := []struct {
		label string
		err   error
		want  int
	}{
		{"nil error", nil, 0},
		{"plain error", plain, 1},
		{"exit error", Exit(42, plain), 42},
		{"exit error with zero code", Exit(0, plain), 0},
		{"wrapped exit error", fmt.Errorf("running command: %w", Exit(17, plain)), 17},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			assert.Equal(t, tt.want, ExitCode(tt.err))
		})
	}
}

func TestExitError(t *testing.T) {
	cause := errors.New("underlying failure")

	err := Exit(3, cause)
	assert.Equal(t, "underlying failure", err.Error())
	assert.ErrorIs(t, err, cause)

	noCause := Exit(4, nil)
	assert.Equal(t, "exit code 4", noCause.Error())
}
