package resp

import (
	"sync"
)

var Command_store = map[string]func([]Client_input) Client_input{
	"PING":      command_ping,
	"SET":       command_set,
	"GET":       command_get,
	"HSET":      command_hset,
	"HGET":      command_hget,
	"HGETALL":   command_hgetall,
	"DEL":       command_del,
	"REDISHELP": command_help,
}

func command_ping(args []Client_input) Client_input {
	if len(args) == 0 {
		return Client_input{Tipe: "string", Str: "PONG"}
	} else if len(args) == 1 {
		return Client_input{Tipe: "bulk", Bulk: args[0].Bulk}
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

var (
	hset_command_hashmap = map[string]map[string]string{}
	hset_command_mutex   = sync.RWMutex{}
)

func command_hset(args []Client_input) Client_input {
	if len(args) != 3 {
		return Client_input{Tipe: "error", Str: "I only take 3 arguments."}
	}

	hset_comm_hash := args[0].Bulk
	hset_comm_key := args[1].Bulk
	hset_comm_value := args[2].Bulk

	hset_command_mutex.Lock()

	_, ok := hset_command_hashmap[hset_comm_hash]
	if !ok {
		hset_command_hashmap[hset_comm_hash] = map[string]string{}
	}
	hset_command_hashmap[hset_comm_hash][hset_comm_key] = hset_comm_value

	hset_command_mutex.Unlock()

	return Client_input{Tipe: "string", Str: "OK"}
}

func command_hget(args []Client_input) Client_input {
	if len(args) != 2 {
		return Client_input{Tipe: "error", Str: "I only take 2 arguments."}
	}

	hget_comm_hash := args[0].Bulk
	hget_comm_key := args[1].Bulk

	hset_command_mutex.RLock()
	hget_comm_value, ok := hset_command_hashmap[hget_comm_hash][hget_comm_key]
	hset_command_mutex.RUnlock()

	if !ok {
		return Client_input{Tipe: "nill"}
	}

	return Client_input{Tipe: "bulk", Bulk: hget_comm_value}
}

func command_hgetall(args []Client_input) Client_input {
	if len(args) != 1 {
		return Client_input{Tipe: "error", Str: "I only take 1 argument."}
	}

	hget_comm_hash := args[0].Bulk

	hset_command_mutex.RLock()

	hget_comm_key, ok := hset_command_hashmap[hget_comm_hash]
	hget_comm_value := []Client_input{}
	for k, v := range hget_comm_key {
		hget_comm_value = append(hget_comm_value, Client_input{Tipe: "bulk", Bulk: k})
		hget_comm_value = append(hget_comm_value, Client_input{Tipe: "bulk", Bulk: v})
	}

	hset_command_mutex.RUnlock()

	if !ok {
		return Client_input{Tipe: "nill"}
	}

	return Client_input{Tipe: "array", Array: hget_comm_value}
}

func command_del(args []Client_input) Client_input {
	if len(args) != 1 {
		return Client_input{Tipe: "error", Str: "I only take 1 argument."}
	}

	del_key := args[0].Bulk
	delete(set_command_hashmap, del_key)

	return Client_input{Tipe: "string", Str: "OK"}
}

func command_help(args []Client_input) Client_input {
	if len(args) != 0 {
		return Client_input{Tipe: "error", Str: "I take 0 arguments."}
	}

	return_value := []Client_input{}
	return_value = append(return_value, Client_input{Tipe: "string", Str: "PING (or) PING <value> (0 or 1 argument)"})
	return_value = append(return_value, Client_input{Tipe: "string", Str: "SET <key> <value> (2 arguments)"})
	return_value = append(return_value, Client_input{Tipe: "string", Str: "GET <key> (1 argument)"})
	return_value = append(return_value, Client_input{Tipe: "string", Str: "HSET <hash> <key> <value> (3 arguments)"})
	return_value = append(return_value, Client_input{Tipe: "string", Str: "HGET <hash> <key> (2 arguments)"})
	return_value = append(return_value, Client_input{Tipe: "string", Str: "HGETALL <hash> (1 argument)"})
	return_value = append(return_value, Client_input{Tipe: "string", Str: "DEL <key> (1 argument)"})

	return Client_input{Tipe: "array", Array: return_value}
}
