package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"fmt"
	"time"
)

type User struct {
	Name string
	Output chan Message
}

type Message struct {
	Username string
	Text string
}

/*Room of a chatserver*/
type ChatServer struct {
	Users map[string]User
	Join chan User
	Leave chan User
	Input chan Message
}

/* how would we send the message?
 * Run method: all the logic of the chat
 * Message in the chat + action in the map of users
 */
func (cs *ChatServer) Run() {
	for {
		select {
			case user := <-cs.Join:
				cs.Users[user.Name] = user
				go func() {
					cs.Input <- Message{
						Username: "SYSTEM",
						Text: fmt.Sprintf("%s joined", user.Name),
					}
				}()
				
			case user := <-cs.Leave:
				delete(cs.Users, user.Name)
				go func() {
					cs.Input <- Message{
						Username: "SYSTEM",
						Text: fmt.Sprintf("%s left", user.Name),
					}
				}()

			case msg := <-cs.Input:
				for _, user := range cs.Users {
					select {
						case user.Output <- msg:
						case <-time.After(time.Second * 10):
					}
					
				}
		}
	}
}

/*****/
func handleConn (chatServer *ChatServer, conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter your username: ")
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

	//write to connection, when it is disconnects, this would return an error
	for msg := range user.Output {
		_, err := io.WriteString(conn, msg.Username+": "+msg.Text+"\n")
		if err != nil {
			break
		}
	}

}

/* 
 * creating the listener
 * defer: execute a task at the end of the enclosing function (main)
 */

func main() {
	server, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer server.Close()

	chatServer := &ChatServer{
		Users: make(map[string]User),
		Join: make(chan User),
		Leave: make(chan User),
		Input: make(chan Message),
	}

	go chatServer.Run()


	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConn(chatServer, conn)
	}
}

/*
goroutine: lightweight thread of execution.
channel: pipes that connect concurrent goroutines.
*/