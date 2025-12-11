# CodeCanvas

轻量级、无依赖、可嵌入的代码画像引擎，面向 DevSecOps 场景，提供标准化的项目技术栈识别能力。

## 核心价值

✅ **快速构建项目画像**：5 秒内输出语言、框架、组件三元组

✅ **安全前置感知**：自动标记高危组件（如 Log4j）及关联 CVE

✅ **规则驱动**：新增框架/组件仅需更新 YAML 配置，无需重新编译

## 典型应用场景

| 场景 | 使用方式 |
|------|---------|
| CI 流水线 | `codecanvas analyze ./src > report.json` |
| SAST 引擎预检 | 根据识别结果动态启用 Java/JS 专用检测规则 |
| 安全运营平台 | 扫描仓库自动标记含 Log4j 的项目 |

## 目录结构

```
CodeCanvas/
├── cmd/
│   └── codecanvas/         # 命令行入口
├── pkg/
│   └── codecanvas/         # 公共 API
├── internal/
│   ├── analyzer/           # 代码分析器
│   └── ruleengine/         # 规则引擎
├── assets/
│   └── rules/              # YAML 规则文件
├── docs/                   # 文档
├── README.md
└── go.mod
```

## 快速开始

### 安装

```bash
go install github.com/codecanvas/codecanvas/cmd/codecanvas@latest
```

### 基本使用

```bash
# 分析目录
codecanvas analyze ./my-project

# 分析单文件
codecanvas analyze ./main.go

# 输出到文件
codecanvas analyze ./my-project > report.json
```

## 输出格式

```json
{
  "code_profile": {
    "path": "./example-spring",
    "total_files": 92,
    "total_lines": 14200,
    "frontend_languages": [],
    "backend_languages": ["Java"],
    "languages": [
      { "name": "Java", "files": 68, "code_lines": 11200, "comment_lines": 900, "blank_lines": 700 },
      { "name": "XML", "files": 10, "code_lines": 800, "comment_lines": 0, "blank_lines": 100 },
      { "name": "Properties", "files": 3, "code_lines": 120, "comment_lines": 10, "blank_lines": 20 }
    ]
  },
  "detection": {
    "frameworks": [
      {
        "name": "Spring Boot",
        "language": "Java",
        "confidence": "high",
        "evidence": "pom.xml contains spring-boot-starter-web"
      }
    ],
    "components": [
      {
        "name": "log4j-checker",
        "version": "2.14.1",
        "language": "Java"
      }
    ]
  },
  "timestamp": "2025-12-11T13:00:00Z",
  "codecanvas_version": "1.0.0"
}
```

## 规则扩展

CodeCanvas 的核心能力由规则驱动。所有规则定义在 `internal/embeds/*.yml` 文件中，编译时自动嵌入二进制文件。

### 规则编写指南

规则采用多文档 YAML 格式，支持以下字段：

- **name**: 框架或组件名称 (如 `Spring Boot`)
- **type**: `framework` (框架) 或 `component` (组件)
- **language**: 关联语言 (如 `Java`, `PHP`)
- **category**: 分类 (`frontend`, `backend`, `desktop`)
- **levels**: 检测策略，支持 L1 (高置信度), L2 (中), L3 (低)

#### 示例规则

```yaml
---
name: Laravel
type: framework
language: PHP
category: backend
levels:
  L1:
    paths: ["artisan"]       # 特征文件路径
    contains: []             # 必须包含的字符串
  L2:
    paths: ["composer.json"]
    contains: ["laravel/framework"]
    extract_version_from_text:  # 版本提取正则
      pattern: '"laravel/framework"\s*:\s*"([^\"]+)"'
```

### 添加规则步骤

1. 修改或新建 `internal/embeds/` 下的 `.yml` 文件。
2. 运行 `go test ./internal/canvas/...` 验证规则加载。
3. 重新构建项目。

## 开发

### 构建

```bash
go build -o codecanvas ./cmd/codecanvas
```

### 测试

```bash
go test -v ./...
```

## 许可证

MIT License
