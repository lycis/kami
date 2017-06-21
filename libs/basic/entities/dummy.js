function $create() {
	this.SetProp("created", true);
}

function $tick() {
	log("debug", "tick tock");
}

function $onShutdown(reason) {
	log("info", "notified that driver will shutdown due to: " + reason);
}
