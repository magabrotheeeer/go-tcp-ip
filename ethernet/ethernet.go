package ethernet // Ethernet II (IEEE 802.3)
import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"

	"github.com/magabrotheeeer/go-tcp-ip/arp"
	"github.com/magabrotheeeer/go-tcp-ip/utils"
)

const (
	Test uint16 = 0xDEAD // тест для нахождения своего frame
	ARP  uint16 = 0x0806 // arp-протокол
	IPv4 uint16 = 0x0800 // ip-протокол
)

type Ethernet struct {
	Tap   *os.File         // дескриптор файла
	Ch    chan []byte      // канал для передачи фреймов
	MyMAC net.HardwareAddr // mac устройства
	Cache *arp.ARPCache    // arp-cache
}

type EthernetFrame struct {
	DstMac    [6]byte // получатель
	SrcMac    [6]byte // отправитель
	EtherType uint16  // ip, test, arp
	Payload   []byte  // если меньше 46 байт, то используется поле-заполнитель для предотвращения коллизий
}

// создает новый ethernet
func NewEthernet(tap *os.File, ch chan []byte, mymac net.HardwareAddr, cache *arp.ARPCache) Ethernet {
	return Ethernet{
		Tap:   tap,
		Ch:    ch,
		MyMAC: mymac,
		Cache: cache,
	}
}

// создает новый ethernetframe
func NewEthernetFrame(dst [6]byte, src [6]byte, protocol string, data []byte) EthernetFrame {
	var pr uint16
	switch protocol {
	case "ip":
		pr = IPv4
	case "arp":
		pr = ARP
	case "test":
		pr = Test
	}

	return EthernetFrame{
		DstMac:    dst,
		SrcMac:    src,
		EtherType: pr,
		Payload:   data,
	}
}

// приводит ethernetframe в слайс байт
func (ef *EthernetFrame) Marshal() []byte {
	frame := []byte{}
	buf := make([]byte, 2)

	binary.BigEndian.PutUint16(buf, ef.EtherType)

	frame = append(frame, ef.DstMac[:]...)
	frame = append(frame, ef.SrcMac[:]...)
	frame = append(frame, buf...)

	if len(ef.Payload) < 46 {
		Padding(&ef.Payload)
	}
	frame = append(frame, ef.Payload...)

	return frame
}

// приводит слайс байт в ethernetframe
func Unmarshal(frame []byte) EthernetFrame {
	var ef EthernetFrame

	copy(ef.DstMac[:], frame[0:6])
	copy(ef.SrcMac[:], frame[6:12])

	buf := frame[12:14]
	ef.EtherType = binary.BigEndian.Uint16(buf)

	ef.Payload = frame[14:]

	return ef
}

// заполняет payload нулями
func Padding(data *[]byte) {
	length := len(*data)
	padLen := 46 - length
	padding := make([]byte, padLen)
	*data = append(*data, padding...)
}

// обрабатывает ethernet frame
func (e *Ethernet) HandleFrame(frame []byte) {
	etherframe := Unmarshal(frame)

	if !bytes.Equal(e.MyMAC, etherframe.DstMac[:]) && !utils.IsBroadcast(etherframe.DstMac[:]) {
		return
	}

	switch etherframe.EtherType {
	case ARP:
		arpPack, err := arp.Unmarshal(etherframe.Payload)
		if err != nil {
			log.Println(err.Error())
			return
		}
		newArp, err := arp.HandleARP(arpPack, e.Cache)

		if newArp != nil {
			respondArp := newArp.Marshal()
			nef := NewEthernetFrame(newArp.DstMac, newArp.SrcMac, "arp", respondArp)
			respondFrame := nef.Marshal()
			e.Ch <- respondFrame
		}

	case IPv4:
		// TODO
	}
}
