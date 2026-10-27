package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/808-d/honya-api/internal/modules/order"
	"github.com/redis/go-redis/v9"
)

const (
	lowStockThreshold = 5
	statsCacheKey     = "dashboard:stats"
	statsCacheTTL     = 30 * time.Minute
)

type Service struct {
	repo *Repository
	rdb  *redis.Client
}

func NewService(repo *Repository, rdb *redis.Client) *Service {
	return &Service{repo: repo, rdb: rdb}
}

func (s *Service) Stats(ctx context.Context) *StatsResponse {
	if cached, ok := s.cachedStats(ctx); ok {
		return cached
	}

	stats := s.loadStats(ctx)
	s.cacheStats(ctx, stats)
	return stats
}

func (s *Service) cachedStats(ctx context.Context) (*StatsResponse, bool) {
	payload, err := s.rdb.Get(ctx, statsCacheKey).Bytes()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Printf("dashboard stats cache get: %v", err)
		}
		return nil, false
	}

	var stats StatsResponse
	if err := json.Unmarshal(payload, &stats); err != nil {
		log.Printf("dashboard stats cache decode: %v", err)
		return nil, false
	}

	return &stats, true
}

func (s *Service) cacheStats(ctx context.Context, stats *StatsResponse) {
	payload, err := json.Marshal(stats)
	if err != nil {
		log.Printf("dashboard stats cache encode: %v", err)
		return
	}

	if err := s.rdb.Set(ctx, statsCacheKey, payload, statsCacheTTL).Err(); err != nil {
		log.Printf("dashboard stats cache set: %v", err)
	}
}

func (s *Service) loadStats(ctx context.Context) *StatsResponse {
	products, err := s.repo.ProductStats(ctx, lowStockThreshold)
	if err != nil {
		panic(err)
	}
	orders, err := s.repo.OrderStats(ctx)
	if err != nil {
		panic(err)
	}
	statuses, err := s.repo.OrderStatusCounts(ctx)
	if err != nil {
		panic(err)
	}
	coupons, err := s.repo.CouponStats(ctx)
	if err != nil {
		panic(err)
	}
	users, err := s.repo.UserStats(ctx)
	if err != nil {
		panic(err)
	}

	orderStats := OrderStats{Total: orders.Total, PaidRevenue: orders.PaidRevenue}
	for _, row := range statuses {
		switch order.OrderStatus(row.Status) {
		case order.OrderStatusPending:
			orderStats.Pending = row.Count
		case order.OrderStatusConfirmed:
			orderStats.Confirmed = row.Count
		case order.OrderStatusProcessing:
			orderStats.Processing = row.Count
		case order.OrderStatusShipping:
			orderStats.Shipping = row.Count
		case order.OrderStatusCompleted:
			orderStats.Completed = row.Count
		case order.OrderStatusCancelled:
			orderStats.Cancelled = row.Count
		case order.OrderStatusRefunded:
			orderStats.Refunded = row.Count
		}
	}

	return &StatsResponse{
		Products: ProductStats{
			Total:    products.Total,
			Active:   products.Active,
			Draft:    products.Draft,
			Archived: products.Archived,
			LowStock: products.LowStock,
		},
		Orders:  orderStats,
		Coupons: CouponStats{Total: coupons.Total, Active: coupons.Active},
		Users:   UserStats{Total: users},
	}
}
