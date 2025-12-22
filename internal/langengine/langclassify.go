package langengine

import (
	"strings"

	"github.com/winezer0/codecanvas/internal/embeds"
	"github.com/winezer0/codecanvas/internal/logging"
	"github.com/winezer0/codecanvas/internal/model"
	"github.com/winezer0/codecanvas/internal/utils"
)

// LangClassify 语言分类器的主结构体
// - rules: 存储所有语言分类规则的映射，键为小写的语言名称
type LangClassify struct {
	ruleMap map[string]model.LangRule
}

// NewLangClassifier 创建一个新的语言分类器实例
// 初始化分类器并加载默认规则
func NewLangClassifier() *LangClassify {
	c := &LangClassify{ruleMap: embeds.EmbeddedLangRules()}
	return c
}

// DetectCategories 检测给定语言的分类（前端/后端/桌面）
// 参数:
// - root: 项目根目录路径
// - langs: 语言信息列表
// 返回值:
// - frontend: 前端语言列表
// - backend: 后端语言列表
// - desktop: 桌面语言列表
// - all: 所有语言列表（去重）
func (c *LangClassify) DetectCategories(root string, langs []model.LangInfo) (frontend, backend, desktop, other, all []string) {
	frontedSet := make(map[string]bool)
	backendSet := make(map[string]bool)
	desktopSet := make(map[string]bool)
	otherSet := make(map[string]bool)
	allSet := make(map[string]bool) // 用于去重所有语言

	deps := readPackageJSONDeps(root)

	for _, li := range langs {
		name := strings.ToLower(li.Name)
		langRule, ok := c.ruleMap[name]
		if ok {
			cat := applyHeuristics(root, langRule, deps)
			switch cat {
			case model.CategoryFrontend:
				frontedSet[li.Name] = true
			case model.CategoryBackend:
				backendSet[li.Name] = true
			case model.CategoryDesktop:
				desktopSet[li.Name] = true
			case model.CategoryOther:
				otherSet[li.Name] = true
			}
		} else {
			otherSet[li.Name] = true
			//// 注意：Desktop 只能通过 rules 判断
			logging.Warnf("Incomplete language category for: %s", name)
		}

		// 所有语言都加入 allSet（去重）
		allSet[li.Name] = true
	}

	// 提取结果（保持顺序无关，若需排序可加 sort.Strings）
	frontend = utils.Mapkeys(frontedSet)
	backend = utils.Mapkeys(backendSet)
	desktop = utils.Mapkeys(desktopSet)
	other = utils.Mapkeys(otherSet)
	all = utils.Mapkeys(allSet)
	return
}
