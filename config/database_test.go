package config

import (
	"net/url"
	"testing"
)

func TestDatabaseDSNEncodesSpecialCharacters(t *testing.T) {
	t.Setenv("DB_HOST", "database")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "rentvibe")
	t.Setenv("DB_USER", "rentvibe")
	t.Setenv("DB_PASSWORD", "p@ss:/?#% word")
	t.Setenv("DB_SSLMODE", "disable")

	dsn := databaseDSN()
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("databaseDSN returned an invalid URL: %v", err)
	}

	password, exists := parsed.User.Password()
	if !exists || password != "p@ss:/?#% word" {
		t.Fatalf("password was not encoded safely: %q", password)
	}
	if parsed.Hostname() != "database" || parsed.Port() != "5432" {
		t.Fatalf("unexpected database host: %s", parsed.Host)
	}
}
