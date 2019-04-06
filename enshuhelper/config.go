package enshuhelper

// Config 转推插件设置
type Config struct {
	EnshuHelper Cfg `toml:"enshuhelper"`
}

// Cfg 转推插件设置
type Cfg struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Module  string `toml:"module"`
}
