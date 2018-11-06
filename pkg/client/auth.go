package client

import (
	"fmt"
	"golang.org/x/net/context"
)

type AuthCreds struct {
	Token string
}

func (c *AuthCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", c.Token),
	}, nil
}

func (c *AuthCreds) RequireTransportSecurity() bool {
	return false
}
