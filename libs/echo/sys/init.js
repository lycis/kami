log("INFO", "****************************************************************\n* Running this library will echo every input back to the user. *\n****************************************************************\n")

// new user function
set_driver_hook(2, newUserToken);
// user input handling
set_driver_hook(3, processUserInput);
// user logs off
set_driver_hook(4, invalidateToken);

// enable REST networking subsystem
var enable = enable_subsystem(0);
if(enable !== true) {
	log("FATAL", "failed to enable REST: " + enable);
} else {
	log("INFO", "REST interface started: "+ enable);
}

// master functions

function newUserToken() {
	log("DEBUG", "new user token was requested. spawning entity");
	shell = spawn("/sys/user.js");
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
    shell = get_entity_by_id(token);
    if(shell === undefined) {
        log("ERROR", "user input provided for invalid token: " + token);
        return false;
    }

    if(shell.GetProp("$path").substring(0, 13) !== "/sys/user.js") {
        log("ERROR", "user input provided for non-shell entity");
        return false;
    }

    call_other(token, "process_input", input);
    return true;
}