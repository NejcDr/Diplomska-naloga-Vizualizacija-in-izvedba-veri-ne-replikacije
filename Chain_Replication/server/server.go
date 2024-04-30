package server

import(
	"fmt"
	"sync"

	"Chain_Replication/storage"
)

var NUM_SERVERS int

var PUT_CHANS []chan storage.Value
var COMMIT_CHANS []chan storage.Value

var GET_INPUT_CHANS []chan storage.Command
var GET_OUTPUT_CHANS []chan storage.Command

var COMMAND_CHANS []chan storage.Command

var TOKEN_CHANS []chan int

func server(id int) {

	fmt.Printf("Servers %d start\n", id)
	s := storage.NewStorage()
	
	for {
		select {
		case value, ok := <-PUT_CHANS[id]:
			if ok {
				_, tokenOK := <-TOKEN_CHANS[id]
				if tokenOK {
					if val, cmd, err := s.Put(value, id); err == nil {
						PUT_CHANS[id+1] <- val
						COMMAND_CHANS[id] <- cmd
					}
				}
			}
		case value, ok := <-COMMIT_CHANS[id+1]:
			if ok {
				_, tokenOK := <-TOKEN_CHANS[id]
				if tokenOK {
					if val, cmd, err := s.Commit(value, id); err == nil {
						COMMIT_CHANS[id] <- val
						COMMAND_CHANS[id] <- cmd
					}
				}
			}
		case command, ok := <-GET_INPUT_CHANS[id]:
			if ok {
				_, tokenOK := <-TOKEN_CHANS[id]
				if tokenOK {
					cmd := s.Get(command, id)
					_, tokenOK := <-TOKEN_CHANS[id]
					if tokenOK {
						GET_OUTPUT_CHANS[id] <- cmd
					}
				}
			}
		case _, tokenOK := <-TOKEN_CHANS[id]:
			if tokenOK {
				continue
			}
		}
	}
	
}

func ServersInit(n_servers int,
	put_chans []chan storage.Value, commit_chans []chan storage.Value,
	get_input_chans []chan storage.Command, get_output_chans []chan storage.Command,
	command_chans []chan storage.Command, token_chans []chan int) {

	NUM_SERVERS = n_servers
	PUT_CHANS = put_chans
	COMMIT_CHANS = commit_chans
	GET_INPUT_CHANS = get_input_chans
	GET_OUTPUT_CHANS = get_output_chans
	COMMAND_CHANS = command_chans
	TOKEN_CHANS = token_chans

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