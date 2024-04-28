package server

import(
	"fmt"
	"time"
	"sync"

	"Chain_Replication/storage"
)

var NUM_SERVERS int

var PUT_CHANS []chan storage.Value
var COMMIT_CHANS []chan storage.Value

var GET_INPUT_CHANS []chan storage.Command
var GET_OUTPUT_CHANS []chan storage.Command

var COMMAND_CHANS []chan storage.Command

func server(id int) {

	fmt.Printf("Servers %d start\n", id)
	s := storage.NewStorage()

	for {
		select {
		case value, ok := <-PUT_CHANS[id]:
			if ok {
				if val, cmd, err := s.Put(value, id); err == nil {
					PUT_CHANS[id+1] <- val
					COMMAND_CHANS[id] <- cmd
				}
			}
		case value, ok := <-COMMIT_CHANS[id+1]:
			if ok {
				if val, cmd, err := s.Commit(value, id); err == nil {
					COMMIT_CHANS[id] <- val
					COMMAND_CHANS[id] <- cmd
				}
			}
		case command, ok := <-GET_INPUT_CHANS[id]:
			if ok {
				cmd := s.Get(command, id)
				GET_OUTPUT_CHANS[id] <- cmd
			}	
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func ServersInit(n_servers int,
	put_chans []chan storage.Value, commit_chans []chan storage.Value,
	get_input_chans []chan storage.Command, get_output_chans []chan storage.Command,
	command_chans []chan storage.Command) {

	NUM_SERVERS = n_servers
	PUT_CHANS = put_chans
	COMMIT_CHANS = commit_chans
	GET_INPUT_CHANS = get_input_chans
	GET_OUTPUT_CHANS = get_output_chans
	COMMAND_CHANS = command_chans

	var wg sync.WaitGroup

	for i := 0; i < NUM_SERVERS; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			server(i)
		}(i)
	}

	wg.Wait()
}