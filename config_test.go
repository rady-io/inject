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

//type ConfigTestSuffixYAML struct {
//	CONF `path:"./resources/application.yaml"`
//}
//
//type ConfigTestTypeYAML struct {
//	CONF `path:"./resources/application.yaml" type:"toml"`
//}

var (
	JSONConf, _ = GetJSONFromAnyFile(DefaultPath, JSON)
	YAMLConf, _ = GetJSONFromAnyFile("./resources/application.yaml", JSON)
)

type ConfigDefaultTest struct {
	Testing
	App *Application
}

func (c *ConfigDefaultTest) TestConfigDefault(t *testing.T) {
	assert.Equal(t, JSONConf, c.App.ConfigFile)
}

type ConfigTypeJSONTest struct {
	Testing
	App *Application
}

func (c *ConfigTypeJSONTest) TestConfigTypeJSON(t *testing.T) {
	assert.Equal(t, YAMLConf, c.App.ConfigFile)
}

func TestConfigFileLoad(t *testing.T) {
	CreateApplication(new(ConfigTestDefault)).RunTest(t, new(ConfigDefaultTest))
	CreateApplication(new(ConfigTestTypeJSON)).RunTest(t, new(ConfigTypeJSONTest))
}
