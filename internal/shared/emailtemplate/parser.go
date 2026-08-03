package emailtemplate

import "regexp"

type Parser interface {
	// Implementation for parsing projects
	ExtractVariables(input string) ([]string, error)
	SanitizeInput(input string) (string, error)
}

type parser struct {
	// Add any dependencies or configurations needed for the parser
}

func NewParser() Parser {
	return &parser{
		// Initialize any dependencies or configurations here
	}
}

var placeholderRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

func (p *parser) ExtractVariables(input string) ([]string, error) {
	matches := placeholderRe.FindAllStringSubmatch(input, -1)
	seen := map[string]bool{}
	var vars []string
	for _, m := range matches {
		if !seen[m[1]] {
			seen[m[1]] = true
			vars = append(vars, m[1])
		}
	}
	return vars, nil
}

func (p *parser) SanitizeInput(input string) (string, error) {
	// Implement the logic to sanitize the input string
	// Return the sanitized string and an error if any
	return "", nil
}
