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
	Test uint16 = 0xDEAD
	ARP  uint16 = 0x0806
	IPv4 uint16 = 0x0800
)

type Ethernet struct {
	Tap   *os.File
	Ch    chan []byte
	MyMAC net.HardwareAddr
	Cache *arp.ARPCache
}

type EthernetFrame struct {
	MacDst    [6]byte
	MacSrc    [6]byte
	EtherType uint16
	Payload   []byte // если меньше 46 байт, то используется поле-заполнитель для предотвращения коллизий
}

func NewEthernet(tap *os.File, ch chan []byte, mymac net.HardwareAddr, cache *arp.ARPCache) Ethernet {
	return Ethernet{
		Tap:   tap,
		Ch:    ch,
		MyMAC: mymac,
		Cache: cache,
	}
}

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
		MacDst:    dst,
		MacSrc:    src,
		EtherType: pr,
		Payload:   data,
	}
}

func (ef *EthernetFrame) Marshal() []byte {
	frame := []byte{}
	buf := make([]byte, 2)

	binary.BigEndian.PutUint16(buf, ef.EtherType)

	frame = append(frame, ef.MacDst[:]...)
	frame = append(frame, ef.MacSrc[:]...)
	frame = append(frame, buf...)

	if len(ef.Payload) < 46 {
		Padding(&ef.Payload)
	}
	frame = append(frame, ef.Payload...)

	return frame
}

func Unmarshal(frame []byte) EthernetFrame {
	var ef EthernetFrame

	copy(ef.MacDst[:], frame[0:6])
	copy(ef.MacSrc[:], frame[6:12])

	buf := frame[12:14]
	ef.EtherType = binary.BigEndian.Uint16(buf)

	ef.Payload = frame[14:]

	return ef
}

func Padding(data *[]byte) {
	length := len(*data)
	padLen := 46 - length
	padding := make([]byte, padLen)
	*data = append(*data, padding...)
}

func (e *Ethernet) HandleFrame(frame []byte) {
	etherframe := Unmarshal(frame)

	if !bytes.Equal(e.MyMAC, etherframe.MacDst[:]) && !utils.IsBroadcast(etherframe.MacDst[:]) {
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
			resData := newArp.Marshal()

			e.Ch <- resData
		}

	case IPv4:
		// TODO
	}

}
