package render

import (
	"image/color"
)

type Info struct {
	Color    color.NRGBA
	Emission uint8
}

var BlockInfo = map[uint16]Info{
	0:   Info{color.NRGBA{0, 0, 0, 0}, 0},           // air
	1:   Info{color.NRGBA{128, 128, 128, 255}, 0},   // stone
	2:   Info{color.NRGBA{108, 169, 66, 255}, 0},    // grass
	3:   Info{color.NRGBA{122, 84, 55, 255}, 0},     // dirt
	4:   Info{color.NRGBA{123, 123, 123, 255}, 0},   // cobblestone
	5:   Info{color.NRGBA{160, 133, 76, 255}, 0},    // oak wood plank
	7:   Info{color.NRGBA{48, 48, 48, 255}, 0},      // bedrock
	9:   Info{color.NRGBA{31, 69, 141, 128}, 0},     // stationary water
	11:  Info{color.NRGBA{213, 140, 33, 255}, 192},  // stationary lava
	13:  Info{color.NRGBA{138, 120, 120, 255}, 0},   // gravel
	17:  Info{color.NRGBA{78, 61, 34, 255}, 0},      // oak wood
	18:  Info{color.NRGBA{55, 148, 11, 255}, 0},     // oak leaves
	20:  Info{color.NRGBA{180, 215, 220, 128}, 0},   // glass
	22:  Info{color.NRGBA{27, 65, 156, 255}, 0},     // lapis lazuli
	24:  Info{color.NRGBA{227, 220, 169, 255}, 0},   // sandstone
	26:  Info{color.NRGBA{162, 34, 31, 255}, 0},     // bed block
	31:  Info{color.NRGBA{91, 41, 1, 255}, 0},       // dead shrub
	35:  Info{color.NRGBA{215, 215, 215, 255}, 0},   // white wool
	37:  Info{color.NRGBA{241, 249, 0, 255}, 0},     // dandelion
	38:  Info{color.NRGBA{210, 1, 2, 255}, 0},       // poppy
	39:  Info{color.NRGBA{146, 109, 84, 255}, 0},    // brown mushroom
	40:  Info{color.NRGBA{227, 10, 10, 255}, 0},     // red mushroom
	42:  Info{color.NRGBA{238, 238, 238, 255}, 0},   // iron block
	44:  Info{color.NRGBA{178, 178, 177, 255}, 0},   // stone slab
	45:  Info{color.NRGBA{125, 67, 51, 255}, 0},     // brick
	47:  Info{color.NRGBA{86, 105, 5, 255}, 0},      // bookshelf
	49:  Info{color.NRGBA{40, 31, 62, 255}, 0},      // obsidian
	50:  Info{color.NRGBA{255, 255, 152, 255}, 192}, // torch
	53:  Info{color.NRGBA{189, 153, 98, 255}, 0},    // oak wood stairs
	54:  Info{color.NRGBA{166, 115, 34, 255}, 0},    // chest
	56:  Info{color.NRGBA{120, 207, 251, 255}, 0},   // diamond ore
	57:  Info{color.NRGBA{136, 230, 226, 255}, 0},   // diamond block
	61:  Info{color.NRGBA{122, 122, 122, 255}, 0},   // furnace
	65:  Info{color.NRGBA{143, 116, 57, 255}, 0},    // ladder
	67:  Info{color.NRGBA{128, 128, 128, 255}, 0},   // cobblestone stairs
	69:  Info{color.NRGBA{125, 98, 60, 255}, 0},     // lever
	70:  Info{color.NRGBA{117, 117, 117, 255}, 0},   // stone pressure plate
	71:  Info{color.NRGBA{206, 206, 206, 255}, 0},   // iron door block
	72:  Info{color.NRGBA{189, 153, 98, 255}, 0},    // wooden pressure plate
	76:  Info{color.NRGBA{255, 217, 0, 255}, 64},    // redstone torch (on)
	77:  Info{color.NRGBA{77, 77, 77, 255}, 0},      // stone button
	79:  Info{color.NRGBA{143, 192, 255, 192}, 0},   // ice
	80:  Info{color.NRGBA{239, 255, 255, 255}, 0},   // snow block
	82:  Info{color.NRGBA{163, 170, 180, 255}, 0},   // clay
	85:  Info{color.NRGBA{140, 114, 71, 255}, 0},    // fence
	89:  Info{color.NRGBA{255, 189, 94, 255}, 192},  // glowstone
	93:  Info{color.NRGBA{147, 147, 147, 255}, 0},   // redstone repeater (off)
	96:  Info{color.NRGBA{146, 109, 53, 255}, 0},    // trapdoor
	98:  Info{color.NRGBA{128, 128, 128, 255}, 0},   // stone brick
	101: Info{color.NRGBA{98, 90, 83, 255}, 0},      // iron bars
	102: Info{color.NRGBA{180, 215, 220, 128}, 0},   // glass pane
	109: Info{color.NRGBA{128, 128, 128, 255}, 0},   // stone brick stairs
	111: Info{color.NRGBA{4, 95, 12, 255}, 0},       // lily pad
	113: Info{color.NRGBA{43, 14, 19, 255}, 0},      // nether brick fence
	114: Info{color.NRGBA{60, 24, 31, 255}, 0},      // nether brick stairs
	116: Info{color.NRGBA{235, 235, 235, 255}, 0},   // enchantment table
	125: Info{color.NRGBA{181, 145, 90, 255}, 0},    // double oak wood slab
	126: Info{color.NRGBA{181, 145, 90, 255}, 0},    // oak wood slab
	128: Info{color.NRGBA{197, 191, 142, 255}, 0},   // sandstone stairs
	130: Info{color.NRGBA{41, 61, 64, 255}, 0},      // ender chest
	133: Info{color.NRGBA{84, 217, 124, 255}, 0},    // emerald block
	135: Info{color.NRGBA{202, 191, 133, 255}, 0},   // birch wood stairs
	137: Info{color.NRGBA{174, 170, 166, 255}, 0},   // command block
	138: Info{color.NRGBA{127, 215, 211, 192}, 0},   // beacon block
	139: Info{color.NRGBA{115, 115, 115, 255}, 0},   // cobblestone wall
	140: Info{color.NRGBA{122, 65, 49, 255}, 0},     // flower pot
	144: Info{color.NRGBA{141, 141, 141, 255}, 0},   // mob head
}
