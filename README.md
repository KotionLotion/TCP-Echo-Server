# TCP Echo Server

TCP echo server written in Go. It implements concurrency, logging, command handling, and other useful features.

## Features

- Concurrent client handling (using goroutines)
- Connection/disconnection logging
- Message logging to client-specific files
- Input trimming and clean echoing
- Graceful client disconnect handling
- Configurable port via command-line flag
- 30-second inactivity timeout
- Message size limit (1024 bytes)
- Special responses for certain messages
- Command protocol (/time, /quit, /echo)

## How to Run

1. Clone this repository
2. Run the server:

go run server.go --port 4000"

- Use commands/ phrases such as:
-   hello
-   goodbye
-   /time
-   /quit
-   /echo
  
## Link:
