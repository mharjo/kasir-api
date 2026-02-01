package repositories

import (
	"context"
	"errors"
	"kasir-api/models"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `SELECT id, name, price, stock, category_id FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var cat pgtype.Int8

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &cat); err != nil {
			return nil, err
		}
		if cat.Valid {
			v := int(cat.Int64)
			p.CategoryID = &v
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *ProductRepository) Create(p *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cat pgtype.Int8
	if p.CategoryID != nil {
		cat = pgtype.Int8{Int64: int64(*p.CategoryID), Valid: true}
	}

	return r.db.QueryRow(ctx,
		`INSERT INTO products (name, price, stock, category_id) VALUES ($1,$2,$3,$4) RETURNING id`,
		p.Name, p.Price, p.Stock, cat,
	).Scan(&p.ID)
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p models.Product
	var cat pgtype.Int8

	err := r.db.QueryRow(ctx,
		`SELECT id, name, price, stock, category_id FROM products WHERE id=$1`,
		id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &cat)

	if err != nil {
		return nil, errors.New("produk belum ada")
	}

	if cat.Valid {
		v := int(cat.Int64)
		p.CategoryID = &v
	}
	return &p, nil
}

func (r *ProductRepository) Update(p *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cat pgtype.Int8
	if p.CategoryID != nil {
		cat = pgtype.Int8{Int64: int64(*p.CategoryID), Valid: true}
	}

	ct, err := r.db.Exec(ctx,
		`UPDATE products SET name=$1, price=$2, stock=$3, category_id=$4 WHERE id=$5`,
		p.Name, p.Price, p.Stock, cat, p.ID,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("produk belum ada")
	}
	return nil
}

func (r *ProductRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ct, err := r.db.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("produk belum ada")
	}
	return nil
}
