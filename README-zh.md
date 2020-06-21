# Crawlab Lite

<p>
  <a href="https://github.com/crawlab-team/crawlab-lite/commits/master" target="_blank">
    <img src="https://img.shields.io/github/last-commit/crawlab-team/crawlab-lite.svg">
  </a>
  <a href="https://github.com/crawlab-team/crawlab-lite/issues?q=is%3Aissue+is%3Aopen+label%3Abug" target="_blank">
    <img src="https://img.shields.io/github/issues/crawlab-team/crawlab-lite/bug.svg?label=bugs&color=red">
  </a>
  <a href="https://github.com/crawlab-team/crawlab-lite/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement" target="_blank">
    <img src="https://img.shields.io/github/issues/crawlab-team/crawlab-lite/enhancement.svg?label=enhancements&color=cyan">
  </a>
  <a href="https://github.com/crawlab-team/crawlab-lite/blob/master/LICENSE" target="_blank">
    <img src="https://img.shields.io/github/license/crawlab-team/crawlab-lite.svg">
  </a>
</p>

中文 | [English](https://github.com/crawlab-team/crawlab-lite#readme)

[Crawlab](https://github.com/crawlab-team/crawlab) 轻量版本, 基于 Golang 的爬虫管理平台，支持任意语言编写的爬虫。

相比较 [Crawlab](https://github.com/crawlab-team/crawlab)，该版本专注于单机上的爬虫管理。

:warning: 目前该版本仍在开发中，代码仅用于体验。

## 运行环境
- Go 1.12+
- Node 8.12+

## 快速开始

1. 克隆仓库

```bash
git clone https://github.com/crawlab-team/crawlab-lite
cd crawlab-lite
```

2. 运行后端

```bash
cd backend
go run main.go
```

3. 运行前端

```bash
cd ..
cd frontend
npm i && npm run serve
```

## 截图

#### 爬虫列表

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-spider-list.png)

#### 任务列表

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-task-list.png)

#### 定时列表

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-schedule-list.png)
