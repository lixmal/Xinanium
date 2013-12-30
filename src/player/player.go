package player

import (
    sf "bitbucket.org/krepa098/gosfml2"
    wm "../worldmap"
    "../animation"
    "../config"
    "../renderer"
)


const PLAYERWIDTH = 31
const PLAYERHEIGHT = 32
//const PLAYERWIDTH = 27
//const PLAYERHEIGHT = 59



type Inventory struct {
}

type Item struct {
}

type Dir struct {
    x, y float32
}

type Player struct {
    Handle string
    Name string
    Sprite *sf.Sprite
    animation.Animation
    Inventory *Inventory
    Item *Item
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
    entityType uint16
    centric bool
}


func New(name string, handle string, centric bool) *Player {
    sprite := sf.NewSprite(config.Conf.Rm.Texture(config.SPRITEDIR + "player2.png"))
    player := &Player{
        Health: 100,
		JumpHeight: 10,
		Speed: 1000,
		Name: name,
		Handle: handle,
		Sprite: sprite,
        entityType: config.LivingEntityPlayer,
        dir: &Dir{0, 1},
        centric: centric,
        Animation: animation.Animation{Sprite: sprite, Stopper: make(chan bool, 1)},
    }
    config.TriggerEvent(&config.EventPlayerNew{Event: config.Event{EType: config.EventTypePlayerNew}, Player: player })

    player.Sprite.SetTextureRect(sf.IntRect{0, 0, PLAYERWIDTH, PLAYERHEIGHT})
//    for i, v := range []string{"N", "NW", "W", "SW", "S", "SE", "E", "NE"} {
    for i, v := range []string{"S", "W", "E", "N"} {
        player.AddAnimation(v, i * 58)
    }

    player.dir.y = 1

    config.Players[handle] = player
    return player
}


func (p *Player) Talk(text string) bool {
    config.TriggerEvent(&config.EventPlayerTalk{Event: config.Event{EType: config.EventTypePlayerTalk}, Player: p, Text: text })
    return true
}
func (p *Player) GetSprite() *sf.Sprite {
    return p.Sprite
}

func (p *Player) Type() uint16 {
    return p.entityType
}

func (p *Player) Move(x, y float32) bool {
    elapsed := float32(renderer.Elapsed)
    newCoords := sf.Vector2f{ x * float32(p.Speed) * elapsed, y * float32(p.Speed) * elapsed }

    // make sure player can actually move
    if p.Speed > 0 && !p.Dead {

        // handle new direction
        if p.dir.x != x || p.dir.y != y {
            p.SetDir(x, y)
        }
        if !p.Collides(newCoords.X, newCoords.Y) {
            config.TriggerEvent(&config.EventPlayerMove{
                Event: config.Event{EType: config.EventTypePlayerMove},
                Player: p,
                NewX: newCoords.X,
                NewY: newCoords.Y,
            })

            // move the sprite
            p.Sprite.Move(newCoords)

            // handle walk animation
            p.FrameCounter++
            if p.FrameCounter >= 3 {
                p.NextFrame()
                p.FrameCounter = 0
            }


            // scroll view
            if config.Conf.Scrolling && p.centric {
                view := config.Conf.Window.GetView()
                view.SetCenter(p.Sprite.GetPosition())
                config.Conf.Window.SetView(view)
            }

            return true
        }
    }
    return false
}

func (p *Player) Jump() bool {
    if !p.InAir && p.JumpHeight > 0 && !p.Dead && !p.Floating {
        config.TriggerEvent(&config.EventPlayerJump{Event: config.Event{EType: config.EventTypePlayerJump}, Player: p })
        p.InAir = true
        go func() {
            println(p.Name, " is jumping")
        }()
        p.InAir = false
        return true
    }
    return false
}

func (p *Player) Hurt(damage int16, damager config.LivingEntity) int16 {
    e := &config.EventPlayerHurt{Event: config.Event{EType: config.EventTypePlayerHurt}, Player: p, Damage: damage, Damager: damager }
    config.TriggerEvent(e)

    if !p.Invincible {
       /* buff, err := sf.NewSoundBufferFromFile("resources/sound/hit.ogg")
        sound := sf.NewSound(buff)
        if err != nil {
        }
        sound.Play()
        */
        health := p.Health
        health -= damage
        if health <= 0 || p.Health < health {
            p.Health = 0
            p.Kill(damager, e)
        } else {
            p.Health = health
        }
    }
    return p.Health
}

