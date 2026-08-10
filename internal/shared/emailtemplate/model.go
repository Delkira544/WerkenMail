package emailtemplate

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/microcosm-cc/bluemonday"
)

type Template struct {
	subject   string
	bodyHTML  string
	bodyText  string
	variables []Variable
}

type Variable struct {
	Key          string
	Type         VariableType // string, number
	DefaultValue *string
}

type VariableType string

const (
	VarTypeString VariableType = "string"
	VarTypeNumber VariableType = "number"
)

var (
	placeholderRe = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)
	openBraceRe   = regexp.MustCompile(`\{\{`)
)

func New(subject, bodyHTML, bodyText string, variables []Variable) *Template {
	return &Template{
		subject:   subject,
		bodyHTML:  bodyHTML,
		bodyText:  bodyText,
		variables: variables,
	}
}

func (t *Template) GetVariables() []Variable {
	return t.variables
}

// SanitizedBodyHTML devuelve el body_html tal como quedó después de Validate() (ya sanitizado).
func (t *Template) SanitizedBodyHTML() string {
	return t.bodyHTML
}

// Validate corre todo el pipeline: texto plano -> sanitización -> extracción -> consistencia -> tipos.
func (t *Template) Validate() error {
	if err := t.validatePlainText(); err != nil {
		return err
	}
	t.sanitize() // muta t.bodyHTML in place

	combined := t.subject + " " + t.bodyHTML + " " + t.bodyText
	if err := validatePlaceholderSyntax(combined); err != nil {
		return err
	}

	extracted := t.extractPlaceholders()
	declared := map[string]Variable{}
	for _, v := range t.variables {
		declared[v.Key] = v
	}
	extractedSet := map[string]bool{}
	for _, k := range extracted {
		extractedSet[k] = true
	}

	// 1) Variables declaradas pero ausentes en subject+body_html+body_text
	for _, v := range t.variables {
		if !extractedSet[v.Key] {
			return fmt.Errorf("variable %q declared but not found in the template", v.Key)
		}
	}
	// 2) Placeholders usados sin declaración
	var missing []string
	for _, k := range extracted {
		if _, ok := declared[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("undeclared placeholder: %v", missing)
	}

	return t.validateVariableTypes()
}

// validatePlainText rechaza HTML en subject/body_text — son campos de texto plano por diseño.
func (t *Template) validatePlainText() error {
	strict := bluemonday.StrictPolicy()
	if strict.Sanitize(t.subject) != t.subject {
		return fmt.Errorf("subject must not contain HTML")
	}
	if strict.Sanitize(t.bodyText) != t.bodyText {
		return fmt.Errorf("body_text must not contain HTML")
	}
	return nil
}

// sanitize limpia body_html contra XSS. No falla si el resultado queda vacío
// (un HTML legítimamente vacío tras sanitizar no es un error).
func (t *Template) sanitize() {
	t.bodyHTML = bluemonday.UGCPolicy().Sanitize(t.bodyHTML)
}

func (t *Template) validateVariableTypes() error {
	for _, v := range t.variables {
		switch v.Type {
		case VarTypeString:
			if v.DefaultValue != nil && *v.DefaultValue == "" {
				return fmt.Errorf("variable %q has an empty default value", v.Key)
			}
		case VarTypeNumber:
			if v.DefaultValue != nil {
				if _, err := strconv.Atoi(*v.DefaultValue); err != nil {
					return fmt.Errorf("variable %q has an invalid default value: %v", v.Key, err)
				}
			}
		default:
			return fmt.Errorf("variable %q has an unknown type: %q", v.Key, v.Type)
		}
	}
	return nil
}

// extractPlaceholders busca {{key}} en subject, body_html y body_text combinados.
func (t *Template) extractPlaceholders() []string {
	text := t.subject + " " + t.bodyHTML + " " + t.bodyText
	matches := placeholderRe.FindAllStringSubmatch(text, -1)
	seen := map[string]bool{}
	var keys []string
	for _, m := range matches {
		if !seen[m[1]] {
			seen[m[1]] = true
			keys = append(keys, m[1])
		}
	}
	return keys
}

func validatePlaceholderSyntax(text string) error {
	opens := len(openBraceRe.FindAllString(text, -1))
	valid := len(placeholderRe.FindAllString(text, -1))
	if opens != valid {
		return fmt.Errorf("malformed placeholder syntax detected — expected \"{{key}}\", got unbalanced or invalid braces")
	}
	return nil
}

func (t *Template) Render(data map[string]string) (string, error) {
	return "", fmt.Errorf("not implemented")
}
