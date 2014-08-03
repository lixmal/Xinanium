package resourcemanager

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"log"
)

type TextureStruct struct {
	m map[string]*sf.Texture
}

/*
type SpriteStruct struct {
    m map[string]*sf.Sprite
}
*/
type ResourceManager struct {
	textures TextureStruct
	//    sprites SpriteStruct
}

func New() *ResourceManager {
	return &ResourceManager{
		textures: TextureStruct{m: map[string]*sf.Texture{}},
		//        sprites:  SpriteStruct{m: map[string]*sf.Sprite{}},
	}
}

func (rm *ResourceManager) Texture(filename string) *sf.Texture {
	tex, ok := rm.textures.m[filename]
	if !ok {
		var err error
		tex, err = sf.NewTextureFromFile(filename, &sf.IntRect{})
		if err != nil || tex == nil {
			log.Fatal("Could not load image '", filename, "':", err)
		}
		rm.textures.m[filename] = tex
	}

	return tex
}

/*func (rm *ResourceManager) Sprite(filename string) *sf.Sprite {
    spr, ok := rm.sprites.m[filename]
    if !ok {
        spr = sf.NewSprite(rm.Texture(filename))
        rm.sprites.m[filename] = spr
    }
    return spr
}
*/
