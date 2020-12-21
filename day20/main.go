package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	reTileHeader *regexp.Regexp
)

type match struct {
	myEdge    int
	otherEdge int
	otherTile int
	flipped   bool
}

type tile struct {
	id      int
	data    [][]bool
	s       int
	matches []match
	x, y, r int
	flip    bool // true if vertical flip before rotate
	placed  bool
}

// Convenience function for flipping arrays
func flipped(l []bool) []bool {
	f := make([]bool, len(l))
	for i := range l {
		f[len(f)-1-i] = l[i]
	}
	return f
}

func equal(l1 []bool, l2 []bool) bool {
	if len(l1) != len(l2) {
		log.Fatal("Problem")
		return false
	}
	for i := range l1 {
		if l1[i] != l2[i] {
			return false
		}
	}
	return true
}

// Only used for the correlator
func (t *tile) getEdge(edgeNum int) []bool {
	switch edgeNum {
	case 0: // Top, left to right
		ret := make([]bool, t.s)
		for i := range t.data {
			ret[i] = t.data[0][i]
		}
		return ret
	case 1: // Right side, top to bottom
		ret := make([]bool, t.s)
		for i := range t.data {
			ret[i] = t.data[i][t.s-1]
		}
		return ret
	case 2: // Bottom, right to left
		ret := make([]bool, t.s)
		for i := range t.data[0] {
			// Bottom row, indexed right to left
			ret[i] = t.data[t.s-1][t.s-1-i]
		}
		return ret
	case 3: // Left, bottom to top
		ret := make([]bool, t.s)
		for i := range t.data {
			ret[i] = t.data[t.s-1-i][0]
		}
		return ret
	}

	log.Fatal("invalid edge number")
	return nil
}

func (t *tile) correlate(m map[int]*tile) {
	// Loop through all other tiles
	for tileNum, otherTile := range m {
		// No need to autocorrelate
		if tileNum == t.id {
			continue
		}
		// Check each one of our edges against the other's
		for edge := 0; edge < 4; edge++ {
			myData := t.getEdge(edge)
			for otherEdge := 0; otherEdge < 4; otherEdge++ {
				otherData := otherTile.getEdge(otherEdge)
				// This is actually the unflipped match
				if equal(myData, flipped(otherData)) {
					t.matches = append(t.matches, match{edge, otherEdge, tileNum, false})
				}
				// If the second array is not flipped, one of the tiles *is*
				if equal(myData, otherData) {
					t.matches = append(t.matches, match{edge, otherEdge, tileNum, true})
				}
			}
		}
	}
}

// Determine the global x,y coordinates, flip, and rotation of each tile
func (t *tile) place(m map[int]*tile, x, y, r int, flip bool) {
	if t.placed {
		return
	}
	t.placed = true
	t.x = x
	t.y = y
	t.r = r
	t.flip = flip
	for _, match := range t.matches {
		var d int  // Global coordinate direction
		var nr int // Rotation to apply to resulting tile
		if !t.flip {
			d = (match.myEdge + t.r) % 4
		} else {
			d = (t.r - match.myEdge) % 4
		}
		newFlip := t.flip != match.flipped
		if !newFlip {
			nr = (d + 2 - match.otherEdge) % 4
		} else {
			nr = (d + 2 + match.otherEdge) % 4
		}
		// Make sure we don't stay negative
		if d < 0 {
			d += 4
		}
		if nr < 0 {
			nr += 4
		}
		switch d {
		case 0: // Up
			m[match.otherTile].place(m, t.x, t.y-1, nr, newFlip)
		case 1: // Right
			m[match.otherTile].place(m, t.x+1, t.y, nr, newFlip)
		case 2: // Down
			m[match.otherTile].place(m, t.x, t.y+1, nr, newFlip)
		case 3: // Left
			m[match.otherTile].place(m, t.x-1, t.y, nr, newFlip)
		default:
			fmt.Println("unsupported...", d)
		}
	}
}

