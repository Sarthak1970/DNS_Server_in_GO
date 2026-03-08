package resolver

import(
	"bytes"
	"encoding/binary"
	"net"
	"strings"
	"fmt"
	"time"
)

func boolToInt(b bool) int {
	if b{
		return 1
	}
	return 0
}

func (c *Client) ipType() (string, error) {
	ip := net.ParseIP(c.IpAddress)

	if ip == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	if ip.To4() != nil {
		return "ipv4", nil
	}

	return "ipv6", nil
}

func IDMatcher(reqID,respID []byte) bool {
	return bytes.Equal(reqID,respID)
}

func appendFromBufferUntilNull(buf *bytes.Buffer) []byte {
	var result []byte

	for {
		b,err:=buf.ReadByte()
		if err!=nil{
			break
		}

		result=append(result,b)

		if b==0{
			break
		}
	}
	return result
}

func parseRData(typ uint16, rdata []byte, messageBuf *bytes.Buffer) (string,error){

	switch typ {

	case 1: 
		ip := net.IP(rdata)
		return ip.String(),nil

	default:
		return "",nil
	}
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
	Questions []Question
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
	IpAddress string
	Port int
}

func (c*Client) Query(message []byte) ([]byte,error){
	ipType,err := c.ipType()
	var addr string
	if err!=nil{
		return nil,fmt.Errorf("Failed to get the IP type: %v",err)
	}

	if ipType=="ipv4"{
		addr=fmt.Sprintf("%s:%d",c.IpAddress,c.Port)
	} else if ipType=="ipv6"{
		addr=fmt.Sprintf("[%s]:%d",c.IpAddress,c.Port)
	}

	conn,err:=net.Dial("udp",addr)
	if err!=nil{
		return nil,fmt.Errorf("Failed to connect to the DNS Server:%v",err)
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5*time.Second))

	_,err=conn.Write(message)
	if err!=nil{
		return nil,fmt.Errorf("failed to send the DNS message: %v",err)
	}

	buf := make([]byte,1024)

	n,err:=conn.Read(buf)

	if err!=nil{
		return nil,fmt.Errorf("Failed to read the response:%v",err)
	}

	response:=buf[:n]

	if !IDMatcher(message[:2],response[:2]){
		return nil,fmt.Errorf("The response ID does not match the request ID")
	}

	return response,nil
}

type ResourceRecord struct{
	Name string
	Type uint16
	Class uint16
	TTL uint32
	RDLength uint16
	RData []byte
	RDataParsed string
}

func ResourceRecordFromBytes(data []byte,messageBufs ...*bytes.Buffer) *ResourceRecord{
	buf := bytes.NewBuffer(data)
	var messageBuf *bytes.Buffer

	if messageBufs != nil{
		messageBuf=messageBufs[0]
	}

	name:=appendFromBufferUntilNull(buf)
	nameLength:=len(name)-1
	decodedName,err:=DecodeName(string(name),messageBuf)

	if err!=nil{
		fmt.Println("Failed to decode the name:%v\n",err)
	}

	typ:=binary.BigEndian.Uint16(data[nameLength : nameLength+2])
	class:=binary.BigEndian.Uint16(data[nameLength+2 : nameLength+4])
    ttl := binary.BigEndian.Uint32(data[nameLength+4 : nameLength+8])
    rdLength := binary.BigEndian.Uint16(data[nameLength+8 : nameLength+10])
    rData := data[nameLength+10 : nameLength+10+int(rdLength)] // 10 is the length of the fields before RData
    rDataParsed, _ := parseRData(typ, rData, messageBuf)

    return &ResourceRecord{
        Name: decodedName,
        Type: typ,
        Class: class,
        TTL: ttl,
        RDLength: rdLength,
        RData: rData,
        RDataParsed: rDataParsed,
    }	
}

func DecodeName(qname string, messageBufs ...*bytes.Buffer) (string, error) {
    encoded := []byte(qname)
    var result bytes.Buffer
    var messageBuf *bytes.Buffer
    if messageBufs != nil {
        messageBuf = messageBufs[0]
    }

    for i := 0; i < len(encoded); {
        length := int(encoded[i])
        if length == 0 {
            break
        }

        if encoded[i]>>6 == 0b11 && messageBuf != nil {
            b := encoded[i+1]
            offset := int(b & 0b11111111)
            messageBytes := messageBuf.Bytes()
            messageBytes = messageBytes[offset:]
            name := appendFromBufferUntilNull(bytes.NewBuffer(messageBytes))
            n, _ := DecodeName(string(name))
            name = []byte(n)
            length = len(name)
            if result.Len() > 0 {
                result.WriteByte('.')
            }
            result.Write(name)
            i += length
            break
        }
        i++

        if i+length > len(encoded) {
            return "", fmt.Errorf("invalid encoded domain name")
        }
        if result.Len() > 0 {
            result.WriteByte('.')
        }
        result.Write(encoded[i : i+length])
        i += length
    }

    return result.String(), nil
}