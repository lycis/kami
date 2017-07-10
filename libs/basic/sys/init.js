log("INFO", "Booting kernel...");
spawn("/entities/dummy.js");

set_driver_hook(2, newUserToken);

var enable = enable_subsystem(0);
if(enable !== true) {
	log("FATAL", "failed to enable REST: " + enable);
} else {
	log("INFO", "REST interface started: "+ enable);
}

function newUserToken() {
	log("DEBUG", "new user token was requested. spawning entity");
	shell = spawn("/shell/user.js");
    log("INFO", "Spawned new user shell. uuid="+shell.GetProp("$uuid"));
	return shell.GetProp("$uuid")
}