package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func wasmHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "wasm_exec.html")
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
		d := r.FormValue("dci")

		content := f + l + m + p + d
		filename := f + l + ".png"

		qrCode, _ := qr.Encode(content, qr.M, qr.Auto)

		// Scale the barcode to 200x200 pixels
		qrCode, _ = barcode.Scale(qrCode, 200, 200)

		file, _ := os.Create(filename)
		defer file.Close()

		png.Encode(file, qrCode)

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
	mux.HandleFunc("/wasm_exec.html", wasmHandler)
	log.Printf("server started")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
