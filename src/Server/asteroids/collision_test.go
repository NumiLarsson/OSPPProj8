package asteroids

import (
	"testing"
)

//CreateFakePlayer creates two player for testing
func CreateFakePlayer() (p1 *Player, p2 *Player) {
	p1 = newPlayer()
	p2 = newPlayer()

	p1.X = 5
	p1.Y = 5

	p2.X = 5
	p2.Y = 5

	p1.ID = 1
	p2.ID = 2

	p1.Alive = true
	p2.Alive = true

	return

}

// CreateFakeAsteroid creates two asteroids for testing
func CreateFakeAsteroids() (a1 *Asteroid, a2 *Asteroid) {

	a1 = newAsteroid()
	a2 = newAsteroid()

	a1.X = 5
	a1.Y = 5

	a2.X = 5
	a2.Y = 5

	a1.Alive = true
	a2.Alive = true

	a1.ID = 1
	a2.ID = 2

	return
}

//TestAsteroidCollision tests when an asteroid collide with a asteroid
func TestAsteroidCollision(t *testing.T) {
	var a1, a2 = CreateFakeAsteroids()
	world := new(World)

	world.Asteroids = append(world.Asteroids, a1)
	world.Asteroids = append(world.Asteroids, a2)

	world.asteroidCollision()

	if world.Asteroids[0].Alive || world.Asteroids[1].Alive {
		t.Error("Two asteroid collided but didn't destroy")
	}

}

//TestPlayerCollision tests when a player collides with a player
//and when ha player collides with a asteroid
func TestPLayerCollision(t *testing.T) {

	var p1, p2 = CreateFakePlayer()
	world := new(World)

	var a1, a2 = CreateFakeAsteroids()

	world.Players = append(world.Players, p1)
	world.Players = append(world.Players, p2)

	world.Asteroids = append(world.Asteroids, a1)
	world.Asteroids = append(world.Asteroids, a2)

	world.playerCollision()

	for _, u1 := range world.Asteroids {

		for _, a := range world.Asteroids {
			if u1.Alive || a.Alive {
				t.Error("One player collided with a asteroid but didn't derstroy")
			}
		}

		for _, u2 := range world.Players {
			if u1.Alive || u2.Alive {
				t.Error("Two Players collided but didn't destroy")
			}

		}
	}

}

//TestIsCollision test the isCollision function in collisionmanager.go
func TestIsCollision(t *testing.T) {
	var t1 = isCollision(1, 1, 1, 1)

	if !t1 {
		t.Error("Collision not detected")
	}

}

//TestCollisionList tests if collisions append to a collisionlist
//and if they after that can be removed
func TestCollisionList(t *testing.T) {
	world := new(World)

	var ITER = 100000

	for i := 1; i <= ITER; i++ {
		world.appendCollision(0, i)
		if len(world.Collisions) != i {
			t.Error("ERROR when trying to append a collision")
		}
	}

	world.removeCollisions()

	if len(world.Collisions) != 0 {
		t.Error("ERROR when trying to remove all collisions")
	}

}
