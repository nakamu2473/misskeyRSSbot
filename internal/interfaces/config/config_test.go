package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadRSSURLs_Numbered(t *testing.T) {
	os.Setenv("RSS_URL_1", "https://example.tld/rss1")
	os.Setenv("RSS_URL_2", "https://example.tld/rss2")
	os.Setenv("RSS_URL_3", "https://example.tld/rss3")
	defer os.Unsetenv("RSS_URL_1")
	defer os.Unsetenv("RSS_URL_2")
	defer os.Unsetenv("RSS_URL_3")

	settings := loadRSSURLs()

	if len(settings) != 3 {
		t.Errorf("expected 3 settings, got %d", len(settings))
	}

	expected := []string{
		"https://example.tld/rss1",
		"https://example.tld/rss2",
		"https://example.tld/rss3",
	}

	for i, s := range settings {
		if s.URL != expected[i] {
			t.Errorf("URL[%d]: expected %s, got %s", i, expected[i], s.URL)
		}
		if s.Filter != false {
			t.Errorf("Filter[%d]: expected false when not set, got %v", i, s.Filter)
		}
	}
}

func TestLoadRSSURLs_NumberedWithGap(t *testing.T) {
	os.Setenv("RSS_URL_1", "https://example.tld/rss1")
	os.Setenv("RSS_URL_2", "https://example.tld/rss2")
	os.Setenv("RSS_URL_4", "https://example.tld/rss4")
	defer os.Unsetenv("RSS_URL_1")
	defer os.Unsetenv("RSS_URL_2")
	defer os.Unsetenv("RSS_URL_4")

	settings := loadRSSURLs()

	if len(settings) != 2 {
		t.Errorf("expected 2 settings, got %d", len(settings))
	}
}

func TestLoadRSSURLs_NoNumbered(t *testing.T) {
	settings := loadRSSURLs()

	if len(settings) != 0 {
		t.Errorf("expected 0 settings, got %d", len(settings))
	}
}

func TestLoadRSSURLs_WithFilter(t *testing.T) {
	os.Setenv("RSS_URL_1", "https://example.tld/rss1")
	os.Setenv("RSS_URL_1_FILTER", "true")
	os.Setenv("RSS_URL_2", "https://example.tld/rss2")
	os.Setenv("RSS_URL_2_FILTER", "false")
	os.Setenv("RSS_URL_3", "https://example.tld/rss3")
	// RSS_URL_3_FILTER は設定しない
	defer os.Unsetenv("RSS_URL_1")
	defer os.Unsetenv("RSS_URL_1_FILTER")
	defer os.Unsetenv("RSS_URL_2")
	defer os.Unsetenv("RSS_URL_2_FILTER")
	defer os.Unsetenv("RSS_URL_3")

	settings := loadRSSURLs()

	if len(settings) != 3 {
		t.Fatalf("expected 3 settings, got %d", len(settings))
	}

	if settings[0].Filter != true {
		t.Errorf("expected Filter[0] to be true, got %v", settings[0].Filter)
	}
	if settings[1].Filter != false {
		t.Errorf("expected Filter[1] to be false, got %v", settings[1].Filter)
	}
	if settings[2].Filter != false {
		t.Errorf("expected Filter[2] to be false when not set, got %v", settings[2].Filter)
	}
}

