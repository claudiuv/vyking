package business

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDistributePrizes(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect the stored procedure call
	mock.ExpectExec("CALL DistributePrizes(?)").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	bl, err := NewBusinessLayer(db)

	// Call the method and assert
	err = bl.DistributePrizes(1)
	assert.NoError(t, err)
}

func TestGetRankings(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect the query
	mock.ExpectQuery(`
        WITH PlayerRankings AS (
            SELECT
                PlayerID,
                Name,
                AccountBalance,
                RANK() OVER (ORDER BY AccountBalance DESC) AS Ranking
            FROM Players
        )
        SELECT * FROM PlayerRankings;
    `).
		WillReturnRows(sqlmock.NewRows([]string{"PlayerID", "Name", "AccountBalance", "Ranking"}).
			AddRow(3, "Alice", 1200.00, 1).
			AddRow(5, "Bob", 800.00, 2))

	bl, _ := NewBusinessLayer(db)

	// Call the method and assert
	rankings, err := bl.GetRankings()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rankings))
}
