package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	InitConfig = `
rady:
  mysql:
    host: localhost
    utf-8: true
  redis:
    host: 127.0.0.1
    port: 6937
  jwt:
    start: 2018-1-30 00:00:00
  server:
    ports:
      - 80
      - 443
    ready:
      - true
      - false
    starts:
      - 2018-1-30 00:00:00
      - 2018-1-30 00:00:00
`
	NoRedis = `
rady:
  mysql:
    host: localhost
    utf-8: true
  jwt:
    start: 2018-1-30 00:00:00
  server:
    ports:
      - 80
      - 443
    ready:
      - true
      - false
    starts:
      - 2018-1-30 00:00:00
      - 2018-1-30 00:00:00
`
	ChangedRedis = `
rady:
  mysql:
    host: localhost
    utf-8: true
  redis:
    host: 127.0.0.1
    port: 1200
  jwt:
    start: 2018-1-30 00:00:00
  server:
    ports:
      - 80
      - 443
    ready:
      - true
      - false
    starts:
      - 2018-1-30 00:00:00
      - 2018-1-30 00:00:00
`
)

type RecallFactoryTest struct {
	Testing
	*RedisComponent
	*Application
}

func (R *RecallFactoryTest) TestRecallFactory(t *testing.T) {
	assert.Equal(t, "127.0.0.1", R.GetHost())
	assert.Equal(t, int64(6937), R.GetPort())
	R.WriteConfigFile(NoRedis)
	R.ReloadValues()
	assert.Equal(t, "localhost", R.GetHost())
	assert.Equal(t, int64(6666), R.GetPort())
	R.WriteConfigFile(ChangedRedis)
	R.ReloadValues()
	assert.Equal(t, "127.0.0.1", R.GetHost())
	assert.Equal(t, int64(1200), R.GetPort())
	R.WriteConfigFile(InitConfig)
}

func TestReloadValues(t *testing.T) {
	CreateTest(new(App)).AddTest(new(RecallFactoryTest)).Test(t)
}
