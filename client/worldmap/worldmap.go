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
    "azul3d.org/lmath.v1"
//    "../renderer"
)
const TILEWIDTH      = 40
const TILEHEIGHT     = 40
const MAPHEADERSIZE  = 1024
const MAPTITLELENGTH = 256
const TILEDIR        = config.RESOURCESDIR + "textures/tiles/"
const TILEEXTENSION  = ".png"
const MAPDIR         = config.RESOURCESDIR + "maps/"
const MAPEXTENSION   = ".dat"

// one mesh for all maps
var   mapMesh        = rm.Mesh(40, 40, 1.0, 1.0)

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
    file, err := os.Open(MAPDIR + filename + MAPEXTENSION)
    checkErr(err)
    defer file.Close()

    decoder := gob.NewDecoder(file)
    var wmap WorldMap
    checkErr(decoder.Decode(&wmap))

    return &wmap
}

func Open(worldmap *WorldMap) {
    var cnt = 0
    for x, v := range worldmap.Tiles {
        for y := range v {
            tileType := worldmap.Tiles[x][y]
            var sprite *gfx.Object
            var tileName string
            switch tileType {
                case 0:
                    tileName = "grass"
                case 1:
                    tileName = "dirt"
                case 2:
                    tileName = "water"
                case 3:
                    tileName = "w_br"
                default:
                    log.Fatal("Unknown tile type: " + fmt.Sprintf("%i", tileType))
            }
            sprite = rm.Sprite(TILEDIR + tileName + TILEEXTENSION, mapMesh)
            // add to drawing
            sprite.SetPos(lmath.Vec3{float64(x * TILEWIDTH), 0, float64(y * TILEHEIGHT)})
            config.ToDraw = append(config.ToDraw, sprite)
            cnt++
        }
    }
    log.Println(cnt)
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
