package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	TUNSETIFF = 0x400454ca // создание или привязка дескриптора к TUN/TAP интерфейсу с заданным именем и параметрами
	IFF_TAP   = 0x0002     // создается или открывается устройство TAP, а не TUN
	IFF_NO_PI = 0x1000     // получаем чистый кадр, без заголовков
)

type ifreq struct {
	Name  [16]byte // имя сетевого интерфейса
	Flags uint16   // параметры интерфейса
	Pad   [22]byte // буфер заполнения для выравнивания
}

func openTAP(name string) (*os.File, error) {
	const op = "cmd.openTap"

	fd, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var req ifreq
	copy(req.Name[:], name)
	req.Flags = IFF_NO_PI | IFF_TAP

	_, _, err = syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(),
		uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if err != nil {
		fd.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return fd, nil
}

func main() {

	fd, err := openTAP("tap0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fd.Close()
	fmt.Println("file is open")

	ch := make(chan []byte)

	// для приема frame
	go func {

	}

	// для отправки
	go func {

	}


	select{}
}
