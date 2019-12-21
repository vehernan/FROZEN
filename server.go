package main //the executable 

// impoting libraries
import (
	"fmt" //package in go for print stdout
	"net/http" //web server!
	"strings" //use str
	"log"
	"bufio"
	"io/ioutil" // implements i/o utility functions
)

type ChatServer struct //
{
	Users map[string]User
	Join chan User //to know when a user has joined
	Leave chan User //to know when user had left
	Input chan string //what has been writen 
}

type User struct
{
	Name string // name of user 
	Output chan Message //channel of message
}

type Message struct
{
	Username string //who wrote it
	Text string //what was put
}

func (cs *ChatServer) Run() //logic needed for process
(
	for (
		select (
		case user := <-cd.Join:
			cs.User[user.Name] = user
			// say "someoned joined"
			go func() { //function that is running concurrently w/the for loop for thr mssg
			cs.Input <- Message{
				UserName: "SYSTEM", //itll tell me thats my message
				Text: fmt.Sprintf("%s joined", user.Name), 
			}
		}()
		case user := <-cs.Leave:
			delete(cs.Users, user.Name)
			//say "someonened left"
			go func() { //function that is running concurrently w/the for loop for thr mssg
				cs.Input <- Message{
					UserName: "SYSTEM",
					Text: fmt.Sprintf("%s left", user.Name),
				}
			}()
		case msg := <-cs.Input:
			for _, u := range cs.Users { //loop through every user and send a message to that user 
				user.Output <- msg // add channel to user called output
			}

		)
	)
)

func handleConn(chatServer *ChatServer, conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter your Username:")
	
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	user := User{
		Name: scanner.Text(),
		Output: make(chan Message),
	}
	chatServer.Join <- user
	defer func() {
		chatServer.Leave <- user
	}()

	//read from connection 
	go func() {
		for scanner.Scan() {
			ln := scanner.Text()
			chatServer.Input <- Message{user.Name, ln}
		}

	}()
	//write to connection
	for msg := range user.Output {
		io.WriteString(conn, msg.Username+ ": " +msg.Text)
	}

func main()
{
	server, er := net.Listen("tcp, 9000") ://error hand, listening to connection
	//tcp is terminal, can be changed to http
	if err := nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()

	chatServer := &ChatServer{ //intitializing it
		Users: make(map[string]User),
		Join: make(chan User), //to know when a user has joined
		Leave: make(chan User), //to know when user had left
		Input: make(chan Message), //what has been writen
	}
	go chat.Server.Run()

	for {
		conn, err := server.Accept()
		if err := nil {
			log.Fatalln(err.Error())
		}
		go handleConn(chatServer, conn)
	}
}
