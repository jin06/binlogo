package service

import (
	"fmt"
	"sync"

	"github.com/go-ldap/ldap/v3"
	"github.com/jin06/binlogo/v2/configs"
	"github.com/sirupsen/logrus"
)

var defaultAuth Authorizer

var initDefaultAuth sync.Once

func DefaultAuth() Authorizer {
	initDefaultAuth.Do(func() {
		authType := configs.Default.Auth.Type
		if authType == "" {
			authType = "none"
		}
		switch authType {
		case "basic":
			{
				defaultAuth = &basicAuth{
					configs.Default.Auth.AuthBasic.UserName,
					configs.Default.Auth.AuthBasic.Password,
				}
			}
		case "ldap":
			{
				defaultAuth = &ldapAuth{
					configs.Default.Auth.AuthLDAP.Addr,
					configs.Default.Auth.AuthLDAP.UserName,
					configs.Default.Auth.AuthLDAP.Password,
					configs.Default.Auth.AuthLDAP.BaseDN,
					configs.Default.Auth.AuthLDAP.IDAttr,
					configs.Default.Auth.AuthLDAP.Attributs,
				}
			}
		case "none":
			defaultAuth = &noneAuth{}
		default:
			panic("unknown auth type")
		}
	})
	return defaultAuth
}

type Authorizer interface {
	Verify(user, psd string) (pass bool)
}

type basicAuth struct {
	username string
	password string
}

func (auth *basicAuth) Verify(user, psd string) (pass bool) {
	if user == auth.username && psd == auth.password {
		return true
	}
	return false
}

type noneAuth struct {
}

func (auth *noneAuth) Verify(user, psd string) (pass bool) {
	logrus.WithField("username", user).WithField("password", psd).Debug("verify")
	return true
}

type ldapAuth struct {
	addr       string
	username   string
	password   string
	baseDN     string
	idAttr     string
	attributes []string
}

func (auth *ldapAuth) Verify(user, psd string) (pass bool) {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	conn, err := ldap.DialURL(auth.addr)
	if err != nil {
		return false
	}
	defer conn.Close()
	if err = conn.Bind(auth.username, auth.password); err != nil {
		return false
	}
	searchRequest := ldap.NewSearchRequest(
		auth.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(%s=%s)", auth.idAttr, ldap.EscapeFilter(user)),
		auth.attributes,
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false
	}
	if len(sr.Entries) != 1 {
		return false
	}
	if err = conn.Bind(sr.Entries[0].DN, psd); err != nil {
		return false
	}
	return true
}
