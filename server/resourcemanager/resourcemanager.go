package resourcemanager

import(
    sf "bitbucket.org/krepa098/gosfml2"
    "log"
    "sync"
)

type TextureStruct struct {
    sync.RWMutex
    m map[string]*sf.Texture
}
/*
type SpriteStruct struct {
    sync.RWMutex
    m map[string]*sf.Sprite
}
*/
type ResourceManager struct{
    textures TextureStruct
//    sprites SpriteStruct
}

func New() *ResourceManager {
    return &ResourceManager {
        textures: TextureStruct{m: map[string]*sf.Texture{}},
//        sprites:  SpriteStruct{m: map[string]*sf.Sprite{}},
    }
}

func (rm *ResourceManager) Texture(filename string) *sf.Texture {
    rm.textures.RLock()
    tex, ok := rm.textures.m[filename]
    rm.textures.RUnlock()
    if !ok {
        var err error
        tex, err = sf.NewTextureFromFile(filename, &sf.IntRect{})
        if err != nil || tex == nil {
            log.Fatal("Could not load image '", filename, "':", err)
        }
        rm.textures.Lock()
        rm.textures.m[filename] = tex
        rm.textures.Unlock()
    }

    return tex
}
/*func (rm *ResourceManager) Sprite(filename string) *sf.Sprite {
    rm.sprites.RLock()
    spr, ok := rm.sprites.m[filename]
    rm.sprites.RUnlock()
    if !ok {
        spr = sf.NewSprite(rm.Texture(filename))
        rm.sprites.Lock()
        rm.sprites.m[filename] = spr
        rm.sprites.Unlock()
    }
    return spr
}
*/




