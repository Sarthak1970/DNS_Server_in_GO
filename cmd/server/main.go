package main

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
	"github.com/joho/godotenv"

	"DNS_Server/internal/cache"
	"DNS_Server/internal/resolver"
	"DNS_Server/internal/blocklist"
	"DNS_Server/internal/aiquery"

)

func main() {

	err:=godotenv.Load()
	if err!=nil{
		fmt.Println(".env file not found")
	}

	addr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8053,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer conn.Close()

	aiClient, err := aiquery.NewGeminiClient()
	if err != nil {
		fmt.Println("Error in Gemini Client:", err)
		return
	}

	fmt.Println("DNS Server running on port:", addr.Port)

	cacheStore := cache.NewCache()
	blocker := blocklist.NewBlockList()

	client := resolver.Client{
		IpAddress: "1.1.1.1",
		Port:      53,
	}

	for {

		buf := make([]byte, 1024)

		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error Reading Packet:", err)
			continue
		}

		fmt.Println("Received Packet from:", remoteAddr.String())

		queryBytes := buf[:n]

		var msg dns.Msg

		err = msg.Unpack(queryBytes)
		if err != nil {
			fmt.Println("Error decoding packet:", err)
			continue
		}

		for _, q := range msg.Question {

			domain := q.Name
			qtype := dns.TypeToString[q.Qtype]

			fmt.Println("Query Domain:", domain)
			fmt.Println("Query Type:", qtype)
		

				if aiquery.IsAIQuery(domain) && q.Qtype == dns.TypeTXT {

					fmt.Println("AI Query Detected:", domain)

					prompt := aiquery.DomainToPrompt(domain)

					ans, err := aiClient.Ask(prompt)
					if err != nil {
						fmt.Println("Gemini error:", err)
						continue
					}

					if len(ans) > 200 {
						ans = ans[:200]
					}

					resp := dns.Msg{}
					resp.SetReply(&msg)

					txt := &dns.TXT{
						Hdr: dns.RR_Header{
							Name:   domain,
							Rrtype: dns.TypeTXT,
							Class:  dns.ClassINET,
							Ttl:    30,
						},
						Txt: []string{ans},
					}

					resp.Answer = append(resp.Answer, txt)

					responseBytes, err := resp.Pack()
					if err != nil {
						fmt.Println("Packing error:", err)
						continue
					}

					_, err = conn.WriteToUDP(responseBytes, remoteAddr)
					if err != nil {
						fmt.Println("Send error:", err)
					}

					fmt.Println("AI response sent")
			}

			if blocker.IsBlocked(domain) {

				fmt.Println("Blocked domain:", domain)

				resp := dns.Msg{}
				resp.SetReply(&msg)
				resp.Rcode = dns.RcodeNameError

				responseBytes, _ := resp.Pack()

				conn.WriteToUDP(responseBytes, remoteAddr)

				continue
			}

			cacheKey := domain + ":" + qtype

			if resp, found := cacheStore.Get(cacheKey); found {

				fmt.Println("Cache HIT:", cacheKey)

				cachedResp := make([]byte, len(resp))
				copy(cachedResp, resp)

				copy(cachedResp[0:2], queryBytes[0:2])

				_, err = conn.WriteToUDP(cachedResp, remoteAddr)
				if err != nil {
					fmt.Println("Send Error:", err)
				}

				continue
			}

			fmt.Println("Cache MISS:", cacheKey)

			upstreamResp, err := client.Query(queryBytes)
			if err != nil {
				fmt.Println("Resolver Error:", err)
				continue
			}

			cacheStore.Set(cacheKey, upstreamResp)

			_, err = conn.WriteToUDP(upstreamResp, remoteAddr)
			if err != nil {
				fmt.Println("Send Error:", err)
			}
		}
	}
}