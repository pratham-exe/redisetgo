package resp

import "sync"

var Command_store = map[string]func([]Client_input) Client_input{
	"PING": command_ping,
	"SET":  command_set,
	"GET":  command_get,
}

func command_ping(args []Client_input) Client_input {
	if len(args) == 0 {
		return Client_input{Tipe: "string", Str: "PONG"}
	} else if len(args) == 1 {
		return Client_input{Tipe: "string", Str: args[0].Bulk}
	}
	return Client_input{Tipe: "error", Str: "I can take max 1 argument."}
}

var (
	set_command_hashmap = map[string]string{}
	set_command_mutex   = sync.RWMutex{}
)

func command_set(args []Client_input) Client_input {
	if len(args) != 2 {
		return Client_input{Tipe: "error", Str: "I only take 2 arguments."}
	}

	set_comm_key := args[0].Bulk
	set_comm_value := args[1].Bulk

	set_command_mutex.Lock()
	set_command_hashmap[set_comm_key] = set_comm_value
	set_command_mutex.Unlock()

	return Client_input{Tipe: "string", Str: "OK"}
}

func command_get(args []Client_input) Client_input {
	if len(args) != 1 {
		return Client_input{Tipe: "error", Str: "I only take 1 argument."}
	}

	get_comm_key := args[0].Bulk

	set_command_mutex.RLock()
	get_comm_value, ok := set_command_hashmap[get_comm_key]
	set_command_mutex.RUnlock()

	if !ok {
		return Client_input{Tipe: "nill"}
	}

	return Client_input{Tipe: "bulk", Bulk: get_comm_value}
}
