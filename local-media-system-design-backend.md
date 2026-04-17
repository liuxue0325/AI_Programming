# 本地视频系统后端详细设计文档

## 1. 项目概述

### 1.1 文档目的

本文档为本地视频系统的后端详细设计文档，详细描述后端系统的架构设计、模块设计、接口设计、数据库设计、核心算法设计等内容，为后端开发提供完整的技术指导和实现规范。

### 1.2 设计依据

本文档基于《本地视频系统概要设计文档》中的后端相关需求进行细化设计，遵循概要设计中定义的技术栈选型、架构风格和功能要求。

### 1.3 范围

本文档涵盖后端系统的以下设计内容：

- 整体架构设计
- 模块职责划分
- API接口详细设计
- 数据库详细设计
- 核心业务模块设计
- 配置管理设计
- 错误处理设计
- 日志设计

## 2. 技术架构

### 2.1 技术栈

| 层级 | 技术 | 版本 | 说明 |
|------|------|------|------|
| 编程语言 | Go | 1.20+ | 高性能编译型语言 |
| Web框架 | Gin | 1.9+ | 轻量级HTTP框架 |
| 数据库 | SQLite | 3.40+ | 零配置嵌入式数据库 |
| 流媒体 | FFmpeg + Nginx | 5.0+ / 1.20+ | 视频转码与流媒体服务 |
| 配置文件 | Viper | 1.18+ | 配置管理库 |
| 日志库 | Zap | 1.26+ | 高性能日志库 |
| ORM | GORM | 1.25+ | Go语言ORM库 |
| 验证器 | Validator | 11.0+ | 参数验证库 |

### 2.2 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/
│   ├── api/
│   │   ├── router/
│   │   │   └── router.go         # 路由配置
│   │   ├── middleware/
│   │   │   ├── cors.go           # 跨域中间件
│   │   │   ├── logger.go         # 日志中间件
│   │   │   ├── recovery.go       # 异常恢复中间件
│   │   │   └── ratelimit.go      # 限流中间件
│   │   └── handler/
│   │       ├── media.go          # 媒体相关处理
│   │       ├── folder.go         # 文件夹相关处理
│   │       ├── scan.go           # 扫描相关处理
│   │       ├── setting.go        # 设置相关处理
│   │       ├── history.go        # 历史记录相关处理
│   │       ├── favorite.go       # 收藏相关处理
│   │       └── websocket.go      # WebSocket处理
│   ├── service/
│   │   ├── media.go              # 媒体服务
│   │   ├── scraper.go            # 刮削服务
│   │   ├── stream.go             # 流媒体服务
│   │   ├── scanner.go            # 扫描服务
│   │   └── setting.go            # 设置服务
│   ├── model/
│   │   ├── media.go              # 媒体模型
│   │   ├── episode.go            # 剧集模型
│   │   ├── history.go            # 历史记录模型
│   │   ├── favorite.go           # 收藏模型
│   │   └── setting.go            # 设置模型
│   ├── repository/
│   │   ├── media.go              # 媒体数据访问
│   │   ├── episode.go            # 剧集数据访问
│   │   ├── history.go            # 历史记录数据访问
│   │   ├── favorite.go           # 收藏数据访问
│   │   └── setting.go            # 设置数据访问
│   ├── pkg/
│   │   ├── database/
│   │   │   └── sqlite.go         # 数据库连接
│   │   ├── config/
│   │   │   └── config.go         # 配置加载
│   │   ├── utils/
│   │   │   ├── path.go           # 路径工具
│   │   │   ├── file.go           # 文件工具
│   │   │   └── hash.go           # 哈希工具
│   │   └── response/
│   │       └── response.go       # 统一响应
│   └── dto/
│       ├── media.go              # 媒体DTO
│       ├── request.go            # 请求DTO
│       └── response.go            # 响应DTO
├── config/
│   └── config.yaml               # 配置文件
├── go.mod
├── go.sum
└── Makefile                      # 构建脚本
```

### 2.3 架构分层

```
┌─────────────────────────────────────────────────────┐
│                    Handler层                         │
│  (HTTP请求处理、参数验证、调用Service、返回响应)      │
├─────────────────────────────────────────────────────┤
│                    Service层                         │
│  (业务逻辑处理、事务管理、业务规则实现)                │
├─────────────────────────────────────────────────────┤
│                   Repository层                       │
│  (数据访问、数据库操作、缓存操作)                      │
├─────────────────────────────────────────────────────┤
│                     Model层                          │
│  (数据结构定义、数据转换)                              │
└─────────────────────────────────────────────────────┘
```

## 3. 数据库设计

### 3.1 数据库配置

```yaml
# config.yaml
database:
  path: ./data/media.db
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 3600
```

### 3.2 数据库表结构

#### 3.2.1 媒体表 (media)

```sql
CREATE TABLE IF NOT EXISTS media (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK(type IN ('movie', 'tv', 'anime')),
    title TEXT NOT NULL,
    original_title TEXT,
    year INTEGER,
    path TEXT NOT NULL UNIQUE,
    poster_path TEXT,
    backdrop_path TEXT,
    overview TEXT,
    rating REAL DEFAULT 0,
    runtime INTEGER DEFAULT 0,
    genres TEXT,
    tmdb_id INTEGER,
    imdb_id TEXT,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_type (type),
    INDEX idx_year (year),
    INDEX idx_rating (rating),
    INDEX idx_tmdb_id (tmdb_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);
```

#### 3.2.2 剧集表 (episodes)

```sql
CREATE TABLE IF NOT EXISTS episodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    series_id INTEGER NOT NULL,
    season_number INTEGER NOT NULL,
    episode_number INTEGER NOT NULL,
    title TEXT,
    path TEXT NOT NULL UNIQUE,
    poster_path TEXT,
    overview TEXT,
    runtime INTEGER DEFAULT 0,
    tmdb_id INTEGER,
    air_date TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES media(id) ON DELETE CASCADE,
    UNIQUE (series_id, season_number, episode_number),
    INDEX idx_series_id (series_id),
    INDEX idx_season (series_id, season_number)
);
```

#### 3.2.3 观看历史表 (watch_history)

```sql
CREATE TABLE IF NOT EXISTS watch_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_id INTEGER,
    episode_id INTEGER,
    progress REAL DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    last_watched TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE CASCADE,
    FOREIGN KEY (episode_id) REFERENCES episodes(id) ON DELETE CASCADE,
    UNIQUE (media_id, episode_id),
    INDEX idx_media_id (media_id),
    INDEX idx_episode_id (episode_id),
    INDEX idx_last_watched (last_watched)
);
```

#### 3.2.4 收藏表 (favorites)

```sql
CREATE TABLE IF NOT EXISTS favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_id INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (media_id) REFERENCES media(id) ON DELETE CASCADE,
    INDEX idx_media_id (media_id)
);
```

#### 3.2.5 设置表 (settings)

```sql
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT,
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 3.2.6 扫描任务表 (scan_tasks)

