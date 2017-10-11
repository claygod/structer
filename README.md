# Structer

Storage for structures (embedded), with the ability to search by field structure.
This library is unsafe.

[![API documentation](https://godoc.org/github.com/claygod/structer?status.svg)](https://godoc.org/github.com/claygod/structer)
[![Go Report Card](https://goreportcard.com/badge/github.com/claygod/structer)](https://goreportcard.com/report/github.com/claygod/structer)

# Usage

An example of using the Structer:
```go
package main

import (
	"github.com/claygod/structer"
)

// Main
func main() {
	s := structer.New()
	...
}
```


# Perfomance

- Add-4              500000	      4452 ns/op
- AddParallel-4      500000	      3104 ns/op
- Find-4           	 2000000	  909 ns/op
- FindParallel-4   	 3000000	  494 ns/op

# API

Methods:
-  *New* - create a new storage for structures


