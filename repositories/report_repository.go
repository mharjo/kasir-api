package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

type BestSeller struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type TodayReport struct {
	TotalRevenue   int        `json:"total_revenue"`
	TotalTransaksi int        `json:"total_transaksi"`
	ProdukTerlaris BestSeller `json:"produk_terlaris"`
}

func (r *ReportRepository) GetReportByDateRange(start, end time.Time) (TodayReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var rep TodayReport

	// total_revenue + total_transaksi
	err := r.db.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) AS total_transaksi
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`, start, end).Scan(&rep.TotalRevenue, &rep.TotalTransaksi)
	if err != nil {
		return rep, err
	}

	// produk_terlaris
	var nama string
	var qty int
	err = r.db.QueryRow(ctx, `
		SELECT p.name, COALESCE(SUM(td.quantity),0) AS qty
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.name
		ORDER BY qty DESC
		LIMIT 1
	`, start, end).Scan(&nama, &qty)

	// Kalau belum ada transaksi hari itu, query ini bisa “no rows”
	if err == nil {
		rep.ProdukTerlaris = BestSeller{Nama: nama, QtyTerjual: qty}
	} else {
		rep.ProdukTerlaris = BestSeller{Nama: "", QtyTerjual: 0}
	}

	return rep, nil
}
