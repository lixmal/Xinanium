package player

import (
//	"../animation"
	"../config"
	"../event"
	"../network"
//	"../renderer"
    rm "../resourcemanager"
//	wm "../worldmap"
	"azul3d.org/gfx.v1"
    "azul3d.org/lmath.v1"
	"time"
    "image"
)

const PLAYERWIDTH           = 31
const PLAYERHEIGHT          = 32
const SPRITEAMNT_VERTICAL   = 3
const SPRITEAMNT_HORIZONTAL = 4
const PLAYEREXTENSION       = ".png"

//const PLAYERWIDTH = 27
//const PLAYERHEIGHT = 59

type Inventory struct {
}

type Item struct {
}

type Dir struct {
	x, y float64
}

type Player struct {
	Handle string
	Name   string
	sprite *gfx.Object
//	animation.Animation
	Inventory  *Inventory
	Item       *Item
	Speed      uint16
	JumpHeight byte
	Health     int16
	Invincible bool
	Invisible  bool
	Walking    bool
	InAir      bool
	Dead       bool
	Floating   bool
	dir        *Dir
	entityType uint16
	centric    bool
}

func New(name string, handle string, centric bool) *Player {
	var sprite *gfx.Object

    // generate one mesh per player, for animations
    mesh := rm.Mesh(
        PLAYERWIDTH,
        PLAYERHEIGHT,
        1.0 / SPRITEAMNT_VERTICAL,
        1.0 / SPRITEAMNT_HORIZONTAL,
    )
	if handle == "vik" {
		network.Send(config.GET_PLAYER_TEX)
		var playertex []byte
		network.Read(&playertex)

        sprite = rm.SpriteFromMemory(&playertex, mesh)
	} else {
        sprite = rm.Sprite(config.SPRITEDIR + "player1" + PLAYEREXTENSION, mesh)
	}

	player := &Player{
		Health:     100,
		JumpHeight: 10,
		Speed:      280,
		Name:       name,
		Handle:     handle,
		sprite:     sprite,
		entityType: config.LivingEntityPlayer,
		dir:        &Dir{0, 1},
		centric:    centric,
//		Animation:  animation.Animation{Sprite: sprite, Stopper: make(chan bool, 1)},
	}
    /*
	if !event.Trigger(&event.PlayerNew{Event: event.New(event.TypePlayerNew), Player: player}) {
		 player.sprite.SetTextureRect(sf.IntRect{0, 0, PLAYERWIDTH, PLAYERHEIGHT})
		//    for i, v := range []string{"N", "NW", "W", "SW", "S", "SE", "E", "NE"} {
		for i, v := range []string{"S", "W", "E", "N"} {
			player.AddAnimation(v, i*58)
		}

		player.dir.y = 1
		config.Players[handle] = player

		return player
	}
    */
    /* delte the following if uncommenting above*/
    player.dir.y = 1
    config.Players[handle] = player
    return player
    /* end */
	return nil
}

func (p *Player) Talk(text string) bool {
	if !event.Trigger(&event.PlayerTalk{Event: event.New(event.TypePlayerTalk), Player: p, Text: text}) {
		return true
	}
	return false
}

func (p *Player) GetSprite() *gfx.Object {
	return p.sprite
}

func (p *Player) Type() uint16 {
	return p.entityType
}

func (p *Player) Move(x, z float64) bool {

	//println("playermove")
	//elapsed := float32(renderer.Elapsed)
	//newCoords := sf.Vector2f{x * float32(p.Speed) * elapsed, y * float32(p.Speed) * elapsed}

	// make sure player can actually move
	if p.Speed > 0 && !p.Dead {
        v := lmath.Vec3{x, 0, z}.MulScalar(config.Conf.Renderer.Clock().Dt() * float64(p.Speed))
        //newCoords := lmath.Vec3{pos.X * v.X, pos.Y, pos.Z * v.Y}
        newCoords := p.sprite.Pos().Add(v)

		if p.Handle == "vik" {
			network.Send(config.PLAYER_MOVE)
			network.Send(config.Dir{x, z})
		}

		// handle new direction
		if p.dir.x != x || p.dir.y != z {
			p.SetDir(x, z)
		}
//		if !p.Collides(newCoords.X, newCoords.Y) {
			//println("playermove event call")
            /*
			if !event.Trigger(&event.PlayerMove{
				Event:  event.New(event.TypePlayerMove),
				Player: p,
				NewX:   newCoords.X,
				NewY:   newCoords.Y,
			}) {
            */

				// move the sprite
                p.sprite.SetPos(newCoords)

                /*
				// handle walk animation
				p.FrameCounter++
				if p.FrameCounter >= 3 {
					p.NextFrame()
					p.FrameCounter = 0
				}
                */

                /*
				// scroll view
				if config.Conf.Scrolling && p.centric {
					view := config.Conf.Window.GetView()
					view.SetCenter(p.sprite.GetPosition())
					config.Conf.Window.SetView(view)
				}
                */

				return true
		//	}
//		}
	}
	return false
}

