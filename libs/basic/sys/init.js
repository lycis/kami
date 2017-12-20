log("INFO", "Booting kernel...");
var dummy = spawn("/entities/dummy.js");

set_driver_hook(2, newUserToken);
set_driver_hook(3, processUserInput);

var enable = enable_subsystem(0);
if(enable !== true) {
	log("FATAL", "failed to enable REST: " + enable);
} else {
	log("INFO", "REST interface started: "+ enable);
}

log("INFO", "destroying dummy with id '" + dummy.GetProp("$uuid")+"'")
destroy(dummy.GetProp("$uuid"));
log("INFO", "dummy destroyed")
if(dummy.IsActive()) {
    log("FATAL", "destroyed entity is still active");
}

function newUserToken() {
	log("DEBUG", "new user token was requested. spawning entity");
	shell = spawn("/shell/user.js");
    log("INFO", "Spawned new user shell. uuid="+shell.GetProp("$uuid"));
	return shell.GetProp("$uuid")
}

function invalidateToken(token) {
    log("DEBUG", "removing token");
    shell = get_entity_by_id(token);
    if(shell === undefined) {
        log("ERROR", "user with the given token is unknown");
        return false;
    }

    destroy(shell.GetProp("$uuid"));
    call_other(token, "SetProp", "disabled", true);
}

function processUserInput(token, input) {
    log("INFO", "processUserInput(" + token + ", " + input + ")");
    shell = get_entity_by_id(token);
    if(shell === undefined) {
        log("ERROR", "user input provided for invalid token: " + token);
        return false;
    }

    log("INFO", "checking shell...");
    if(shell.GetProp("$path").substring(0, 6) !== "/shell") {
        log("ERROR", "user input provided for non-shell entity");
        return false;
    }

    log("INFO", "forwarding user input. token = " + token + " input="+input);

    call_other(token, "process_input", input);

    return true;
}