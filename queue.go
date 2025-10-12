package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Message map[string]interface{}
type Queue struct {
	messages []Message
	lock     sync.Mutex
}

var queues = make(map[string]*Queue)
var queuesLock sync.Mutex

func getQueue(name string) *Queue {
	queuesLock.Lock()
	defer queuesLock.Unlock()

	if _, exists := queues[name]; !exists {
		queues[name] = &Queue{messages: []Message{}}
	}
	return queues[name]
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Path[len("/publish/"):]
	q := getQueue(queueName)

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	q.lock.Lock()
	q.messages = append(q.messages, msg)
	q.lock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"queued": len(q.messages),
	})
}

func consumeHandler(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Path[len("/consume/"):]
	q := getQueue(queueName)

	q.lock.Lock()
	defer q.lock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if len(q.messages) > 0 {
		msg := q.messages[0]
		q.messages = q.messages[1:]
		json.NewEncoder(w).Encode(msg)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "empty"})
	}
}

func main() {
	http.HandleFunc("/publish/", publishHandler)
	http.HandleFunc("/consume/", consumeHandler)

	fmt.Println("Queue service running on :5001")
	if err := http.ListenAndServe(":5001", nil); err != nil {
		panic(err)
	}
}