```sql
CREATE TABLE IF NOT EXISTS scan_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    folder_path TEXT NOT NULL,
    status TEXT DEFAULT 'pending' CHECK(status IN ('pending', 'scanning', 'completed', 'failed')),
    total_files INTEGER DEFAULT 0,
    processed_files INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_folder_path (folder_path)
);
```

### 3.3 数据模型定义

```go
// internal/model/media.go
type MediaType string

const (
    MediaTypeMovie MediaType = "movie"
    MediaTypeTV    MediaType = "tv"
    MediaTypeAnime MediaType = "anime"
)

type MediaStatus string

const (
    MediaStatusPending   MediaStatus = "pending"
    MediaStatusScraped   MediaStatus = "scraped"
    MediaStatusFailed    MediaStatus = "failed"
)

type Media struct {
    ID            int64       `gorm:"primaryKey;autoIncrement" json:"id"`
    Type          MediaType   `gorm:"type:text;not null;index" json:"type"`
    Title         string      `gorm:"type:text;not null" json:"title"`
    OriginalTitle string      `gorm:"type:text" json:"original_title,omitempty"`
    Year          int         `gorm:"type:integer;index" json:"year,omitempty"`
    Path          string      `gorm:"type:text;not null;uniqueIndex" json:"path"`
    PosterPath    string      `gorm:"type:text" json:"poster_path,omitempty"`
    BackdropPath  string      `gorm:"type:text" json:"backdrop_path,omitempty"`
    Overview      string      `gorm:"type:text" json:"overview,omitempty"`
    Rating        float64     `gorm:"type:real;default:0;index" json:"rating,omitempty"`
    Runtime       int         `gorm:"type:integer;default:0" json:"runtime,omitempty"`
    Genres        string      `gorm:"type:text" json:"genres,omitempty"`
    TmdbID        int64       `gorm:"type:integer;index" json:"tmdb_id,omitempty"`
    ImdbID        string      `gorm:"type:text" json:"imdb_id,omitempty"`
    Status        MediaStatus `gorm:"type:text;default:'pending';index" json:"status"`
    Seasons       []Season    `gorm:"foreignKey:SeriesID" json:"seasons,omitempty"`
    CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

type Season struct {
    ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
    SeriesID     int64      `gorm:"not null;index" json:"series_id"`
    SeasonNumber int        `gorm:"not null" json:"season_number"`
    Title        string     `gorm:"type:text" json:"title,omitempty"`
    PosterPath   string     `gorm:"type:text" json:"poster_path,omitempty"`
    Episodes     []Episode  `gorm:"foreignKey:SeriesID,SeasonNumber" json:"episodes,omitempty"`
    CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Episode struct {
    ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    SeriesID      int64     `gorm:"not null;index:idx_series_id" json:"series_id"`
    SeasonNumber  int       `gorm:"not null" json:"season_number"`
    EpisodeNumber int       `gorm:"not null" json:"episode_number"`
    Title         string    `gorm:"type:text" json:"title,omitempty"`
    Path          string    `gorm:"type:text;not null;uniqueIndex" json:"path"`
    PosterPath    string    `gorm:"type:text" json:"poster_path,omitempty"`
    Overview      string    `gorm:"type:text" json:"overview,omitempty"`
    Runtime       int       `gorm:"type:integer;default:0" json:"runtime,omitempty"`
    TmdbID        int64     `gorm:"type:integer" json:"tmdb_id,omitempty"`
    AirDate       string    `gorm:"type:text" json:"air_date,omitempty"`
    CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

## 4. API接口设计

### 4.1 接口规范

#### 4.1.1 请求格式

- Content-Type: application/json
- 字符编码: UTF-8
- 请求方法: RESTful风格

#### 4.1.2 响应格式

```json
{
    "code": 0,
    "message": "success",
    "data": {}
}
```

#### 4.1.3 响应码定义

| 响应码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 资源不存在 |
| 1003 | 权限不足 |
| 1004 | 资源已存在 |
| 2001 | 服务器内部错误 |
| 2002 | 数据库错误 |
| 2003 | 文件操作错误 |
| 2004 | 外部服务错误 |

### 4.2 媒体接口

#### 4.2.1 获取媒体列表

```
GET /api/media
```

**请求参数**

| 参数名 | 类型 | 位置 | 必填 | 说明 |
|--------|------|------|------|------|
| type | string | query | 否 | 媒体类型 (movie/tv/anime) |
| year | int | query | 否 | 发行年份 |
| genre | string | query | 否 | 题材 |
| status | string | query | 否 | 刮削状态 |
| keyword | string | query | 否 | 搜索关键字 |
| page | int | query | 否 | 页码 (默认1) |
| page_size | int | query | 否 | 每页数量 (默认20) |
| sort | string | query | 否 | 排序字段 (created_at/rating/year) |
| order | string | query | 否 | 排序方向 (asc/desc) |

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "type": "movie",
                "title": "肖申克的救赎",
                "original_title": "The Shawshank Redemption",
                "year": 1994,
                "path": "/media/movies/Shawshark.Redemption.1994.1080p.BluRay.x264.mkv",
                "poster_path": "/posters/1.jpg",
                "backdrop_path": "/backdrops/1.jpg",
                "overview": "一场谋杀案使银行家安迪蒙冤入狱...",
                "rating": 9.7,
                "runtime": 142,
                "genres": "剧情,犯罪",
                "status": "scraped"
            }
        ],
        "total": 100,
        "page": 1,
        "page_size": 20,
        "total_pages": 5
    }
}
```

