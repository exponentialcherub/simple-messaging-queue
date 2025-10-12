# simple-messaging-queue
A simple messaging queue to publish and consume messages

# Requirements

Go installed - see https://go.dev/doc/install

# To run

go run queue.go

# Endpoints

/publish/{queue_name}
Body = {message}

/consume/{queue_name}
Return Body = {message}
