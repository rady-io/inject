package rady

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/tidwall/gjson"
)

const (
	FormatTime = "2006-01-02 15:04:05"
)

var (
	TrueTime, _ = time.Parse(FormatTime, "2018-1-30 00:00:00")
)

type (
	ValueInjectRoot struct {
		CONF `path:"./resources/application.yaml"`
	}

	ValueInjectTest struct {
		Testing
		RedisPortInt      int64                   `value:"rady.redis.port"`
		RedisPortIntPtr   *int64                  `value:"rady.redis.port"`
		RedisPortUint     uint64                  `value:"rady.redis.port"`
		RedisPortUintPtr  *uint64                 `value:"rady.redis.port"`
		RedisPortFloat    float64                 `value:"rady.redis.port"`
		RedisPortFloatPtr *float64                `value:"rady.redis.port"`
		RedisPortStr      string                  `value:"rady.redis.port"`
		RedisPortStrPtr   *string                 `value:"rady.redis.port"`
		MysqlUtf8         bool                    `value:"rady.mysql.utf-8"`
		MysqlUtf8Ptr      *bool                   `value:"rady.mysql.utf-8"`
		JWTStartTime      time.Time               `value:"rady.jwt.start"`
		JWTStartTimePtr   *time.Time              `value:"rady.jwt.start"`
		RadyConfig        map[string]gjson.Result `value:"rady"`
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

func TestValueInjection(t *testing.T) {
	CreateApplication(new(ValueInjectRoot)).PrepareTest().AddTest(new(ValueInjectTest)).Test(t)
}