#### 4.2.2 获取媒体详情

```
GET /api/media/:id
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "type": "movie",
        "title": "肖申克的救赎",
        "original_title": "The Shawshank Redemption",
        "year": 1994,
        "path": "/media/movies/Shawshark.Redemption.1994.1080p.BluRay.x264.mkv",
        "poster_path": "/posters/1.jpg",
        "backdrop_path": "/backdrops/1.jpg",
        "overview": "一场谋杀案使银行家安迪蒙冤入狱...",
        "rating": 9.7,
        "runtime": 142,
        "genres": "剧情,犯罪",
        "tmdb_id": 278,
        "imdb_id": "tt0111161",
        "status": "scraped",
        "seasons": [],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 4.2.3 获取媒体播放信息

```
GET /api/media/:id/play
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "media_id": 1,
        "type": "movie",
        "play_url": "http://localhost:8090/stream/1.m3u8",
        "subtitles": [
            {
                "id": 1,
                "language": "zh-CN",
                "name": "简体中文",
                "url": "http://localhost:8080/api/media/1/subtitles/1"
            },
            {
                "id": 2,
                "language": "en",
                "name": "English",
                "url": "http://localhost:8080/api/media/1/subtitles/2"
            }
        ],
        "transcode_needed": true,
        "quality": "1080p",
        "bitrate": 5000000
    }
}
```

#### 4.2.4 获取剧集列表

```
GET /api/media/:id/episodes
```

**请求参数**

| 参数名 | 类型 | 位置 | 必填 | 说明 |
|--------|------|------|------|------|
| season | int | query | 否 | 季数 |

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "series_id": 2,
        "seasons": [
            {
                "season_number": 1,
                "title": "第1季",
                "poster_path": "/posters/2/season1.jpg",
                "episodes": [
                    {
                        "id": 10,
                        "season_number": 1,
                        "episode_number": 1,
                        "title": "Pilot",
                        "path": "/media/tv/show/S01E01.mkv",
                        "overview": "首播集...",
                        "runtime": 45,
                        "status": "watched"
                    }
                ]
            }
        ]
    }
}
```

