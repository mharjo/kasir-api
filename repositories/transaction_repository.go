package repositories

import (
	"context"
	"fmt"
	"kasir-api/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	totalAmount := 0
	details := make([]models.TransactionDetail, 0, len(items))

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity harus > 0 (product_id=%d)", item.ProductID)
		}

		var productName string
		var price int
		var stock int

		// ✅ Lock row agar stok aman (race-free)
		err := tx.QueryRow(ctx,
			`SELECT name, price, stock FROM products WHERE id = $1 FOR UPDATE`,
			item.ProductID,
		).Scan(&productName, &price, &stock)
		if err != nil {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("stok tidak cukup untuk %s (stok=%d, qty=%d)", productName, stock, item.Quantity)
		}

		subtotal := price * item.Quantity
		totalAmount += subtotal

		// Update stok
		_, err = tx.Exec(ctx,
			`UPDATE products SET stock = stock - $1 WHERE id = $2`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert transaction header
	var transactionID int
	err = tx.QueryRow(ctx,
		`INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id`,
		totalAmount,
	).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// ✅ TASK 3 FIX (Best practice):
	// Insert details + ambil id detail pakai RETURNING id
	for i := range details {
		details[i].TransactionID = transactionID

		var detailID int
		err = tx.QueryRow(ctx,
			`INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal,
		).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
