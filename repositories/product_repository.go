package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	rows, err := r.db.Query(`SELECT id, name, price, stock, category_id FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var categoryID sql.NullInt64

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID); err != nil {
			return nil, err
		}
		if categoryID.Valid {
			v := int(categoryID.Int64)
			p.CategoryID = &v
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *ProductRepository) Create(p *models.Product) error {
	var category sql.NullInt64
	if p.CategoryID != nil {
		category = sql.NullInt64{Int64: int64(*p.CategoryID), Valid: true}
	}

	return r.db.QueryRow(
		`INSERT INTO products (name, price, stock, category_id) VALUES ($1,$2,$3,$4) RETURNING id`,
		p.Name, p.Price, p.Stock, category,
	).Scan(&p.ID)
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	var categoryID sql.NullInt64

	err := r.db.QueryRow(
		`SELECT id, name, price, stock, category_id FROM products WHERE id=$1`,
		id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID)

	if err == sql.ErrNoRows {
		return nil, errors.New("produk belum ada")
	}
	if err != nil {
		return nil, err
	}
	if categoryID.Valid {
		v := int(categoryID.Int64)
		p.CategoryID = &v
	}
	return &p, nil
}

func (r *ProductRepository) Update(p *models.Product) error {
	var category sql.NullInt64
	if p.CategoryID != nil {
		category = sql.NullInt64{Int64: int64(*p.CategoryID), Valid: true}
	}

	res, err := r.db.Exec(
		`UPDATE products SET name=$1, price=$2, stock=$3, category_id=$4 WHERE id=$5`,
		p.Name, p.Price, p.Stock, category, p.ID,
	)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("produk belum ada")
	}
	return nil
}

func (r *ProductRepository) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("produk belum ada")
	}
	return nil
}
