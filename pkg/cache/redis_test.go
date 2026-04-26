package cache

import (
	"testing"

	"github.com/xiaochangtongxue/my-gin/pkg/config"
)

func TestRedisAddr(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.RedisConfig
		want string
	}{
		{
			name: "separate host and port",
			cfg:  config.RedisConfig{Host: "127.0.0.1", Port: 6379},
			want: "127.0.0.1:6379",
		},
		{
			name: "host already includes port",
			cfg:  config.RedisConfig{Host: "redis.internal:6380", Port: 6379},
			want: "redis.internal:6380",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redisAddr(&tt.cfg); got != tt.want {
				t.Fatalf("redisAddr() = %q, want %q", got, tt.want)
			}
		})
	}
}
