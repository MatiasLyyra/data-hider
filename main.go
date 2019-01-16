package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"

	"github.com/matiaslyyra/data-hider/lsb"
)

func main() {
	mode := flag.String("mode", "hide", "Values: hide / reveal")
	inFile := flag.String("in", "", "Input image")
	outFile := flag.String("out", "", "File path to output file. Ignored in reveal mode.")
	text := flag.String("text", "", "Tet to hide. Ignored in reveal mode.")
	flag.Parse()
	reader, err := os.Open(*inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	srcImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	switch *mode {
	case "hide":
		writer, err := os.Create(*outFile)
		if err != nil {
			log.Fatal(err)
		}
		defer writer.Close()
		img := lsb.Hide(srcImage, []byte(*text))
		png.Encode(writer, img)
		fmt.Printf("Wrote %v bytes\n", len(*text))
	case "reveal":
		data := lsb.Reveal(srcImage)
		fmt.Printf("Revealed text:\n%v\n", string(data))
	}
}
