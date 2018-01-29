package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestGetModeEnv(t *testing.T) {
	assert.Equal(t, GetModeEnv(), "")
	os.Setenv(ModeEnv, TestMod)
	assert.Equal(t, GetModeEnv(), TestMod)
	ResetEnv(ModeEnv)
}

func TestIsTestMode(t *testing.T) {
	assert.False(t, IsTestMode())
	os.Setenv(ModeEnv, TestMod)
	assert.True(t, IsTestMode())
	ResetEnv(ModeEnv)
}

func TestGetConfigFileByMode(t *testing.T) {
	testSets := [][]string {
		{"test", TestMod, "test.test"},
		{"test.yaml", TestMod, "test.test.yaml"},
		{"test.json", TestMod, "test.test.json"},
	}

	for _, testCase := range testSets {
		assert.Equal(t, GetConfigFileByMode(testCase[0]), testCase[0])
		os.Setenv(ModeEnv, testCase[1])
		assert.Equal(t, GetConfigFileByMode(testCase[0]), testCase[2])
		ResetEnv(ModeEnv)
	}
}

func TestIsAutoRollback(t *testing.T) {
	testSets := map[string]bool {
		"true": true,
		"ok": false,
	}

	for env, ok := range testSets {
		assert.False(t, IsAutoRollback())
		os.Setenv(AutoRollbackEnv, env)
		assert.Equal(t, IsAutoRollback(), ok)
		ResetEnv(AutoRollbackEnv)
	}
}