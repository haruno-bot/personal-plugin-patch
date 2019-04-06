package retweet

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/haruno-bot/haruno/coolq"
	"github.com/haruno-bot/haruno/logger"
)

// Retweet 转推插件
type Retweet struct {
	coolq.PbPlugin
	name      string
	version   string
	module    string
	broadcast map[string][]int64
	imageRoot string
}

// Name 插件名称
func (_plugin *Retweet) Name() string {
	return fmt.Sprintf("%s@%s", _plugin.name, _plugin.version)
}

// Module 插件名称
func (_plugin *Retweet) Module() string {
	return _plugin.module
}

func removeRepeatedString(arr []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	for _, val := range arr {
		if m[val] {
			continue
		}
		n = append(n, val)
		m[val] = true
	}
	return n
}

func removeRepeatedInteger(arr []int64) []int64 {
	m := make(map[int64]bool)
	n := make([]int64, 0)
	for _, val := range arr {
		if m[val] {
			continue
		}
		n = append(n, val)
		m[val] = true
	}
	return n
}

func (_plugin *Retweet) loadConfig() error {
	cfg := new(Config)
	_, err := toml.DecodeFile("config.toml", cfg)
	if err != nil {
		return err
	}
	pcfg := cfg.Retweet
	_plugin.name = pcfg.Name
	_plugin.module = pcfg.Module
	_plugin.imageRoot = pcfg.ImageRoot
	_plugin.version = pcfg.Version
	_plugin.broadcast = make(map[string][]int64)
	// 确定广播组
	for _, broadcast := range pcfg.Broadcast {
		accounts := make([]string, 0)
		toGroupNums := removeRepeatedInteger(broadcast.GroupNums)
		if broadcast.Account != "" {
			accounts = append(accounts, broadcast.Account)
		}
		for _, account := range broadcast.Accounts {
			accounts = append(accounts, account)
		}
		accounts = removeRepeatedString(accounts)
		for _, account := range accounts {
			if _plugin.broadcast[account] == nil {
				_plugin.broadcast[account] = make([]int64, 0)
			}
			_plugin.broadcast[account] = append(_plugin.broadcast[account], toGroupNums...)
			_plugin.broadcast[account] = removeRepeatedInteger(_plugin.broadcast[account])
		}
	}
	return nil
}

func handleAvatar(id, name, url string, groupNums []int64) {
	cqMsg := coolq.NewMessage()
	section := coolq.NewTextSection(fmt.Sprintf("%s\n更新了头像\n", name))
	cqMsg = coolq.AddSection(cqMsg, section)
	logger.Field("Plugin retweet").Infof("头像地址 = %s", url)
	section = coolq.NewImageSection(url)
	cqMsg = coolq.AddSection(cqMsg, section)
	cqMsgTxt := string(coolq.Marshal(cqMsg))
	logger.Field("Plugin retweet").Infof("向酷Q发送 -> %s", cqMsgTxt)
	for _, groupID := range groupNums {
		coolq.Client.SendGroupMsg(groupID, cqMsgTxt)
	}
	logger.Field("Plugin retweet").Successf("成功转发了一条来自%s(%s)的头像更新信息", name, id)
}

// Callback 回调接口
func (_plugin *Retweet) Callback(data string) {
	logger.Field("Plugin Retweet").Successf("%s", data)
}

// Callback 回调接口
// func (_plugin *Retweet) Callback(data string) {
// 	logger.Field("Plugin retweet").Successf("%s", data)
// 	msg := new(TweetMsg)
// 	err := json.Unmarshal([]byte(data), msg)
// 	if err != nil {
// 		logger.Field("Plugin retweet").Errorf("%s", err.Error())
// 		return
// 	}
// 	if !coolq.Client.IsAPIOk() {
// 		if msg.Cmd == "1" || msg.Cmd == "2" {
// 			logger.Field("Plugin retweet").Errorf("一条来自%s的消息被弄丢了(因为api连接没有准备好)", msg.FromName)
// 		}
// 		return
// 	}
// 	groupNums := _plugin.broadcast[msg.FromID]
// 	if len(groupNums) == 0 {
// 		return
// 	}
// 	switch msg.Cmd {
// 	case "1": // 推文
// 		// 检测是否有头像
// 		if msg.Avatar != "" {
// 			go handleAvatar(msg.FromID, msg.FromName, fmt.Sprintf("%s%s", _plugin.imageRoot, msg.Avatar), groupNums)
// 		}
// 		cqMsg := coolq.NewMessage()
// 		section := coolq.NewTextSection(msg.Text)
// 		cqMsg = coolq.AddSection(cqMsg, section)
// 		for _, img := range msg.Imgs {
// 			imgSrc := fmt.Sprintf("%s%s", _plugin.imageRoot, img)
// 			log.Printf("包含图片：%s\n", imgSrc)
// 			section = coolq.NewImageSection(imgSrc)
// 			cqMsg = coolq.AddSection(cqMsg, section)
// 		}
// 		cqMsgTxt := string(coolq.Marshal(cqMsg))
// 		logger.Field("Plugin retweet").Infof("向酷Q发送 -> %s", cqMsgTxt)
// 		for _, groupID := range groupNums {
// 			coolq.Client.SendGroupMsg(groupID, cqMsgTxt)
// 		}
// 		logger.Field("Plugin retweet").Successf("成功转发了一条来自%s(%s)的推文", msg.FromName, msg.FromID)
// 	case "2": // 头像
// 		handleAvatar(msg.FromID, msg.FromName, fmt.Sprintf("%s%s", _plugin.imageRoot, msg.Avatar), groupNums)
// 	}
// }

// Load 插件加载
func (_plugin *Retweet) Load() error {
	err := _plugin.loadConfig()
	if err != nil {
		return err
	}
	return nil
}

// Instance 转推插件实体
var Instance = &Retweet{}
