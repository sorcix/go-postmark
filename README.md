Go Postmark library
===================

This library should provide everything you need to send e-mails using the Postmark API.

Important: Work in progress! You shouldn't rely on this for important e-mails. As I currently don't handle the API response, the server.Send methods could change in the future in order to return more information.

Works
------
 * Sending text and HTML e-mail without attachments.
 * Errors based on HTTP response codes.
 
Roadmap
--------
 * Attachments
 * Use Url Fetch API when running on Google App Engine.
 * Batch sending
 * Handling Postmark API response and errors.
 
Example
--------

Sending a simple text/plain e-mail:
```go
server := postmark.NewServer("api-key-here")
	
err := server.SendSimpleText("signature@example.com","recipient@example.com", "Example e-mail", "Hello there, this is an example e-mail!")

if err != nil {
	fmt.Println(err)
}
```

You can also use the Message object directly if you need more options:
```go
server := postmark.NewServer("api-key-here")

message := postmark.NewMessage()
message.To = "signature@example.com"
message.From = "recipient@example.com"
message.ReplyTo = "replyto@example.com"
message.Tag = "signup"
message.TextBody = "Hello there, this is an example e-mail!")

err := server.Send(message)

if err != nil {
	fmt.Println(err)
}
```
