package idworker

import "github.com/google/uuid"

type IdentifierGenerator interface {
	// NextId 生成id
	NextId() (int64, error)

	// NextUUID 生成uuid
	NextUUID() string
}

type DefaultIdentifierGenerator struct {
	snowflake *Snowflake
}

func (d *DefaultIdentifierGenerator) NextId() (int64, error) {
	return d.snowflake.NextId()
}

func (d *DefaultIdentifierGenerator) NextUUID() string {
	return uuid.NewString()
}

func NewIdentifierGenerator() IdentifierGenerator {
	return &DefaultIdentifierGenerator{
		snowflake: NewSnowflake(""),
	}
}
