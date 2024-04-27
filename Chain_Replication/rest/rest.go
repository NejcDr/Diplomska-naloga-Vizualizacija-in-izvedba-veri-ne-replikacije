package rest

import(
	"fmt"
	"time"
	"Chain_Replication/storage"
)

func WriteOutGet(id int, output chan storage.Command) {
	/*
	for cmd := range output {
		fmt.Printf("%s\n", cmd.Command)
		for _, arg := range cmd.Arguments {
			fmt.Printf("%s\n", arg)
		}
		fmt.Println()
	}
	*/
	for {
		select {
		case cmd, ok := <-output:
			if ok {
				fmt.Printf("%s\n", cmd.Command)
				for i := 0; i < len(cmd.Arguments); i++ {
					fmt.Printf("%s\n", cmd.Arguments[i])
				}
				fmt.Println()
			}
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func WriteOutCommand(id int, output chan storage.Command) {
	/*
	for cmd := range output {
		fmt.Printf("%s\n", cmd.Command)
		for _, arg := range cmd.Arguments {
			fmt.Printf("%s\n", arg)
		}
		fmt.Println()
	}
	*/
	for {
		select {
		case cmd, ok := <-output:
			if ok {
				fmt.Printf("%s\n", cmd.Command)
				for i := 0; i < len(cmd.Arguments); i++ {
					fmt.Printf("%s\n", cmd.Arguments[i])
				}
				fmt.Println()
			}
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func Rest(num_of_servers int,
	head_input_chan chan storage.Value, head_output_chan chan storage.Value,
	get_input_chans []chan storage.Command, get_output_chans []chan storage.Command,
	command_chans []chan storage.Command) {

	fmt.Println("Rest start")

	value1 := storage.Value{Key: "Jabolka", Value: "10"}
	value2 := storage.Value{Key: "Hruska", Value: "5"}
	value3 := storage.Value{Key: "Jabolka", Value: "20"}

	readAll := storage.Command{Command: "readAll", Arguments: []string{}}
	read := storage.Command{Command: "read", Arguments: []string{"Jabolka"}}

	for i := 0; i < num_of_servers; i++ {
		go WriteOutGet(i, get_output_chans[i])
		go WriteOutCommand(i, command_chans[i])
	}

	time.Sleep(1 * time.Second)
	head_input_chan <- value1
	<-head_output_chan
	fmt.Println("WRITE_DONE!")

	time.Sleep(1 * time.Second)
	get_input_chans[num_of_servers - 1] <- readAll

	time.Sleep(1 * time.Second)
	head_input_chan <- value2
	<-head_output_chan
	fmt.Println("WRITE_DONE!")

	time.Sleep(1 * time.Second)
	get_input_chans[num_of_servers - 1] <- readAll

	time.Sleep(1 * time.Second)
	head_input_chan <- value3
	<-head_output_chan
	fmt.Println("WRITE_DONE!")

	time.Sleep(1 * time.Second)
	get_input_chans[num_of_servers - 1] <- readAll

	time.Sleep(1 * time.Second)
	get_input_chans[num_of_servers - 1] <- read
}