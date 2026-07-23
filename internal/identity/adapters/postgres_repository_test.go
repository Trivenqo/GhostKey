package adapters

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostgresIdentityRepository_FindByID_Success(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostgresIdentityRepository{DB: db}

	// Expected Database Row
	expectedID := "123"
	expectedMetadataJSON := `{"role": "admin", "active": true}`
	
	rows := sqlmock.NewRows([]string{"id", "metadata"}).
		AddRow(expectedID, []byte(expectedMetadataJSON))

	// Expect the query to be called
	query := `SELECT id, metadata FROM identities WHERE id = \$1`
	mock.ExpectQuery(query).WithArgs(expectedID).WillReturnRows(rows)

	// Execute
	identity, err := repo.FindByID(context.Background(), expectedID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, identity)
	assert.Equal(t, expectedID, identity.ID)
	assert.Equal(t, "admin", identity.Metadata["role"])
	assert.Equal(t, true, identity.Metadata["active"])

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostgresIdentityRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PostgresIdentityRepository{DB: db}

	// Expect query but return no rows
	query := `SELECT id, metadata FROM identities WHERE id = \$1`
	mock.ExpectQuery(query).WithArgs("999").WillReturnError(sql.ErrNoRows)

	identity, err := repo.FindByID(context.Background(), "999")

	assert.Error(t, err)
	assert.Equal(t, "identity not found", err.Error())
	assert.Nil(t, identity)
}