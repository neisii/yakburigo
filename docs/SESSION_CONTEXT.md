# Yakburigo (약버리GO) - Session Context

**Purpose**: Provide complete context for new Claude sessions to continue the Yakburigo project.

**Last Updated**: 2025-10-22

---

## Project Overview

**Goal**: 폐의약품 수거 가능한 우체통, 약국, 보건소 등의 위치를 지도에 표시하여 사용자가 쉽게 찾을 수 있도록 돕는 웹 애플리케이션 개발

**Tech Stack**:
- Backend: Go (Gin Framework)
- Database: PostgreSQL + PostGIS
- ORM: GORM
- Map API: Kakao Map or Naver Map API
- Deployment: Docker

**Project Status**: 0% Complete / Planning Phase

**Methodology**: AI-DLC (see `docs/ai-dlc.txt`)

---

## Communication Guidelines

1. **Tone**: Informal (반말) between user and Claude
2. **Documentation**: 
   - SESSION_CONTEXT.md: English (for Claude sessions)
   - Other docs: Korean (for user reading)
   - Style: Fact-based (현상-원인-결과), avoid emojis, use numbering/bullets
3. **Documentation Structure**:
   - **PROGRESS.md**: 전체 Phase 간결한 요약 (Phase별 5-10줄)
   - **PHASE_N_PLAN.md**: Phase 시작 전 상세 계획 (요구사항, 기술적 접근, 예상 이슈)
   - **PHASE_N_SUMMARY.md**: Phase 완료 후 상세 회고 (구현 내용, 트러블슈팅, 코드 예시)
   - **SESSION_CONTEXT.md**: Claude 세션용 스냅샷 (현재 상태, 다음 작업)
4. **Phase Documentation Templates**:
   - Planning: Use `docs/templates/PHASE_PLAN_TEMPLATE.md`
   - Summary: Use `docs/templates/PHASE_SUMMARY_TEMPLATE.md`
5. **Progress Tracking**: Always update this file after completing tasks
6. **Approval**: Ask user approval when requirements are ambiguous
7. **Commit Convention**: Follow Conventional Commits with Scope (see `docs/COMMIT_CONVENTION.md`)

---

## Architecture

### System Design
```
[Frontend Layer]
    ↓
[Backend API Layer (Go + Gin)]
    ↓
[Service Layer] ← Business logic
    ↓
[Repository Layer] ← GORM + PostGIS
    ↓
[PostgreSQL + PostGIS] ← Spatial data
```

**Key Design Principles**:
- Clean Architecture: Domain, Repository, Service, Handler 계층 분리
- RESTful API 설계
- PostGIS 활용한 공간 쿼리 최적화
- 환경변수 기반 설정 관리

### Technology Decisions
- **Go + Gin**: 경량화되고 빠른 백엔드 API 서버 구축
- **PostgreSQL + PostGIS**: 공간 데이터 처리 및 거리 기반 검색
- **GORM**: 타입 안정성과 생산성 향상
- **Docker**: 개발 환경 일관성 및 배포 용이성

---

## Implemented Features

### Completed Phases

None yet - Project is in Planning Phase

---

### In Progress

#### Phase 1: 프로젝트 구조 및 기본 설정
**Status**: ⚠️ Planning

**Planned Features**:
- [ ] Go 프로젝트 초기화 (go.mod)
- [ ] 디렉토리 구조 생성 (Clean Architecture)
- [ ] Docker Compose 설정 (PostgreSQL + PostGIS)
- [ ] 기본 환경변수 설정
- [ ] Gin 서버 Hello World

**Reference**: `docs/PHASE_1_PLAN.md`

---

### Not Implemented

1. **Location API**: 위치 기반 수거함 조회 API
2. **Search API**: 주소 검색 API
3. **Detail API**: 특정 위치 상세 정보
4. **Disposal Guide API**: 폐의약품 배출 가이드
5. **Database Migration**: 초기 데이터 삽입
6. **Frontend**: 지도 기반 UI

**Reference**: `yakburigo_ai_prompt.md`

---

## File Structure

```
yakburigo/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Domain entities
│   │   └── location.go
│   ├── repository/                 # Data access layer
│   │   └── location_repository.go
│   ├── service/                    # Business logic
│   │   └── location_service.go
│   ├── handler/                    # HTTP handlers
│   │   ├── location_handler.go
│   │   └── disposal_guide_handler.go
│   ├── dto/                        # Request/Response DTOs
│   │   └── location_dto.go
│   └── middleware/                 # Middlewares
│       ├── cors.go
│       └── logger.go
├── pkg/
│   ├── database/                   # DB connection
│   │   └── postgres.go
│   └── response/                   # Common response format
│       └── response.go
├── config/
│   └── config.go                   # Configuration management
├── migrations/                     # DB migrations
│   ├── 001_create_locations_table.sql
│   └── 002_seed_sample_data.sql
├── scripts/
│   └── seed_data.json              # Sample data
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── go.mod
├── go.sum
├── docs/
│   ├── SESSION_CONTEXT.md          # This file
│   ├── PHASE_1_PLAN.md
│   ├── PROGRESS.md
│   └── templates/
└── README.md
```

**Important Files** (To Be Created):
- `cmd/server/main.go`: Gin server initialization
- `internal/domain/location.go`: Location entity with GORM tags
- `pkg/database/postgres.go`: PostgreSQL connection with PostGIS support
- `config/config.go`: Environment variable management

---

## Data Models

