package monster

import (
    sf "bitbucket.org/krepa098/gosfml2"
    "github.com/stevedonovan/luar"
    "math/rand"
    "time"
    "log"
    "../animation"
    wm "../worldmap"
    "../renderer"
    "../config"
    "../event"
)

const RESOURCESDIR = "resources/"
//const MONSTERHEIGHT = 90
//const MONSTERWIDTH =90
const MONSTERHEIGHT = 32
const MONSTERWIDTH = 31

type Inventory struct {
}


type Duration float64

type Dir struct {
    x, y float32
}

type Monster struct {
    Name string
    Sprite *sf.Sprite
    animation.Animation
    Inventory *Inventory
    Speed uint16
    JumpHeight uint8
    Health int16
    Invincible bool
    Invisible bool
    Walking bool
    InAir bool
    Dead bool
    Floating bool
    dir *Dir
    ID int64
    entityType uint16
    handler *luar.LuaObject
}

// monsterxy.type == mosnteremptystruct.type

func New(montype string) *Monster {
    id := rand.Int63()
    sprite := sf.NewSprite(config.Conf.Rm.Texture(config.SPRITEDIR + "player1.png"))
    monster := &Monster{
        Health: 100,
        ID: id,
		JumpHeight: 1,
		Speed: 1000,
		Name: montype,
		Sprite: sprite,
        dir: &Dir{0, 1},
        entityType: config.LivingEntityMonster,
        Animation: animation.Animation{Sprite: sprite, Stopper: make(chan bool, 1)},
    }

    // trigger event and check if it was cancelled
    if !event.Trigger(&event.MonsterNew{Event: event.New(event.TypeMonsterNew), Monster: monster }) {
        sprite.SetTextureRect(sf.IntRect{0, 0, MONSTERWIDTH, MONSTERHEIGHT});
        monster.dir.y = 1

        config.Monsters[id] = monster
        return monster
    }
    return nil
}

func (m *Monster) Handler() *luar.LuaObject {
    return m.handler
}
func (m *Monster) SetHandler(fn interface{}) bool {
    if fn, ok := fn.(*luar.LuaObject); ok && fn.Type == "function" {
        m.handler = fn
        return true
    }
    log.Println("Failed to add handler to monster", m)
    return false
}

func (m *Monster) GetSprite() *sf.Sprite {
    return m.Sprite
}

func (m *Monster) Type() uint16 {
    return m.entityType
}

func (m *Monster) Move(x, y float32) bool {
    elapsed := float32(renderer.Elapsed)
    newCoords := sf.Vector2f{ x * float32(m.Speed) * elapsed, y * float32(m.Speed) * elapsed }

    if m.Speed > 0 && !m.Dead {
        // handle new direction
        if m.dir.x != x || m.dir.y != y {
            m.SetDir(x, y)
        }
        if !m.Collides(newCoords.X, newCoords.Y) {
            if !event.Trigger(&event.MonsterMove{
                Event: event.New(event.TypeMonsterMove),
                Monster: m,
                NewX: newCoords.X,
                NewY: newCoords.Y,
            }) {
                // move the sprite
                m.Sprite.Move(newCoords)

                // handle walk animation
                m.FrameCounter++
                if m.FrameCounter >= 3 {
                //   m.NextFrame()
                    m.FrameCounter = 0
                }

                return true
            }
        }
    }
    return false
}

func (m *Monster) Hurt(damage int16, damager config.LivingEntity) int16 {
    e := &event.MonsterHurt{ Event: event.New(event.TypeMonsterHurt), Monster: m, Damage: damage, Damager: damager }
    if !m.Invincible && !event.Trigger(e) {
    /* buff, err := sf.NewSoundBufferFromFile("resources/sound/hit.ogg")
        sound := sf.NewSound(buff)
        if err != nil {
        }
        sound.Play()
        */
        health := m.Health
        health -= damage
        if health <= 0 || m.Health < health {
            m.Health = 0
            m.Kill(damager, e)
        } else {
            m.Health = health
        }
        m.Invincible = true
        go func() {
            time.Sleep(time.Second)
            m.Invincible = false
        }()
    }
    return m.Health
}

