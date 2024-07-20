package player

import (
	"github.com/dfirebaugh/cube/pkg/component"
	"github.com/dfirebaugh/cube/pkg/message"
)

type Player struct {
	component.Position
	component.Size
	VX, VY, VZ float32
	RX, RY, RZ float32
	H, W, L    float32
	broker     message.MessageBus
}

func New(x, y, z float32, broker message.MessageBus) *Player {
	player := &Player{
		Position: component.Position{
			X: x,
			Y: y,
			Z: z,
		},
		Size: component.Size{
			H: 2,
			W: 1,
			L: 1,
		},
		broker: broker,
	}
	go player.listenForMovement(broker)
	return player
}

func (p *Player) listenForMovement(broker message.MessageBus) {
	sub := broker.Subscribe()
	for msg := range sub {
		if msg.GetTopic() == "PlayerMove" {
			payload := msg.GetPayload().([3]float32)
			p.ApplyVelocity(component.Velocity{X: payload[0], Y: payload[1], Z: payload[2]})
		}
	}
}

func (p *Player) ApplyVelocity(v component.Velocity) {
	if !p.shouldApplyVelocity(v) {
		return
	}
	p.VX = v.X
	p.VY = v.Y
	p.VZ = v.Z
}

func (p *Player) shouldApplyVelocity(v component.Velocity) bool {
	return v.X != 0 || v.Y != 0 || v.Z != 0
}

func (p *Player) Update() {
	// Update player position based on velocity and other factors
	p.X += p.VX
	p.Y += p.VY
	p.Z += p.VZ

	// Check if the player hits the ground
	if p.Y <= 0 {
		p.Y = 0
		p.VY = 0 // Stop the downward velocity
	}

	// Apply friction or other physics effects
	p.VX *= 0.9
	p.VY *= 0.9
	p.VZ *= 0.9

	// Publish the updated position
	p.broker.Publish(message.Message{
		Topic:     "PlayerMove",
		Requestor: "player",
		Payload:   [3]float32{p.X, p.Y, p.Z},
	})
}

func (p *Player) Jump() {
	p.VY += 10 // set a positive velocity to jump up
}