### Location
```go
type Location struct {
    ID              uint      `gorm:"primaryKey"`
    Name            string    `gorm:"size:200;not null"`
    Type            string    `gorm:"size:50;not null"` // 'postbox', 'pharmacy', 'health_center'
    Address         string    `gorm:"size:500;not null"`
    Latitude        float64   `gorm:"type:decimal(10,8);not null"`
    Longitude       float64   `gorm:"type:decimal(11,8);not null"`
    OperatingHours  string    `gorm:"size:200"`
    Contact         string    `gorm:"size:50"`
    Location        string    `gorm:"type:geography(POINT,4326)"` // PostGIS
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

**Usage**: Core entity for storing disposal locations

---

## API Endpoints

### Location Service

**Base URL**: `/api/v1`

**Endpoints**:
```bash
# Get nearby locations
GET /locations/nearby?latitude=37.5665&longitude=126.9780&radius=1000&type=all

# Search by address
GET /locations/search?query=서울시+강남구&limit=20

# Get location detail
GET /locations/:id

# Get disposal guide
GET /disposal-guide
```

**Authentication**: None (public API)

**Rate Limits**: TBD

---

## Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=yakburigo

# Server
SERVER_PORT=8080
GIN_MODE=debug

# Map API (선택)
KAKAO_MAP_API_KEY=your_key
NAVER_MAP_CLIENT_ID=your_client_id
NAVER_MAP_CLIENT_SECRET=your_secret
```

**Setup**:
```bash
cp .env.example .env
# Edit .env with your values
```

---

## User Decisions (All Final)

### 1. 개발 우선순위
- **Decision**: Phase 1은 프로젝트 구조 및 기본 설정에 집중
- **Rationale**: Clean Architecture 기반 구조를 먼저 확립해야 이후 API 개발이 용이함
- **Date**: 2025-10-22

---

### 2. Map API 선택 (Pending)
- **Decision**: TBD by user
- **Options**: Kakao Map vs Naver Map vs Backend-only first
- **Date**: 2025-10-22

---

### 3. Sample Data Source (Pending)
- **Decision**: TBD by user
- **Options**: Auto-generated vs Public Data API
- **Date**: 2025-10-22

---

### 4. PostGIS Implementation Level (Pending)
- **Decision**: TBD by user
- **Options**: Full PostGIS GEOGRAPHY type vs Simple DECIMAL with Haversine formula
- **Date**: 2025-10-22

---

## Development Workflow

### Getting Started (Phase 1+)
```bash
cd yakburigo

# Initialize Go project
go mod init github.com/neisii/yakburigo

# Start PostgreSQL + PostGIS
docker-compose up -d

# Run server
go run cmd/server/main.go
```

### Running Tests (Phase 2+)
```bash
go test ./...
```

### Building
```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

---

## Testing Strategy

### Unit Tests (Phase 2+)
- **Framework**: Go testing package
- **Coverage Target**: 80%

**Key Test Files**:
- TBD

---

### E2E Tests (Phase 3+)
- **Framework**: TBD
- **Test Scenarios**: API endpoint tests

---

## Known Issues

None yet - Project is in Planning Phase

---

## Dependencies

### Core Dependencies (Phase 1)
```json
{
  "github.com/gin-gonic/gin": "latest",
  "gorm.io/gorm": "latest",
  "gorm.io/driver/postgres": "latest",
  "github.com/joho/godotenv": "latest"
}
```

### Dev Dependencies
```json
{
  "github.com/stretchr/testify": "latest"
}
```

**Update Policy**: Use latest stable versions

---

## Git Workflow

### Branch Strategy
- `main`: Production-ready code
- Feature branches as needed

### Commit Convention
**Format**: `<type>(<scope>): <subject>`

**Types**: feat, fix, docs, style, refactor, test, chore

**Scope**: `yakburigo` (main backend project)

**Reference**: `docs/COMMIT_CONVENTION.md`

**Examples**:
```bash
feat(yakburigo): initialize Go project structure
feat(yakburigo): add PostgreSQL connection with PostGIS
feat(yakburigo): implement nearby locations API
docs(yakburigo): update SESSION_CONTEXT with Phase 1 status
```

---

## Useful Commands

```bash
# Development
go run cmd/server/main.go        # Start dev server
go build -o bin/server cmd/server/main.go  # Build binary

# Testing
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -cover ./...             # With coverage

# Database
docker-compose up -d             # Start PostgreSQL
docker-compose down              # Stop PostgreSQL
docker-compose logs postgres     # View logs

# Dependencies
go mod tidy                      # Clean dependencies
go mod download                  # Download dependencies
```

---

## Quick Reference

### Common Tasks

#### Add New API Endpoint
1. Define domain entity in `internal/domain/`
2. Create repository in `internal/repository/`
3. Implement service in `internal/service/`
4. Add handler in `internal/handler/`
5. Register route in `cmd/server/main.go`
6. Write tests
7. Update documentation

#### Add Database Migration
1. Create SQL file in `migrations/`
2. Test with PostgreSQL
3. Document in PHASE summary

---

## Resources

### Documentation
- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [PostGIS Documentation](https://postgis.net/documentation/)

### External Links
- [Kakao Map API](https://apis.map.kakao.com/)
- [Naver Map API](https://www.ncloud.com/product/applicationService/maps)
- [Go Best Practices](https://go.dev/doc/effective_go)

---

## Notes for New Sessions

### What to Check First
1. [ ] Read this SESSION_CONTEXT.md
2. [ ] Check latest PHASE_N_PLAN.md
3. [ ] Review open issues in Known Issues section
4. [ ] Check PROGRESS.md for overall status
5. [ ] Verify environment: `docker-compose ps`

### Important Reminders
- Ask user approval before implementing each Phase
- Follow AI-DLC methodology: Plan → Clarify → Implement
- Update SESSION_CONTEXT.md after major changes
- Document all decisions and issues

---

**Last Session**: 2025-10-22  
**Next Steps**: Create PHASE_1_PLAN.md and get user approval  
**Status**: Planning Phase
