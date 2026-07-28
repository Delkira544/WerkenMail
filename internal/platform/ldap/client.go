package ldap

import (
	"fmt"

	"github.com/Delkira544/rakiduam/config"
	"github.com/go-ldap/ldap/v3"
)

type Client struct {
	conn   *ldap.Conn
	config config.LDAPConfig
}

// NewClient crea y devuelve una conexión LDAP lista para usar
func NewClient(cfg config.LDAPConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	var conn *ldap.Conn
	var err error

	if cfg.UseTLS {
		conn, err = ldap.DialTLS("tcp", addr, nil)
	} else {
		conn, err = ldap.Dial("tcp", addr)
	}
	if err != nil {
		return nil, fmt.Errorf("error conectando a LDAP: %w", err)
	}

	if err := conn.Bind(cfg.BindDN, cfg.BindPass); err != nil {
		conn.Close()
		return nil, fmt.Errorf("error en bind LDAP: %w", err)
	}

	return &Client{conn: conn, config: cfg}, nil
}

func (c *Client) Conn() *ldap.Conn {
	return c.conn
}

func (c *Client) Addr() string {
	return fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
