package config
import (
    "time"
    "sync"
	"azul3d.org/gfx/window.v2"
    "github.com/aarzilli/golua/lua"
	"azul3d.org/gfx.v1"
)

const GAMETITLE = "Xinanium"
const SPRITEDIR = "resources/textures/spritesheets/"
const RESOURCESDIR = "resources/"
const TICKS = 50

// Entity Types
const (
    LivingEntityPlayer = iota
    LivingEntityMonster
)
// actions
const (
    PLAYER_MOVE uint64 = iota
    GET_PLAYER
    GET_PLAYER_TEX
    PLAYER_LOGIN
)

var ToDraw []*gfx.Object

type Coord uint64

type Dir struct {
    X, Y float64
}

type Entity interface {
    Position() (float64, float64)
    SetPosition(x, y float64) bool
    Type() uint16
    Remove() bool
    GetSprite() *gfx.Object
}

type LivingEntity interface {
    Entity
    Move(float64, float64) bool
    Collides(float64, float64) bool
    Dir() (float64, float64)
    SetDir(float64, float64) bool
    Talk(string) bool
    Hurt(int16, LivingEntity) int16
}

type Gameconfig struct {
    ScreenWidth int
    ScreenHeight int
    BitDepth uint
    GameTitle string
    GameActive bool
    Scrolling bool
    Window window.Window
    TextMode bool
    Connected bool
    Renderer gfx.Renderer
}

var Conf = &Gameconfig{
        ScreenWidth:     800,
        ScreenHeight:    600,
        BitDepth:        32,
        GameTitle:       GAMETITLE,
        Scrolling:       true,
}

var GameTicker = time.Tick(time.Second / TICKS)

var Players = map[string]LivingEntity{}
var Monsters = map[int64]LivingEntity{}


type luaState struct {
    sync.RWMutex
    State *lua.State
}

//var Lua *lua.State
var Lua = &luaState{}