#### 4.2.5 获取字幕列表

```
GET /api/media/:id/subtitles
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "media_id": 1,
        "subtitles": [
            {
                "id": 1,
                "language": "zh-CN",
                "name": "简体中文",
                "format": "srt",
                "path": "/subtitles/1_zh.srt"
            }
        ]
    }
}
```

#### 4.2.6 下载字幕

```
GET /api/media/:id/subtitles/:subtitle_id
```

**响应**: 返回字幕文件流

#### 4.2.7 删除媒体

```
DELETE /api/media/:id
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": null
}
```

### 4.3 文件夹接口

#### 4.3.1 获取文件夹列表

```
GET /api/folders
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "folders": [
            {
                "path": "/media/movies",
                "name": "电影",
                "type": "movie",
                "media_count": 50,
                "last_scan": "2024-01-01T00:00:00Z"
            }
        ]
    }
}
```

#### 4.3.2 添加媒体文件夹

```
POST /api/folders
```

**请求体**

```json
{
    "path": "/media/movies",
    "name": "电影",
    "type": "movie"
}
```

#### 4.3.3 删除媒体文件夹

```
DELETE /api/folders/:id
```

### 4.4 扫描接口

#### 4.4.1 触发媒体扫描

```
POST /api/scan
```

**请求体**

```json
{
    "folder_path": "/media/movies",
    "force": false
}
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "task_id": 1,
        "status": "pending"
    }
}
```

#### 4.4.2 获取扫描进度

```
GET /api/scan/:task_id
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "task_id": 1,
        "folder_path": "/media/movies",
        "status": "scanning",
        "total_files": 100,
        "processed_files": 45,
        "progress": 45.0,
        "started_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 4.4.3 取消扫描任务

```
DELETE /api/scan/:task_id
```

### 4.5 刮削接口

#### 4.5.1 手动刮削媒体

```
POST /api/media/:id/scrape
```

**请求体**

```json
{
    "source": "tmdb"
}
```

#### 4.5.2 批量刮削

```
POST /api/scrape/batch
```

**请求体**

```json
{
    "media_ids": [1, 2, 3],
    "source": "tmdb"
}
```

### 4.6 设置接口

#### 4.6.1 获取系统设置

```
GET /api/settings
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "library": {
            "media_folders": ["/media/movies", "/media/tv"],
            "scan_on_startup": true,
            "auto_scrape": true
        },
        "scraper": {
            "tmdb_api_key": "xxx",
            "preferred_language": "zh-CN",
            "auto_download_poster": true,
            "auto_download_backdrop": true
        },
        "stream": {
            "default_quality": "1080p",
            "transcode_enabled": true,
            "hardware_acceleration": true,
            "advertisement_enabled": false
        },
        "system": {
            "theme": "dark",
            "language": "zh-CN",
            "log_level": "info"
        }
    }
}
```

#### 4.6.2 更新系统设置

```
PUT /api/settings
```

**请求体**

```json
{
    "library.media_folders": ["/media/movies", "/media/tv"],
    "scraper.preferred_language": "zh-CN"
}
```

### 4.7 历史记录接口

#### 4.7.1 获取观看历史

```
GET /api/history
```

**请求参数**

| 参数名 | 类型 | 位置 | 必填 | 说明 |
|--------|------|------|------|------|
| page | int | query | 否 | 页码 |
| page_size | int | query | 否 | 每页数量 |

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "media_id": 1,
                "episode_id": null,
                "media": {
                    "id": 1,
                    "title": "肖申克的救赎",
                    "poster_path": "/posters/1.jpg"
                },
                "progress": 0.75,
                "completed": false,
                "last_watched": "2024-01-01T00:00:00Z"
            }
        ],
        "total": 50
    }
}
```

#### 4.7.2 记录观看进度

```
POST /api/history
```

**请求体**

```json
{
    "media_id": 1,
    "episode_id": null,
    "progress": 0.75,
    "completed": false
}
```

#### 4.7.3 清空观看历史

```
DELETE /api/history
```

### 4.8 收藏接口

#### 4.8.1 获取收藏列表

```
GET /api/favorites
```

