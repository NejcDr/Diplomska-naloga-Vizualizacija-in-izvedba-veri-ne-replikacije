package main

import(
	"fmt"
	"flag"
	"time"
	"sync"

	"Chain_Replication/storage"
	"Chain_Replication/server"
	"Chain_Replication/rest"
)

func main() {
	var nPtr int
	flag.IntVar(&nPtr, "n", 5, "total number of replication servers.")
	flag.Parse()

	fmt.Printf("Start %d servers!\n", nPtr)

	//nPtr := flag.Int("n", 5, "total number of replication servers.")

	put_chans := make([]chan storage.Value, nPtr + 1)
	commit_chans := make([]chan storage.Value, nPtr + 1)

	get_input_chans := make([]chan storage.Command, nPtr)
	get_output_chans := make([]chan storage.Command, nPtr)

	command_chans := make([]chan storage.Command, nPtr)

	for i := 0; i < nPtr; i++ {
		put_chans[i] = make(chan storage.Value, 10)
		commit_chans[i] = make(chan storage.Value, 10)
		get_input_chans[i] = make(chan storage.Command, 10)
		get_output_chans[i] = make(chan storage.Command, 10)
		command_chans[i] = make(chan storage.Command, 10)
	}

	tail_chan := make(chan storage.Value, 10)
	put_chans[nPtr] = tail_chan;
	commit_chans[nPtr] = tail_chan;

	//fmt.Println("Servers start")

	var wg sync.WaitGroup

	for i := 0; i < nPtr; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			server.Server(i, nPtr, put_chans[i], put_chans[i+1], commit_chans[i+1], commit_chans[i], get_input_chans[i], get_output_chans[i], command_chans[i])
		}(i)
		//go server.Server(i, put_chans[i], put_chans[i+1], commit_chans[i+1], commit_chans[i], get_input_chans[i], get_output_chans[i], command_chans[i])
	}

	time.Sleep(3 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		rest.Rest(nPtr, put_chans[0], commit_chans[0], get_input_chans, get_output_chans, command_chans)
	}()

	wg.Wait()

	select{}

	fmt.Println("END")
	//go rest.Rest(nPtr, put_chans[0], commit_chans[0], get_input_chans, get_output_chans, command_chans)
}