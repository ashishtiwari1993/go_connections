# go_connections
A central library that contains different functions for different connections like MySQL, Redis, etc. I will keep updating this file. Any Contribution would appreciate :)

# Usage

## git clone
`git clone https://github.com/ashishtiwari1993/go_connections.git`

## Set configs in `configs.json`
`vim go_connections/configs.json`

## Sample `main.go`

```go
package main

import c "YourProjectPath/go_connections"

var configs map[string]string

func main() {

	// Access all value of configs.json
	configs = c.Configs
	
	m := c.ConnectMysql()
	m.Close()
	
	r := c.ConnectRedis()
	r.Close()
}
```
