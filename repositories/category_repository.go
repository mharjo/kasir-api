package repositories

import (
	"context"
	"errors"
	"kasir-api/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `SELECT id, name, COALESCE(description,'') FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *CategoryRepository) Create(c *models.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.QueryRow(ctx,
		`INSERT INTO categories (name, description) VALUES ($1,$2) RETURNING id`,
		c.Name, c.Description,
	).Scan(&c.ID)
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c models.Category
	err := r.db.QueryRow(ctx,
		`SELECT id, name, COALESCE(description,'') FROM categories WHERE id=$1`,
		id,
	).Scan(&c.ID, &c.Name, &c.Description)

	if err != nil {
		return nil, errors.New("category belum ada")
	}
	return &c, nil
}

func (r *CategoryRepository) Update(c *models.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ct, err := r.db.Exec(ctx,
		`UPDATE categories SET name=$1, description=$2 WHERE id=$3`,
		c.Name, c.Description, c.ID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("category belum ada")
	}
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ct, err := r.db.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("category belum ada")
	}
	return nil
}
