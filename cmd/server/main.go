package main

import (
	"fmt"
	"net"

	"github.com/miekg/dns"

	"DNS_Server/internal/cache"
	"DNS_Server/internal/resolver"
)

func main(){

	addr := &net.UDPAddr{
		IP: net.IPv4(0,0,0,0),
		Port:8053,
	}

	conn,err := net.ListenUDP("udp",addr)

	if err != nil{
		fmt.Println("Failed to start server:",err)
		return
	}

	defer conn.Close()

	fmt.Println("DNS Server running on port:",addr.Port)

	cacheStore := cache.NewCache()

	client := resolver.Client{
		IpAddress:"8.8.8.8",
		Port:53,
	}

	for{

		buf := make([]byte,1024)

		n,remoteAddr,err := conn.ReadFromUDP(buf)

		if err != nil{
			fmt.Println("Error Reading Packet:",err)
			continue
		}

		fmt.Println("Received Packet from:",remoteAddr.String())

		queryBytes := buf[:n]

		var msg dns.Msg

		err = msg.Unpack(queryBytes)

		if err != nil{
			fmt.Println("Error decoding packet:",err)
			continue
		}

		for _,q := range msg.Question{

			domain := q.Name
			qtype := dns.TypeToString[q.Qtype]

			cacheKey := domain + ":" + qtype

			fmt.Println("Query Domain:",domain)
			fmt.Println("Query Type:",qtype)

			if resp,found := cacheStore.Get(cacheKey); found{

				fmt.Println("Cache HIT:",cacheKey)

				_,err = conn.WriteToUDP(resp,remoteAddr)

				if err != nil{
					fmt.Println("Send Error:",err)
				}

				continue
			}

			fmt.Println("Cache MISS:",cacheKey)

			upstreamResp,err := client.Query(queryBytes)

			if err != nil{
				fmt.Println("Resolver Error:",err)
				continue
			}

			cacheStore.Set(cacheKey,upstreamResp)

			_,err = conn.WriteToUDP(upstreamResp,remoteAddr)

			if err != nil{
				fmt.Println("Send Error:",err)
			}
		}
	}
}