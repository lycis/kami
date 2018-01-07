function $create() {}

function process_input(input) {
    send_user_event(this.GetProp("$uuid"), input);
}