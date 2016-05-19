package asteroids

import (
	"fmt"
	"testing"
)

//Just to keep fmt package imported
var fakeFmtPrint bool = false

func startAsteroidManagerForTest() *asteroidManager {
	manager := newAsteroidManager()
	conn := MakeConnection()

	// Go-routine to be able to read confirmation
	go manager.init(conn.FlipConnection())

	if fakeFmtPrint != false {
		fmt.Println("Just to keep fmt package imported")
	}

	<-conn.read

	return manager
}

func TestNewAsteroidManager(t *testing.T) {

	manager := newAsteroidManager

	if manager == nil {
		t.Error("Manager not created")
	}
}

func TestInitAsteroidManager(t *testing.T) {

	manager := newAsteroidManager()
	conn := MakeConnection()

	// Go-routine to be able to read confirmation
	go manager.init(conn.FlipConnection())

	response := <-conn.read

	if manager.input != conn.write {
		t.Error("Channels incorrect")
	}

	// Testing hardcoded lines in init
	if response.action != "a.manager_ready" || response.result != 200 {
		t.Error("Response incorrect")
	}

	if manager.yMax != 100 {
		t.Error("yMax incorrect")
	}

	if manager.xMax != 100 {
		t.Error("xMax incorrect")
	}

	if manager.maxRoids != 20 {
		t.Error("yMax incorrect")
	}
}

func TestNewAsteroidInManager(t *testing.T) {

	manager := startAsteroidManagerForTest()

	if len(manager.asteroids) != 0 {
		t.Error("Asteroids should be 0 in length")
	}

	manager.newAsteroid()

	if manager.asteroids[0] == nil {
		t.Error("Asteroids appended incorrect")
	}

	if len(manager.asteroids) != 1 {
		t.Error("Asteroids should be 1 in length")
	}

	manager.newAsteroid()

	if len(manager.asteroids) != 2 {
		t.Error("Asteroids should be 2 in length")
	}
}

func TestShouldSpawn(t *testing.T) {

	manager := startAsteroidManagerForTest()

	i := 0
	for i < 20 {
		manager.asteroids = append(manager.asteroids, newAsteroid())
		i++
	}

	if len(manager.asteroids) != 20 {
		t.Error("Asteroids should be 20 in length")
	}

	i = 0

	for i < 100 {
		if manager.shouldSpawn() != false {
			t.Error("Random spawn should not return true if len is 20")
		}
		i++
	}

	manager.removeAsteroid(len(manager.asteroids) - 1)

	for i < 100 {
		if manager.shouldSpawn() {
			if len(manager.asteroids) >= 20 {
				t.Error("Random spawn should not return true if len is 20+ ")
			}
		}
		i++
	}
}

func TestGetAsteroids(t *testing.T) {

	manager := startAsteroidManagerForTest()
	asteroid := newAsteroid()
	manager.asteroids = append(manager.asteroids, asteroid)

	copyList := manager.getAsteroids()

	if manager.asteroids[0] != copyList[0] {
		t.Error("Get asteroids pointer incorrect")
	}
}

func TestRemoveAsteroid(t *testing.T) {

	manager := startAsteroidManagerForTest()
	i := 0
	for i < 20 {
		manager.asteroids = append(manager.asteroids, newAsteroid())
		if len(manager.asteroids) != i+1 {
			t.Error("Remove append incorrect")
		}
		i++
	}

	for i > 0 {
		manager.removeAsteroid(i - 1)
		if len(manager.asteroids) != i-1 {
			t.Error("Remove incorrect")
		}
		i--
	}
}

func TestResumeAsteroid(t *testing.T) {
	manager := startAsteroidManagerForTest()
	asteroid := newAsteroid()
	manager.asteroids = append(manager.asteroids, asteroid)
	conn := MakeConnection()

	asteroid.input = conn.read

	go manager.resumeAsteroids()

	response := <-conn.read

	if response.action != "a.manager_ok" || response.result != 0 {
		t.Error("Response incorrect")
	}

}

func TestGetNextID(t *testing.T) {
	manager := startAsteroidManagerForTest()

	ID := manager.nextID

	if manager.getNextID() != ID {
		t.Error("Next ID incorrect")
	}

	if manager.getNextID() != ID+1 {
		t.Error("Next ID incorrect")
	}
}
