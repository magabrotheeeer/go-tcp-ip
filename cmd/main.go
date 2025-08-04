package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"unsafe"

	"github.com/magabrotheeeer/go-tcp-ip/arp"
	"github.com/magabrotheeeer/go-tcp-ip/ethernet"
	"golang.org/x/sys/unix"
)

const (
	TUNSETIFF        = 0x400454ca // создание или привязка дескриптора к TUN/TAP интерфейсу с заданным именем и параметрами
	IFF_TAP          = 0x0002     // создается или открывается устройство TAP, а не TUN
	IFF_NO_PI        = 0x1000     // получаем чистый кадр, без заголовков
)

type ifreq struct {
	Name  [16]byte // имя сетевого интерфейса
	Flags uint16   // параметры интерфейса
	Pad   [22]byte // буфер заполнения для выравнивания
}

// открывает и настраивает TAP интерфейс
func openTAP(name string) (*os.File, error) {
	const op = "cmd.openTap"

	fd, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var req ifreq
	copy(req.Name[:], name)
	req.Flags = IFF_NO_PI | IFF_TAP

	_, _, err = syscall.Syscall(
		unix.SYS_IOCTL,
		fd.Fd(),
		uintptr(TUNSETIFF),
		uintptr(unsafe.Pointer(&req)),
	)
	if err != nil {
		fd.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return fd, nil
}

func main() {

	fd, err := openTAP("tap0")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer fd.Close()
	log.Println("file is open")

	iface, err := net.InterfaceByName("tap0")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ch := make(chan []byte)

	arpCache := arp.NewARPCache()

	ether := ethernet.NewEthernet(fd, ch, iface.HardwareAddr, arpCache)

	// для приема frame
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := fd.Read(buf)
			if err != nil {
				log.Println("tap read error:", err)
				continue
			}
			go ether.HandleFrame(buf[:n])
		}

	}()

	// для отправки
	go func() {
		for frame := range ether.Ch {
			_, err := fd.Write(frame)
			if err != nil {
				log.Println("TAP write error:", err)
				continue
			}
		}
	}()

	select {}
}
