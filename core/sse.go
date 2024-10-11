package core

import (
	"fmt"
	"log"
	"net/http"
)

func (ctrl *DynamicController) Sse(path string, sseFnc SseFnc) {
	b := &SseBroker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}
	b.Start()

	ctrl.Handler(path, b)
	go sseFnc(b)
}

type SseBroker struct {
	// Create a map of Clients, the keys of the map are the channels
	// over which we can push messages to attached Clients.  (The values
	// are just booleans and are meaningless.)
	//
	Clients map[chan string]bool

	// Channel into which new clients can be pushed
	//
	NewClients chan chan string

	// Channel into which disconnected clients should be pushed
	//
	DefunctClients chan chan string

	// Channel into which Messages are pushed to be broadcast out
	// to attahed clients.
	//
	Messages chan string
}

type SseFnc func(broker *SseBroker)

func (b *SseBroker) Start() {
	// Start a goroutine
	//
	go func() {

		// Loop endlessly
		//
		for {

			// Block until we receive from one of the
			// three following channels.
			select {

			case s := <-b.NewClients:

				// There is a new client attached and we
				// want to start sending them messages.
				b.Clients[s] = true
				log.Println("Added new client")

			case s := <-b.DefunctClients:

				// A client has dettached and we want to
				// stop sending them messages.
				delete(b.Clients, s)
				close(s)

				log.Println("Removed client")

			case msg := <-b.Messages:

				// There is a new message to send.  For each
				// attached client, push the new message
				// into the client's message channel.
				for s := range b.Clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.Clients))
			}
		}
	}()
}

func (b *SseBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Make sure that the writer supports flushing.
	//
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a new channel, over which the broker can
	// send this client messages.
	messageChan := make(chan string)

	// Add this client to the map of those that should
	// receive updates
	b.NewClients <- messageChan

	// Listen to the closing of the http connection via the Request.Context
	notify := r.Context().Done()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		b.DefunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Don't close the connection, instead loop endlessly.
	for {

		// Read from our messageChan.
		msg, open := <-messageChan

		if !open {
			// If our messageChan was closed, this means that the client has
			// disconnected.
			break
		}

		// Write to the ResponseWriter, `w`.
		fmt.Fprintf(w, "data: Message: %s\n\n", msg)

		// Flush the response.  This is only possible if
		// the repsonse supports streaming.
		f.Flush()
	}

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}
