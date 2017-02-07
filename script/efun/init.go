package efun

import "github.com/lycis/kami/script"

func init() {
	script.ExposeFunction("include", CreateIncludeEfun)
	script.ExposeFunction("call_other", createCallOther)
}
