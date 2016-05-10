package main

import (
	"fmt"
	"time"
)

// Change state by shifting x bits
type World int

// channels struct used to implement a structured way to handle multiple
// write/read channels for session
type channels struct {
	server    chan (Data)
	players   chan (Data)
	asteroids chan (Data)
}

// session struct stores info regarding players,session managers,
// read/write channels etc.
type session struct {
	players         int
	maxPlayers      int
	world           World
	asteroids       []*asteroid
	asteroidManager *asteroidManager
	listenerManager *ListenerManager
	// For external communication
	write channels
	read  channels
}

// loop is the sessions ....TODO
func (session *session) loop() {

	for {

		select {
		case data := <-session.read.server:

			if data.action == "poke" {

				// Check if theres room inside the session
				if session.players < session.maxPlayers {
					session.write.server <- Data{"Session has room", 200}
				} else {
					session.write.server <- Data{"Session full", -1}
				}

			} else {

				// Spawn a new player
				var port = session.listenerManager.NewObject()
				session.players++
				session.write.server <- Data{"Session: player created", port}
			}

		// Send response back to server
		case userdata := <-session.read.players:

			fmt.Printf("Session: Read from manager %s\n", userdata.action)
			session.write.server <- userdata

		default:
			// Nothing
		}

	}

}

// Session …
func Session(serverConn *Connection, startPort int, players int) {

	session := new(session)
	session.maxPlayers = players

	session.write.server = serverConn.write
	session.read.server = serverConn.read
	session.asteroids = make([]*asteroid, 0)

	// CREATE MANAGERS
	// TODO: Loopify
	session.write.server <- Data{"Session created", 0}
	session.createManagers(startPort)

	// RESPOND TO SERVER
	//

	session.loop()

}

// createManagers sets up managers and their respective connections to/from session
func (session *session) createManagers(startPort int) {

	toPlayers, fromPlayers := makeConnection()
	session.write.players = toPlayers.write
	session.read.players = toPlayers.read

	toAsteroids, _ := makeConnection()
	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	session.asteroidManager = newAsteroidManager()
	session.listenerManager = newListenerManager()

	go session.asteroidManager.loop(toAsteroids, session.asteroids)
	go session.listenerManager.loop(fromPlayers, session.maxPlayers, startPort)

	time.Sleep(250 * time.Millisecond)

	//var port = session.listenerManager.NewObject()
	//fmt.Println("create manager Player created", port)
	//session.players++
	//session.write.server <- Data{"Session: response to server", port}

	//go createAsteroidManager(toAsteroids, session.asteroids)

}
