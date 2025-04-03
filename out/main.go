package out

import (
	"database/sql"

	"github.com/zachklingbeil/factory"
)

type Output struct {
	Factory *factory.Factory
	Db      *sql.DB
}