**响应示例**

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "media_id": 1,
                "media": {
                    "id": 1,
                    "title": "肖申克的救赎",
                    "poster_path": "/posters/1.jpg",
                    "rating": 9.7
                },
                "created_at": "2024-01-01T00:00:00Z"
            }
        ],
        "total": 10
    }
}
```

#### 4.8.2 添加收藏

```
POST /api/favorites
```

**请求体**

```json
{
    "media_id": 1
}
```

#### 4.8.3 删除收藏

```
DELETE /api/favorites/:media_id
```

### 4.9 WebSocket接口

#### 4.9.1 实时通知

```
WS /api/ws
```

**消息格式**

```json
{
    "type": "scan_progress",
    "data": {
        "task_id": 1,
        "progress": 45.0
    }
}
```

**消息类型**

| 类型 | 说明 |
|------|------|
| scan_progress | 扫描进度更新 |
| scrape_progress | 刮削进度更新 |
| media_updated | 媒体信息更新 |
| notification | 系统通知 |

## 5. 核心模块设计

### 5.1 媒体服务 (MediaService)

#### 5.1.1 模块职责

- 媒体文件管理
- 媒体信息存储与查询
- 播放控制
- 字幕管理

#### 5.1.2 接口定义

```go
type MediaService interface {
    // 媒体管理
    GetMediaList(ctx context.Context, req *MediaListRequest) (*MediaListResponse, error)
    GetMediaByID(ctx context.Context, id int64) (*Media, error)
    GetMediaByPath(ctx context.Context, path string) (*Media, error)
    DeleteMedia(ctx context.Context, id int64) error
    
    // 剧集管理
    GetEpisodes(ctx context.Context, seriesID int64, season int) ([]Episode, error)
    GetEpisodeByID(ctx context.Context, id int64) (*Episode, error)
    
    // 播放信息
    GetPlayInfo(ctx context.Context, mediaID int64, episodeID *int64) (*PlayInfo, error)
    
    // 字幕管理
    GetSubtitles(ctx context.Context, mediaID int64) ([]Subtitle, error)
    DownloadSubtitle(ctx context.Context, mediaID int64, subtitleID int64) ([]byte, error)
    
    // 搜索
    SearchMedia(ctx context.Context, keyword string) ([]Media, error)
}
```

#### 5.1.3 核心算法

**文件名解析算法**

```go
func ParseMediaInfo(filename string) (title, year, quality string, is3D bool) {
    re := regexp.MustCompile(`(?i)(.+?)[.\s]+(19|20)\d{2}[.\s]`)
    if matches := re.FindStringSubmatch(filename); len(matches) > 2 {
        title = matches[1]
        year = matches[2]
    }
    
    qualityPattern := regexp.MustCompile(`(?i)(2160p|1080p|720p|480p|4k|hd)`)
    if match := qualityPattern.FindString(filename); match != "" {
        quality = strings.ToUpper(match)
    }
    
    is3D = regexp.MustCompile(`(?i)3d|half-ou|half-su`).MatchString(filename)
    return
}
```

### 5.2 扫描服务 (ScannerService)

#### 5.2.1 模块职责

- 媒体文件扫描
- 文件分类识别
- 新增文件检测
- 扫描任务管理

#### 5.2.2 接口定义

```go
type ScannerService interface {
    // 扫描管理
    StartScan(ctx context.Context, folderPath string, force bool) (*ScanTask, error)
    GetScanProgress(ctx context.Context, taskID int64) (*ScanProgress, error)
    CancelScan(ctx context.Context, taskID int64) error
    
    // 文件扫描
    ScanFolder(ctx context.Context, folderPath string) ([]MediaFile, error)
    ScanFile(ctx context.Context, filePath string) (*MediaFile, error)
    
    // 文件过滤
    IsMediaFile(filename string) bool
    GetMediaType(filename string) (MediaType, error)
}
```

#### 5.2.3 核心算法

**媒体文件识别算法**

```go
var videoExtensions = map[string]bool{
    ".mkv": true, ".mp4": true, ".avi": true, ".mov": true,
    ".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
    ".mpg": true, ".mpeg": true, ".3gp": true,
}

var moviePatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)(.+?)[.\s]+(19|20)\d{2}[.\s].*\.(mkv|mp4|avi)`),
    regexp.MustCompile(`(?i)(.+?)[.\s]+(2160p|1080p|720p|480p)[.\s].*\.(mkv|mp4|avi)`),
}

var tvPatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)S(\d{1,2})E(\d{1,2})`),
    regexp.MustCompile(`(?i)Season[.\s]*(\d{1,2})[.\s]*Episode[.\s]*(\d{1,2})`),
    regexp.MustCompile(`(?i)(\d{1,2})x(\d{1,2})`),
}

func (s *ScannerService) IsMediaFile(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    return videoExtensions[ext]
}

func (s *ScannerService) GetMediaType(filename string) (MediaType, error) {
    basename := filepath.Base(filename)
    
    for _, pattern := range tvPatterns {
        if pattern.MatchString(basename) {
            return MediaTypeTV, nil
        }
    }
    
    for _, pattern := range moviePatterns {
        if pattern.MatchString(basename) {
            return MediaTypeMovie, nil
        }
    }
    
    return "", ErrUnknownMediaType
}
```

### 5.3 刮削服务 (ScraperService)

#### 5.3.1 模块职责

- 元数据获取
- TMDB API集成
- 封面下载
- 媒体信息匹配

#### 5.3.2 接口定义

```go
type ScraperService interface {
    // 刮削
    ScrapeMedia(ctx context.Context, mediaID int64, source string) error
    ScrapeBatch(ctx context.Context, mediaIDs []int64, source string) error
    ScrapeAll(ctx context.Context, force bool) error
    
    // TMDB
    SearchTMDB(ctx context.Context, query string, mediaType MediaType) ([]TMDBSearchResult, error)
    GetTMDBDetails(ctx context.Context, tmdbID int64, mediaType MediaType) (*TMDBDetails, error)
    
    // 下载
    DownloadPoster(ctx context.Context, url string, mediaID int64) (string, error)
    DownloadBackdrop(ctx context.Context, url string, mediaID int64) (string, error)
}
```

#### 5.3.3 TMDB API集成

```go
type TMDBClient struct {
    APIKey    string
    BaseURL   string
    Language  string
    HTTPClient *http.Client
}

