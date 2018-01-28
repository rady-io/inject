package test_orm

import (
	"rady"
	"testing"
)

func TestOrm(t *testing.T) {
	rady.CreateApplication(new(OrmRoot)).RunTest(t, new(OrmTest))
}
