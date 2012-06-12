package postmark

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

// General constants
const (
	empty = "" // An empty string.
)

// Header represents an SMTP header.
type Header struct {
	Name  string // Name of the SMTP header
	Value string // Value of the SMTP header
}

// NewHeader creates a new SMTP header.
func NewHeader(name string, value string) (h *Header) {
	h = new(Header)
	h.Name = name
	h.Value = value
	return h
}

// Message represents an e-mail message.
type Message struct {
	From     string    // The sender e-mail address.
	To       string    // The recipient e-mail address.
	Cc       string    // Comma separated CC e-mail addresses (optional).
	Bcc      string    // Comma separated BCC e-mail addresses (optional).
	Tag      string    // Categorize outgoing e-mail using a tag (optional).
	Subject  string    // Subject line.
	HtmlBody string    // Body in text/html format.
	TextBody string    // Body in text/plain format.
	ReplyTo  string    // The reply-to e-mail address (optional).
	Headers  []*Header // SMTP headers (optional).
}

// NewMessage creates a new e-mail message.
func NewMessage() *Message {
	return new(Message)
}

// AppendTo adds a new recipient.
func (m *Message) AppendTo(recipient string) {
	if len(m.To) > 0 {
		m.To += ", " + recipient
		return
	}
	m.To = recipient
}

// AppendCc adds a new CC address.
func (m *Message) AppendCc(recipient string) {
	if len(m.Cc) > 0 {
		m.Cc += ", " + recipient
		return
	}
	m.To = recipient
}

// AppendBcc adds a new BCC address.
func (m *Message) AppendBcc(recipient string) {
	if len(m.Bcc) > 0 {
		m.Bcc += ", " + recipient
		return
	}
	m.To = recipient
}

// Server represents a Postmark server.
type Server struct {
	APIKey         string // The API Key for this server.
	DefaultFrom    string // The default sender e-mail address.
	DefaultReplyTo string // The default reply-to e-mail address.
}

// Variables shared between servers
var (
	client     *http.Client // The HTTP client used to communicate with the API.
	clientOnce sync.Once    // Make sure we initialize the client only once.
)

// NewServer creates a new Postmark server.
func NewServer(key string) (s *Server) {
	s = new(Server)
	s.APIKey = key
	return s
}

// Postmark API constants
const (
	pmContentType = "application/json"                 // The Content-type used on the API (application/json)
	pmAPIAddress  = "http://api.postmarkapp.com/email" // The API endpoint
	pmRequestType = "POST"                             // The request type (POST)
	headerAccept  = "Accept"                           // The Accept http header name.
	headerContent = "Content-Type"                     // The Content-type http header name.
	headerToken   = "X-Postmark-Server-Token"          // The postmark server token header name.
)

// Postmark API HTTP response codes
const (
	responseSuccess       = 200 // E-mail sent!
	responseUnauthorized  = 401 // Missing or wrong Server API key
	responseUnprocessable = 422 // Invalid message?
	responseInternalError = 500 // Internal
)

// Errors
var (
	ErrorUnknownResponse     = errors.New("The Postmark API returned an unknown HTTP response code.")
	ErrorUnauthorized        = errors.New("The server API key is missing or incorrect.")
	ErrorUnprocessable       = errors.New("The message is invalid. Check if all required fields are filled correctly.")
	ErrorInternalServerError = errors.New("The Postmark server has an internal server error.")
)

// SendSimple is a shortcut for sending simple messages with both a text/html and a text/plain body.
func (s *Server) SendSimple(from string, to string, subject string, textBody string, htmlBody string) error {

	m := NewMessage()
	m.From = from
	m.To = to
	m.Subject = subject
	m.HtmlBody = htmlBody
	m.TextBody = textBody

	return s.Send(m)

}

// SendSimpleHtml is a shortcut for sending simple messages with a text/html body.
func (s *Server) SendSimpleHtml(from string, to string, subject string, body string) error {
	return s.SendSimple(from, to, subject, body, empty)
}

// SendSimpleText is a shortcut for sending simple messages with a text/plain body.
func (s *Server) SendSimpleText(from string, to string, subject string, body string) error {
	return s.SendSimple(from, to, subject, empty, body)
}

// Send connects to the postmark API and sends a message.
func (s *Server) Send(message *Message) error {

	// Initialize HTTP client
	clientOnce.Do(func() {
		client = new(http.Client)
	})

	// Encode data
	data, err := json.Marshal(&message)

	if err != nil {
		return err
	}

	// Set up request.
	request, err := http.NewRequest(pmRequestType, pmAPIAddress, bytes.NewReader(data))

	if err != nil {
		return err
	}

	// Set headers
	request.Header.Add(headerAccept, pmContentType)
	request.Header.Add(headerContent, pmContentType)
	request.Header.Add(headerToken, s.APIKey)

	// Execute request
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	switch response.StatusCode {

	default:
		return ErrorUnknownResponse
	case responseSuccess:
		return nil
	case responseUnauthorized:
		return ErrorUnauthorized
	case responseUnprocessable:
		return ErrorUnprocessable
	case responseInternalError:
		return ErrorInternalServerError

	}

	// Unreachable?
	return nil

}
