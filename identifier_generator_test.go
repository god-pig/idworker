package idworker

import (
	"testing"
)

func TestNewIdentifierGenerator(t *testing.T) {
	generate := NewIdentifierGenerator()
	t.Log(generate.NextId())
	t.Log(generate.NextUUID())
}
