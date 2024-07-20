package component

type Position struct {
	X, Y, Z float32
}

func (p Position) ApplyVelocity(v Velocity) {

}

type Velocity struct {
	X, Y, Z float32
}

type Size struct {
	H, W, L int
}
