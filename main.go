package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"

	_ "image/jpeg"

	wf "EMPTY_TEMPLATE/webframework"
)

/*
const (
	WM_COUNTS = 4
)*/

func main() {

	exeDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	logrus.SetFormatter(&wf.ErrorFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	errorFile, _ := os.Create(exeDir + "/error.txt")
	logrus.SetOutput(errorFile)

	cfg, err := ini.Load(exeDir + "/wm.ini")
	if err != nil {
		wf.Fatal(err)
	}
	WM_OFFSET, err := cfg.Section("").Key("wm_offset").Int()
	if err != nil {
		wf.Fatal(err)
	}
	OPACITY, err := cfg.Section("").Key("opacity").Float64()
	if err != nil {
		wf.Fatal(err)
	}
	ANGLE, err := cfg.Section("").Key("angle").Float64()
	if err != nil {
		wf.Fatal(err)
	}

	ioutil.WriteFile("debug.txt", []byte(os.Args[1]), os.ModeExclusive)

	targetImg, err := imaging.Open(os.Args[1])
	if err != nil {
		wf.Fatal(err)
	}
	wmImg, err := imaging.Open("wm.png")
	if err != nil {
		wf.Fatal(err)
	}

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
		wf.Fatal(err)
	}

	png.Encode(out, rgba)
}
