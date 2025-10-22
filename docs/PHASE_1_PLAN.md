# Phase 1 Implementation Plan

**작성일**: 2025-10-22  
**목표**: Go + Gin 백엔드 프로젝트 구조 및 기본 환경 설정 완료

---

## 📋 Phase 1 목표

### 주요 작업
1. **Go 프로젝트 초기화 및 디렉토리 구조 생성**
2. **Docker Compose로 PostgreSQL + PostGIS 환경 구성**
3. **기본 Gin 서버 구현 (Hello World)**
4. **환경변수 관리 설정**
5. **데이터베이스 연결 및 테스트**

### 예상 소요 시간
- Go 프로젝트 초기화: 0.5시간
- 디렉토리 구조 생성: 0.5시간
- Docker Compose 설정: 1시간
- Gin 서버 구현: 1시간
- 데이터베이스 연결: 1시간
- **총 예상 시간**: 4시간

---

## 🎯 작업 1: Go 프로젝트 초기화

### 1.1 go.mod 생성

**설명**: Go 모듈 초기화 및 프로젝트 루트 설정

**구현 방식**:
- 모듈 이름: `github.com/neisii/yakburigo`
- Go 버전: 1.21+ (최신 안정 버전)

**명령어**:
```bash
cd /Users/neisii/Development/yakburigo
go mod init github.com/neisii/yakburigo
```

### 1.2 핵심 의존성 추가

**필수 패키지**:
```bash
# Web framework
go get github.com/gin-gonic/gin

# Database
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Environment variables
go get github.com/joho/godotenv

# Testing (dev dependency)
go get github.com/stretchr/testify
```

### 1.3 구현 체크리스트

- [ ] `go mod init` 실행
- [ ] 핵심 의존성 설치
- [ ] `go.mod` 파일 생성 확인
- [ ] `go.sum` 파일 자동 생성 확인

---

## 🎯 작업 2: 디렉토리 구조 생성

### 2.1 Clean Architecture 기반 구조

**설명**: 계층 분리를 통한 유지보수성 확보

**디렉토리 구조**:
```
yakburigo/
├── cmd/
│   └── server/
│       └── main.go                 # 애플리케이션 진입점
├── internal/                       # 내부 패키지 (외부 import 불가)
│   ├── domain/                     # 도메인 엔티티
│   │   └── location.go
│   ├── repository/                 # 데이터 접근 계층
│   │   └── location_repository.go
│   ├── service/                    # 비즈니스 로직
│   │   └── location_service.go
│   ├── handler/                    # HTTP 핸들러
│   │   ├── location_handler.go
│   │   └── disposal_guide_handler.go
│   ├── dto/                        # 요청/응답 DTO
│   │   └── location_dto.go
│   └── middleware/                 # 미들웨어
│       ├── cors.go
│       └── logger.go
├── pkg/                            # 외부 공유 가능 패키지
│   ├── database/                   # DB 연결
│   │   └── postgres.go
│   └── response/                   # 공통 응답 포맷
│       └── response.go
├── config/
│   └── config.go                   # 환경변수 관리
├── migrations/                     # SQL 마이그레이션
├── scripts/                        # 유틸리티 스크립트
├── docs/                           # 문서
├── .env.example                    # 환경변수 예시
├── .gitignore                      # Git ignore 파일
├── docker-compose.yml              # Docker 설정
├── Dockerfile                      # Go 앱 컨테이너
├── go.mod
└── README.md
```

### 2.2 구현 체크리스트

**디렉토리 생성**:
- [ ] `cmd/server/` 생성
- [ ] `internal/domain/` 생성
- [ ] `internal/repository/` 생성
- [ ] `internal/service/` 생성
- [ ] `internal/handler/` 생성
- [ ] `internal/dto/` 생성
- [ ] `internal/middleware/` 생성
- [ ] `pkg/database/` 생성
- [ ] `pkg/response/` 생성
- [ ] `config/` 생성
- [ ] `migrations/` 생성
- [ ] `scripts/` 생성

**파일 생성**:
- [ ] `.gitignore` (Go 기본 + .env)
- [ ] `.env.example` (환경변수 템플릿)
- [ ] `README.md` (프로젝트 설명)