func (c *TMDBClient) SearchMovie(ctx context.Context, query string) ([]TMDBSearchResult, error) {
    url := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s&language=%s",
        c.BaseURL, c.APIKey, url.QueryEscape(query), c.Language)
    
    resp, err := c.HTTPClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result TMDBResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Results, nil
}

func (c *TMDBClient) GetMovieDetails(ctx context.Context, movieID int64) (*TMDBMovieDetails, error) {
    url := fmt.Sprintf("%s/movie/%d?api_key=%s&language=%s&append_to_response=credits",
        c.BaseURL, movieID, c.APIKey, c.Language)
    
    resp, err := c.HTTPClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result TMDBMovieDetails
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

### 5.4 流媒体服务 (StreamService)

#### 5.4.1 模块职责

- 视频转码
- HLS流生成
- 自适应码率
- 流媒体分发

#### 5.4.2 接口定义

```go
type StreamService interface {
    // 播放
    GetStreamURL(ctx context.Context, mediaID int64, episodeID *int64, quality string) (string, error)
    Transcode(ctx context.Context, inputPath string, outputPath string, options *TranscodeOptions) error
    
    // 转码管理
    StartTranscodeTask(ctx context.Context, mediaID int64, episodeID *int64, options *TranscodeOptions) (string, error)
    GetTranscodeProgress(ctx context.Context, taskID string) (float64, error)
    CancelTranscodeTask(ctx context.Context, taskID string) error
    
    // 清理
    CleanupOldStreams(ctx context.Context) error
}
```

#### 5.4.3 转码配置

```go
type TranscodeOptions struct {
    Quality      string  // 720p, 1080p, 4k
    VideoCodec   string  // h264, h265, vp9
    AudioCodec   string  // aac, mp3, flac
    VideoBitrate string  // 2000k, 5000k
    AudioBitrate string  // 128k, 320k
    Resolution   string  // 1920x1080
    HardwareAccel bool   // 硬件加速
}

var QualityPresets = map[string]TranscodeOptions{
    "4k": {
        VideoCodec:   "h265",
        VideoBitrate: "15000k",
        Resolution:   "3840x2160",
    },
    "1080p": {
        VideoCodec:   "h264",
        VideoBitrate: "5000k",
        Resolution:   "1920x1080",
    },
    "720p": {
        VideoCodec:   "h264",
        VideoBitrate: "2500k",
        Resolution:   "1280x720",
    },
    "480p": {
        VideoCodec:   "h264",
        VideoBitrate: "1000k",
        Resolution:   "854x480",
    },
}
```

### 5.5 设置服务 (SettingService)

#### 5.5.1 模块职责

- 系统配置管理
- 配置验证
- 配置持久化

#### 5.5.2 接口定义

```go
type SettingService interface {
    // 获取设置
    GetSettings(ctx context.Context) (*SystemSettings, error)
    GetSetting(ctx context.Context, key string) (string, error)
    
    // 更新设置
    UpdateSettings(ctx context.Context, settings map[string]interface{}) error
    UpdateSetting(ctx context.Context, key string, value interface{}) error
    
    // 重置
    ResetSettings(ctx context.Context) error
    
    // 验证
    ValidateSettings(ctx context.Context, settings *SystemSettings) error
}
```

#### 5.5.3 配置结构

```go
type SystemSettings struct {
    Library LibrarySettings `json:"library"`
    Scraper ScraperSettings `json:"scraper"`
    Stream  StreamSettings   `json:"stream"`
    System  SystemSettings   `json:"system"`
}

type LibrarySettings struct {
    MediaFolders     []string `json:"media_folders"`
    ScanOnStartup    bool     `json:"scan_on_startup"`
    AutoScrape       bool     `json:"auto_scrape"`
    ExcludePatterns  []string `json:"exclude_patterns"`
    FileSizeLimit    int64    `json:"file_size_limit"`
}

type ScraperSettings struct {
    TMDBAPIKey          string   `json:"tmdb_api_key"`
    PreferredLanguage  string   `json:"preferred_language"`
    AutoDownloadPoster  bool     `json:"auto_download_poster"`
    AutoDownloadBackdrop bool    `json:"auto_download_backdrop"`
    AutoDownloadSubtitle bool    `json:"auto_download_subtitle"`
    AdditionalSources   []string `json:"additional_sources"`
}

type StreamSettings struct {
    DefaultQuality      string `json:"default_quality"`
    TranscodeEnabled    bool   `json:"transcode_enabled"`
    HardwareAccel       bool   `json:"hardware_acceleration"`
    AdvertisementEnabled bool  `json:"advertisement_enabled"`
    MaxConcurrent       int    `json:"max_concurrent_streams"`
}

type SystemSettings struct {
    Theme        string `json:"theme"`
    Language     string `json:"language"`
    LogLevel     string `json:"log_level"`
    EnableUPNP   bool   `json:"enable_upnp"`
    EnableSharing bool  `json:"enable_sharing"`
}
```

## 6. 中间件设计

### 6.1 CORS中间件

```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

### 6.2 日志中间件

```go
func Logger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()

        logger.Info("request",
            zap.Int("status", status),
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.String("query", query),
            zap.String("ip", c.ClientIP()),
            zap.Duration("latency", latency),
        )
    }
}
```

### 6.3 限流中间件

```go
type RateLimiter struct {
    limiter  *rate.Limiter
    bursts   int
}

