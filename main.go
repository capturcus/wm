package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/go-ini/ini"

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
	WM_OFFSET, err := cfg.Section("").Key("wm_offset").Int()
	detonate(err)
	OPACITY, err := cfg.Section("").Key("opacity").Float64()
	detonate(err)
	ANGLE, err := cfg.Section("").Key("angle").Float64()
	detonate(err)

	targetImg, err := imaging.Open(os.Args[1])
	detonate(err)
	wmImg, err := imaging.Open("wm.png")
	detonate(err)

	rgba := image.NewRGBA(targetImg.Bounds())

	mask := image.NewUniform(color.Alpha{uint8(255 * OPACITY)})

	rotatedWm := imaging.Rotate(wmImg, ANGLE, color.Transparent)

	resizedWm := imaging.Resize(rotatedWm, targetImg.Bounds().Max.X, 0, imaging.Lanczos)

	draw.Draw(rgba, targetImg.Bounds(), targetImg, image.Point{0, 0}, draw.Src)

	for i := -resizedWm.Bounds().Max.Y; i < targetImg.Bounds().Max.Y/WM_OFFSET; i++ {
		draw.DrawMask(rgba, image.Rectangle{image.Point{0, WM_OFFSET * i}, targetImg.Bounds().Max}, resizedWm, image.Point{0, 0}, mask, image.Point{0, 0}, draw.Over)
	}

	extension := filepath.Ext(os.Args[1])
	name := os.Args[1][0 : len(os.Args[1])-len(extension)]

	out, err := os.Create(name + ".wm.png")
	if err != nil {
		fmt.Println(err)
	}

	png.Encode(out, rgba)
}
