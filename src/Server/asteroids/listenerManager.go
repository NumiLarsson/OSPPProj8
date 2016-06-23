package asteroids

import (
	"fmt"
	"os"
	"time"
)

//ListenerManager is used as a struct to basically emulate an object
type ListenerManager struct {
	xMax           int
	yMax           int
	maxPlayers     int
	currentPlayers int
	currentPort    int
	nextID         int
	input          chan Data
	output         chan Data
	listeners      []*Listener
	players        []*Player
	lastWorld	   World
}

//Update is used to send all the changes in the world to the client, 
//rather than sending the entire gamestate every frame.
type Update struct {
	World *World
}

// loop is where the listenerManager spinns waiting for tick message from session,
// once received it handles collisions from last tick and collect all new positions
func (manager *ListenerManager) loop(sessionConn *Connection, maxPlayers int, startPort int) {

	manager.init(sessionConn, maxPlayers, startPort)

	for {
		select {

		case msg := <-manager.input:

			if msg.action == "session.tick" {				
				//manager.print()
				for _, character := range manager.players {
					//If the character is currently alive, award points
					if character.Alive {
						character.Points++;
					}
				}
				// TODO: below correct way to use handle ??
				//We should move this to use manager.listener.player positions instead.
				manager.players = manager.collectPlayerPositions()
				// Send update + world to players
			}

			if msg.action == terminate {

				DebugPrint(fmt.Sprintln("[LIST.MAN] Dead {TODO: Close socket & kill listeners}"))
				manager.output <- Data{terminated, ok}
				return

			}
		}
	}
}

//Kill is used to terminate the listenerManager.
func (manager *ListenerManager) kill() {

	go func() {
		manager.input <- Data{terminate, request}
	}()

}

// newAsteroidsManager creates a new asteroid manager
func newListenerManager() *ListenerManager {

	DebugPrint(fmt.Sprintln("[LIST.MAN] Created"))
	return new(ListenerManager)

}

// init initiates the listenerManager with a cap on maxPlayers
// connected and maxPlayers numbers of ports in a row from firstPort
func (manager *ListenerManager) init(sessionConn *Connection,
	maxPlayers int, firstPort int) {

	// TODO fix hardcoded variables
	manager.xMax = 400
	manager.yMax = 225
	//World size in x = horizontal & y = longitude.
	
	manager.maxPlayers = maxPlayers
	manager.nextID = 1
	manager.currentPort = firstPort
	manager.input = sessionConn.read
	manager.output = sessionConn.write

	//manager.listeners = make([]*Listener, 0)
	sessionConn.write <- Data{"l.manager_ready", ok}
}

// newPlayer creates a new listener for the listener manager, used to connect to a new player.
func (manager *ListenerManager) newPlayer() (int, *Player) {

	DebugPrint(fmt.Sprintln("[LIST.MAN] Creating new object in listener manager"))

	//Creation of the listener and listener-player
	listener := newListener()
	listener.init(manager.currentPort)

	listener.player = newPlayer()
	listener.player.init(manager.getNextID(), manager.xMax, manager.yMax)

	//Insert in the managers lists
	manager.listeners = append(manager.listeners, listener)
	manager.players = append(manager.players, listener.player)

	manager.incrementCurrentPlayers()

	go listener.startUpListener()

	return manager.getNextPort(), listener.player
}

// getNextID returns the id to be used and sets the next value
func (manager *ListenerManager) getNextPort() int {
	defer func() { manager.currentPort++ }()
	return manager.currentPort
}

// incrementCurrentPlayers increments currentPlayers in the manager.
// This is a function so that we can use defer() in a better looking way.
func (manager *ListenerManager) incrementCurrentPlayers() {
	manager.currentPlayers++
}

// getNextID returns the id to be used and sets the next value
//This is the alternative way to use defer(), but it's less readable.
func (manager *ListenerManager) getNextID() int {
	defer func() { manager.nextID++ }()
	return manager.nextID
}

// collectPlayerPositions collect all player positions and return an array of them
func (manager *ListenerManager) collectPlayerPositions() []*Player {

	var playerList []*Player
	for _, listener := range manager.listeners {

		player := listener.getPlayer()
		playerList = append(playerList, player)
	}

	return playerList
}

// getPlayers returns an array of players
func (manager *ListenerManager) getPlayers() []*Player {
	return manager.players
}

// sendToClient broadcasts world-info to every listener
func (manager *ListenerManager) sendToClient(world *World) {
	AllDead := true
	for _, character := range manager.getPlayers() {
		if character.Alive || character.Lives > 0 {
			AllDead = false
		}
	}
		
	if (AllDead) {
		fmt.Println("Range on manager:", len( manager.listeners))		
		for variant, listener := range manager.listeners {
			fmt.Println("Sendind endgame to:", variant)
			go listener.WriteEndGame(world)			
		}
		shutdown();
	}
	
	for _, listener := range manager.listeners {
		go listener.Write(world)
	}
}
func shutdown() {
	time.Sleep(time.Second * 4)
	os.Exit(0)
}

// print is used to print all players that have been in a collision
func (manager *ListenerManager) print() {

	var list []int

	for _, player := range manager.players {
		if !player.isAlive() {
			list = append(list, player.ID)
		}
	}

	if len(list) > 0 {
		DebugPrint(fmt.Sprintln("[LIST.MAN] Collision:", list))
	}
}