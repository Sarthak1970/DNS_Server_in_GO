# DNS-GPT

DNS-GPT is an experimental **AI-enhanced DNS server written in Go** that works as a normal DNS resolver while also responding intelligently to **non-domain queries using an LLM**.

Traditional DNS servers resolve domain names to IP addresses. DNS-GPT extends this idea by adding an **AI layer** that detects queries which are not valid domains and returns conversational responses, effectively turning DNS into a **CLI-style chatbot interface**.

---

## Features:

- **Recursive DNS Resolution** – Handles standard DNS queries.
- **Caching Mechanism** – Improves performance by storing previously resolved queries.
- **Domain Blocklist Filtering** – Blocks unwanted or malicious domains.
- **AI Query Processing** – Detects non-domain queries and generates conversational responses.
- **CLI Chatbot Behavior** – Tools like `dig` or `nslookup` can trigger AI responses.

---

## How It Works?

1. A DNS query is received by the server.
2. The system checks:
   - If the domain exists in the **blocklist** → request is blocked.
   - If the query is a **valid domain** → resolved through recursive lookup.
   - If the query is **not a domain** → forwarded to the AI layer.
3. The AI generates a response which is returned to the client.

---

## 🛠 Tech Stack

- **Golang**
- **DNS Protocol**
- **LLM APIs**
- **Caching**
- **Domain Blocklists**

---


## References-
This project was inspired by experiments with DNS servers and the idea of "toying with DNS beyond simple domain resolution."

1)https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4

2)https://harshagarwal29.hashnode.dev/building-a-dns-resolver-in-golang-a-step-by-step-guide

3)https://www.youtube.com/watch?v=ANmFZ8rbmnc

---

## 🤝 Contributing & Ideas

Contributions, suggestions, and experimental ideas are welcome!

If you'd like to improve DNS-GPT, you can:
- Open an **issue** to discuss new features or ideas
- Submit a **pull request** with improvements
- Suggest **AI integrations, DNS optimizations, or caching strategies**

Interesting areas to explore:
- Latency Optimization
- Improved DNS caching mechanisms
- Additional DNS record support
- Better AI response formatting
- Security and abuse prevention

Feel free to fork the repository and experiment. Creative hacks and unconventional ideas are always encouraged !!
