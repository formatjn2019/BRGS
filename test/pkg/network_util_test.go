package test

import (
	"BRGS/pkg/utils"
	"testing"
)

func TestGetIp(t *testing.T) {
	println(utils.GotLocalIP())
}
