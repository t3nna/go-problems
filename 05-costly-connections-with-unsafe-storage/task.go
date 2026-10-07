package main

import "sync"

type Connection interface {
	// Need call Connect before Send
	// Take time to connect
	Connect()

	// Every connection should be disconnected after use
	// Take time to disconnect
	Disconnect()

	Send(req string) (string, error)
}

type ConnectionCreator interface {
	// Create new connection
	// Will return error if there is more than maxConn
	NewConnection() (Connection, error)
}

type Saver interface {
	// Saves data to unsafe storage
	// WILL CORRUPT DATA on concurrent save
	Save(data string)
}

// SendAndSave should send all requests concurrently using at most `maxConn` simultaneous connections.
// Responses must be saved using Saver.Save.
// Be careful: Saver.Save is not safe for concurrent use.
func SendAndSave(creator ConnectionCreator, saver Saver, requests []string, maxConn int) {

	var wg sync.WaitGroup

	wg.Add(maxConn)

	reqCh, resCh := make(chan string, len(requests)), make(chan string, len(requests))

	for _, v := range requests {
		reqCh <- v
	}
	close(reqCh)

	for range maxConn {
		go func() {

			conn, err := creator.NewConnection()
			if err != nil {
				return
			}

			conn.Connect()
			defer conn.Disconnect()

			for c := range reqCh {
				v, err := conn.Send(c)
				if err != nil {
					break
				}
				resCh <- v
			}

			wg.Done()
		}()

	}
	go func() {
		wg.Wait()
		close(resCh)
	}()
	for v := range resCh {
		saver.Save(v)
	}

}
