# 약버리고 (Yakburigo)

폐의약품 수거함 위치 안내 서비스 백엔드 API

## 기술 스택

- Go 1.25+
- Gin 웹 프레임워크
- PostgreSQL + PostGIS
- GORM
- Docker

## 시작하기

```bash
# PostgreSQL 시작
docker-compose up -d

# 환경변수 설정
cp .env.example .env

# 서버 실행
go run cmd/server/main.go

# 헬스 체크
curl http://localhost:8080/health
```

## API 엔드포인트

- GET /health - 헬스 체크
- GET /api/v1/ping - Ping 테스트

## 문서

- [세션 컨텍스트](docs/SESSION_CONTEXT.md)
- [Phase 1 계획](docs/PHASE_1_PLAN.md)
