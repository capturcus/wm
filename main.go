package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"github.com/nfnt/resize"

	_ "image/jpeg"
)

/*
const (
	WM_COUNTS = 4
)*/

func detonate(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg, err := ini.Load("wm.ini")
	detonate(err)
	WM_COUNTS, err := cfg.Section("").Key("num_wms").Int()
	detonate(err)
	OPACITY, err := cfg.Section("").Key("opacity").Float64()
	detonate(err)
	targetFile, err := os.Open(os.Args[1])
	detonate(err)
	wmFile, err := os.Open("wm.png")
	detonate(err)
	targetImg, _, err := image.Decode(targetFile)
	detonate(err)
	wmImg, _, err := image.Decode(wmFile)
	detonate(err)
	rgba := image.NewRGBA(targetImg.Bounds())

	mask := image.NewUniform(color.Alpha{uint8(255 * OPACITY)})

	resizedWm := resize.Resize(uint(targetImg.Bounds().Max.X), 0, wmImg, resize.Lanczos3)

	draw.Draw(rgba, targetImg.Bounds(), targetImg, image.Point{0, 0}, draw.Src)

	wmYDt := (targetImg.Bounds().Max.Y - resizedWm.Bounds().Max.Y) / (WM_COUNTS - 1)

	for i := 0; i < WM_COUNTS; i++ {
		draw.DrawMask(rgba, image.Rectangle{image.Point{0, wmYDt * i}, targetImg.Bounds().Max}, resizedWm, image.Point{0, 0}, mask, image.Point{0, 0}, draw.Over)
	}

	extension := filepath.Ext(os.Args[1])
	name := os.Args[1][0 : len(os.Args[1])-len(extension)]

	out, err := os.Create(name + ".wm.png")
	if err != nil {
		fmt.Println(err)
	}

	png.Encode(out, rgba)
}
