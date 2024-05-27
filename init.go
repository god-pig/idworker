package idworker

var (
	defaultIdentifierGenerator = NewIdentifierGenerator()
)

func NextId() (int64, error) {
	return defaultIdentifierGenerator.NextId()
}

func NextUUID() string {
	return defaultIdentifierGenerator.NextUUID()
}
