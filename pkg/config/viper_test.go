package config

import "testing"

func TestInitBindsImportantEnvironmentVariables(t *testing.T) {
	t.Setenv("APP_JWT_SECRET", "12345678901234567890123456789012")
	t.Setenv("APP_DATABASE_NAME", "env_database")
	t.Setenv("APP_DATABASE_HOST", "db.internal")
	t.Setenv("APP_REDIS_HOST", "redis.internal")
	t.Setenv("APP_REDIS_PORT", "6380")
	t.Setenv("APP_MODE", "test")

	cfg, err := Init("../../configs/config.yaml")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if cfg.Server.Mode != "test" {
		t.Fatalf("Server.Mode = %q, want test", cfg.Server.Mode)
	}
	if cfg.Database.Database != "env_database" {
		t.Fatalf("Database.Database = %q, want env_database", cfg.Database.Database)
	}
	if cfg.Database.Host != "db.internal" {
		t.Fatalf("Database.Host = %q, want db.internal", cfg.Database.Host)
	}
	if cfg.Redis.Host != "redis.internal" {
		t.Fatalf("Redis.Host = %q, want redis.internal", cfg.Redis.Host)
	}
	if cfg.Redis.Port != 6380 {
		t.Fatalf("Redis.Port = %d, want 6380", cfg.Redis.Port)
	}
}

func TestInitRejectsShortJWTSecret(t *testing.T) {
	t.Setenv("APP_JWT_SECRET", "short-secret")

	if _, err := Init("../../configs/config.yaml"); err == nil {
		t.Fatal("Init() error = nil, want jwt secret length error")
	}
}
