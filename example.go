package main

import (
	"./postmark"
	"fmt"
)

func main() {

	server := postmark.NewServer("api-key-here")

	err := server.SendSimpleText("signature@example.com", "recipient@example.com", "Example e-mail", "Hello there, this is an example e-mail!")

	if err != nil {
		fmt.Println(err)
	}

}
