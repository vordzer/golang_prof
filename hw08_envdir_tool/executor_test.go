package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	t.Run("empty cmd", func(t *testing.T) {
		result := RunCmd(make([]string, 0), Environment{})
		require.Equal(t, -1, result)
	})

	t.Run("echo cmd", func(t *testing.T) {
		cmd := make([]string, 0)
		cmd = append(cmd, "echo 1")
		result := RunCmd(cmd, Environment{})
		require.Equal(t, 0, result)
	})
}
