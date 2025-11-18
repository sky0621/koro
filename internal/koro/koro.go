package koro

type Koro struct {
	X, Y  float64
	Speed float64
}

func New(x, y float64) *Koro {
	return &Koro{
		X:     x,
		Y:     y,
		Speed: 2,
	}
}

func (k *Koro) MoveLeft() {
	k.X -= k.Speed
}

func (k *Koro) MoveRight() {
	k.X += k.Speed
}

func (k *Koro) MoveUp() {
	k.Y -= k.Speed
}

func (k *Koro) MoveDown() {
	k.Y += k.Speed
}
