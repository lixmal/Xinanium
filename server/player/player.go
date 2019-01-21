package player

import (
    commonnet "../../common/network"
    "../network"
    "azul3d.org/lmath.v1"
    "log"
    "math"
)

const PLAYERWIDTH = 31
const PLAYERHEIGHT = 32

type Inventory struct {
}

type Item struct {
}

type Dir struct {
    X, Y int8
}

type Player struct {
    Handle     string
    name       string
    Speed      uint16
    JumpHeight uint8
    Health     int16
    Invincible bool
    Invisible  bool
    Walking    bool
    InAir      bool
    Dead       bool
    Floating   bool
    entityType uint16
    Dir        Dir
    Pos        lmath.Vec3
    *network.Client
}

func New(handle string) *Player {
    player := &Player{
        Health:     100,
        JumpHeight: 10,
        Speed:      1000,
        name:       handle,
        Handle:     handle,
        entityType: 0,
        Dir:        Dir{0, 1},
        Pos:        lmath.Vec3{0, 0, 0},
    }

    return player
}

func (p *Player) Move(v *lmath.Vec3) error {
    maxv := lmath.Vec3{10, 10, 10}

    // restrain movement
    if (lmath.Vec3{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)}.Sub(maxv).AnyGreater(lmath.Vec3{0, 0, 0})) {
        log.Println(p.Handle+": Movement to fast:", v.Sub(maxv))
    } else {
        p.Pos = p.Pos.Add(*v)
    }

    p.Dir.X = int8(v.X)
    p.Dir.Y = int8(v.Z)
    log.Println(p)
    log.Println(p.Dir.Name())
    return nil
}

func (p *Player) Remove() {
    p.Client.Player = nil
    p.Client = nil
}

func (p *Player) SetPosition(v *lmath.Vec3) {
    p.Pos = *v
    p.WriteAction(commonnet.PLAYER_POS, v)
}

func (d *Dir) Name() string {
    /*
                 y
                /|\
                 |
            –––––0–––––>x
                 |
                 |
    */
    x, y := d.X, d.Y
    if x > 0 && y > 0 {
        return "NorthEast"
    } else if x < 0 && y > 0 {
        return "NorthWest"
    } else if x > 0 && y < 0 {
        return "SouthEast"
    } else if x < 0 && y < 0 {
        return "SouthWest"
    } else if x > 0 {
        return "East"
    } else if x < 0 {
        return "West"
    } else if y < 0 {
        return "South"
    } else if y > 0 {
        return "North"
    }
    return ""
}

func (d *Dir) String() string {
    return d.Name()
}

func (p *Player) Collides(float64, float64) bool { return false }
func (p *Player) Talk(string) bool               { return false }
func (p *Player) Position() (float64, float64)   { return 0, 0 }
func (p *Player) Type() uint16                   { return 0 }
func (p *Player) Name() string {
    return p.name
}
