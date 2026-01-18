package config

var current *Config

// SetCurrent はロード済みの設定を保持する
func SetCurrent(cfg *Config) {
	current = cfg
}

// Current は現在の設定を取得する
func Current() *Config {
	if current == nil {
		return &Config{}
	}
	return current
}
