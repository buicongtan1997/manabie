package cmd

import (
	"github.com/buicongtan1997/manabie/pkg/protocol/rest"
)

func RunRestServer() error {
	// run HTTP gateway
	return rest.RunServer()
}