func TestGetNumberedEnvInt(t *testing.T) {
	os.Setenv("TEST_1", "100")
	os.Setenv("TEST_2", "invalid")
	defer os.Unsetenv("TEST_1")
	defer os.Unsetenv("TEST_2")

	tests := []struct {
		name         string
		prefix       string
		index        int
		defaultValue int
		expected     int
	}{
		{"valid value", "TEST", 1, 50, 100},
		{"invalid value", "TEST", 2, 50, 50},
		{"not exists", "TEST", 3, 50, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNumberedEnvInt(tt.prefix, tt.index, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestLoadConfig_NumberedRSSURLs(t *testing.T) {
	os.Setenv("MISSKEY_HOST", "test.example.tld")
	os.Setenv("AUTH_TOKEN", "test_token")
	os.Setenv("RSS_URL_1", "https://example.tld/rss1")
	os.Setenv("RSS_URL_2", "https://example.tld/rss2")

	defer os.Unsetenv("MISSKEY_HOST")
	defer os.Unsetenv("AUTH_TOKEN")
	defer os.Unsetenv("RSS_URL_1")
	defer os.Unsetenv("RSS_URL_2")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.MisskeyHost != "test.example.tld" {
		t.Errorf("expected MisskeyHost 'test.example.tld', got '%s'", cfg.MisskeyHost)
	}

	if len(cfg.RSSURL) != 2 {
		t.Errorf("expected 2 RSS URLs, got %d", len(cfg.RSSURL))
	}

	if cfg.LocalOnly != false {
		t.Errorf("expected LocalOnly to be false by default, got %v", cfg.LocalOnly)
	}
}

func TestLoadConfig_WithFilterSettings(t *testing.T) {
	os.Setenv("MISSKEY_HOST", "test.example.tld")
	os.Setenv("AUTH_TOKEN", "test_token")
	os.Setenv("RSS_URL_1", "https://example.tld/rss1")
	os.Setenv("RSS_URL_1_FILTER", "true")
	os.Setenv("RSS_URL_2", "https://example.tld/rss2")

	defer os.Unsetenv("MISSKEY_HOST")
	defer os.Unsetenv("AUTH_TOKEN")
	defer os.Unsetenv("RSS_URL_1")
	defer os.Unsetenv("RSS_URL_1_FILTER")
	defer os.Unsetenv("RSS_URL_2")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.RSSURL) != 2 {
		t.Fatalf("expected 2 RSS settings, got %d", len(cfg.RSSURL))
	}

	if cfg.RSSURL[0].URL != "https://example.tld/rss1" {
		t.Errorf("expected URL 'https://example.tld/rss1', got '%s'", cfg.RSSURL[0].URL)
	}
	if cfg.RSSURL[0].Filter != true {
		t.Errorf("expected Filter to be true, got %v", cfg.RSSURL[0].Filter)
	}
	if cfg.RSSURL[1].Filter != false {
		t.Errorf("expected Filter to be false when not set, got %v", cfg.RSSURL[1].Filter)
	}
}

func TestLoadConfig_LocalOnlyTrue(t *testing.T) {
	os.Setenv("MISSKEY_HOST", "test.example.tld")
	os.Setenv("AUTH_TOKEN", "test_token")
	os.Setenv("RSS_URL_1", "https://example.tld/rss1")
	os.Setenv("LOCAL_ONLY", "true")

	defer os.Unsetenv("MISSKEY_HOST")
	defer os.Unsetenv("AUTH_TOKEN")
	defer os.Unsetenv("RSS_URL_1")
	defer os.Unsetenv("LOCAL_ONLY")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.LocalOnly != true {
		t.Errorf("expected LocalOnly to be true, got %v", cfg.LocalOnly)
	}
}

func TestLoadConfig_NoRSSURLs(t *testing.T) {
	os.Setenv("MISSKEY_HOST", "test.example.tld")
	os.Setenv("AUTH_TOKEN", "test_token")

	defer os.Unsetenv("MISSKEY_HOST")
	defer os.Unsetenv("AUTH_TOKEN")

	_, err := LoadConfig()
	if err == nil {
		t.Error("expected error when no RSS URLs are configured, got nil")
	}
}

func TestConfig_GetCacheCleanupInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval int
		expected time.Duration
	}{
		{"default 24 hours", 24, 24 * time.Hour},
		{"custom 48 hours", 48, 48 * time.Hour},
		{"1 hour", 1, 1 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{CacheCleanupInterval: tt.interval}
			result := cfg.GetCacheCleanupInterval()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConfig_GetCacheRetentionPeriod(t *testing.T) {
	tests := []struct {
		name     string
		days     int
		expected time.Duration
	}{
		{"default 7 days", 7, 7 * 24 * time.Hour},
		{"custom 14 days", 14, 14 * 24 * time.Hour},
		{"1 day", 1, 24 * time.Hour},
		{"30 days", 30, 30 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{CacheRetentionDays: tt.days}
			result := cfg.GetCacheRetentionPeriod()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConfig_IsPersistentCache(t *testing.T) {
	tests := []struct {
		name     string
		dbPath   string
		expected bool
	}{
		{"empty path", "", false},
		{"with path", "./cache.db", true},
		{"absolute path", "/var/data/cache.db", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{CacheDBPath: tt.dbPath}
			result := cfg.IsPersistentCache()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
