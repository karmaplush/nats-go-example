package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

func senderServer(nc *nats.Conn) {

	senderMux := http.NewServeMux()
	senderMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			"/send\n"+
				"Path for sending simple message (current unix time) in queue\n"+
				"Use ?message= query param for sending time and some text message",
		)
	})

	// Simple handler for sending message in queue
	senderMux.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {

		unixNow := time.Now().Unix()
		unixNowStr := strconv.Itoa(int(unixNow))

		messageQuery := r.URL.Query().Get("message")

		if messageQuery != "" {
			messageToPublish := fmt.Sprintf("%s, message: %s", unixNowStr, messageQuery)

			nc.Publish("messages", []byte(messageToPublish))
			fmt.Fprintf(w, "Unix time and \"%s\" message was was sended to queue", messageQuery)
		} else {
			nc.Publish("messages", []byte(unixNowStr))
			fmt.Fprintf(w, "Unix time was was sended to queue")
		}

	})

	sender := &http.Server{Addr: ":8080", Handler: senderMux}
	sender.ListenAndServe()
}

func listenerServer(nc *nats.Conn) {

	// Some kind of "storage" (for testing purposes only)
	messagesStorage := []string{}

	// Subscribe to the "messages" subject and store message in storage
	_, err := nc.Subscribe("messages", func(m *nats.Msg) {
		messagesStorage = append(messagesStorage, string(m.Data))
	})
	if err != nil {
		panic(err)
	}

	listenerMux := http.NewServeMux()
	listenerMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			"/listen\n"+
				"Path for reading messages from storage\n"+
				"Every /listen call will read one message from storage",
		)
	})

	// Simple handler for listening messages from queue
	listenerMux.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {
		if len(messagesStorage) > 0 {

			message := messagesStorage[0]
			messagesStorage = messagesStorage[1:]

			fmt.Fprintf(
				w,
				"%s\nMessages in queue allowed: %d",
				message,
				len(messagesStorage),
			)

		} else {
			fmt.Fprintln(w, "No messages was found in queue")
		}
	})

	listener := &http.Server{Addr: ":9090", Handler: listenerMux}
	listener.ListenAndServe()
}

func main() {

	// Connect to NATS server for publishing messages
	ncPublisher, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}
	defer ncPublisher.Close()

	// Connect to NATS server for subscribing to messages
	ncSubscriber, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}
	defer ncSubscriber.Close()

	go senderServer(ncPublisher)
	go listenerServer(ncSubscriber)

	fmt.Println(
		"HTTP servers at :8080 (sender) & :9090 (listener) was started (Shutdown it ungracefully!)",
	)
	select {}
}
