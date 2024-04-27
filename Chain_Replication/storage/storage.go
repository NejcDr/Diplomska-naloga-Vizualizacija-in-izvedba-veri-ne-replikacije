package storage

import(
	"errors"
	"sync"
	"fmt"
)

type Command struct {
	Command string
	Arguments []string
}

type Value struct {
	Key string 
	Value string
	Version int
	Commited bool
}

type Storage struct {
	dict map[string]Value
	mu sync.RWMutex
}

var ErrorNotFound = errors.New("not found")

func NewStorage() *Storage{
	dict := make(map[string]Value)
	return &Storage {
		dict: dict,
	}
}

func (s *Storage) Get(command Command, id int) Command {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var args []string
	var return_command string
	
	if command.Command == "readAll" {
		for k, v := range s.dict {
			if v.Commited {
				str := fmt.Sprintf("Key: %s | Value: %s | Version: %d", k, v.Value, v.Version)
				args = append(args, str)
			}
		}
		return_command = fmt.Sprintf("returnedAll_%d", id)
	} else if command.Command == "read" {
		k := command.Arguments[0]
		if v, ok := s.dict[k]; ok {
			if v.Commited {
				str := fmt.Sprintf("Key: %s | Value: %s | Version: %d", k, v.Value, v.Version)
				args = append(args, str)
				return_command = fmt.Sprintf("returned_%d", id)
			}
		} 
	}

	return Command{Command: return_command, Arguments: args}
}

func (s *Storage) Put(value Value, id int) (Value, Command, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var version int

	if v, ok := s.dict[value.Key]; ok {
		version = v.Version
	} else {
		version = 0
	}

	value.Commited = false
	value.Version = version + 1
	s.dict[value.Key] = value

	var args []string
	str := fmt.Sprintf("Key: %s | Value: %s | Version: %d", value.Key, value.Value, value.Version)
	args = append(args, str)

	return Value{Key: value.Key, Value: value.Value, Version: value.Version}, 
		Command{Command: fmt.Sprintf("put_%d", id), Arguments: args}, 
		nil
}

func (s *Storage) Commit(value Value, id int) (Value, Command, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := s.dict[value.Key]; ok && v.Version <= value.Version {
		value.Commited = true
		s.dict[value.Key] = value

		var args []string
		str := fmt.Sprintf("Key: %s | Value: %s | Version: %d", value.Key, value.Value, value.Version)
		args = append(args, str)

		return Value{Key: value.Key, Value: value.Value, Version: value.Version}, 
			Command{Command: fmt.Sprintf("commited_%d", id), Arguments: args}, 
			nil
	}

	return Value{}, Command{}, ErrorNotFound
}