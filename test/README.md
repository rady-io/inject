## Test framework

Add RunTest(interface{}, *testing.T) for Application

Then rady can integrate go test

```go
// someapp_test.go
func TestSomeApp(t *testing.T){
    rady.CreateApplication(new(Root)).WithTest(t).RunTest(new(SomeAppTest))
}
```

However, what does SomeAppTest look like?

```go
type SomeAppTest struct {
    rady.Test
    client *test.RadyClient
    *OtherTest
    SomeValue *string `value:"rady.app.some_key"`
}

func (r *SomeAppTest) TestSomeValue(t *testing.T) {
    assert.Equel(t, *r.SomeValue, "some string")
}

func (r *SomeAppTest) TestApiSome(t *testing.T) {
    response := client.Get("/api/some").AssertOk()
    response.AssertHeader("content-type", "application/json").AssertBody([]bytes(`
    {
        "msg": "ok"
        "data": {
            "name": "xixi"
        }
    }
    `))
    response.AssertJson("data.name", "xixi")
    
    body := response.Return()
}
```

Then run 

```bash
go test
```


