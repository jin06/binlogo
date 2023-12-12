package service

import (
	"fmt"
	"sync"

	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var defaultAuth Authorizer

var initDefaultAuth sync.Once

func DefaultAuth() Authorizer {
	initDefaultAuth.Do(func() {
		switch viper.GetString("auth.authorizer.type") {
		case "basic":
			{
				defaultAuth = &basicAuth{
					username: viper.GetString("auth.authorizer.basic.username"),
					password: viper.GetString("auth.authorizer.basic.password"),
				}
			}
		case "ldap":
			{
				defaultAuth = &ldapAuth{
					addr:       viper.GetString("auth.authorizer.ldap.addr"),
					username:   viper.GetString("auth.authorizer.ldap.username"),
					password:   viper.GetString("auth.authorizer.ldap.password"),
					baseDN:     viper.GetString("auth.authorizer.ldap.baseDN"),
					idAttr:     viper.GetString("auth.authorizer.ldap.idAttr"),
					attributes: viper.GetStringSlice("auth.authorizer.ldap.attributs"),
				}
			}
		case "none":
			fallthrough
		default:
			defaultAuth = &noneAuth{}
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
