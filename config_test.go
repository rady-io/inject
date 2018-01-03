package rhapsody

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type ConfigTestDefault struct {
	CONF
}

type ConfigTestTypeJSON struct {
	CONF `path:"./resources/application.yaml" type:"json"`
}

type ConfigTestSuffixYAML struct {
	CONF `path:"./resources/application.yaml"`
}

type ConfigTestTypeYAML struct {
	CONF `path:"./resources/application.yaml" type:"toml"`
}

var (
	JSONConf, _ = GetJSONFromAnyFile(DefaultPath, JSON)
	JSONFromYAML = "{\"rhapsody\":{\"db\":{\"type\":\"mysql\"},\"redis\":{\"host\":\"127.0.0.1\",\"port\":6937}}}"
	YAMLConf, _ = GetJSONFromAnyFile("./resources/application.yaml", JSON)
)

func TestConfigDefault(t *testing.T)  {
	app := CreateApplication(new(ConfigTestDefault))
	app.Run()
	t.Logf("conf:\n %s", app.ConfigFile)
	assert.Equal(t, JSONConf, app.ConfigFile)
}

func TestConfigTypeJSON(t *testing.T)  {
	app := CreateApplication(new(ConfigTestTypeJSON))
	app.Run()
	t.Logf("conf:\n %s", app.ConfigFile)
	assert.Equal(t, YAMLConf, app.ConfigFile)
}

func TestConfigSuffixYAML(t *testing.T)  {
	app := CreateApplication(new(ConfigTestSuffixYAML))
	app.Run()
	t.Logf("conf:\n %s", app.ConfigFile)
	assert.Equal(t, JSONFromYAML, app.ConfigFile)
}

func TestConfigTypeYAML(t *testing.T)  {
	app := CreateApplication(new(ConfigTestTypeYAML))
	app.Run()
	t.Logf("conf:\n %s", app.ConfigFile)
	assert.Equal(t, JSONFromYAML, app.ConfigFile)
}

