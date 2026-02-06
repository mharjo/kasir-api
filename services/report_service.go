package services

import (
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) HariIni() (repositories.TodayReport, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	return s.repo.GetReportByDateRange(start, end)
}

func (s *ReportService) Range(start, end time.Time) (repositories.TodayReport, error) {
	return s.repo.GetReportByDateRange(start, end)
}
