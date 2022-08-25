package idworker

import (
	"testing"
)

func TestNewSnowFlake(t *testing.T) {
	snow := NewSnowflake("")
	t.Log(snow.NextId())
}
