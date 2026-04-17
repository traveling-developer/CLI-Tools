//go:build windows && amd64

package main

import _ "embed"

//go:embed defuddle/defuddle-windows-x64.exe
var defuddleBinary []byte
