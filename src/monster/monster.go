package monster

import (
    sf "bitbucket.org/krepa098/gosfml2"
    "../animation"
    wm "../worldmap"
    "math/rand"
    "../renderer"
    "../config"
    "time"
)

const RESOURCESDIR = "resources/"
const MONSTERHEIGHT = 90
const MONSTERWIDTH = 90

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
    JumpHeight byte
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
}

func New(montype string) *Monster {
    id := rand.Int63()
    sprite := sf.NewSprite(config.Conf.Rm.Texture(config.SPRITEDIR + "monster.png"))
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
    config.TriggerEvent(&config.EventMonsterNew{Event: config.Event{EType: config.EventTypeMonsterNew}, Monster: monster })

    sprite.SetTextureRect(sf.IntRect{0, 0, 90, 90});
    monster.dir.y = -1



    config.Monsters[id] = monster
    return monster
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
            config.TriggerEvent(&config.EventMonsterMove{
                Event: config.Event{EType: config.EventTypeMonsterMove},
                Monster: m,
                NewX: newCoords.X,
                NewY: newCoords.Y,
            })

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
    return false
}

func (m *Monster) Hurt(damage int16, damager config.LivingEntity) int16 {
    e := &config.EventMonsterHurt{Event: config.Event{EType: config.EventTypeMonsterHurt}, Monster: m, Damage: damage, Damager: damager }
    config.TriggerEvent(e)

    if !m.Invincible {
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
    }
    return m.Health
}

func (m *Monster) Kill(killer config.LivingEntity, hurtEvent *config.EventMonsterHurt) bool {
    config.TriggerEvent(&config.EventMonsterKilled{Event: config.Event{EType: config.EventTypeMonsterKilled}, Monster: m, Killer: killer, HurtEvent: hurtEvent })
    return true
}


func (m *Monster) Remove() bool {
    _, ok := config.Monsters[m.ID]
    if !ok {
        return false
    }
    config.TriggerEvent(&config.EventMonsterRemoved{Event: config.Event{EType: config.EventTypeMonsterRemoved}, Monster: m })
    delete(config.Monsters, m.ID)
    // TODO: memory leak!?
    return true
}

func (m *Monster) SetPosition(x, y float32) {
    config.TriggerEvent(&config.EventMonsterChangedPosition{Event: config.Event{EType: config.EventTypeMonsterChangedPosition}, Monster: m, NewX: x, NewY: y })
    m.Sprite.SetPosition(sf.Vector2f{x, y})
}

func (m *Monster) Position() (float32, float32) {
    pos := m.Sprite.GetPosition()
    return pos.X, pos.Y
}

func (m *Monster) Run() {
    go func() {
        ticker := time.Tick(time.Second / 50)
        for _ = range ticker {
            if !m.Dead && config.Conf.GameActive {
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
    }()
}
func (m *Monster) Talk(text string) bool {
    config.TriggerEvent(&config.EventMonsterTalk{Event: config.Event{EType: config.EventTypeMonsterTalk}, Monster: m, Text: text })
    return true
}
func (p *Monster) Dir() (float32, float32) {
    return p.dir.x, p.dir.y
}

func (m *Monster) SetDir(x, y float32) {
    rect := m.Sprite.GetTextureRect()
    if x == 1 {
        rect.Top = 90 * 1
        rect.Width = -90
        rect.Left = 90
    } else if x == -1 {
        rect.Top = 90 * 1
        rect.Width = 90
        rect.Left = 0
    } else if y == 1 {
        rect.Top = 90 * 2
        rect.Width = 90
        rect.Left = 0
    } else if y == -1 {
        rect.Top = 0
        rect.Width = 90
        rect.Left = 0
    }
    config.TriggerEvent(&config.EventMonsterChangedDirection{Event: config.Event{EType: config.EventTypeMonsterChangedDirection}, Monster: m, NewDirX: x, NewDirY: y})

    m.Sprite.SetTextureRect(rect)
    m.dir.x = x
    m.dir.y = y
}


const MONSTERFEET = MONSTERWIDTH
func (m *Monster) Collides(x, y float32) bool {
    bounds := m.Sprite.GetGlobalBounds()
    bounds.Left += x
    bounds.Top += y + MONSTERHEIGHT
    worldmap := config.Conf.CurrentMap
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


