package main

import (
	"log"
	"net/http"
	"net/url"
	"syscall/js"

	"github.com/dennwc/dom"
)

type writer dom.Element

// Write implements io.Writer.

func (d writer) Write(p []byte) (n int, err error) {
	node := dom.GetDocument().CreateElement("div")
	node.SetTextContent(string(p))
	(*dom.Element)(&d).AppendChild(node)
	return len(p), nil
}

func main() {

	t := dom.GetDocument().GetElementById("target")
	i := dom.GetDocument().GetElementById("qrcode")

	f := js.Global().Get("document").Call("getElementById", "primary_player").Get("value").String()
	l := js.Global().Get("document").Call("getElementById", "primary_dci").Get("value").String()
	m := js.Global().Get("document").Call("getElementById", "secondary_player").Get("value").String()
	p := js.Global().Get("document").Call("getElementById", "secondary_dci").Get("value").String()
	r := js.Global().Get("document").Call("getElementById", "rarity").Get("value").String()
	//d := js.Global().Get("document").Call("getElementById", "dci").Get("value").String()

	filename := "./qr/" + f + ".png"

	_, err := http.PostForm("./index.html",
		//url.Values{"first": {f}, "last": {l}, "mail": {m}, "phone": {p}, "dci": {d}})
		url.Values{"first": {f}, "last": {l}, "mail": {m}, "phone": {p}, "rarity": {r}})
	if err != nil {
		log.Fatal(err)
	}

	i.SetAttribute("src", filename)

	logger := log.New((*writer)(t), "", log.LstdFlags)
	logger.Print("QR code is ready: " + filename)

}
