package world

import (
	"github.com/Mischanix/anvil-go"
)

type blockArray [16][16][16]uint16

type World struct {
	region *anvil.Level

	loaded map[struct{ x, y, z int }]*blockArray
}

func New(path string) *World {
	var w World
	w.region = anvil.New(path)
	w.loaded = make(map[struct{ x, y, z int }]*blockArray)
	return &w
}

func (w *World) Get(bx, by, bz int) uint16 {
	cx := bx >> 4
	cy := by >> 4
	cz := bz >> 4

	if cy < 0 || cy > 15 {
		return 0
	}

	offx := bx - cx<<4
	offy := by - cy<<4
	offz := bz - cz<<4

	if secData, ok := w.loaded[struct{ x, y, z int }{cx, cy, cz}]; ok {
		if secData == nil {
			return 0
		} else {
			return secData[offx][offy][offz]
		}
	}

	// load all sections of the chunk

	// set to nil to denote attempted, we re-set these if needed
	for y := 0; y <= 15; y++ {
		w.loaded[struct{ x, y, z int }{cx, y, cz}] = nil
	}

	chunk, err := w.region.Chunk(cx, cz)
	if err != nil {
		return 0
	}

	level := chunk.Data.At("Level")
	sections := level.Compound().At("Sections").List()

	for i := 0; i < int(sections.Length()); i++ {
		section := sections.At(int32(i)).Compound()
		sectionIndex := section.At("Y").Byte()

		blocksArray := section.At("Blocks").ByteArray()
		var addArray []byte
		addTag := section.At("Add")
		if addTag != nil {
			addArray = addTag.ByteArray()
		}

		var secData *blockArray = new(blockArray)

		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				for z := 0; z < 16; z++ {
					if addArray == nil {
						secData[x][y][z] = uint16(blocksArray[(y*16+z)*16+x])
					} else {
						pos := (y*16+z)*16 + x

						part := addArray[pos>>1]
						if pos%2 == 1 {
							part >>= 4
						}
						part &= 0x0F

						secData[x][y][z] = uint16(blocksArray[pos]) + uint16(part)<<8
					}
				}
			}
		}

		w.loaded[struct{ x, y, z int }{cx, int(sectionIndex), cz}] = secData
	}

	if secData, ok := w.loaded[struct{ x, y, z int }{cx, cy, cz}]; ok {
		if secData == nil {
			return 0
		} else {
			return secData[offx][offy][offz]
		}
	} else {
		panic("loaded a chunk but it wasn't loaded")
	}
}
