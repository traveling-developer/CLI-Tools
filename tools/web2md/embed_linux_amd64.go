//go:build linux && amd64

package main

import _ "embed"

//go:embed defuddle/defuddle-linux-x64
var defuddleBinary []byte
