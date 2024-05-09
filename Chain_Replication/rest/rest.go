package rest

import(
	"fmt"
	"time"
	"sync"
	"encoding/json"
    "log"
    "net/http"
	"strconv"

    "github.com/gorilla/mux"
	"Chain_Replication/storage"
)

var NUM_SERVERS int

var TOKEN_ID int

var INPUT_HEAD_CHAN chan storage.Value
var OUTPUT_HEAD_CHAN chan storage.Value

var GET_INPUT_CHANS []chan storage.Command
var GET_OUTPUT_CHANS []chan storage.Command

var COMMAND_CHANS []chan storage.Command

var TOKEN_CHANS []chan int

/*
func WriteOutGet(id int, output chan storage.Command) {
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
*/

func WriteOutCommand(id int) {
	for {
		select {
		case cmd, ok := <-COMMAND_CHANS[id]:
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

func readMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var readInCommand storage.Command
	_ = json.NewDecoder(r.Body).Decode(&readInCommand)

	server, err := strconv.Atoi(readInCommand.Arguments[0])
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}

	GET_INPUT_CHANS[server] <- readInCommand

	readOutCommand, ok := <-GET_OUTPUT_CHANS[server]
	if ok {
		json.NewEncoder(w).Encode(readOutCommand)
		return
	}

	json.NewEncoder(w).Encode(nil)
}

func insertMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inputValue storage.Value
	_ = json.NewDecoder(r.Body).Decode(&inputValue)

	INPUT_HEAD_CHAN <- inputValue

	outValue, ok := <-OUTPUT_HEAD_CHAN
	if ok {
		json.NewEncoder(w).Encode(outValue)
		return
	}

	json.NewEncoder(w).Encode(nil)
}


func clockMethod(w http.ResponseWriter, r *http.Request) {
	TOKEN_ID++

	for i := 0; i < NUM_SERVERS; i++ {
		TOKEN_CHANS[i] <- TOKEN_ID
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

func restReciver() {
	r := mux.NewRouter()

	r.HandleFunc("/read", readMethod).Methods("GET")
	r.HandleFunc("/insert", insertMethod).Methods("POST")
	r.HandleFunc("/clock", clockMethod).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func Rest(num_of_servers int,
	head_input_chan chan storage.Value, head_output_chan chan storage.Value,
	get_input_chans []chan storage.Command, get_output_chans []chan storage.Command,
	command_chans []chan storage.Command, token_chans []chan int) {

	fmt.Println("Rest start")

	NUM_SERVERS = num_of_servers
	TOKEN_ID = 0

	INPUT_HEAD_CHAN = head_input_chan
	OUTPUT_HEAD_CHAN = head_output_chan
	GET_INPUT_CHANS = get_input_chans
	GET_OUTPUT_CHANS = get_output_chans
	COMMAND_CHANS = command_chans
	TOKEN_CHANS = token_chans

	var wg sync.WaitGroup

	for i := 0; i < NUM_SERVERS; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			WriteOutCommand(i)
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		restReciver()
	}()

	wg.Wait()

	/*
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
	*/
}