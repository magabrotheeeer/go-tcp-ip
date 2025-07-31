package main

import (
	"github.com/magabrotheeeer/go-tcp-ip/ethernet"
)

func main() {

	data := []byte("Hello world, its works!")

	src := [6]byte{0, 0, 0, 0, 0, 0}
	dest := [6]byte{0, 0, 0, 0, 0, 0}

	ether := ethernet.New(src, dest, "test", data)
	frame := ether.Marshal()
	_ = frame

}
