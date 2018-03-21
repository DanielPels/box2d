package box2d_test

import (
	"testing"
	"github.com/DanielPels/box2d"
	"math"
	"math/rand"
)

func BenchmarkTest(b *testing.B) {

	fps := 1.0 / 60.0

	// Define the gravity vector.
	gravity := box2d.MakeB2Vec2(0.0, 0.0)

	// Construct a world object, which will hold and simulate the rigid bodies.
	world := box2d.MakeB2World(gravity)

	world.SetContactFilter(&box2d.B2ContactFilter{})

	var PlayerDef box2d.B2BodyDef
	var PlayerShape box2d.B2CircleShape
	var PlayerFixture box2d.B2FixtureDef

	PlayerDef = box2d.MakeB2BodyDef()
	PlayerDef.AllowSleep = false
	PlayerDef.FixedRotation = true
	PlayerDef.LinearDamping = 5
	PlayerDef.Type = box2d.B2BodyType.B2_dynamicBody

	PlayerShape = box2d.MakeB2CircleShape()
	PlayerShape.M_radius = 0.5

	PlayerFixture = box2d.MakeB2FixtureDef()
	PlayerFixture.Shape = &PlayerShape
	PlayerFixture.Density = 20

	var WallDef box2d.B2BodyDef
	var WallShape box2d.B2PolygonShape
	var WallFixture box2d.B2FixtureDef

	WallDef = box2d.MakeB2BodyDef()
	WallDef.AllowSleep = true
	WallDef.Type = box2d.B2BodyType.B2_staticBody

	WallShape = box2d.MakeB2PolygonShape()

	WallFixture = box2d.MakeB2FixtureDef()
	WallFixture.Shape = &WallShape

	// maak een box

	WallShape.SetAsBox(0.5, 5)
	WallDef.Position.Set(-5, 0)
	leftWall := world.CreateBody(&WallDef)
	leftWall.CreateFixtureFromDef(&WallFixture)

	WallShape.SetAsBox(0.5, 5)
	WallDef.Position.Set(5, 0)
	rightWall := world.CreateBody(&WallDef)
	rightWall.CreateFixtureFromDef(&WallFixture)

	WallShape.SetAsBox(5, 0.5)
	WallDef.Position.Set(0, -5)
	topWall := world.CreateBody(&WallDef)
	topWall.CreateFixtureFromDef(&WallFixture)

	WallShape.SetAsBox(5, 0.5)
	WallDef.Position.Set(0, 5)
	bottomWall := world.CreateBody(&WallDef)
	bottomWall.CreateFixtureFromDef(&WallFixture)

	// maak x aantal cirkels

	players := make([]*playerCircle, 0, 100)

	for i := 0; i < 100; i++ {
		PlayerDef.Position.Set(math.Cos(math.Pi*2)*rand.Float64(), math.Sin(math.Pi*2)*rand.Float64())
		b := world.CreateBody(&PlayerDef)
		b.CreateFixtureFromDef(&PlayerFixture)

		players = append(players, &playerCircle{
			body: b,
			vec:  box2d.B2Vec2{X: math.Cos(math.Pi*2) * rand.Float64(), Y: math.Sin(math.Pi*2) * rand.Float64()},
		})
	}

	b.ReportAllocs()
	b.ResetTimer()

	// begin world step bench mark

	for i := 0; i < b.N; i++ {
		world.Step(fps, 1, 1)
		for p := 0; p < len(players)-1; p++ {
			players[p].body.SetLinearVelocity(players[p].vec)
		}

		b.StopTimer()
		p1 := box2d.B2Vec2{X: (math.Cos(math.Pi*2) * rand.Float64() * 5), Y: (math.Sin(math.Pi*2) * rand.Float64() * 5)}
		p2 := box2d.B2Vec2{X: math.Cos(math.Pi*2) * rand.Float64() * 5, Y: math.Sin(math.Pi*2) * rand.Float64() * 5}
		b.StartTimer()

		world.RayCast(func(fixture *box2d.B2Fixture, point box2d.B2Vec2, normal box2d.B2Vec2, fraction float64) float64 {
			return 1
		}, p1, p2)
	}
}

type playerCircle struct {
	body *box2d.B2Body
	vec  box2d.B2Vec2
}
