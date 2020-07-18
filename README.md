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

[中文](https://github.com/crawlab-team/crawlab-lite/blob/master/README-zh.md) | English

Lite version of [Crawlab](https://github.com/crawlab-team/crawlab), golang-based web crawler management platform, supporting crawlers in any language.

Compared with [Crawlab](https://github.com/crawlab-team/crawlab), this lite version focuses on crawler management on a single machine, it runs independent of any external database and removes a lot of non-essential features.

:warning: This version is still in early development and some features may be unstable.

## Quick Start

#### Docker Compose

1. Create `docker-compose.yml` in any directory as follows:

```yaml
version: '3'
services:
  master:
    image: zkqiang/crawlab-lite:latest
    container_name: master
    ports:
      - "8080:8080"
    volumes:
      - "./data:/app/data"  # persistent volume
```

2. Run the command in this directory:

```bash
docker-compose up -d
```

3. Visit `http://localhost:8080`

#### Source Code

1. Clone repository

```bash
git clone https://github.com/crawlab-team/crawlab-lite
cd crawlab-lite
```

2. Run backend

```bash
cd backend
go run main.go
```

3. Run frontend

```bash
cd ../frontend
npm i && npm run serve
```

4. Visit `http://localhost:8080`

## Screenshot

#### Spider List

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-spider-list.png)

#### Task List

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-task-list.png)

#### Schedule List

![](https://github.com/crawlab-team/crawlab-docs/blob/master/assets/images/lite-schedule-list.png)
