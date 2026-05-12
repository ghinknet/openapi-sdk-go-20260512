package openapi

import (
	"fmt"
	"runtime"
)

const Endpoint = "https://api.gh.ink/v3.1"

var Version = [3]int{3, 0, 2}
var UserAgent = fmt.Sprintf(
	"GhinkOpenAPISDK-20260512-Go/%d.%d.%d (%s; %s)",
	Version[0], Version[1], Version[2],
	runtime.GOOS, runtime.GOARCH,
)
