package config
import (
    sf "bitbucket.org/krepa098/gosfml2"
    wm "../worldmap"
    rm "../resourcemanager"
    "time"
    "github.com/aarzilli/golua/lua"
    "github.com/stevedonovan/luar"
    "sort"
    "log"
)

const GAMETITLE = "WeoutHibot"
const SPRITEDIR = "resources/textures/spritesheets/"
const TICKS = 50

type Coord uint64

type Entity interface {
    Position() (float32, float32)
    SetPosition(x, y float32)
    GetSprite() *sf.Sprite
    Type() uint16
}

type LivingEntity interface {
    Entity
    Move(float32, float32) bool
    Collides(float32, float32) bool
    Dir() (float32, float32)
    SetDir(float32, float32)
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
    CurrentMap *wm.WorldMap
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

var Lua *lua.State





// Eventhandling

var eventQueue = make(chan Eventer, 1000)
var EventRegistry = make(map[EventType][]*EventHandler)

type EventHandler struct {
    Priority uint16
    Handler *luar.LuaObject
}

type Eventer interface {
    Type() EventType
}

func RegisterEvent(eType float64, fn interface{}, priority float64) bool {
    e := EventType(eType)
    if v, ok := fn.(*luar.LuaObject); ok {
        EventRegistry[e] = append(EventRegistry[e], &EventHandler{Handler: v, Priority: uint16(priority)})
        sort.Sort(byPriority(EventRegistry[e]))
        return true
    }
    log.Fatal("Couldn't add Event: ", fn)
    return false
}

func RegisterTicker(ticks float64, fn interface{}, priority float64) bool {
    bl := RegisterEvent(EventTypeTicker * ticks + 1000, fn, priority)
    go makeTicker(ticks)
    return bl
}

func makeTicker(ticks float64) {
    ticker := time.Tick(time.Second / TICKS * time.Duration(ticks))
    controlChan := make(chan bool, 1)
    active := true
    tickid := EventTypeTicker * ticks + 1000

    for _ = range ticker {
        select {
        case v := <-controlChan:
            if v == true {
                return
            }
            active = !active && true
        default:
        }
        if active {
            TriggerEvent(&EventTicker{Event: Event{EType: EventType(tickid)}, control: controlChan})
        }
    }
}

// TODO: auto adjust buffer to map
func TriggerEvent(event Eventer) () {
    select {
    case eventQueue <- event:
    default:
    }
}

func InitEventQueue() {
    go func() {
        var err error
        for e := range eventQueue {
            for _, eHandler := range EventRegistry[e.Type()] {
                _, err = eHandler.Handler.Call(e)
                if err != nil {
                    log.Println("error in function call of event type '", e.Type(), "':")
                    log.Println(err)
                }
            }
        }
    }()
}

// sorting
type byPriority []*EventHandler

func (p byPriority) Len() int           { return len(p) }
func (p byPriority) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byPriority) Less(i, j int) bool { return p[i].Priority < p[j].Priority }




// Entity Types
const (
    LivingEntityPlayer = iota
    LivingEntityMonster
)

// Event Types

type EventType uint16
const (
    EventTypePlayerNew = iota
    EventTypePlayerMove
    EventTypePlayerCollision
    EventTypePlayerJump
    EventTypePlayerTalk
    EventTypePlayerChangedDirection
    EventTypePlayerChangedPosition
    EventTypePlayerHurt
    EventTypePlayerKilled
    EventTypePlayerRemoved
    EventTypeMonsterNew
    EventTypeMonsterMove
    EventTypeMonsterCollision
    EventTypeMonsterJump
    EventTypeMonsterTalk
    EventTypeMonsterChangedDirection
    EventTypeMonsterChangedPosition
    EventTypeMonsterHurt
    EventTypeMonsterKilled
    EventTypeMonsterRemoved
    EventTypeKeyPressed
    EventTypeTicker
)

type Event struct {
    EType EventType
}

func (e *Event) Type() EventType {
    return e.EType
}


type EventPlayerNew struct {
    Event
    Player LivingEntity
}

type EventPlayerMove struct {
    Event
    Player LivingEntity
    NewX float32
    NewY float32
}

type EventPlayerCollision struct {
    Event
    Player LivingEntity
    What LivingEntity
    X float32
    Y float32
}

type EventPlayerJump struct {
    Event
    Player LivingEntity
}

type EventPlayerTalk struct {
    Event
    Player LivingEntity
    Text string
}

type EventPlayerChangedDirection struct {
    Event
    Player LivingEntity
    NewDirX float32
    NewDirY float32
}

type EventPlayerChangedPosition struct {
    Event
    Player LivingEntity
    NewX float32
    NewY float32
}

type EventPlayerHurt struct {
    Event
    Player LivingEntity
    Damager Entity
    Damage int16
}

type EventPlayerKilled struct {
    Event
    Player LivingEntity
    Killer Entity
    HurtEvent *EventPlayerHurt
}

type EventPlayerRemoved struct {
    Event
    Player LivingEntity
}

type EventMonsterNew struct {
    Event
    Monster LivingEntity
}

type EventMonsterMove struct {
    Event
    Monster LivingEntity
    NewX float32
    NewY float32
}

type EventMonsterCollision struct {
    Event
    Monster LivingEntity
    What LivingEntity
    X float32
    Y float32
}

type EventMonsterJump struct {
    Event
    Monster LivingEntity
}

type EventMonsterTalk struct {
    Event
    Monster LivingEntity
    Text string
}

type EventMonsterChangedDirection struct {
    Event
    Monster LivingEntity
    NewDirX float32
    NewDirY float32
}

type EventMonsterChangedPosition struct {
    Event
    Monster LivingEntity
    NewX float32
    NewY float32
}

type EventMonsterHurt struct {
    Event
    Monster LivingEntity
    Damager Entity
    Damage int16
}

type EventMonsterKilled struct {
    Event
    Monster LivingEntity
    Killer Entity
    HurtEvent *EventMonsterHurt
}

type EventMonsterRemoved struct {
    Event
    Monster LivingEntity
}

type EventKeyPressed struct {
    Event
    KeyCode uint16
}

type EventTicker struct {
    Event
    control chan bool
}

func (e *EventTicker) Stop() {
    select {
    case e.control <- true:
    default:
    }
}
func (e *EventTicker) PauseResume() {
    select {
    case e.control <- false:
    default:
    }
}
