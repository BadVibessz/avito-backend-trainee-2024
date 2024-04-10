package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

func New() *Hasher {
	return &Hasher{}
}

func (h *Hasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (h *Hasher) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
