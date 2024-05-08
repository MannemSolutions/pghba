// hba provides a programmatic interface to the rules described by the PostgreSQL pg_hba.conf file using abstractions
// for connection types, databases, users, addresses, authentication methods and complete rules.
package hba

import (
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

// TODO Is this copying a structure?
func InitLogger(logger *zap.SugaredLogger) {
	log = logger
}
