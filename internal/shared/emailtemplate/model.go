package emailtemplate

import (
	"fmt"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

type Template struct {
	rawHTML   string
	variables []Variable
	doc       *html.Node // resultado del parseo, guardado tras Parse()
}

type Variable struct {
	Key          string
	Type         VariableType // string, url, number, etc.
	Required     bool
	DefaultValue *string
}

type VariableType string

const (
	VarTypeString VariableType = "string"
	VarTypeNumber VariableType = "number"
)

func New(rawHTML string, variables []Variable) *Template {
	return &Template{
		rawHTML:   rawHTML,
		variables: variables,
	}
}

func (t *Template) Parse() error {
	doc, err := html.ParseFragment(strings.NewReader(t.rawHTML), nil)
	if err != nil {
		return err
	}
	t.doc = &html.Node{
		Type: html.ElementNode,
		Data: "div",
	}
	for _, n := range doc {
		t.doc.AppendChild(n)
	}
	return nil
}

func (t *Template) Validate() error {
	extracted, err := t.ExtractVariables()
	if err != nil {
		return err
	}
	declared := map[string]Variable{}
	for _, v := range t.variables {
		declared[v.Key] = v
	}
	extractedSet := map[string]bool{}
	for _, k := range extracted {
		extractedSet[k] = true
	}
	// 1) Variables declaradas pero ausentes en el HTML
	for _, v := range t.variables {
		if !extractedSet[v.Key] {
			return fmt.Errorf("variable %q declarated but not found in the template", v.Key)
		}
	}
	// 2) Placeholders extraídos sin declaración
	var missing []string
	for _, k := range extracted {
		if _, ok := declared[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("undeclared placeholder: %v", missing)
	}
	return nil
}

func (t *Template) Sanitize() (string, error) {
	p := bluemonday.UGCPolicy()
	clean := p.Sanitize(t.rawHTML)
	if clean == "" {
		return "", fmt.Errorf("Error to sanitized html")
	}
	return clean, nil
}

func (t *Template) GetVariables() []Variable {
	return t.variables
}

func (t *Template) ExtractVariables() ([]string, error) {
	matches := placeholderRe.FindAllStringSubmatch(t.rawHTML, -1)
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

func (t *Template) Render(data map[string]string) (string, error) {
	return "", nil
}
