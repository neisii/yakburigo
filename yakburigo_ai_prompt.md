# 약버리고(藥버리GO) - 폐의약품 수거함 위치 안내 서비스

## 프로젝트 개요
폐의약품 수거 가능한 우체통, 약국, 보건소 등의 위치를 지도에 표시하여 사용자가 쉽게 찾을 수 있도록 돕는 웹 애플리케이션을 개발합니다.

## 핵심 요구사항

### 1. 기술 스택
- **백엔드**: Go (Gin 프레임워크)
- **데이터베이스**: PostgreSQL + PostGIS (공간 데이터 처리)
- **ORM**: GORM
- **지도 API**: 카카오맵 또는 네이버 지도 API
- **배포**: Docker 컨테이너

### 2. 주요 기능

#### 2.1 지도 기반 위치 표시
- 사용자 현재 위치를 중심으로 지도 표시
- 다음 유형별로 다른 마커 아이콘 사용:
  - 폐의약품 수거 가능 우체통 (빨간색 우체통 아이콘)
  - 약국 (초록색 십자가 아이콘)
  - 보건소/주민센터 (파란색 건물 아이콘)

#### 2.2 위치 기반 검색
- 현재 GPS 위치 기준 가까운 순으로 정렬
- 주소 입력으로 검색 가능
- 반경 500m, 1km, 3km, 5km 필터링 옵션

#### 2.3 상세 정보 제공
각 마커 클릭 시 표시할 정보:
- 장소명
- 주소
- 현재 위치로부터의 거리
- 운영 시간 (해당되는 경우)
- 길찾기 버튼 (카카오맵/네이버 지도 앱 연동)

#### 2.4 폐의약품 배출 방법 안내
- 하단 또는 별도 페이지에 배출 가능한 의약품 종류 안내
- 배출 방법 및 주의사항 텍스트 제공
- 이미지 또는 아이콘으로 시각화

### 3. 데이터베이스 스키마

```sql
-- 위치 정보 테이블
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- 'postbox', 'pharmacy', 'health_center'
    address VARCHAR(500) NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    operating_hours VARCHAR(200),
    contact VARCHAR(50),
    location GEOGRAPHY(POINT, 4326),  -- PostGIS 공간 인덱스용
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 공간 인덱스 생성
CREATE INDEX idx_location_geography ON locations USING GIST(location);
```

### 4. API 엔드포인트 설계

#### 4.1 GET /api/locations/nearby
현재 위치 기준 가까운 수거함 조회

**요청 파라미터:**
```json
{
  "latitude": 37.5665,
  "longitude": 126.9780,
  "radius": 1000,  // 미터 단위, optional (기본값: 3000)
  "type": "all"    // 'all', 'postbox', 'pharmacy', 'health_center'
}
```

**응답:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "서울중앙우체국 우체통",
      "type": "postbox",
      "address": "서울특별시 종로구 세종대로 54",
      "latitude": 37.5665,
      "longitude": 126.9780,
      "distance": 250.5,  // 미터 단위
      "operating_hours": "24시간",
      "contact": null
    }
  ],
  "count": 15
}
```

#### 4.2 GET /api/locations/search
주소 검색

**요청 파라미터:**
```json
{
  "query": "서울시 강남구",
  "limit": 20
}
```

#### 4.3 GET /api/locations/:id
특정 위치 상세 정보

**응답:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "서울중앙우체국 우체통",
    "type": "postbox",
    "address": "서울특별시 종로구 세종대로 54",
    "latitude": 37.5665,
    "longitude": 126.9780,
    "operating_hours": "24시간",
    "contact": null,
    "disposal_info": "일반 의약품, 연고, 물약 등 배출 가능"
  }
}
```

#### 4.4 GET /api/disposal-guide
폐의약품 배출 방법 안내 조회

**응답:**
```json
{
  "success": true,
  "data": {
    "acceptable_items": [
      "복용하지 않은 알약",
      "사용하지 않은 연고",
      "물약 (용기째)"
    ],
    "unacceptable_items": [
      "주사기",
      "의료기기"
    ],
    "instructions": "의약품은 포장을 벗기지 말고 그대로 배출해주세요."
  }
}
```

### 5. 프론트엔드 요구사항

#### 5.1 메인 화면 구성
- 상단: 검색바 (주소 입력)
- 중앙: 전체 화면 지도
- 하단: 필터 버튼 (우체통/약국/보건소), 반경 선택 버튼
- 우측 하단: 현재 위치로 이동 버튼

#### 5.2 반응형 디자인
- 모바일 우선 설계
- 태블릿, 데스크탑 대응

#### 5.3 사용자 경험
- 지도 로딩 시 스켈레톤 UI 표시
- 마커 클러스터링 적용 (마커 100개 이상 시)
- 길찾기 클릭 시 카카오맵/네이버 지도 앱으로 딥링크

### 6. 성능 요구사항
- API 응답 시간: 500ms 이내
- 지도 초기 로딩: 2초 이내
- 동시 접속자 100명 처리 가능

### 7. 배포 및 인프라
- Docker Compose로 로컬 개발 환경 구성
- PostgreSQL + PostGIS 컨테이너
- Go 애플리케이션 컨테이너
- Nginx 리버스 프록시 (선택사항)

### 8. 초기 데이터 준비
- 서울시 폐의약품 수거함 샘플 데이터 (최소 50개 위치)
- CSV 또는 JSON 형태로 준비
- 마이그레이션 스크립트로 초기 데이터 삽입

### 9. 개발 우선순위

**Phase 1 (MVP):**
1. 기본 지도 표시 및 마커 표시
2. 현재 위치 기반 가까운 수거함 조회 API
3. 마커 클릭 시 상세 정보 팝업

**Phase 2:**
1. 주소 검색 기능
2. 필터링 (유형별, 반경별)
3. 길찾기 연동

**Phase 3 (선택):**
1. 폐의약품 배출 가이드 페이지
2. 즐겨찾기 기능
3. 관리자 페이지 (위치 정보 CRUD)

## 개발 시 주의사항

1. **환경변수 관리**: 지도 API 키는 반드시 환경변수로 관리
2. **에러 처리**: 지도 API 호출 실패 시 적절한 fallback 처리
3. **CORS 설정**: 프론트엔드에서 API 호출 가능하도록 설정
4. **로깅**: 요청/응답 로깅으로 디버깅 용이하게
5. **보안**: SQL Injection 방지 (GORM 사용 시 자동 처리됨)

## 참고 자료
- 카카오맵 API 문서: https://apis.map.kakao.com/
- 네이버 지도 API 문서: https://www.ncloud.com/product/applicationService/maps
- PostGIS 공식 문서: https://postgis.net/documentation/
- Gin 프레임워크: https://gin-gonic.com/docs/

## 테스트 시나리오

### 시나리오 1: 현재 위치로 수거함 찾기
1. 사용자가 앱 접속
2. GPS 권한 허용
3. 자동으로 현재 위치 주변 지도 표시
4. 가까운 수거함 마커 확인

### 시나리오 2: 주소로 수거함 찾기
1. 상단 검색바에 "강남역" 입력
2. 해당 지역 지도로 이동
3. 주변 수거함 마커 표시

### 시나리오 3: 특정 수거함으로 길찾기
1. 마커 클릭
2. 상세 정보 팝업 확인
3. "길찾기" 버튼 클릭
4. 카카오맵/네이버 지도 앱 실행