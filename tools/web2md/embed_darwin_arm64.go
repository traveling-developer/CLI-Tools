//go:build darwin && arm64

package main

import _ "embed"

//go:embed defuddle/defuddle-darwin-arm64
var defuddleBinary []byte
