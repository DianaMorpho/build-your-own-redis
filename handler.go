package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) != 0 {
		return Value{typ: "array", array: args}
	}
	return Value{typ: "string", str: "PONG"}
}

var SETs = map[string]string{}
var SETsMutex = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "string", str: "Invalid request"}
	}
	key := args[0].bulk
	value := args[1].bulk
	SETsMutex.Lock()
	SETs[key] = value
	SETsMutex.Unlock()
	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "string", str: "Invalid request"}
	}
	key := args[0].bulk
	SETsMutex.RLock()
	value, ok := SETs[key]
	SETsMutex.RUnlock()
	if !ok {
		return Value{typ: "string", str: "Key not found"}
	}
	return Value{typ: "string", str: value}
}

var HSETS = map[string]map[string]string{}
var HSETsMutex = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "string", str: "Invalid request"}
	}
	key := args[0].bulk
	field := args[1].bulk
	value := args[2].bulk
	HSETsMutex.Lock()
	if _, ok := HSETS[key]; !ok {
		HSETS[key] = map[string]string{}
	}
	HSETS[key][field] = value
	HSETsMutex.Unlock()
	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "string", str: "Invalid request"}
	}
	key := args[0].bulk
	field := args[1].bulk
	HSETsMutex.RLock()
	value, ok := HSETS[key][field]
	HSETsMutex.RUnlock()
	if !ok {
		return Value{typ: "string", str: "Field not found"}
	}
	return Value{typ: "string", str: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "string", str: "Invalid request"}
	}
	key := args[0].bulk
	HSETsMutex.RLock()
	fields := HSETS[key]
	HSETsMutex.RUnlock()
	var res []Value
	for field, value := range fields {
		res = append(res, Value{typ: "array", array: []Value{
			{typ: "string", str: field},
			{typ: "string", str: value},
		}})
	}
	return Value{typ: "array", array: res}
}
