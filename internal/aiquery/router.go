package aiquery

import "strings"

func IsAIQuery(domain string) bool {

	domain = strings.TrimSuffix(domain, ".")

	if strings.HasPrefix(domain, "ask.") {
		return true
	}

	if domain == "time.now" {
		return true
	}

	return false
}

func DomainToPrompt(domain string) string {

	domain = strings.TrimSuffix(domain, ".")

	if strings.HasPrefix(domain, "ask.") {

		prompt := strings.TrimPrefix(domain, "ask.")
		prompt = strings.ReplaceAll(prompt, "-", " ")

		return prompt
	}

	if domain == "time.now" {
		return "tell me the current time"
	}

	return domain
}