package util

import (
	"fmt"
	"runtime"
)

var UserAgent = fmt.Sprintf("%s/%s golang/%v", Name, Version, runtime.Version())
