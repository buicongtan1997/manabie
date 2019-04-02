package main

import (
	_ "github.com/buicongtan1997/manabie/pkg/configs"
	_ "github.com/buicongtan1997/manabie/pkg/database"
	"fmt"
	"github.com/buicongtan1997/manabie/pkg/cmd"
	"os"
)

func main() {
	if err := cmd.RunRestServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
