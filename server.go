package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dchest/uniuri"

	//"github.com/rs/xid"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func mergePng(path1 string, path2 string, path3 string, x int, y int) {
	image1, err := os.Open(path1)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	first, err := png.Decode(image1)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer image1.Close()

	image2, err := os.Open(path2)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	second, err := png.Decode(image2)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer image2.Close()

	//offset := image.Pt(890, 350)
	offset := image.Pt(x, y)
	b := first.Bounds()
	image3 := image.NewRGBA(b)
	draw.Draw(image3, b, first, image.ZP, draw.Src)
	draw.Draw(image3, second.Bounds().Add(offset), second, image.ZP, draw.Over)

	third, err := os.Create(path3)
	if err != nil {
		log.Fatalf("failed to create: %s", err)
	}
	png.Encode(third, image3)
	defer third.Close()
}

func renderText(text string) {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	dc := gg.NewContext(550, 100)
	dc.SetFontFace(face)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(text, 275, 50, 0.5, 0.5)
	dc.SavePNG("out.png")

}

func wasmHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		f := r.FormValue("first")
		l := r.FormValue("last")
		m := r.FormValue("mail")
		p := r.FormValue("phone")
		//d := r.FormValue("dci")
		//	x := xid.New()
		x := uniuri.New()

		content := f + " " + l + " " + " " + m + " " + p + x
		filename := "./qr/" + f + ".png"

		qrCode, _ := qr.Encode(content, qr.M, qr.Auto)

		// Scale the barcode to 270x270 pixels
		qrCode, _ = barcode.Scale(qrCode, 270, 270)

		file, _ := os.Create(filename)
		defer file.Close()

		log.Printf("new code: " + filename + " id: " + x)

		png.Encode(file, qrCode)

		//rare := "rare.png"
		//mergePng(rare, filename)

		rarity := r.FormValue("rarity") + ".png"
		mergePng(rarity, filename, filename, 650, 1250)

		renderText(f)
		mergePng(filename, "out.png", filename, 50, 1400)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func main() {

	files, err := filepath.Glob("*")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(files) // contains a list of all files in the current directory

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/index.html", wasmHandler)
	log.Printf("server started")
	log.Fatal(http.ListenAndServe(":3000", mux))

}
