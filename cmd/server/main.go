package main 

import (
	"fmt"
	"net"
	// "sync"
	"github.com/miekg/dns"
)

func main(){
	addr := &net.UDPAddr{ //struct representing UDP endpoint
		IP:net.IPv4(0,0,0,0),
		Port:8053, //53 requires root privilege
		Zone:"",
	}
	conn,err:=net.ListenUDP("udp",addr)

	if err!=nil{
		fmt.Println("Failed to start server:",err)
		return
	}
	fmt.Println("DNS Server running on port :",addr.Port)
	defer conn.Close()

	for{
		buf:=make([]byte,1024)
		
		n, remoteAddr,err:=conn.ReadFromUDP(buf)
		if err!=nil{
			fmt.Println("Error Reading Packet",err)
			continue
		}
		fmt.Println("Received Packet from:", remoteAddr.String())
		// fmt.Println("Message says:", string(buf[:n]))

		var msg dns.Msg

		err=msg.Unpack(buf[:n])
		if err!=nil{
			fmt.Println("Error in Decoding packet:",err)
			continue
		}

		for _, q := range msg.Question {
			fmt.Println("Query Domain:", q.Name)
			fmt.Println("Query Type:", dns.TypeToString[q.Qtype])
		}

		resp:=dns.Msg{}
		resp.SetReply(&msg)

		for _,q := range msg.Question {
			fmt.Println("Query Domain:",q.Name)

			if q.Qtype == dns.TypeA{
				rr,err:=dns.NewRR(q.Name+" 60 IN A 1.2.3.4") //rr=Resource Record
				if err!=nil{
					fmt.Println("RR error",err)
					continue
				}
				resp.Answer=append(resp.Answer,rr)
			}
		}

		responseBytes,err:=resp.Pack()  //converting message to bytes
		if err!=nil{
			fmt.Println("Packing error",err)
			continue
		}

		_,err=conn.WriteToUDP(responseBytes,remoteAddr)
		if err!=nil{
			fmt.Println("Send Error:",err)
		}

	}
}