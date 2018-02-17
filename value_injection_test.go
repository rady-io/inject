package rady

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
	"time"
)

const (
	FormatTime   = "2006-01-02 15:04:05"
	ReloadConfig = `
rady:
  mysql:
    host: localhost
    utf-8: false
  redis:
    host: 127.0.0.1
    port: 6937
  jwt:
    start: 2018-1-29 00:00:00
  server:
    ports:
      - 8080
      - 443
    ready:
      - false
      - false
    starts:
      - 2018-1-29 00:00:00
      - 2018-1-30 00:00:00
`
)

var (
	TrueTime, _         = time.Parse(FormatTime, "2018-1-30 00:00:00")
	ReloadedTrueTime, _ = time.Parse(FormatTime, "2018-1-29 00:00:00")
)

type (
	ValueInjectRoot struct {
		CONF `path:"./resources/application.yaml"`
	}

	ValueInjectTest struct {
		Testing
		*Application
		RedisPortInt         int64                    `value:"rady.redis.port"`
		RedisPortIntPtr      *int64                   `value:"rady.redis.port"`
		RedisPortUint        uint64                   `value:"rady.redis.port"`
		RedisPortUintPtr     *uint64                  `value:"rady.redis.port"`
		RedisPortFloat       float64                  `value:"rady.redis.port"`
		RedisPortFloatPtr    *float64                 `value:"rady.redis.port"`
		RedisPortStr         string                   `value:"rady.redis.port"`
		RedisPortStrPtr      *string                  `value:"rady.redis.port"`
		MysqlUtf8            bool                     `value:"rady.mysql.utf-8"`
		MysqlUtf8Ptr         *bool                    `value:"rady.mysql.utf-8"`
		JWTStartTime         time.Time                `value:"rady.jwt.start"`
		JWTStartTimePtr      *time.Time               `value:"rady.jwt.start"`
		RadyConfig           map[string]gjson.Result  `value:"rady"`
		RadyConfigPtr        *map[string]gjson.Result `value:"rady"`
		ServePortsResults    []gjson.Result           `value:"rady.server.ports"`
		ServePortsInt        []int64                  `value:"rady.server.ports"`
		ServePortsUint       []uint64                 `value:"rady.server.ports"`
		ServePortsFloat      []float64                `value:"rady.server.ports"`
		ServePortsStr        []string                 `value:"rady.server.ports"`
		IfPortsReady         []bool                   `value:"rady.server.ready"`
		ServerPowerTimes     []time.Time              `value:"rady.server.starts"`
		ServePortsResultsPtr *[]gjson.Result          `value:"rady.server.ports"`
		ServePortsIntPtr     *[]int64                 `value:"rady.server.ports"`
		ServePortsUintPtr    *[]uint64                `value:"rady.server.ports"`
		ServePortsFloatPtr   *[]float64               `value:"rady.server.ports"`
		ServePortsStrPtr     *[]string                `value:"rady.server.ports"`
		IfPortsReadyPtr      *[]bool                  `value:"rady.server.ready"`
		ServerPowerTimesPtr  *[]time.Time             `value:"rady.server.starts"`
	}
)

func (v *ValueInjectTest) TestElemOrPtrInject(t *testing.T) {
	assert.Equal(t, int64(6937), v.RedisPortInt)
	assert.Equal(t, int64(6937), *v.RedisPortIntPtr)
	assert.Equal(t, uint64(6937), v.RedisPortUint)
	assert.Equal(t, uint64(6937), *v.RedisPortUintPtr)
	assert.Equal(t, float64(6937), v.RedisPortFloat)
	assert.Equal(t, float64(6937), *v.RedisPortFloatPtr)
	assert.Equal(t, "6937", v.RedisPortStr)
	assert.Equal(t, "6937", *v.RedisPortStrPtr)
	assert.True(t, v.MysqlUtf8)
	assert.True(t, *v.MysqlUtf8Ptr)
	assert.Equal(t, TrueTime, v.JWTStartTime)
	assert.Equal(t, TrueTime, *v.JWTStartTimePtr)
}

func (v *ValueInjectTest) TestReloadValues(t *testing.T) {
	v.WriteConfigFile(ReloadConfig)
	v.ReloadValues()
	assert.False(t, *v.MysqlUtf8Ptr)
	assert.Equal(t, ReloadedTrueTime, *v.JWTStartTimePtr)
	v.WriteConfigFile(InitConfig)
	v.ReloadValues()
}

func (v *ValueInjectTest) TestMapInject(t *testing.T) {
	for Key, Value := range v.RadyConfig {
		switch Key {
		case "mysql":
			for key, value := range Value.Map() {
				switch key {
				case "host":
					assert.Equal(t, "localhost", value.String())
				case "utf-8":
					assert.True(t, value.Bool())
				default:
				}
			}
		case "redis":
			for key, value := range Value.Map() {
				switch key {
				case "host":
					assert.Equal(t, "127.0.0.1", value.String())
				case "utf-8":
					assert.Equal(t, int64(6937), value.Int())
				default:
				}
			}
		case "jwt":
			assert.Equal(t, TrueTime, Value.Time())
		default:

		}
	}
}

func (v *ValueInjectTest) TestReloadedMap(t *testing.T) {
	v.WriteConfigFile(ReloadConfig)
	v.ReloadValues()
	for Key, Value := range *v.RadyConfigPtr {
		switch Key {
		case "mysql":
			for key, value := range Value.Map() {
				switch key {
				case "host":
					assert.Equal(t, "localhost", value.String())
				case "utf-8":
					assert.False(t, value.Bool())
				default:
				}
			}
		case "redis":
			for key, value := range Value.Map() {
				switch key {
				case "host":
					assert.Equal(t, "127.0.0.1", value.String())
				case "utf-8":
					assert.Equal(t, int64(6937), value.Int())
				default:
				}
			}
		case "jwt":
			assert.Equal(t, ReloadedTrueTime, Value.Time())
		default:

		}
	}
	v.WriteConfigFile(InitConfig)
	v.ReloadValues()
}

func (v *ValueInjectTest) TestArrays(t *testing.T) {
	assert.Equal(t, int64(80), v.ServePortsInt[0])
	assert.Equal(t, int64(443), v.ServePortsInt[1])
}

func (v *ValueInjectTest) TestReloadedArrays(t *testing.T) {
	v.WriteConfigFile(ReloadConfig)
	v.ReloadValues()
	assert.Equal(t, int64(8080), (*v.ServePortsIntPtr)[0])
	assert.Equal(t, int64(443), (*v.ServePortsIntPtr)[1])
	v.WriteConfigFile(InitConfig)
	v.ReloadValues()
}

func TestValueInjection(t *testing.T) {
	CreateTest(new(ValueInjectRoot)).AddTest(new(ValueInjectTest)).Test(t)
}
