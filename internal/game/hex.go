package game

import "math"

type Axial struct {
	Q, R int
}

type Hex struct {
	Axial  Axial
	Center struct{ X, Y float64 }
	Radius float64
	TerrID int // -1 if unassigned
}

type HexGrid struct {
	Hexes map[Axial]*Hex
}

func NewHexGrid() *HexGrid {
	g := &HexGrid{Hexes: make(map[Axial]*Hex)}
	height := HexRadius * 2
	rowStep := height * 3 / 4
	colStep := math.Sqrt(3) / 2 * height

	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			ax := Axial{Q: x - (y-(y&1))/2, R: y}
			posX := colStep * float64(x)
			if y%2 != 0 {
				posX += colStep / 2
			}
			h := &Hex{
				Axial:  ax,
				Radius: HexRadius,
				TerrID: -1,
			}
			h.Center.X = posX
			h.Center.Y = rowStep * float64(y)
			g.Hexes[ax] = h
		}
	}
	return g
}

func (g *HexGrid) Neighbor(h *Hex, dir int) *Hex {
	d := Directions[dir%6]
	ax := Axial{Q: h.Axial.Q + d[0], R: h.Axial.R + d[1]}
	return g.Hexes[ax]
}

func (g *HexGrid) PickHex(px, py float64) *Hex {
	cubeX := (px*math.Sqrt(3)/3 - py/3) / HexRadius
	cubeZ := py * 2 / 3 / HexRadius
	cubeY := -cubeX - cubeZ
	return g.CubeRound(cubeX, cubeY, cubeZ)
}

func (g *HexGrid) CubeRound(x, y, z float64) *Hex {
	rx := int(math.Round(x))
	ry := int(math.Round(y))
	rz := int(math.Round(z))

	xDiff := math.Abs(float64(rx) - x)
	yDiff := math.Abs(float64(ry) - y)
	zDiff := math.Abs(float64(rz) - z)

	switch {
	case xDiff > yDiff && xDiff > zDiff:
		rx = -ry - rz
	case yDiff > zDiff:
		ry = -rx - rz
	default:
		rz = -rx - ry
	}

	return g.Hexes[Axial{Q: rx, R: rz}]
}

func HexVertices(cx, cy, radius float64) [7]struct{ X, Y float64 } {
	var verts [7]struct{ X, Y float64 }
	for i := 0; i < 7; i++ {
		angle := (60.0*float64(i) - 30) * math.Pi / 180
		verts[i].X = cx + radius*math.Cos(angle)
		verts[i].Y = cy + radius*math.Sin(angle)
	}
	return verts
}
