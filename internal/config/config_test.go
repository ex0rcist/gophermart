package config

import (
	"os"
	"testing"
	"time"

	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewDefault(t *testing.T) {
	cfg, err := NewDefault(nil)
	assert.NoError(t, err)

	assert.Equal(t, "file://internal/storage/migrations", cfg.DB.MigrationsSource)
	assert.Equal(t, "0.0.0.0:8080", cfg.Server.Address)
	assert.Equal(t, 5*time.Second, cfg.Server.Timeout)
	assert.Equal(t, "0.0.0.0:8181", cfg.Accrual.Address)
	assert.Equal(t, 5*time.Second, cfg.Accrual.RefillInterval)
	assert.Equal(t, 5*time.Second, cfg.Accrual.Timeout)
}

func TestConfigFromEnv(t *testing.T) {
	// устанавливаем переменные окружения
	os.Setenv("DATABASE_URI", "postgres://user:password@localhost/dbname")
	os.Setenv("RUN_ADDRESS", "127.0.0.1:8080")
	os.Setenv("APP_KEY", "mysecretkey")
	os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "127.0.0.1:8181")

	defer os.Clearenv() // очищаем окружение после теста

	cfg, err := ConfigFromEnv(&Config{})
	assert.NoError(t, err)

	assert.Equal(t, "postgres://user:password@localhost/dbname", cfg.DB.DSN)
	assert.Equal(t, "127.0.0.1:8080", cfg.Server.Address)
	assert.Equal(t, entities.Secret("mysecretkey"), cfg.Server.Secret)
	assert.Equal(t, "127.0.0.1:8181", cfg.Accrual.Address)
}

func TestConfigFromFlags(t *testing.T) {
	// устанавливаем флаги
	os.Args = []string{
		"test",
		"--database=postgres://user:password@localhost/dbname",
		"--accrual-address=127.0.0.1:8181",
		"--gophermart-address=127.0.0.1:8080",
	}

	cfg, err := ConfigFromFlags(&Config{
		Server: Server{
			Secret: entities.Secret(""),
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "postgres://user:password@localhost/dbname", cfg.DB.DSN)
	assert.Equal(t, "127.0.0.1:8080", cfg.Server.Address)
	assert.Equal(t, "127.0.0.1:8181", cfg.Accrual.Address)
	assert.NotEmpty(t, cfg.Server.Secret) // должен сгенерироваться, если пустой
}

func TestValidateConfig(t *testing.T) {
	validConfig := &Config{
		DB: DB{
			DSN: "postgres://user:password@localhost/dbname",
		},
		Server: Server{
			Address: "127.0.0.1:8080",
		},
		Accrual: Accrual{
			Address: "127.0.0.1:8181",
		},
	}

	err := validateConfig(validConfig)
	assert.NoError(t, err)
}

func TestValidateConfig_InvalidServerAddress(t *testing.T) {
	invalidConfig := &Config{
		DB: DB{
			DSN: "postgres://user:password@localhost/dbname",
		},
		Server: Server{
			Address: "invalid_address",
		},
		Accrual: Accrual{
			Address: "127.0.0.1:8181",
		},
	}

	err := validateConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing URL")
}

func TestValidateConfig_EmptyDSN(t *testing.T) {
	invalidConfig := &Config{
		DB: DB{
			DSN: "",
		},
		Server: Server{
			Address: "127.0.0.1:8080",
		},
		Accrual: Accrual{
			Address: "127.0.0.1:8181",
		},
	}

	err := validateConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty DB DSN")
}
