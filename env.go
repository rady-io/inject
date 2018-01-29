package rady

import (
	"os"
	"strings"
	"fmt"
)

const (
	ModeEnv = "RADY_MODE"
	TestMod = "test"
)

func GetModeEnv() string {
	return os.Getenv(ModeEnv)
}

func IsTestMode() bool {
	return GetModeEnv() == TestMod
}

func GetConfigFileByMode(filePath string) string {
	mode := GetModeEnv()
	if mode == "" {
		return filePath
	}
	index := strings.LastIndexByte(filePath, os.PathSeparator)
	return fmt.Sprintf("%s%s.%s", filePath[:index+1], mode, filePath[index+1:])
}
