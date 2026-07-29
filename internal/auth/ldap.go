package auth

import (
	"fmt"
	"log"

	ldapclient "github.com/Delkira544/rakiduam/internal/platform/ldap"
	"github.com/go-ldap/ldap/v3"
)

type LDAPGateway interface {
	Authenticate(username, password string) (*User, error)
	FindByUsername(username string) (*User, error)
}

type ldapGateway struct {
	client *ldapclient.Client
	baseDN string
}

func NewLDAPRepository(client *ldapclient.Client, baseDN string) LDAPGateway {
	return &ldapGateway{client: client, baseDN: baseDN}
}

func (r *ldapGateway) Authenticate(username, password string) (*User, error) {
	user, err := r.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	// Bind con las credenciales del usuario para validar password
	userConn, err := ldap.Dial("tcp", r.client.Addr())
	if err != nil {
		return nil, fmt.Errorf("error conectando: %w", err)
	}
	defer userConn.Close() // se cierra siempre, sin afectar la conexión admin

	err = userConn.Bind(user.DN, password)
	if err != nil {
		return nil, nil // credenciales inválidas, no es error de sistema
	}

	return user, nil
}

func (r *ldapGateway) FindByUsername(username string) (*User, error) {
	searchRequest := ldap.NewSearchRequest(
		r.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(username)),
		[]string{"dn", "cn", "mail", "uid", "gidNumber"},
		nil,
	)

	result, err := r.client.Conn().Search(searchRequest)
	if err != nil {
		log.Printf("error lo, %w", err)
		return nil, fmt.Errorf("error buscando usuario: %w", err)
	}
	if len(result.Entries) == 0 {
		log.Printf("error ")
		return nil, fmt.Errorf("usuario no encontrado")
	}

	entry := result.Entries[0]
	return &User{
		Username: entry.GetAttributeValue("uid"),
		FullName: entry.GetAttributeValue("cn"),
		Email:    entry.GetAttributeValue("mail"),
		Role:     ToRole(entry.GetAttributeValue("gidNumber")),
		DN:       entry.DN,
	}, nil
}
