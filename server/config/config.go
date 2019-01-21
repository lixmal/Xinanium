package config

import (
    "azul3d.org/lmath.v1"
    "code.google.com/p/gcfg"
    "log"
)

const (
    SPRITEDIR = "resources/textures/spritesheets/"
    // SPRITEDIR       = filepath.Dir(filepath.FromSlash("resources/textures/spritesheets")) + "/"
    SPRITEEXTENSION = ".png"
)

type Entity interface {
    Position() (float64, float64)
    SetPosition(*lmath.Vec3)
    Type() uint16
    Remove()
}

type LivingEntity interface {
    Entity
    Move(*lmath.Vec3) error
    Collides(float64, float64) bool
    Talk(string) bool
    //    Hurt(int16, LivingEntity) int16
    Name() string
}

type network struct {
    Server struct {
        Host string
        Port string
    }
}

var Network network

func init() {
    if err := gcfg.ReadFileInto(&Network, "server.conf"); err != nil {
        log.Fatal("Unable to read config file server.conf: ", err)
    }
}
