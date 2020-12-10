package artnet

import "fmt"

type ArtDmx struct {
	ArtHeader
	ProtVer ProtVer

	Sequence uint8
	Physical uint8
	SubUni   uint8
	Net      uint8
	Length   uint16
}

const Port = 6454
const MTU = 1500

var ARTNET = [8]byte{'A', 'r', 't', '-', 'N', 'e', 't', 0}

type OpCode struct {
	Lo uint8
	Hi uint8
}
type ProtVer struct {
	Hi uint8
	Lo uint8
}

func (protVer ProtVer) ToUint() uint {
	return uint(protVer.Hi<<8) + uint(protVer.Lo)
}

func (protVer ProtVer) IsCompatible(otherVersion ProtVer) bool {
	return protVer.ToUint() == otherVersion.ToUint()
}

type ArtHeader struct {
	ID     [8]byte
	OpCode OpCode
}

func (h ArtHeader) Header() ArtHeader {
	return h
}

type ArtPacket interface {
	Header() ArtHeader
}

var ProtVer14 = ProtVer{0, 14}
var (
	OpDmx = OpCode{0x00, 0x50}
)

func DecodeMac(mac [6]byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}
