package config
import (
    sf "bitbucket.org/krepa098/gosfml2"
    "time"
    "sync"
    "github.com/aarzilli/golua/lua"
    rm "../resourcemanager"
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

type Coord uint64

type Entity interface {
    Position() (float32, float32)
    SetPosition(x, y float32) bool
    GetSprite() *sf.Sprite
    Type() uint16
    Remove() bool
}

type LivingEntity interface {
    Entity
    Move(float32, float32) bool
    Collides(float32, float32) bool
    Dir() (float32, float32)
    SetDir(float32, float32) bool
    Talk(string) bool
    Hurt(int16, LivingEntity) int16
}

type Gameconfig struct {
    Rm *rm.ResourceManager
    ContextSettings sf.ContextSettings
    ScreenWidth uint
    ScreenHeight uint
    BitDepth uint
    GameTitle string
    GameActive bool
    Scrolling bool
    Window *sf.RenderWindow
    TextMode bool
    Connected bool
}

var Conf = &Gameconfig{
        Rm: rm.New(),
        ContextSettings: sf.ContextSettingsDefault(),
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
