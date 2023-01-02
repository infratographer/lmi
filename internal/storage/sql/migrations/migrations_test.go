package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/stretchr/testify/assert"

	"github.com/infratographer/lmi/internal/storage/sql/migrations"
)

func TestMigrations(t *testing.T) {
	t.Parallel()

	ts, crdberr := testserver.NewTestServer()
	assert.NoError(t, crdberr)
	defer ts.Stop()

	dbdialect := "postgres"
	dbConn, dbopenerr := sql.Open(dbdialect, ts.PGURL().String())
	assert.NoError(t, dbopenerr, "failed to open db connection")

	err := migrations.Migrate(dbConn)
	assert.NoError(t, err, "failed to run migrations")
}
