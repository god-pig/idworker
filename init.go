package idworker

var (
	defaultIdentifierGenerator = NewIdentifierGenerator()
)

func NextId() (string, error) {
	return defaultIdentifierGenerator.NextId()
}

func NextUUID() string {
	return defaultIdentifierGenerator.NextUUID()
}
