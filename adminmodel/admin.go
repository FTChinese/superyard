package adminmodel

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// Admin let administrators to manage staff
type Admin struct {
	DB *sql.DB
}

var adminLogger = log.WithFields(log.Fields{
	"package":  "adminmodel",
	"resource": "Admin",
})

// VIPRoster list all vip account on ftchinese.com
func (m Admin) VIPRoster() {

}
