package business

import (
	"database/sql"
	"fmt"
	"vyking/settings"
)

// BusinessLayer contains the database instance and business methods
type BusinessLayer struct {
	DB *sql.DB
}

// NewBusinessLayer initializes the business layer and sets up the database connection
func NewBusinessLayer(db *sql.DB) (*BusinessLayer, error) {

	return &BusinessLayer{DB: db}, nil
}

func GetDatabaseConnection() (*sql.DB, error) {
	config, err := settings.GetDatabaseConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

// Close closes the database connection
func (b *BusinessLayer) Close() error {
	if b.DB != nil {
		return b.DB.Close()
	}
	return nil
}

// DistributePrizes handles prize distribution logic
func (b *BusinessLayer) DistributePrizes(tournamentID int) error {
	_, err := b.DB.Exec("CALL DistributePrizes(?)", tournamentID)
	if err != nil {
		return err
	}
	return nil
}

// GetRankings fetches player rankings based on account balance
func (b *BusinessLayer) GetRankings() ([]map[string]interface{}, error) {
	rows, err := b.DB.Query(`
        WITH PlayerRankings AS (
            SELECT
                PlayerID,
                Name,
                AccountBalance,
                RANK() OVER (ORDER BY AccountBalance DESC) AS Ranking
            FROM Players
        )
        SELECT * FROM PlayerRankings;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []map[string]interface{}
	for rows.Next() {
		var playerID int
		var name string
		var accountBalance float64
		var ranking int
		if err := rows.Scan(&playerID, &name, &accountBalance, &ranking); err != nil {
			return nil, err
		}
		rankings = append(rankings, map[string]interface{}{
			"PlayerID":       playerID,
			"Name":           name,
			"AccountBalance": accountBalance,
			"Ranking":        ranking,
		})
	}
	return rankings, nil
}
