log("INFO", "Booting kernel...");
spawn("/entities/dummy.js");
set_driver_hook(1, on_wr);

function on_wr() {
	log("info", "world is running");
	shutdown("expected end");
}
