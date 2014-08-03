package worldmap

import (
	"../config"
	"../renderer"
	sf "bitbucket.org/krepa098/gosfml2"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	MAPDIR         = filepath.Dir(filepath.FromSlash("resources/maps")) + "/"
	MAPEXTENSION   = ".dat"
	TILEWIDTH      = 40
	TILEHEIGHT     = 40
	MAPHEADERSIZE  = 1024
	MAPTITLELENGTH = 256
)

type WorldMap struct {
	Tiles   [][]byte
	Name    string
	Width   float32
	Height  float32
	Version float32
	Script  string
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
		Tiles:   tiles,
		Width:   width,
		Height:  height,
		Version: 0.1,
		Script:  string(script),
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

func ReadOpen(filename string) *WorldMap {
	worldmap := Read(filename)
	Current = worldmap
	return worldmap
}

func init() {
	//  Write(MAPDIR + "/gobmap.dat", 40, 40)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var Current *WorldMap
