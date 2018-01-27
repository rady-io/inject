package test_orm

import (
	"testing"
	"rady"
)

func TestOrm(t *testing.T)  {
	rady.CreateApplication(new(OrmRoot)).RunTest(t, new(OrmTest))
}
