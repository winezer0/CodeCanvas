// Package analyzer 提供了 CodeCanvas 的代码画像分析功能。
package analyzer

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/winezer0/codecanvas/internal/langengine"
	"github.com/winezer0/codecanvas/internal/model"
)

// CodeAnalyzer 实现代码画像分析功能。
type CodeAnalyzer struct{}

// NewCodeAnalyzer 创建一个新的代码分析器实例。
func NewCodeAnalyzer() *CodeAnalyzer {
	return &CodeAnalyzer{}
}

// AnalysisTask 定义一个分析任务
type AnalysisTask struct {
	Path    string
	LangDef *LangDefine
}

// AnalysisResult 定义分析结果
type AnalysisResult struct {
	LangName string
	Stats    FileStats
	Err      error
}

var (
	extToLanguage  = make(map[string]*LangDefine)
	fileToLanguage = make(map[string]*LangDefine)
)

// AnalyzeCodeProfile 分析给定路径下的代码库并返回代码画像和文件索引。
func (a *CodeAnalyzer) AnalyzeCodeProfile(projectPath string) (*model.CodeProfile, *model.FileIndex, error) {
	// 获取绝对路径
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, nil, err
	}
	// 初始化langMap
	initFileInfoToLangeMap()

	// 初始化文件索引
	fileIndex := model.NewFileIndex(absPath)

	// 准备并发处理
	workers := runtime.NumCPU() / 4
	tasks := make(chan AnalysisTask, workers)
	results := make(chan AnalysisResult, workers)
	var wg sync.WaitGroup

	// 启动 Worker Pool
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				stats, err := CountFileStats(task.Path, task.LangDef)
				results <- AnalysisResult{
					LangName: task.LangDef.Name,
					Stats:    stats,
					Err:      err,
				}
			}
		}()
	}

	// 启动结果收集协程
	stats := make(map[string]*model.LangSummary)
	var errorFiles int
	done := make(chan struct{})
	go func() {
		for res := range results {
			if res.Err != nil {
				errorFiles++
				continue
			}
			summary, ok := stats[res.LangName]
			if !ok {
				summary = &model.LangSummary{Name: res.LangName}
				stats[res.LangName] = summary
			}
			summary.Count++
			summary.Code += res.Stats.Code
			summary.Comment += res.Stats.Comment
			summary.Blank += res.Stats.Blank
		}
		close(done)
	}()

	// 遍历目录并分发任务
	err = filepath.WalkDir(absPath, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			// 如果无法访问文件/目录，跳过
			return nil
		}
		if dirEntry.IsDir() {
			// 跳过隐藏目录，如 .git
			if strings.HasPrefix(dirEntry.Name(), ".") && dirEntry.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// 计算相对路径并添加到索引 (保持在主协程，无需锁)
		relPath, _ := filepath.Rel(absPath, path)
		// 统一使用 "/" 作为路径分隔符
		relPath = filepath.ToSlash(relPath)
		fileIndex.AddFile(relPath, dirEntry.Name(), filepath.Ext(dirEntry.Name()))

		// 识别语言
		langDef := parseFileLangByLangeMap(path, dirEntry)
		if langDef != nil {
			// 分发任务
			tasks <- AnalysisTask{
				Path:    path,
				LangDef: langDef,
			}
		}
		return nil
	})

	close(tasks)   // 停止发送任务
	wg.Wait()      // 等待所有 Worker 完成
	close(results) // 停止发送结果
	<-done         // 等待结果收集完成

	if err != nil {
		return nil, nil, err
	}

	codeProfile := convertToCodeProfile(absPath, stats, errorFiles)
	return codeProfile, fileIndex, nil
}

// convertToCodeProfile converts statistics to CodeCanvas CodeProfile.
func convertToCodeProfile(absPath string, stats map[string]*model.LangSummary, errorFiles int) *model.CodeProfile {

	profile := &model.CodeProfile{
		Path:              absPath,
		TotalFiles:        0,
		TotalLines:        0,
		ErrorFiles:        errorFiles,
		FrontendLanguages: []string{},
		BackendLanguages:  []string{},
		LanguageInfos:     []model.LangInfo{},
	}

	// 将统计表转换为切片
	var summaries []model.LangSummary
	for _, summary := range stats {
		summaries = append(summaries, *summary)
	}

	for _, stat := range summaries {
		langInfo := model.LangInfo{
			Name:         stat.Name,
			Files:        int(stat.Count),
			CodeLines:    int(stat.Code),
			CommentLines: int(stat.Comment),
			BlankLines:   int(stat.Blank),
		}

		// Add to profile
		profile.LanguageInfos = append(profile.LanguageInfos, langInfo)
		profile.TotalFiles += langInfo.Files
		profile.TotalLines += langInfo.CodeLines + langInfo.CommentLines + langInfo.BlankLines
	}

	// 进行语言信息优化
	langClassifier := langengine.NewLangClassifier()
	frontend, backend, desktop, other, allLang := langClassifier.DetectCategories(absPath, profile.LanguageInfos)
	profile.FrontendLanguages = frontend
	profile.BackendLanguages = backend
	profile.DesktopLanguages = desktop
	profile.OtherLanguages = other
	profile.Languages = allLang
	profile.ExpandLanguages = langengine.ExpandLanguages(allLang)
	return profile
}

// initFileInfoToLangeMap 初始化语言映射
func initFileInfoToLangeMap() {
	for i := range LangDefines {
		lang := &LangDefines[i]
		for _, ext := range lang.Extensions {
			extToLanguage[ext] = lang
		}
		for _, name := range lang.Filenames {
			fileToLanguage[name] = lang
		}
	}
}

// parseFileLangByLangeMap
func parseFileLangByLangeMap(path string, dirEntry os.DirEntry) *LangDefine {
	ext := filepath.Ext(path)
	langDef := extToLanguage[strings.ToLower(ext)]
	if langDef == nil {
		langDef = fileToLanguage[dirEntry.Name()]
	}
	return langDef
}