func (m *Monster) Kill(killer config.LivingEntity, hurtEvent *event.MonsterHurt) bool {
    if !event.Trigger(&event.MonsterKilled{Event: event.New(event.TypeMonsterKilled), Monster: m, Killer: killer, HurtEvent: hurtEvent }) {
        m.Dead = true
        m.Remove()
        // TODO Set dead animation/sprite
        return true
    }
    return false
}


func (m *Monster) Remove() bool {
    if _, ok := config.Monsters[m.ID]; ok {
        if !event.Trigger(&event.MonsterRemoved{Event: event.New(event.TypeMonsterRemoved), Monster: m }) {
            // TODO: memory leak!?
            delete(config.Monsters, m.ID)
            return true
        }
    }
    return false
}

func (m *Monster) SetPosition(x, y float32) bool {
    if !event.Trigger(&event.MonsterChangedPosition{Event: event.New(event.TypeMonsterChangedPosition), Monster: m, NewX: x, NewY: y }) {
        m.Sprite.SetPosition(sf.Vector2f{x, y})
        return true
    }
    return false
}

func (m *Monster) Position() (float32, float32) {
    pos := m.Sprite.GetPosition()
    return pos.X, pos.Y
}


// add stop after removal
func (m *Monster) Run() {
    go func() {
        ticker := time.Tick(time.Second / config.TICKS)
        if m.handler == nil {
            for _ = range ticker {
                if m.Dead { return }
                if config.Conf.GameActive {
                    if (rand.Intn(100) < 98) {
                        m.Move(m.dir.x, m.dir.y)
                    } else {
                        m.SetDir(float32(rand.Intn(3) - 1), float32(rand.Intn(3) - 1))
                        if m.dir.x != 0 && m.dir.y != 0 {
                            m.Move(m.dir.x, m.dir.y)
                        }
                    }
                }
            }
        } else{
            for _ = range ticker {
                if m.Dead { return }
                config.Lua.Lock()
                _, err := m.handler.Call(m)
                config.Lua.Unlock()
                if err != nil {
                    log.Println("error in function call in monster run of", m)
                }
            }
        }
    }()
}

func (m *Monster) Talk(text string) bool {
    if !event.Trigger(&event.MonsterTalk{Event: event.New(event.TypeMonsterTalk), Monster: m, Text: text }) {
        return true
    }
    return false
}

func (p *Monster) Dir() (float32, float32) {
    return p.dir.x, p.dir.y
}

func (m *Monster) SetDir(x, y float32) bool {
    rect := m.Sprite.GetTextureRect()
    if x == 1 {
        rect.Top =  MONSTERHEIGHT* 1
        rect.Width = -MONSTERHEIGHT
        rect.Left = MONSTERHEIGHT
    } else if x == -1 {
        rect.Top = MONSTERHEIGHT * 3
        rect.Width = MONSTERHEIGHT
        rect.Left = 0
    } else if y == 1 {
        rect.Top = MONSTERHEIGHT * 2
        rect.Width = MONSTERHEIGHT
        rect.Left = 0
    } else if y == -1 {
        rect.Top = 0
        rect.Width = MONSTERHEIGHT
        rect.Left = 0
    }
    if !event.Trigger(&event.MonsterChangedDirection{Event: event.New(event.TypeMonsterChangedDirection), Monster: m, NewDirX: x, NewDirY: y}) {
        m.Sprite.SetTextureRect(rect)
        m.dir.x = x
        m.dir.y = y
        return true
    }
    return false
}


const MONSTERFEET = MONSTERWIDTH
func (m *Monster) Collides(x, y float32) bool {
    bounds := m.Sprite.GetGlobalBounds()
    bounds.Left += x
    bounds.Top += y + MONSTERHEIGHT
    worldmap := wm.Current
    if bounds.Left < 0 || bounds.Left + MONSTERWIDTH > float32(worldmap.Width * wm.TILEWIDTH) || bounds.Top - MONSTERFEET < 0 || bounds.Top > float32(worldmap.Height * wm.TILEHEIGHT) {
        return true
    }
    // TODO: Scan for nearby entities
    bounds.Height -= MONSTERFEET + MONSTERHEIGHT
    for _, player := range config.Players {
        if collision, _ := bounds.Intersects(player.GetSprite().GetGlobalBounds()); collision {
            return true
        }
    }

    for _, mon := range config.Monsters {
        if mon != m {
            if collision, _ := bounds.Intersects(mon.GetSprite().GetGlobalBounds()); collision {
                return true
            }
        }
    }

    return false
}


