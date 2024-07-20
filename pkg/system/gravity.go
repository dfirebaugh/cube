package system

import "github.com/dfirebaugh/cube/pkg/component"

type Gravity struct{}

type HasMass interface {
	ApplyVelocity(component.Velocity)
}

func (g Gravity) Apply(e HasMass) {
	e.ApplyVelocity(component.Velocity{Y: -.1})
}
