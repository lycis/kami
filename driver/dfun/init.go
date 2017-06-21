package dfun

import "github.com/lycis/kami/script"

func init() {
	script.ExposeFunction("spawn", create_dfun_spawn)
	script.ExposeFunction("get_entity_by_id", create_dfun_get_entity_by_id)
	script.ExposeFunction("set_driver_hook", create_dfun_set_driver_hook)
	script.ExposeFunction("shutdown", create_shutdown)
}
