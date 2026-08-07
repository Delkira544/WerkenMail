package emailtemplate

import (
	"fmt"
	"strconv"
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

	err = t.ValidateVariableTypes()
	if err != nil {
		return err
	}
	return nil
}

func (t *Template) ValidateVariableTypes() error {
	for _, v := range t.variables {
		switch v.Type {
		case VarTypeString:
			if v.DefaultValue != nil && *v.DefaultValue == "" {
				return fmt.Errorf("variable %q has an empty default value", v.Key)
			}
		case VarTypeNumber:
			if v.DefaultValue != nil {
				_, err := strconv.Atoi(*v.DefaultValue)
				if err != nil {
					return fmt.Errorf("variable %q has an invalid default value: %v", v.Key, err)
				}
			}
		default:
			return fmt.Errorf("variable %q has an unknown type: %q", v.Key, v.Type)
		}
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
