//go:build linux && arm64

package main

import _ "embed"

//go:embed defuddle/defuddle-linux-arm64
var defuddleBinary []byte
