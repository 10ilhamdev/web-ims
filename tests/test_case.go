package tests

import (
	"github.com/goravel/framework/testing"

	"ims/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
