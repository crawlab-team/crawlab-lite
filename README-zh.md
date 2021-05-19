# Crawlab Lite

<p>
  <a href="https://hub.docker.com/r/zkqiang/crawlab-lite" target="_blank">
    <img src="https://img.shields.io/docker/pulls/zkqiang/crawlab-lite?label=pulls&logo=docker">
  </a>
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

相比较 [Crawlab](https://github.com/crawlab-team/crawlab)，该版本专注于单机上的爬虫管理，平台运行不依赖任何的外部数据库，去除了大量非必要功能。

:warning: 目前该版本仍在前期开发中，部分功能可能不稳定。

## 快速开始

#### 通过 Docker Compose 运行

1. 在任意目录下创建 `docker-compose.yml`，内容如下：

```yaml
version: '3'
services:
  master:
    image: zkqiang/crawlab-lite:latest
    container_name: master
    ports:
      - "8080:8080"
    volumes:
      - "./data:/app/data"  # 数据持久化的挂载
```

2. 在目录下运行命令：

```bash
docker-compose up -d
```

3. 访问 `http://localhost:8080`

#### 源代码运行

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
cd ../frontend
npm i && npm run serve
```

4. 访问 `http://localhost:8080`

## 截图

#### 爬虫列表

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-spider-list.png)

#### 任务列表

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-task-list.png)

#### 定时列表

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-schedule-list.png)

## 与 Crawlab 比较

| | Crawlab Lite | Crawlab |
| :---: | :---: | :---: |
| 跨语言爬虫 | ✅ | ✅ |
| 多节点部署 | ❌ | ✅ |
| 定时任务 | ✅ | ✅ |
| 查看日志 | ✅ | ✅ |
| 爬虫版本管理 | ✅ | ❌ |
| 数据统计 | ❌ | ✅ |
| 消息通知 | ❌ | ✅ |
| 在线编辑 | ❌ | ✅ |
| 可配置爬虫 | ❌ | ✅ |
| SDK | ❌ | ✅ |