func NewRateLimiter(rps float64, bursts int) *RateLimiter {
    return &RateLimiter{
        limiter:  rate.NewLimiter(rate.Limit(rps), bursts),
        bursts:   bursts,
    }
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !rl.limiter.Allow() {
            c.JSON(429, gin.H{
                "code":    1005,
                "message": "rate limit exceeded",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## 7. 错误处理设计

### 7.1 错误定义

```go
var (
    ErrNotFound           = errors.New("resource not found")
    ErrInvalidParameter   = errors.New("invalid parameter")
    ErrUnauthorized       = errors.New("unauthorized")
    ErrForbidden          = errors.New("forbidden")
    ErrMediaFileNotFound  = errors.New("media file not found")
    ErrScrapeFailed       = errors.New("scrape failed")
    ErrTranscodeFailed     = errors.New("transcode failed")
    ErrDatabaseError      = errors.New("database error")
    ErrExternalAPIError   = errors.New("external API error")
)

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}
```

### 7.2 错误处理中间件

```go
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            var appErr *AppError
            switch {
            case errors.Is(err.Err, ErrNotFound):
                appErr = &AppError{Code: 1002, Message: err.Error()}
            case errors.Is(err.Err, ErrInvalidParameter):
                appErr = &AppError{Code: 1001, Message: err.Error()}
            case errors.Is(err.Err, ErrDatabaseError):
                appErr = &AppError{Code: 2002, Message: "database error"}
                logger.Error("database error", zap.Error(err.Err))
            default:
                appErr = &AppError{Code: 2001, Message: "internal server error"}
                logger.Error("internal error", zap.Error(err.Err))
            }
            
            c.JSON(appErr.Code, appErr)
        }
    }
}
```

## 8. 日志设计

### 8.1 日志配置

```go
var logger *zap.Logger

func InitLogger(level string) {
    var zapLevel zapcore.Level
    switch level {
    case "debug":
        zapLevel = zapcore.DebugLevel
    case "info":
        zapLevel = zapcore.InfoLevel
    case "warn":
        zapLevel = zapcore.WarnLevel
    case "error":
        zapLevel = zapcore.ErrorLevel
    default:
        zapLevel = zapcore.InfoLevel
    }

    config := zap.Config{
        Level: zap.NewAtomicLevelAt(zapLevel),
        Encoding: "json",
        EncoderConfig: zapcore.EncoderConfig{
            TimeKey:        "timestamp",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "caller",
            MessageKey:     "message",
            StacktraceKey: "stacktrace",
        },
        OutputPaths:      []string{"stdout", "./logs/app.log"},
        ErrorOutputPaths: []string{"stderr", "./logs/error.log"},
    }

    logger, _ = config.Build()
}
```

### 8.2 日志格式

```json
{
    "timestamp": "2024-01-01T00:00:00.000Z",
    "level": "info",
    "caller": "handler/media.go:45",
    "message": "media list retrieved",
    "request_id": "uuid",
    "user_id": 1,
    "duration": 125
}
```

## 9. 配置管理

### 9.1 配置文件

```yaml
# config/config.yaml
app:
  name: media-server
  host: 0.0.0.0
  port: 8080
  mode: release
  context_timeout: 60

database:
  path: ./data/media.db
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 3600

scanner:
  max_depth: 10
  exclude_patterns:
    - ".*"
    - "__MACOSX"
    - "Thumbs.db"
  file_size_limit: 10737418240

scraper:
  tmdb_api_key: ${TMDB_API_KEY}
  tmdb_base_url: https://api.themoviedb.org/3
  language: zh-CN
  timeout: 30
  retry_count: 3

stream:
  hls_dir: /tmp/hls
  default_quality: 1080p
  transcode_enabled: true
  hardware_acceleration: false
  max_concurrent: 3

log:
  level: info
  path: ./logs
  max_size: 100
  max_age: 30
  max_backups: 7

cors:
  allow_origins:
    - "*"
  allow_methods:
    - GET
    - POST
    PUT
    DELETE
    OPTIONS
  allow_headers:
    - "*"

rate_limit:
  rps: 100
  bursts: 200
```

### 9.2 配置加载

```go
type Config struct {
    App       AppConfig       `mapstructure:"app"`
    Database  DatabaseConfig  `mapstructure:"database"`
    Scanner   ScannerConfig   `mapstructure:"scanner"`
    Scraper   ScraperConfig   `mapstructure:"scraper"`
    Stream    StreamConfig    `mapstructure:"stream"`
    Log       LogConfig       `mapstructure:"log"`
    CORS      CORSConfig      `mapstructure:"cors"`
    RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")
    
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

## 10. 性能优化

### 10.1 数据库优化

- 连接池配置
- 索引优化
- 批量操作
- 查询缓存

```go
// 索引优化
func (r *MediaRepository) CreateIndexes(db *gorm.DB) error {
    return db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_media_type ON media(type);
        CREATE INDEX IF NOT EXISTS idx_media_year ON media(year);
        CREATE INDEX IF NOT EXISTS idx_media_rating ON media(rating);
        CREATE INDEX IF NOT EXISTS idx_media_status ON media(status);
        CREATE INDEX IF NOT EXISTS idx_media_tmdb_id ON media(tmdb_id);
        CREATE INDEX IF NOT EXISTS idx_episodes_series ON episodes(series_id, season_number);
    `).Error
}
```

### 10.2 并发处理

- 扫描任务并发
- 刮削任务并发
- 转码任务池

```go
type WorkerPool struct {
    workers int
    tasks   chan func() error
    results chan error
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers: workers,
        tasks:   make(chan func() error, workers*2),
        results: make(chan error, workers*2),
    }
}

func (wp *WorkerPool) Start() {
    var wg sync.WaitGroup
    for i := 0; i < wp.workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for task := range wp.tasks {
                wp.results <- task()
            }
        }()
    }
    wg.Wait()
    close(wp.results)
}
```

### 10.3 缓存策略

- 媒体信息缓存
- TMDB响应缓存
- 缩略图缓存

```go
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

type InMemoryCache struct {
    data map[string]cacheItem
    mu   sync.RWMutex
}

type cacheItem struct {
    value      interface{}
    expireAt    time.Time
}
```

## 11. 部署设计

### 11.1 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| API_PORT | API服务端口 | 8080 |
| DATABASE_PATH | 数据库路径 | ./data/media.db |
| MEDIA_FOLDER | 媒体文件夹 | /media |
| TMDB_API_KEY | TMDB API密钥 | - |
| LOG_LEVEL | 日志级别 | info |

### 11.2 启动流程

```go
func main() {
    cfg := config.LoadConfig("config/config.yaml")
    logger := zap.NewProduction()
    defer logger.Sync()
    
    db, err := database.NewSQLite(cfg.Database)
    if err != nil {
        logger.Fatal("failed to connect database", zap.Error(err))
    }
    
    if err := database.Migrate(db); err != nil {
        logger.Fatal("failed to migrate database", zap.Error(err))
    }
    
    repo := repository.NewRepository(db)
    svc := service.NewService(repo, cfg)
    handler := handler.NewHandler(svc, logger)
    
    router := router.SetupRouter(handler, cfg, logger)
    
    srv := &http.Server{
        Addr:    fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port),
        Handler: router,
    }
    
    if err := srv.ListenAndServe(); err != nil {
        logger.Fatal("server failed", zap.Error(err))
    }
}
```

## 12. 测试设计

### 12.1 单元测试

```go
func TestParseMediaInfo(t *testing.T) {
    tests := []struct {
        name     string
        filename string
        expected struct {
            title    string
            year     string
            quality  string
        }
    }{
        {
            name:     "Standard movie filename",
            filename: "The.Shawshank.Redemption.1994.1080p.BluRay.x264.mkv",
            expected: struct {
                title    string
                year     string
                quality  string
            }{
                title:   "The Shawshank Redemption",
                year:    "1994",
                quality: "1080P",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            title, year, quality, _ := ParseMediaInfo(tt.filename)
            if title != tt.expected.title {
                t.Errorf("expected title %s, got %s", tt.expected.title, title)
            }
        })
    }
}
```

### 12.2 集成测试

```go
func TestMediaAPI(t *testing.T) {
    router := SetupTestRouter()
    
    t.Run("Get media list", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "/api/media", nil)
        resp := httptest.NewRecorder()
        router.ServeHTTP(resp, req)
        
        assert.Equal(t, 200, resp.Code)
    })
}
```
