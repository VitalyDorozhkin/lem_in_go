package main

import (
	"fmt"
	"github.com/VitalyDorozhkin/lem_in_go/datastruct"
	"image"
	"image/color"
	"image/gif"
	"os"
)

var (
	white    = color.RGBA{f, f, f, f}
	black    = color.RGBA{0, 0, 0, f}
	redL     = color.RGBA{f, 127, 127, f}
	red      = color.RGBA{f, 0, 0, f}
	redD     = color.RGBA{127, 0, 0, f}
	greenL   = color.RGBA{127, f, 127, f}
	green    = color.RGBA{0, f, 0, f}
	greenD   = color.RGBA{0, 127, 0, f}
	blueL    = color.RGBA{127, 127, f, f}
	blue     = color.RGBA{0, 0, f, f}
	blueDark = color.RGBA{0, 0, 127, f}
	bg = color.RGBA{
		R: 0x45,
		G: 0x47,
		B: 0x54,
		A: f,
	}
)

func drawThickLine(img *image.Paletted, x1, y1, x2, y2 int, col color.Color, col2 color.Color, thickness int) {
	thickness += 1 - thickness%2
	for x := 0 - thickness/2; x <= thickness/2; x++ {
		for y := 0 - thickness/2; y <= thickness/2; y++ {
			Bresenham(img, x1+x, y1+y, x2+x, y2+y, col2)
		}
	}
	Bresenham(img, x1, y1, x2, y2, col)
}

func drawThickPoint(img *image.Paletted, x1, y1 int, col color.Color, thickness int) {
	thickness += 1 - thickness%2
	for x := 0 - thickness/2; x <= thickness/2; x++ {
		for y := 0 - thickness/2; y <= thickness/2; y++ {
			img.Set(x1+x, y1+y, col)
		}
	}
}

func Bresenham(img *image.Paletted, x1, y1, x2, y2 int, col color.Color) {
	var dx, dy, e, slope int

	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	if dy < 0 {
		dy = -dy
	}
	c := 1
	if y1 >= y2 {
		c = -c
	}
	if dx > dy {
		dy, e, slope = 2*dy, dx, 2*dx
		for ; dx != 0; dx-- {
			img.Set(x1, y1, col)
			x1++
			e -= dy
			if e < 0 {
				y1 += c
				e += slope
			}
		}
	} else {
		dx, e, slope = 2*dx, dy, 2*dy
		for ; dy != 0; dy-- {
			img.Set(x1, y1, col)
			y1 += c
			e -= dx
			if e < 0 {
				x1++
				e += slope
			}
		}
	}
	img.Set(x2, y2, col)
}

func fill(img *image.Paletted, col color.Color) {
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			img.Set(x, y, col)
		}
	}
}

var f = uint8(0xff)
var colors1 = map[string]color.Color{
	"white":  color.RGBA{f, f, f, f},
	"black":  color.RGBA{0, 0, 0, f},
	"red":    color.RGBA{f, 0, 0, f},
	"red2":   color.RGBA{127, 0, 0, f},
	"green":  color.RGBA{0, f, 0, f},
	"green2": color.RGBA{0, 127, 0, f},
	"blue":   color.RGBA{0, 0, f, f},
}

func drawNodes(img *image.Paletted, nodes map[string]*datastruct.Node, thickness int) {
	for _, node := range nodes {
		if node.Status == "start" || node.Status == "end" {
			drawThickPoint(img, node.X, node.Y, red, 2*thickness)
		}
		drawThickPoint(img, node.X, node.Y, redD, thickness)
	}
}

func drawLinks(img *image.Paletted, links []*datastruct.Link, thickness int) {
	for _, link := range links {
		drawThickLine(img, link.NodeStart.X, link.NodeStart.Y, link.NodeEnd.X, link.NodeEnd.Y, green, greenL, thickness)
	}
}

func drawAnt(img *image.Paletted, start *datastruct.Node, end *datastruct.Node, progress float64) {
	dx := end.X - start.X
	x := start.X + int(float64(dx)*progress)
	dy := end.Y - start.Y
	y := start.Y + int(float64(dy)*progress)
	drawThickPoint(img, x, y, red, 5)
}

func main() {
	ants, graph, stepList := datastruct.NewReadedGraph()
	var side = 1200
	datastruct.MoveGraph(graph, side, 60)

	var palette = []color.Color{
		white,
		black,
		redL,
		red,
		redD,
		greenL,
		green,
		greenD,
		blueL,
		blue,
		blueDark,
	}
	var images []*image.Paletted
	var delays []int

	frames := 30

	imgPer := image.NewPaletted(image.Rect(0, 0, side, side), palette)
	fill(imgPer, blueL)
	drawLinks(imgPer, graph.Links, 3)
	drawNodes(imgPer, graph.Nodes, 9)
	for _, steps := range stepList {
		println("step")
		for frame := 0; frame < frames; frame++ {
			pix := make([]uint8, len(imgPer.Pix), len(imgPer.Pix))
			copy(pix, imgPer.Pix)
			img := image.NewPaletted(image.Rect(0, 0, side, side), palette)
			img.Pix = pix
			for _, step := range steps {
				drawAnt(img, ants[step.LeminNumber-1], step.NodeEnd, float64(frame)/float64(frames))
				if frame+1 == frames {
					ants[step.LeminNumber-1] = step.NodeEnd
				}
			}
			images = append(images, img)
			delays = append(delays, 0)
		}
	}

	f, err := os.OpenFile("rgb.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image:           images,
		Delay:           delays,
		LoopCount:       0,
		Disposal:        nil,
		Config:          image.Config{},
		BackgroundIndex: 0,
	})
}
