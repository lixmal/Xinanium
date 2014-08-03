package player

import (
	"../animation"
	"../config"
	"../event"
	"../network"
	"../renderer"
	wm "../worldmap"
	sf "bitbucket.org/krepa098/gosfml2"
	"log"
	"time"
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
	Name   string
	sprite *sf.Sprite
	animation.Animation
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
	var sprite *sf.Sprite
	var err error
	if handle == "vik" {
		network.Send(config.GET_PLAYER_TEX)
		var playertex []byte
		network.Read(&playertex)
		tex, err := sf.NewTextureFromMemory(playertex, nil)
		if err != nil {
			log.Println(err)
			// returns nil???
			//    return nil
		}
		sprite, err = sf.NewSprite(tex)
	} else {
		sprite, err = sf.NewSprite(config.Conf.Rm.Texture(config.SPRITEDIR + "player1.png"))
	}

	_ = err
	player := &Player{
		Health:     100,
		JumpHeight: 10,
		Speed:      1000,
		Name:       name,
		Handle:     handle,
		sprite:     sprite,
		entityType: config.LivingEntityPlayer,
		dir:        &Dir{0, 1},
		centric:    centric,
		Animation:  animation.Animation{Sprite: sprite, Stopper: make(chan bool, 1)},
	}
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
	return nil
}

func (p *Player) Talk(text string) bool {
	if !event.Trigger(&event.PlayerTalk{Event: event.New(event.TypePlayerTalk), Player: p, Text: text}) {
		return true
	}
	return false
}

func (p *Player) GetSprite() *sf.Sprite {
	return p.sprite
}

func (p *Player) Type() uint16 {
	return p.entityType
}

func (p *Player) Move(x, y float32) bool {
	//println("playermove")
	elapsed := float32(renderer.Elapsed)
	newCoords := sf.Vector2f{x * float32(p.Speed) * elapsed, y * float32(p.Speed) * elapsed}

	// make sure player can actually move
	if p.Speed > 0 && !p.Dead {
		if p.Handle == "vik" {
			network.Send(config.PLAYER_MOVE)
			network.Send(config.Dir{x, y})
		}

		// handle new direction
		if p.dir.x != x || p.dir.y != y {
			p.SetDir(x, y)
		}
		if !p.Collides(newCoords.X, newCoords.Y) {
			//println("playermove event call")
			if !event.Trigger(&event.PlayerMove{
				Event:  event.New(event.TypePlayerMove),
				Player: p,
				NewX:   newCoords.X,
				NewY:   newCoords.Y,
			}) {

				// move the sprite
				p.sprite.Move(newCoords)

				// handle walk animation
				p.FrameCounter++
				if p.FrameCounter >= 3 {
					p.NextFrame()
					p.FrameCounter = 0
				}

				// scroll view
				if config.Conf.Scrolling && p.centric {
					view := config.Conf.Window.GetView()
					view.SetCenter(p.sprite.GetPosition())
					config.Conf.Window.SetView(view)
				}

				return true
			}
		}
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

func (p *Player) SetPosition(x, y float32) bool {
	if !event.Trigger(&event.PlayerChangedPosition{Event: event.New(event.TypePlayerChangedPosition), Player: p, NewX: x, NewY: y}) {
		p.sprite.SetPosition(sf.Vector2f{x, y})
		return true
	}
	return false
}

func (p *Player) Position() (float32, float32) {
	pos := p.sprite.GetPosition()
	return pos.X, pos.Y
}

func (p *Player) Dir() (float32, float32) {
	return p.dir.x, p.dir.y
}

func (p *Player) SetDir(x, y float32) bool {
	rect := p.sprite.GetTextureRect()
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
	if !event.Trigger(&event.PlayerChangedDirection{Event: event.New(event.TypePlayerChangedDirection), Player: p, NewDirX: x, NewDirY: y}) {
		p.sprite.SetTextureRect(rect)
		p.dir.x = x
		p.dir.y = y
		return true
	}
	return false
}

const PLAYERFEET = PLAYERWIDTH

func (p *Player) Collides(x, y float32) bool {
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
	return false
}
