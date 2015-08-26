package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/mgutz/ansi"
	"github.com/nfnt/resize"
	. "gopkg.in/godo.v1"
)

func resizeImage(inDir string, name string, size uint, outDir string) {
	file, err := os.Open(inDir + name)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Thumbnail(size, size, img, resize.Lanczos3)

	out, err := os.Create(outDir + name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}

func tasks(p *Project) {
	// Given a json list of files and a source path and out path, will rename the files in the order specified.
	p.Task("rename", func(c *Context) {
		// the path to a json file of the form `{files:[]}` where the contents of the array
		// are a list of files in the order to be renamed numerically
		sourceFile := c.Args.MustString("files", "f")
		// Path where the files exist
		inPath := c.Args.MustString("inpath", "i")
		// Outpath should be a path with a trailing slash
		outPath := c.Args.MustString("out", "o")
		os.MkdirAll(outPath, 0777)
		source, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			msg := "Error reading file: " + err.Error()
			ansi.Color(msg, "red")
		}
		files := map[string][]string{}
		json.Unmarshal(source, &files)
		for i := 0; i < len(files["files"]); i++ {
			imgSrc, imgErr := ioutil.ReadFile(inPath + files["files"][i])
			if imgErr != nil {
				msg := "Error reading image: " + imgErr.Error()
				ansi.Color(msg, "red")
				return
			}
			ioutil.WriteFile(outPath+strconv.Itoa(i)+".jpg", imgSrc, 0777)
		}
	})

	// Resizes images
	p.Task("resize", func(c *Context) {
		dir := c.Args.MustString("dir", "d")
		width := c.Args.MustInt("size", "s")
		outDir := c.Args.MustString("out", "o")

		fmt.Println(width)

		files, _ := ioutil.ReadDir(dir)
		for _, f := range files {
			fmt.Println(f.Name())
			resizeImage(dir, f.Name(), uint(width), outDir)
		}
	})
}

func main() {
	Godo(tasks)
}
