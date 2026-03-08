package resolver

import(
	"bytes"
	"encoding/binary"
	"net"
	"stirngs"
)

func boolToInt(b bool) int {
	if b{
		return 1
	}
	return 0
}

type Header struct{
	ID uint16
	Flags uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func (h *Header) toBytes() []byte{
	buf := new(bytes.Buffer)

    binary.Write(buf, binary.BigEndian, h.ID)
    binary.Write(buf, binary.BigEndian, h.Flags)
    binary.Write(buf, binary.BigEndian, h.QDCount)
    binary.Write(buf, binary.BigEndian, h.ANCount)
    binary.Write(buf, binary.BigEndian, h.NSCount)
    binary.Write(buf, binary.BigEndian, h.ARCount)

	return buf.Bytes()
}

type HeaderFlag struct{
	QR bool
	Opcode uint8
	AA bool
	TC bool
	RD bool
	RA bool
	Z uint8
	RCode uint8
}

func (hf *HeaderFlag) GenerateFlags() uint16{
	qr:=uint16(boolToInt(hf.QR))
	opcode:=uint16(hf.Opcode)
	aa:=uint16(boolToInt(hf.AA))
	tc:=uint16(boolToInt(hf.TC))
	rd:=uint16(boolToInt(hf.RD))
	ra:=uint16(boolToInt(hf.RA))
	z:=uint16(hf.Z)
	rcode:=uint16(hf.RCode)

    return uint16(qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode)  //OR Operator to stack the bits in 16 bits
}

type Question struct{
	Name string
	QName string
	Qtype uint16
	QClass uint16
}

func encodeName(name string) []byte {

	parts := strings.Split(name,".")
	buf := []byte{}

	for _, part := range parts {
		buf = append(buf, byte(len(part)))
		buf = append(buf, []byte(part)...)
	}

	buf = append(buf, 0)

	return buf
}

func (q* Question) ToBytes() []byte{
	buf := new(bytes.Buffer)
	buf.Write([]byte(q.QName))
	binary.Write(buf,binary.BigEndian,q.Qtype)
	binary.Write(buf,binary.BigEndian,q.QClass)
	return buf.Bytes()
}

type DNSMessage struct{
	Header Header
	Question []Question
	Answers []ResourceRecord
	AuthorityRRs []ResourceRecord
	AdditionalRRs []ResourceRecord
}

//creates a new DNSMessage to send over network
func NewDNSMessage(header Header,questions []Question,records ...[]ResourceRecord) *DNSMessage{ //Viriadic Parameter-Gives Flexibility-Since Query has 0 ans but response has many
	answers:=make([]ResourceRecord,0)
	authorityRRs:=make([]ResourceRecord,0)
	additionalRRs:=make([]ResourceRecord,0)

	if len(records)>0{
		answers=records[0]
	}
    if len(records) > 1 {
        authorityRRs = records[1]
    }

    if len(records) > 2 {
        additionalRRs = records[2]
    }

	return &DNSMessage{
		Header:header,
		Questions:questions,
		Answers:answers,
		AuthorityRRs:authorityRRs,
		AdditionalRRs:additionalRRs,
	}
}

type Client struct{
	ipAddress string
	port int
}

func (c*Client) Query(message []byte) ([]byte,error){
	ipType,err := c.ipType()
	var addr string
	if err!=nil{
		return nil,fmt.Errorf("Failed to get the IP type: %v",err)
	}

	if ipType=="ipv4"{
		addr:=fmt.Sprintf("%s:%d",c.ipAddress,c.port)
	} else if ipType=="ipv6"{
		addr=fmt.Sprint
	}
}