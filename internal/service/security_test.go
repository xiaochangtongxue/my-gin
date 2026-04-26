package service

import (
	"context"
	"testing"
	"time"

	"github.com/xiaochangtongxue/my-gin/pkg/cache"
)

func TestSecurityServiceBlacklistsIPOnlyAfterThreshold(t *testing.T) {
	ctx := context.Background()
	memoryCache := cache.NewMemoryCache()
	defer memoryCache.Close()

	svc := NewSecurityService(memoryCache, SecurityConfig{
		MaxAttempts:       1,
		LockDuration:      time.Minute,
		Window:            time.Minute,
		BlacklistDuration: time.Hour,
	})

	if err := svc.RecordFailure(ctx, "127.0.0.1", "10001"); err != nil {
		t.Fatalf("RecordFailure() error = %v", err)
	}

	blacklisted, err := svc.IsIPBlacklisted(ctx, "127.0.0.1")
	if err != nil {
		t.Fatalf("IsIPBlacklisted() error = %v", err)
	}
	if blacklisted {
		t.Fatal("IP should not be blacklisted after one blacklist counter increment")
	}

	for i := 0; i < 9; i++ {
		if err := svc.RecordFailure(ctx, "127.0.0.1", "10001"); err != nil {
			t.Fatalf("RecordFailure() error = %v", err)
		}
	}

	blacklisted, err = svc.IsIPBlacklisted(ctx, "127.0.0.1")
	if err != nil {
		t.Fatalf("IsIPBlacklisted() error = %v", err)
	}
	if !blacklisted {
		t.Fatal("IP should be blacklisted after ten blacklist counter increments")
	}
}
