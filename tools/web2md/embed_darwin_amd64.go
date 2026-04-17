//go:build darwin && amd64

package main

import _ "embed"

//go:embed defuddle/defuddle-darwin-x64
var defuddleBinary []byte