// Not the prettiest copy function, but it gets the job done
func (t *tile) transformData(out *[][]bool, xOffset, yOffset int) int {
	nPounds := 0
	for y := 1; y < t.s-1; y++ {
		for x := 1; x < t.s-1; x++ {
			// Step 1: map global coordinate to flipped and rotated coordinate
			var lx, ly int
			if !t.flip {
				switch t.r {
				case 0:
					lx = x
					ly = y
				case 1:
					lx = y
					ly = t.s - 1 - x
				case 2:
					lx = t.s - 1 - x
					ly = t.s - 1 - y
				case 3:
					lx = t.s - 1 - y
					ly = x
				}
			} else {
				switch t.r {
				case 0: // 1 flip
					lx = t.s - 1 - x
					ly = y
				case 1: // 2 flip
					lx = t.s - 1 - y
					ly = t.s - 1 - x
				case 2: // 3 flip
					lx = x
					ly = t.s - 1 - y
				case 3: // 0 flip
					lx = y
					ly = x
				}
			}
			// Step 2: output
			(*out)[y+yOffset-1][x+xOffset-1] = t.data[ly][lx]

			// For conveniencem count `#`s
			if t.data[ly][lx] {
				nPounds++
			}
		}
	}
	return nPounds
}

func checkForSeaMonster(a [][]bool, x, y int) bool {
	sm1 := "                  # "
	sm2 := "#    ##    ##    ###"
	sm3 := " #  #  #  #  #  #   "
	for sx, c := range sm1 {
		if c == '#' && !a[y][x+sx] {
			return false
		}
	}
	for sx, c := range sm2 {
		if c == '#' && !a[y+1][x+sx] {
			return false
		}
	}
	for sx, c := range sm3 {
		if c == '#' && !a[y+2][x+sx] {
			return false
		}
	}
	return true
}

func init() {
	reTileHeader = regexp.MustCompile(`^Tile ([0-9]+):$`)
}

func main() {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}

	tileNum := -1
	var t *tile
	tiles := make(map[int]*tile)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines
		if len(line) == 0 {
			continue
		}
		// First, check if regex matches a tile header
		m := reTileHeader.FindSubmatch(scanner.Bytes())
		if m != nil {
			// Save off the previous tile
			if tileNum != -1 {
				t.s = len(t.data)
				tiles[tileNum] = t
			}
			n, err := strconv.Atoi(string(m[1]))
			if err != nil {
				log.Fatal(err)
			}

			tileNum = n
			t = &tile{n, make([][]bool, 0), 0, nil, 0, 0, 0, false, false}
		} else {
			// Read in data as an array of booleans
			row := make([]bool, len(line))
			for i, c := range line {
				row[i] = c == '#'
			}
			t.data = append(t.data, row)
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	// Save off the final tile
	t.s = len(t.data)
	tiles[tileNum] = t

	// Run the correlation step
	for _, t := range tiles {
		t.correlate(tiles)
	}

	// Find all the corner tiles (tiles with only 2 matches)
	// and multiply their IDs
	product := 1
	for _, t := range tiles {
		if len(t.matches) == 2 {
			product *= t.id
		}
	}
	fmt.Println("Part 1 answer:", product)

	// Attempt all combinations of flips and rotations
	for fr := 0; fr < 8; fr++ {
		doAFlip := fr >= 4
		rot := fr % 4
		// Ensure all tiles are reset
		for _, tile := range tiles {
			tile.placed = false
		}
		// Determine tile positions with different initial states
		// each loop
		tiles[tileNum].place(tiles, 0, 0, rot, doAFlip)
		minX := tiles[tileNum].x
		maxX := tiles[tileNum].x
		minY := tiles[tileNum].y
		maxY := tiles[tileNum].y
		for _, tile := range tiles {
			if tile.x > maxX {
				maxX = tile.x
			}
			if tile.x < minX {
				minX = tile.x
			}
			if tile.y > maxY {
				maxY = tile.y
			}
			if tile.y < minY {
				minY = tile.y
			}
		}

		sz := (tiles[tileNum].s - 2) * (maxX - minX + 1)

		// Make an array to hold the full image
		bigData := make([][]bool, sz)
		for i := range bigData {
			bigData[i] = make([]bool, sz)
		}

		// Render tiles out onto the large array
		nPounds := 0
		nMonsters := 0
		for _, tile := range tiles {
			nPounds += tile.transformData(&bigData, (tile.x-minX)*(tile.s-2), (tile.y-minY)*(tile.s-2))
		}

		// Look for sea monsters
		for y := 0; y < sz-3; y++ {
			for x := 0; x < sz-20; x++ {
				found := checkForSeaMonster(bigData, x, y)
				if found {
					nMonsters++
					nPounds -= 15
				}
			}
		}
		if nMonsters > 0 {
			fmt.Println("Sea Monsters:", nMonsters, "Total roughness:", nPounds)
			break
		}
	}

}
