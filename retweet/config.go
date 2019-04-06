package retweet

// Config 转推插件设置
type Config struct {
	Retweet Cfg `toml:"retweet"`
}

// Cfg 转推插件设置
type Cfg struct {
	Name      string      `toml:"name"`
	Version   string      `toml:"version"`
	Broadcast []Broadcast `toml:"broadcast"`
	Module    string      `toml:"module"`
	ImageRoot string      `toml:"imageRoot"`
}

// Broadcast 消息广播配置
type Broadcast struct {
	Account   string   `toml:"account"`
	Accounts  []string `toml:"accounts"`
	GroupNums []int64  `toml:"groupNums"`
}
