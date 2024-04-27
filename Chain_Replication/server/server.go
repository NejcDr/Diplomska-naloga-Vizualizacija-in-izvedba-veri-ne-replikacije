package server

import(
	"fmt"
	"time"

	"Chain_Replication/storage"
)

func Put(id int, s *storage.Storage, input chan storage.Value, output chan storage.Value, command chan storage.Command) {
	//fmt.Printf("%d PUT\n", id)
	/*
	for value := range input {
		fmt.Println("PUT")
		if cmd, err := s.Put(value, id); err == nil {
			output <- value
			command <- cmd
		}
	}
	*/
	for {
		select {
		case value, ok := <-input:
			if ok {
				//fmt.Printf("%d read from input head\n", id)
				//fmt.Printf("%d send to output head\n", id)
				if val, cmd, err := s.Put(value, id); err == nil {
					output <- val
					//fmt.Printf("%d send to output head\n", id)
					command <- cmd
				}
			}
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func Commit(id int, s *storage.Storage, input chan storage.Value, output chan storage.Value, command chan storage.Command) {
	//fmt.Printf("%d COMMIT\n", id)
	/*
	for value := range input {
		fmt.Println("COMMIT")
		if cmd, err := s.Commit(value, id); err == nil {
			output <- value
			command <- cmd
		}
	}
	*/
	for {
		select {
		case value, ok := <-input:
			if ok {
				//fmt.Printf("%d read from input tail\n", id)
				//fmt.Printf("%d send to output tail\n", id)
				if val, cmd, err := s.Commit(value, id); err == nil {
					output <- val
					//fmt.Printf("%d send to output tail\n", id)
					command <- cmd
				}
			}
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func Get(id int, s *storage.Storage, input chan storage.Command, output chan storage.Command) {
	//fmt.Printf("%d GET\n", id)
	/*
	for command := range input {
		fmt.Println("GET")
		return_command := s.Get(command, id)
		output <- return_command
	}
	*/
	for {
		select {
		case command, ok := <-input:
			if ok {
				//fmt.Printf("%d read from input get\n", id)
				return_command := s.Get(command, id)
				output <- return_command
				//fmt.Printf("%d send to output get\n", id)
			}	
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}	
}

func Server(id int, n_servers int,
	input_head_chan chan storage.Value, output_head_chan chan storage.Value, 
	input_tail_chan chan storage.Value, output_tail_chan chan storage.Value,
	input_read_chan chan storage.Command, output_read_chan chan storage.Command,
	command_output_chan chan storage.Command) {

	fmt.Printf("Servers %d start\n", id)

	s := storage.NewStorage()

	go Put(id, s, input_head_chan, output_head_chan, command_output_chan)
	go Commit(id, s, input_tail_chan, output_tail_chan, command_output_chan)
	go Get(id, s, input_read_chan, output_read_chan)

	select {}
	
	/*
	if id == (n_servers - 1) {
		for {
			select {
			case value, ok := <-input_head_chan:
				if ok {
					if cmd, err := s.Put(value, id); err == nil {
						command_output_chan <- cmd
						if cmd, err := s.Commit(value, id); err == nil {
							output_tail_chan <- value
							command_output_chan <- cmd
						}
					}
				}
			case command, ok := <-input_read_chan:
				if ok {
					return_command := s.Get(command, id)
					output_read_chan <- return_command
				}
			}
		}

	} else {
		for {
			select {
			case value, ok := <-input_head_chan:
				if ok {
					if cmd, err := s.Put(value, id); err == nil {
						output_head_chan <- value
						command_output_chan <- cmd
					}
				}
			case value, ok := <-input_tail_chan:
				if ok {
					if cmd, err := s.Commit(value, id); err == nil {
						output_tail_chan <- value
						command_output_chan <- cmd
					}
				}
			case command, ok := <-input_read_chan:
				if ok {
					return_command := s.Get(command, id)
					output_read_chan <- return_command
				}
			}
		}
	}
	*/
}