---

## 🎯 작업 3: Docker Compose 설정

### 3.1 PostgreSQL + PostGIS 컨테이너

**설명**: 로컬 개발 환경용 데이터베이스 구성

**docker-compose.yml**:
```yaml
version: '3.8'

services:
  postgres:
    image: postgis/postgis:15-3.3
    container_name: yakburigo-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: yakburigo
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

**특징**:
- PostGIS 3.3 포함된 PostgreSQL 15 이미지 사용
- 볼륨 마운트로 데이터 영속성 보장
- healthcheck로 컨테이너 상태 확인
- migrations 디렉토리를 초기화 스크립트로 마운트

### 3.2 구현 체크리스트

- [ ] `docker-compose.yml` 작성
- [ ] `docker-compose up -d` 실행 테스트
- [ ] `docker-compose ps` 로 컨테이너 상태 확인
- [ ] PostgreSQL 접속 테스트: `docker exec -it yakburigo-postgres psql -U postgres -d yakburigo`
- [ ] PostGIS 확장 설치 확인: `SELECT PostGIS_version();`

---

## 🎯 작업 4: 기본 Gin 서버 구현

### 4.1 main.go 작성

**파일**: `cmd/server/main.go`

**기능**:
- Gin 서버 초기화
- 기본 라우트 설정 (헬스 체크)
- 환경변수에서 포트 읽기
- Graceful shutdown 지원

**구현 예시**:
```go
package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // 환경변수 로드
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Gin 모드 설정
    ginMode := os.Getenv("GIN_MODE")
    if ginMode == "" {
        ginMode = gin.DebugMode
    }
    gin.SetMode(ginMode)

    // 라우터 초기화
    router := gin.Default()

    // 헬스 체크 엔드포인트
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "service": "yakburigo-api",
        })
    })

    // API v1 그룹
    v1 := router.Group("/api/v1")
    {
        v1.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message": "pong",
            })
        })
    }

    // 서버 시작
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting server on port %s...", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

### 4.2 .env.example 작성

**파일**: `.env.example`

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=yakburigo

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Map API Keys (Optional)
KAKAO_MAP_API_KEY=
NAVER_MAP_CLIENT_ID=
NAVER_MAP_CLIENT_SECRET=
```

### 4.3 구현 체크리스트

- [ ] `cmd/server/main.go` 작성
- [ ] `.env.example` 작성
- [ ] `.env` 파일 생성 (`.env.example` 복사)
- [ ] `go run cmd/server/main.go` 실행 테스트
- [ ] `curl http://localhost:8080/health` 응답 확인
- [ ] `curl http://localhost:8080/api/v1/ping` 응답 확인

---

## 🎯 작업 5: 데이터베이스 연결 구현

### 5.1 config.go 작성

**파일**: `config/config.go`

**기능**:
- 환경변수에서 설정 읽기
- 구조체로 설정 관리
- 검증 로직

**구현 예시**:
```go
package config

import (
    "fmt"
    "os"
)

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

type ServerConfig struct {
    Port string
}

func Load() (*Config, error) {
    cfg := &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            DBName:   getEnv("DB_NAME", "yakburigo"),
        },
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
    }

    return cfg, nil
}

func (c *DatabaseConfig) DSN() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        c.Host, c.Port, c.User, c.Password, c.DBName,
    )
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 5.2 postgres.go 작성

**파일**: `pkg/database/postgres.go`

**기능**:
- GORM + PostgreSQL 연결
- PostGIS 확장 설치 확인
- 연결 풀 설정

**구현 예시**:
```go
package database

import (
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }

    // PostGIS 확장 설치 확인
    var result string
    if err := db.Raw("SELECT PostGIS_version()").Scan(&result).Error; err != nil {
        log.Println("Warning: PostGIS not installed or not available")
    } else {
        log.Printf("PostGIS version: %s", result)
    }

    // 연결 풀 설정
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)

    log.Println("Database connection established")
    return db, nil
}
```

### 5.3 main.go에 DB 연결 추가

**수정 사항**:
```go
package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/neisii/yakburigo/config"
    "github.com/neisii/yakburigo/pkg/database"
)

