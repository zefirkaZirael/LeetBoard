package domain

import "flag"

var (
	HelpFlag = flag.Bool("help", false, "Flag to show information about usage of project")
	Port     = flag.Int("port", 8080, "Server launch port")
)
