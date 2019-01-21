package worldmap

import(
    "log"
    "encoding/gob"
    "os"
    "io/ioutil"
    "fmt"
    "../config"
    rm "../resourcemanager"
    "azul3d.org/gfx.v1"
//    "../renderer"
)
const TILEWIDTH = 40
const TILEHEIGHT = 40
const MAPHEADERSIZE = 1024
const MAPTITLELENGTH = 256

type WorldMap struct {
    Tiles [][]byte
    Name string
    Width float32
    Height float32
    Version float32
    Script string
    Objects []config.Entity
}

func (w *WorldMap) RunScripts() bool {
    state := config.Lua.State
    /*state.LoadFile(w.Script)
    state.NewTable
    state.SetField(-1, )
    // set env for map first; push env, inherit global env, set 'self, call
    state.Setfenv*/
    if state.DoString(w.Script) != nil {
        return false
    }
    return true
}

func Write(filename string, width, height float32) *WorldMap {
    file, err := os.Create(filename)
    checkErr(err)
    defer file.Close()

    tiles := make([][]byte, int32(height))
    for i := range tiles {
        tiles[i] = make([]byte, int32(width))
    }
    encoder := gob.NewEncoder(file)
    script, err := ioutil.ReadFile("test.lua")
    checkErr(err)
    wmap := &WorldMap{
        Tiles: tiles,
        Width: width,
        Height: height,
        Version: 0.1,
        Script: string(script),
    }

    checkErr(encoder.Encode(wmap))
    return wmap
}

func Read(filename string) *WorldMap {
    file, err := os.Open(filename)
    checkErr(err)
    defer file.Close()

    decoder := gob.NewDecoder(file)
    var wmap WorldMap
    checkErr(decoder.Decode(&wmap))

    return &wmap
}

func Open(worldmap *WorldMap) {
    for x, v := range worldmap.Tiles {
        for y := range v {
            tileType := worldmap.Tiles[x][y]
            var sprite *gfx.Object
            switch tileType {
                case 0:
                    sprite = rm.Sprite(config.RESOURCESDIR + "textures/tiles/grass.png")
                case 1:
                    sprite = rm.Sprite(config.RESOURCESDIR + "textures/tiles/dirt.png")
                case 2:
                    sprite = rm.Sprite(config.RESOURCESDIR + "textures/tiles/water.png")
                case 3:
                    sprite = rm.Sprite(config.RESOURCESDIR + "textures/tiles/w_br.png")
            }
            if sprite != nil {
                // add to drawing
            //    sprite.SetPosition(sf.Vector2f{float32(x * TILEWIDTH), float32(y * TILEHEIGHT)})
            //    renderer.ToDraw = append(renderer.ToDraw, sprite)
            } else {
                log.Fatal("Unknown tile type: " + fmt.Sprintf("%i", tileType))
            }
        }
    }
    Current = worldmap
}

func ReadOpen(filename string) *WorldMap {
    worldmap := Read(filename)
    Open(worldmap)
    return worldmap
}

func init() {
  //  Write("resources/maps/gobmap.dat", 40, 40)
}


func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

var Current *WorldMap
