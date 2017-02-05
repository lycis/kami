package efun

import "github.com/lycis/kami/script"

func init() {
	script.ExposeFunction("include", CreateIncludeEfun)
}
