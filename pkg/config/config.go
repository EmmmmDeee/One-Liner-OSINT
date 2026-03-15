package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// General settings
	Verbose   bool
	Output    string
	Workers   int
	Timeout   time.Duration
	NoColor   bool
	CacheDir  string
	UserAgent string

	// API Keys
	APIKeys APIKeys

	// Rate limiting
	RateLimit RateLimit

	// Search settings
	Search SearchConfig

	// Output settings
	OutputConfig OutputConfig
}

// APIKeys holds API credentials for various services
type APIKeys struct {
	HaveIBeenPwned  string
	HunterIO        string
	Shodan          string
	VirusTotal      string
	GoogleCSE       string
	GoogleCSECX     string
	BingSearch      string
	EmailRepIO      string
	IPData          string
	IPGeolocation   string
	WhoisXMLAPI     string
	GitHub          string
	Twitter         string
	Facebook        string
	Instagram       string
	LinkedIn        string
}

// RateLimit configuration for API requests
type RateLimit struct {
	RequestsPerSecond int
	BurstSize         int
	RetryAttempts     int
	RetryDelay        time.Duration
}

// SearchConfig for search operations
type SearchConfig struct {
	MaxResults      int
	MaxDepth        int
	FollowRedirects bool
	UserAgents      []string
	ProxyList       []string
	UseProxy        bool
	DNSServers      []string
}

// OutputConfig for output formatting
type OutputConfig struct {
	Format       string
	Pretty       bool
	IncludeEmpty bool
	Timestamp    bool
	SaveToFile   bool
	OutputDir    string
}

// LoadConfig loads configuration from viper
func LoadConfig() *Config {
	return &Config{
		Verbose: viper.GetBool("verbose"),
		Output:  viper.GetString("output"),
		Workers: viper.GetInt("workers"),
		Timeout: time.Duration(viper.GetInt("timeout")) * time.Second,
		NoColor: viper.GetBool("no-color"),
		CacheDir: viper.GetString("cache.dir"),
		UserAgent: viper.GetString("user_agent"),

		APIKeys: APIKeys{
			HaveIBeenPwned:  viper.GetString("api_keys.haveibeenpwned"),
			HunterIO:        viper.GetString("api_keys.hunter_io"),
			Shodan:          viper.GetString("api_keys.shodan"),
			VirusTotal:      viper.GetString("api_keys.virustotal"),
			GoogleCSE:       viper.GetString("api_keys.google_cse"),
			GoogleCSECX:     viper.GetString("api_keys.google_cse_cx"),
			BingSearch:      viper.GetString("api_keys.bing_search"),
			EmailRepIO:      viper.GetString("api_keys.emailrep_io"),
			IPData:          viper.GetString("api_keys.ipdata"),
			IPGeolocation:   viper.GetString("api_keys.ipgeolocation"),
			WhoisXMLAPI:     viper.GetString("api_keys.whoisxmlapi"),
			GitHub:          viper.GetString("api_keys.github"),
			Twitter:         viper.GetString("api_keys.twitter"),
			Facebook:        viper.GetString("api_keys.facebook"),
			Instagram:       viper.GetString("api_keys.instagram"),
			LinkedIn:        viper.GetString("api_keys.linkedin"),
		},

		RateLimit: RateLimit{
			RequestsPerSecond: viper.GetInt("rate_limit.requests_per_second"),
			BurstSize:         viper.GetInt("rate_limit.burst_size"),
			RetryAttempts:     viper.GetInt("rate_limit.retry_attempts"),
			RetryDelay:        time.Duration(viper.GetInt("rate_limit.retry_delay")) * time.Second,
		},

		Search: SearchConfig{
			MaxResults:      viper.GetInt("search.max_results"),
			MaxDepth:        viper.GetInt("search.max_depth"),
			FollowRedirects: viper.GetBool("search.follow_redirects"),
			UserAgents:      viper.GetStringSlice("search.user_agents"),
			ProxyList:       viper.GetStringSlice("search.proxy_list"),
			UseProxy:        viper.GetBool("search.use_proxy"),
			DNSServers:      viper.GetStringSlice("search.dns_servers"),
		},

		OutputConfig: OutputConfig{
			Format:       viper.GetString("output.format"),
			Pretty:       viper.GetBool("output.pretty"),
			IncludeEmpty: viper.GetBool("output.include_empty"),
			Timestamp:    viper.GetBool("output.timestamp"),
			SaveToFile:   viper.GetBool("output.save_to_file"),
			OutputDir:    viper.GetString("output.output_dir"),
		},
	}
}

// SetDefaults sets default configuration values
func SetDefaults() {
	// General defaults
	viper.SetDefault("workers", 10)
	viper.SetDefault("timeout", 30)
	viper.SetDefault("cache.dir", ".osint_cache")
	viper.SetDefault("user_agent", "OSINT-Tool/1.0 (https://github.com/EmmmmDeee/One-Liner-OSINT)")

	// Rate limiting defaults
	viper.SetDefault("rate_limit.requests_per_second", 10)
	viper.SetDefault("rate_limit.burst_size", 20)
	viper.SetDefault("rate_limit.retry_attempts", 3)
	viper.SetDefault("rate_limit.retry_delay", 2)

	// Search defaults
	viper.SetDefault("search.max_results", 100)
	viper.SetDefault("search.max_depth", 3)
	viper.SetDefault("search.follow_redirects", true)
	viper.SetDefault("search.user_agents", []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})
	viper.SetDefault("search.use_proxy", false)
	viper.SetDefault("search.dns_servers", []string{"8.8.8.8", "1.1.1.1"})

	// Output defaults
	viper.SetDefault("output.format", "text")
	viper.SetDefault("output.pretty", true)
	viper.SetDefault("output.include_empty", false)
	viper.SetDefault("output.timestamp", true)
	viper.SetDefault("output.save_to_file", false)
	viper.SetDefault("output.output_dir", "./output")
}

func init() {
	SetDefaults()
}
