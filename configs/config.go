package configs

import (
	"net"
)

// ENV global environment
var ENV Env

// NodeName current node's name
var NodeName string

// NodeIP current node's ip
var NodeIP net.IP

// NodePort current node's port
var NodePort int