func main() {
    // 환경변수 로드
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // 설정 로드
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 데이터베이스 연결
    db, err := database.NewPostgresDB(cfg.Database.DSN())
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    log.Println("Database connected successfully")

    // Gin 모드 설정
    gin.SetMode(gin.DebugMode)

    // 라우터 초기화
    router := gin.Default()

    // 헬스 체크 (DB 포함)
    router.GET("/health", func(c *gin.Context) {
        sqlDB, _ := db.DB()
        if err := sqlDB.Ping(); err != nil {
            c.JSON(503, gin.H{
                "status": "error",
                "message": "database unhealthy",
            })
            return
        }

        c.JSON(200, gin.H{
            "status": "ok",
            "service": "yakburigo-api",
            "database": "connected",
        })
    })

    // 서버 시작
    log.Printf("Starting server on port %s...", cfg.Server.Port)
    if err := router.Run(":" + cfg.Server.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

### 5.4 구현 체크리스트

- [ ] `config/config.go` 작성
- [ ] `pkg/database/postgres.go` 작성
- [ ] `cmd/server/main.go` 업데이트 (DB 연결 추가)
- [ ] Docker PostgreSQL 컨테이너 실행 확인
- [ ] 서버 시작 시 DB 연결 로그 확인
- [ ] `/health` 엔드포인트에서 DB 상태 확인

---

## 🎯 작업 6: .gitignore 및 README.md 작성

### 6.1 .gitignore

**파일**: `.gitignore`

```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment variables
.env
.env.local
.env.*.local

# IDEs
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log
logs/

# Database
*.db
*.sqlite
```

### 6.2 README.md

**파일**: `README.md`

```markdown
# 약버리고 (Yakburigo)

폐의약품 수거함 위치 안내 서비스 백엔드 API

## 기술 스택

- **언어**: Go 1.21+
- **웹 프레임워크**: Gin
- **데이터베이스**: PostgreSQL + PostGIS
- **ORM**: GORM
- **컨테이너**: Docker

## 시작하기

### 사전 요구사항

- Go 1.21+
- Docker & Docker Compose

### 설치 및 실행

1. 저장소 클론
```bash
git clone https://github.com/neisii/yakburigo.git
cd yakburigo
```

2. 환경변수 설정
```bash
cp .env.example .env
# .env 파일을 편집하여 필요한 값 입력
```

3. PostgreSQL 시작
```bash
docker-compose up -d
```

4. 의존성 설치
```bash
go mod download
```

5. 서버 실행
```bash
go run cmd/server/main.go
```

6. 헬스 체크
```bash
curl http://localhost:8080/health
```

## API 엔드포인트

### Phase 1 (완료)
- `GET /health` - 헬스 체크
- `GET /api/v1/ping` - API 응답 테스트

### Phase 2+ (예정)
- `GET /api/v1/locations/nearby` - 주변 수거함 조회
- `GET /api/v1/locations/search` - 주소 검색
- `GET /api/v1/locations/:id` - 상세 정보
- `GET /api/v1/disposal-guide` - 배출 가이드

## 프로젝트 구조

```
yakburigo/
├── cmd/server/          # 애플리케이션 진입점
├── internal/            # 내부 패키지
│   ├── domain/          # 도메인 엔티티
│   ├── repository/      # 데이터 접근
│   ├── service/         # 비즈니스 로직
│   ├── handler/         # HTTP 핸들러
│   ├── dto/             # DTO
│   └── middleware/      # 미들웨어
├── pkg/                 # 공유 패키지
│   ├── database/        # DB 연결
│   └── response/        # 응답 포맷
├── config/              # 설정
├── migrations/          # DB 마이그레이션
└── docs/                # 문서
```

## 개발

### 테스트 실행
```bash
go test ./...
```

### 빌드
```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## 문서

- [세션 컨텍스트](docs/SESSION_CONTEXT.md)
- [Phase 1 계획](docs/PHASE_1_PLAN.md)
- [커밋 컨벤션](docs/COMMIT_CONVENTION.md)

## 라이선스

MIT
```

### 6.3 구현 체크리스트

- [ ] `.gitignore` 작성
- [ ] `README.md` 작성
- [ ] Git 초기화: `git init`
- [ ] 첫 커밋: `git add . && git commit -m "chore(yakburigo): initialize project structure"`

---

## 📊 작업 우선순위

### High Priority (필수)
1. ✅ Go 프로젝트 초기화
2. ✅ 디렉토리 구조 생성
3. ✅ Docker Compose 설정
4. ✅ Gin 서버 구현
5. ✅ 데이터베이스 연결

### Medium Priority (권장)
6. 📄 README.md 작성
7. 📄 .gitignore 작성

### Low Priority (선택)
8. 🔧 Dockerfile 작성 (Phase 2로 이동 가능)

---

## ✅ 완료 조건

### Phase 1 완료 기준
- [ ] Go 프로젝트 초기화 완료 (`go.mod` 존재)
- [ ] 모든 디렉토리 구조 생성 완료
- [ ] Docker Compose로 PostgreSQL 실행 가능
- [ ] Gin 서버가 정상적으로 시작됨
- [ ] `/health` 엔드포인트가 200 OK 응답
- [ ] `/health` 엔드포인트가 DB 연결 상태 확인
- [ ] `.env` 파일로 환경변수 관리
- [ ] 문서화 완료 (`README.md`, `SESSION_CONTEXT.md` 업데이트)

### 검증 체크리스트
- [ ] `docker-compose up -d` 성공
- [ ] `docker-compose ps` 에서 postgres 컨테이너 healthy 상태
- [ ] `go run cmd/server/main.go` 에러 없이 실행
- [ ] `curl http://localhost:8080/health` 응답:
  ```json
  {
    "status": "ok",
    "service": "yakburigo-api",
    "database": "connected"
  }
  ```
- [ ] 로그에 "Database connected successfully" 출력
- [ ] 로그에 PostGIS 버전 정보 출력

---

## 🚀 시작 절차

### 1. 사용자 승인 대기

현재 이 문서를 읽고 승인해주세요:
- 프로젝트 구조가 적절한가?
- 기술 스택 선택이 적합한가?
- Phase 1 범위가 적정한가?

### 2. 미결정 사항 확인

다음 질문들에 답변이 필요합니다:

**Q1. 지도 API 선택**
- 옵션 A: 카카오맵 API 사용
- 옵션 B: 네이버 지도 API 사용
- 옵션 C: Phase 1에서는 백엔드만 구현 (지도 API는 Phase 2+에서)

**Q2. 초기 샘플 데이터**
- 옵션 A: 50개 위치 데이터를 임의로 생성 (서울 주요 지역)
- 옵션 B: 공공데이터 포털 API에서 실제 데이터 가져오기
- 옵션 C: Phase 1에서는 데이터 없이 구조만 (데이터는 Phase 2+에서)

**Q3. PostGIS 활용 수준**
- 옵션 A: PostGIS GEOGRAPHY 타입 + ST_Distance 함수 사용 (정확한 거리 계산)
- 옵션 B: DECIMAL 타입 + 하버사인 공식 (간단한 구현)

**Q4. Git 저장소**
- 로컬 저장소만 사용할까요?
- 아니면 GitHub 등 원격 저장소를 생성할까요?

### 3. Phase 1 시작

승인 및 답변 후 즉시 구현 시작:
```bash
cd /Users/neisii/Development/yakburigo
# 작업 시작...
```

---

## 📝 참고 사항

### 관련 문서
- `docs/ai-dlc.txt` - AI-DLC 방법론
- `yakburigo_ai_prompt.md` - 전체 기획서
- `docs/COMMIT_CONVENTION.md` - 커밋 컨벤션
- `docs/SESSION_CONTEXT.md` - 세션 컨텍스트

### 이전 Phase 이슈
- 없음 (첫 번째 Phase)

### 기술 스택 버전
- Go: 1.21+
- PostgreSQL: 15
- PostGIS: 3.3
- Gin: latest
- GORM: latest

### 개발 환경
- OS: macOS (Darwin 24.6.0)
- 작업 디렉토리: `/Users/neisii/Development/yakburigo`

---

**작성일**: 2025-10-22  
**상태**: 승인 대기  
**다음 단계**: 사용자 승인 후 구현 시작
