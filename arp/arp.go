package arp

import (
	"encoding/binary"
	"fmt"

	"github.com/magabrotheeeer/go-tcp-ip/utils"
)

const (
	ARPRequest = 1 // arp-запрос
	ARPReply   = 2 // arp-ответ
)

type ARPPackage struct {
	HardwareType    uint16 // Ethernet или Wifi
	ProtocolType    uint16 // ip
	HardwareAddrLen uint8
	ProtocolAddrLen uint8
	Operation       uint16 // arp-запрос или arp-ответ
	SrcMac          [6]byte
	SrcIp           [4]byte
	DstMac          [6]byte
	DstIp           [4]byte
}

// создает новый arp пакет
func NewARPPackage(operation uint16, srcmac [6]byte, srcip [4]byte, dstmac [6]byte, dstip [4]byte) ARPPackage {
	return ARPPackage{
		HardwareType:    0x0001, // Ethernet
		ProtocolType:    0x0800, // IP
		HardwareAddrLen: 6,
		ProtocolAddrLen: 4,
		Operation:       operation,
		SrcMac:          srcmac,
		SrcIp:           srcip,
		DstMac:          dstmac,
		DstIp:           dstip,
	}
}

// приводит arppackage в слайс байт
func (ap *ARPPackage) Marshal() []byte {
	pack := make([]byte, 28)

	buf := make([]byte, 2)

	binary.BigEndian.PutUint16(buf, ap.HardwareType)
	copy(pack[0:2], buf)

	binary.BigEndian.PutUint16(buf, ap.ProtocolType)
	copy(pack[2:4], buf)

	pack[4] = ap.HardwareAddrLen
	pack[5] = ap.ProtocolAddrLen

	binary.BigEndian.PutUint16(buf, ap.Operation)
	copy(pack[6:8], buf)

	copy(pack[8:14], ap.SrcMac[:])
	copy(pack[14:18], ap.SrcIp[:])
	copy(pack[18:24], ap.DstMac[:])
	copy(pack[24:28], ap.DstIp[:])

	return pack
}

// приводит слайс байт в arppackage
func Unmarshal(data []byte) (ARPPackage, error) {

	if len(data) < 28 {
		return ARPPackage{}, fmt.Errorf("data more then 28 bytes")
	}

	var ap ARPPackage

	buf := data[0:2]
	ap.HardwareType = binary.BigEndian.Uint16(buf)

	buf = data[2:4]
	ap.ProtocolType = binary.BigEndian.Uint16(buf)

	ap.HardwareAddrLen = data[4]
	ap.ProtocolAddrLen = data[5]

	buf = data[6:8]
	ap.Operation = binary.BigEndian.Uint16(buf)

	copy(ap.SrcMac[:], data[8:14])
	copy(ap.SrcIp[:], data[14:18])
	copy(ap.DstMac[:], data[18:24])
	copy(ap.DstIp[:], data[24:28])

	return ap, nil
}

// обработка arppackage в зависимости от поля operation
func HandleARP(pack ARPPackage, cache *ARPCache) (*ARPPackage, error) {
	var res ARPPackage
	var err error
	switch pack.Operation {
	case ARPRequest:
		res, err = pack.BuildReply()
		if err != nil {
			return nil, err
		}
		return &res, nil
	case ARPReply:
		pack.UpdateData(cache)
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown operation %d", pack.Operation)
	}

}

// создает arppackage-reply
func (ap ARPPackage) BuildReply() (ARPPackage, error) {
	targetIP := ap.DstIp

	list, err := utils.GetInterfaceIPs("tar0")
	if err != nil {
		return ARPPackage{}, err
	}
	if !utils.IsMyIP(targetIP, list) {
		return ARPPackage{}, nil
	}
	res := NewARPPackage(2, ap.DstMac, ap.DstIp, ap.SrcMac, ap.SrcIp)

	return res, nil
}

// обновляет ip и mac в arp-cache
func (ap ARPPackage) UpdateData(cache *ARPCache) {
	senderIp := ap.SrcIp
	senderMac := ap.SrcMac

	stringSenderIp := fmt.Sprintf("%d.%d.%d.%d",
		senderIp[0], senderIp[1], senderIp[2], senderIp[3])

	cache.Add(stringSenderIp, senderMac[:])

}
