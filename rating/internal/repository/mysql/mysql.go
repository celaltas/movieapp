package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

type Repository struct {
	db *sql.DB
}

func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:password@/movieexample")
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id, value FROM ratings WHERE record_id = ? AND record_type = ?", recordID, recordType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ratings []model.Rating
	for rows.Next() {
		var user_id string
		var value int32
		if err := rows.Scan(&user_id, &value); err != nil {
			return nil, err
		}
		ratings = append(ratings, model.Rating{UserID: model.UserID(user_id), Value: model.RatingValue(value)})

	}
	if len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return ratings, nil

}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO ratings(record_id, record_type, user_id, value) VALUES (?, ?, ?, ?)", recordID, recordType, rating.UserID, rating.Value)
	return err
}
