package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase16ToBase2(t *testing.T) {
	asst := assert.New(t)

	base2, _ := Base16ToBase2("06")
	asst.Equal("110", base2, "test Base16ToBase2() failed")
}
