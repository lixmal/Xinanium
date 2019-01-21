package config

const HashCost = 0
const (
    HandleMinLength = 3
    HandleMaxLength = 128
)

type Credentials struct {
    Handle   string `db:"handle"`
    Password []byte `db:"password"`
}