func (p *Player) Kill(killer config.LivingEntity, hurtEvent *config.EventPlayerHurt) bool {
    config.TriggerEvent(&config.EventPlayerKilled{Event: config.Event{EType: config.EventTypePlayerKilled}, Player: p, Killer: killer, HurtEvent: hurtEvent })
    p.Dead = true
    return true
}

func (p *Player) Remove() bool {
    if _, ok := config.Players[p.Handle]; !ok {
        return false
    }
    config.TriggerEvent(&config.EventPlayerRemoved{Event: config.Event{EType: config.EventTypePlayerRemoved}, Player: p })
    delete(config.Players, p.Handle)
    return true
}

func (p *Player) SetPosition(x, y float32) {
    config.TriggerEvent(&config.EventPlayerChangedPosition{Event: config.Event{EType: config.EventTypePlayerChangedPosition}, Player: p, NewX: x, NewY: y })
    p.Sprite.SetPosition(sf.Vector2f{x, y})
}

func (p *Player) Position() (float32, float32) {
    pos := p.Sprite.GetPosition()
    return pos.X, pos.Y
}

func (p *Player) Dir() (float32, float32) {
    return p.dir.x, p.dir.y
}
func (p *Player) SetDir(x, y float32) {
    rect := p.Sprite.GetTextureRect()
/*    if x == 1 && y == 1 {
        rect.Top = PLAYERHEIGHT * 5
    } else if x == -1 && y == 1 {
        rect.Top = PLAYERHEIGHT * 3
    } else if x == 1 && y == -1 {
        rect.Top = PLAYERHEIGHT * 7
    } else if x == -1 && y == -1 {
        rect.Top = PLAYERHEIGHT * 1
    } else if x == 1 {
        rect.Top = PLAYERHEIGHT * 6
    } else if x == -1 {
        rect.Top = PLAYERHEIGHT * 2
    } else if y == 1 {
        rect.Top = PLAYERHEIGHT * 4
    } else if y == -1 {
        rect.Top = 0
    }
    */
    if x == 1 {
        rect.Top = PLAYERHEIGHT * 2
    } else if x == -1 {
        rect.Top = PLAYERHEIGHT * 1
    } else if y == 1 {
        rect.Top = 0
    } else if y == -1 {
        rect.Top = PLAYERHEIGHT * 3
    }
    config.TriggerEvent(&config.EventPlayerChangedDirection{Event: config.Event{EType: config.EventTypePlayerChangedDirection}, Player: p, NewDirX: x, NewDirY: y})

    p.Sprite.SetTextureRect(rect)
    p.dir.x = x
    p.dir.y = y
}



const PLAYERFEET = PLAYERWIDTH
func (p *Player) Collides(x,y float32) bool {
    bounds := p.Sprite.GetGlobalBounds()
    bounds.Left += x
    bounds.Top += y + PLAYERHEIGHT
    worldmap := config.Conf.CurrentMap

    var success bool
    var what config.LivingEntity
    if bounds.Left < 0 || bounds.Left + PLAYERWIDTH > float32(worldmap.Width * wm.TILEWIDTH) || bounds.Top - PLAYERFEET < 0 || bounds.Top > float32(worldmap.Height * wm.TILEHEIGHT) {
        success = true
    }
    // TODO: Scan for nearby entities
    bounds.Height -= PLAYERFEET + PLAYERHEIGHT
    if !success {
        for _, player := range config.Players {
            if (player != p) {
                if collision, _ := bounds.Intersects(player.GetSprite().GetGlobalBounds()); collision {
                    success = true
                    what = player
                    break
                }
            }
        }
    }
    if !success {
        for _, mon := range config.Monsters {
            if collision, _ := bounds.Intersects(mon.GetSprite().GetGlobalBounds()); collision {
                success = true
                what = mon
                break
            }
        }
    }
    if success {
        config.TriggerEvent(
            &config.EventPlayerCollision{
                Player: p,
                Event: config.Event{EType: config.EventTypePlayerCollision},
                What: what,
                X: bounds.Left + x,
                Y: bounds.Top + y,
            },
        )
    }
    return success
}
