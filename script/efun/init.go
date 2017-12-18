package efun

import "gitlab.com/lycis/kami/script"

func init() {
	script.ExposeFunction("include", CreateIncludeEfun)
	script.ExposeFunction("call_other", createCallOther)
	script.ExposeFunction("grant_privilege", create_grant_privilege)
	script.ExposeFunction("log", create_log)
}
