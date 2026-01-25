package config

var (
	currentConifg *Config
	repoRootPath  string
)

// ロード済みの設定を保持する
func SetCurrentConfig(cfg *Config) {
	currentConifg = cfg
}

// CurrentConfig は現在の設定を取得する
func CurrentConfig() *Config {
	if currentConifg == nil {
		return &Config{}
	}
	return currentConifg
}

// リポジトリのルートパスをセット
func SetCurrentRepoRoot(path string) {
	repoRootPath = path
}

func CurrentRepoRootPath() string {
	return repoRootPath
}
