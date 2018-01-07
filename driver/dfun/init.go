package dfun

import "gitlab.com/lycis/kami/script"

func init() {
	script.ExposeFunction("spawn", create_dfun_spawn)
	script.ExposeFunction("get_entity_by_id", create_dfun_get_entity_by_id)
	script.ExposeFunction("set_driver_hook", create_dfun_set_driver_hook)
	script.ExposeFunction("shutdown", create_shutdown)
	script.ExposeFunction("enable_subsystem", create_enable_subsystem)
	script.ExposeFunction("disable_subsystem", create_disable_subsystem)
	script.ExposeFunction("destroy", createDfunDestroy)
	script.ExposeFunction("send_user_event", createDfunSendUserEvent)
}
