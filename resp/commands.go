package resp

var Command_store = map[string]func([]Client_input) Client_input{
	"PING": command_ping,
}

func command_ping(args []Client_input) Client_input {
	if len(args) == 0 {
		return Client_input{Tipe: "string", Str: "PONG"}
	}
	return Client_input{Tipe: "string", Str: args[0].Bulk}
}
