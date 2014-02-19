package animation

import (
	sf "bitbucket.org/krepa098/gosfml2"
    "time"
)

const PLAYERWIDTH = 31
const PLAYERHEIGHT = 32

type Animation struct {
    anims map[string][]*sf.IntRect
    Stopper chan bool
    Sprite *sf.Sprite
    FrameCounter uint8
//    top int
}

func (a *Animation) AddAnimation(name string, top int) {
    if a.anims == nil {
        a.anims = make(map[string][]*sf.IntRect)
    }
    size := a.Sprite.GetTexture().GetSize()
    for i := 0; i < int(size.X); i += PLAYERWIDTH {
        a.anims[name] = append(a.anims[name], &sf.IntRect{i, top, PLAYERWIDTH, PLAYERHEIGHT})
    }
}

func (a *Animation) LoopAnimation(index string, dur uint8) {
    ticker := time.Tick(time.Second / time.Duration(dur))
    go func() {
        for {
            for _, v := range a.anims[index] {
                select {
                case <-a.Stopper: println("got true"); return
                default:
                }
                <-ticker
                a.Sprite.SetTextureRect(*v)
            }
        }
    }()
}

func (a *Animation) StopAnimation() {
    select {
    case a.Stopper<- true:
    default:
    }
}

// TODO: Refactor
func (a *Animation) PlayAnimation(index string, dur uint8) {
    go func() {
        ticker := time.Tick(time.Second / time.Duration(dur))
        for _, v := range a.anims[index]{
            <-ticker
            a.Sprite.SetTextureRect(*v)
        }
    }()
}

func (a *Animation) NextFrame() {
    rect := a.Sprite.GetTextureRect()
    defer func() { a.Sprite.SetTextureRect(rect) }()
    if (rect.Left >= 62) {
        rect.Left = 0
        return
    }
    rect.Left += PLAYERWIDTH
}

func (a *Animation) PreviousFrame() {
    rect := a.Sprite.GetTextureRect()
    defer func() { a.Sprite.SetTextureRect(rect) }()
    if (rect.Left == 0) {
        rect.Left = 62
        return
    }
    rect.Left -= PLAYERWIDTH
}