func (p *Player) Jump() bool {
	if !p.InAir && p.JumpHeight > 0 && !p.Dead && !p.Floating && !event.Trigger(&event.PlayerJump{Event: event.New(event.TypePlayerJump), Player: p}) {
		p.InAir = true
		go func() {
			//println(p.Name, " is jumping")
		}()
		p.InAir = false
		return true
	}
	return false
}

func (p *Player) Hurt(damage int16, damager config.LivingEntity) int16 {
	e := &event.PlayerHurt{Event: event.New(event.TypePlayerHurt), Player: p, Damage: damage, Damager: damager}
	if !p.Invincible && !event.Trigger(e) {
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
		p.Invincible = true
		go func() {
			time.Sleep(time.Second)
			p.Invincible = false
		}()
	}
	return p.Health
}

func (p *Player) Kill(killer config.LivingEntity, hurtEvent *event.PlayerHurt) bool {
	if !event.Trigger(&event.PlayerKilled{Event: event.New(event.TypePlayerKilled), Player: p, Killer: killer, HurtEvent: hurtEvent}) {
		p.Dead = true
		return true
		// TODO Set dead animation/sprite
	}
	return false
}

func (p *Player) Remove() bool {
	if _, ok := config.Players[p.Handle]; ok {
		if !event.Trigger(&event.PlayerRemoved{Event: event.New(event.TypePlayerRemoved), Player: p}) {
			delete(config.Players, p.Handle)
			return true
		}
	}
	return false
}

func (p *Player) SetPosition(x, y float64) bool {
//	if !event.Trigger(&event.PlayerChangedPosition{Event: event.New(event.TypePlayerChangedPosition), Player: p, NewX: x, NewY: y}) {
//		p.sprite.SetPosition(sf.Vector2f{x, y})
		return true
//	}
	return false
}

func (p *Player) Position() (float64, float64) {
//	pos := p.sprite.GetPosition()
    pos := struct{X, Y float64}{0, 0}
	return pos.X, pos.Y
}

func (p *Player) Dir() (float64, float64) {
	return p.dir.x, p.dir.y
}

func (p *Player) SetDir(x, y float64) bool {
    /*^
	if x == 1 && y == 1 {
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
    var top float32
	if x == 1 {
		top = 2
	} else if x == -1 {
		top = 1
	} else if y == 1 {
		top = 0
	} else if y == -1 {
		top = 3
	}
    p.sprite.Textures[0].Bounds = image.Rect(0, int(PLAYERHEIGHT * top), PLAYERWIDTH, PLAYERHEIGHT)
   //  fmt.Println(p.sprite.Textures[0].Bounds)
   // fmt.Printf("%+#v\n", p.sprite.Meshes[0])
    /*
    p.sprite.Textures[0] = rm.TexCoords(
        PLAYERWIDTH,
        PLAYERHEIGHT * top,
        1.0 / SPRITEAMNT_VERTICAL,
        1.0 / SPRITEAMNT_HORIZONTAL,
    )
    */

    return true
    /*
	if !event.Trigger(&event.PlayerChangedDirection{Event: event.New(event.TypePlayerChangedDirection), Player: p, NewDirX: x, NewDirY: y}) {
		p.sprite.SetTextureRect(rect)
		p.dir.x = x
		p.dir.y = y
		return true
	}
    */
	return false
}

const PLAYERFEET = PLAYERWIDTH

func (p *Player) Collides(x, y float64) bool {
    /*
	//println("player collides")
	bounds := p.sprite.GetGlobalBounds()
	bounds.Left += x
	bounds.Top += y + PLAYERHEIGHT
	worldmap := wm.Current

	var success bool
	var what config.LivingEntity
	if bounds.Left < 0 || bounds.Left+PLAYERWIDTH > float32(worldmap.Width*wm.TILEWIDTH) || bounds.Top-PLAYERFEET < 0 || bounds.Top > float32(worldmap.Height*wm.TILEHEIGHT) {
		success = true
	}
	// TODO: Scan for nearby entities
	bounds.Height -= PLAYERFEET + PLAYERHEIGHT
	if !success {
		for _, player := range config.Players {
			if player != p {
				if collision, _ := bounds.Intersects(player.sprite.GetGlobalBounds()); collision {
					success = true
					what = player
					break
				}
			}
		}
	}
	if !success {
		for _, mon := range config.Monsters {
			if collision, _ := bounds.Intersects(mon.sprite.GetGlobalBounds()); collision {
				success = true
				what = mon
				break
			}
		}
	}
	//println("player collides event call")
	if success && !event.Trigger(
		&event.PlayerCollision{
			Player: p,
			Event:  event.New(event.TypePlayerCollision),
			What:   what,
			X:      bounds.Left + x,
			Y:      bounds.Top + y,
		},
	) {
		return true
	}
    */
	return false
}
