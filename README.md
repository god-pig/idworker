# idworker

分布式高效有序ID

## Usage

```go
package main

import (
	"fmt"
	"github.com/god-pig/idworker"
)

func main() {
	generate := NewIdentifierGenerator()
	fmt.Println(generate.NextId())
	// -> 1562674057830715394

	fmt.Println(generate.NextUUID())
	// -> f3cbb52b-c41e-4c88-b1d3-d64a27debaea
}
```