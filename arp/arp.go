package arp

import (
	"encoding/binary"
	"fmt"
)

const (
	ARPRequest = 1
	ARPReply   = 2
)

type ARPPackage struct {
	HardwareType    uint16
	ProtocolType    uint16
	HardwareAddrLen uint8
	ProtocolAddrLen uint8
	Operation       uint16
	SrcMac          [6]byte
	SrcIp           [4]byte
	DstMac          [6]byte
	DstIp           [4]byte
}

func New(operation uint16, srcmac [6]byte, srcip [4]byte, dstmac [6]byte, dstip [4]byte) ARPPackage {
	return ARPPackage{
		HardwareType:    0x0001,
		ProtocolType:    0x0800,
		HardwareAddrLen: 6,
		ProtocolAddrLen: 4,
		Operation:       operation,
		SrcMac:          srcmac,
		SrcIp:           srcip,
		DstMac:          dstmac,
		DstIp:           dstip,
	}
}

func (ap ARPPackage) Marshal() []byte {
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

func (ap ARPPackage) CheckOperation() error {
	switch ap.Operation {
	case ARPRequest:
		Request()
	case ARPReply:
		Reply()
	default:
		return fmt.Errorf("unknown operation")
	}
	return nil
}

func Request() ARPPackage {

}

func Reply() ARPPackage {

}
