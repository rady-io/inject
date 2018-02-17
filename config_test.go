package rady

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ConfigTestDefault struct {
}

type ConfigTestTypeJSON struct {
	CONF `path:"./resources/application.yaml" type:"json"`
}

type ConfigTestSuffixYAML struct {
	CONF `path:"./resources/application.yaml"`
}

type ConfigTestSuffixJSON struct {
	CONF `path:"./resources/application.json"`
}

type ConfigTestTypeTomlSuffixYAML struct {
	CONF `path:"./resources/application.yaml" type:"toml"`
}

type ConfigTestTypeTomlSuffixJSON struct {
	CONF `path:"./resources/application.json" type:"toml"`
}

var (
	YAMLToJSONConfig, _ = GetJSONFromAnyFile(DefaultPath, DefaultConfType)
	YAMLConf, _         = GetJSONFromAnyFile("./resources/application.yaml", JSON)
)

type ConfigDefaultTest struct {
	Testing
	App *Application
}

func (c *ConfigDefaultTest) TestConfigDefault(t *testing.T) {
	assert.Equal(t, YAMLToJSONConfig, c.App.ConfigFile)
}

type ConfigTypeJSONTest struct {
	Testing
	App *Application
}

func (c *ConfigTypeJSONTest) TestConfigTypeJSON(t *testing.T) {
	assert.Equal(t, YAMLConf, c.App.ConfigFile)
}

func TestConfigFileLoad(t *testing.T) {
	CreateTest(new(ConfigTestDefault)).AddTest(new(ConfigDefaultTest)).Test(t)
	CreateTest(new(ConfigTestTypeJSON)).AddTest(new(ConfigTypeJSONTest)).Test(t)
}

func TestConfigFilePath(t *testing.T) {
	roots := []interface{}{
		new(ConfigTestDefault),
		new(ConfigTestTypeJSON),
		new(ConfigTestSuffixYAML),
		new(ConfigTestSuffixJSON),
		new(ConfigTestTypeTomlSuffixYAML),
		new(ConfigTestTypeTomlSuffixJSON),
	}

	results := [][]string{
		{"./resources/application.conf", YAML},
		{"./resources/application.yaml", JSON},
		{"./resources/application.yaml", YAML},
		{"./resources/application.json", JSON},
		{"./resources/application.yaml", YAML},
		{"./resources/application.json", JSON},
	}
	for i, root := range roots {
		path, Type := CreateApplication(root).GetRealConfigPathAndType()
		assert.Equal(t, results[i][0], path)
		assert.Equal(t, results[i][1], Type)
	}
}
