package blocklist

import (
	"strings"
)

type BlockList struct {
	Domains map[string]bool
}

func NewBlockList() *BlockList {

	domains := map[string]bool{
		"ads.google.com.": true,
		"tracker.facebook.com.": true,
		"doubleclick.net.": true,
		"go.dev":true,
	}

	return &BlockList{
		Domains: domains,
	}
}

func (b *BlockList) IsBlocked(domain string) bool {

	domain = strings.ToLower(domain)

	if b.Domains[domain] {
		return true
	}

	return false
}