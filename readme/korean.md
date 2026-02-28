# TesselBox - 한국어 README
## 육각형 복셀 게임

[![Open Source Helpers](https://www.codetriage.com/tesselstudio/tesselbox-game/badges/users.svg)](https://www.codetriage.com/tesselstudio/tesselbox-game)

*Terraria*에서 영감을 받은 2D 샌드박스 어드벤처 게임이지만 **육각형 그리드**에서 구축되었습니다.

세계를 탐험하고, 자원을 채굴하고, 구조물을 건설하고, 아이템을 만들고, 적과 싸우고 생존하세요 — 모두 아름다운 육각형 타일에서.

## 게임 기능

### ✅ **완전한 기능**
- **육각형 세계 생성** - 생물군계가 있는 절차적 생성 세계
- **채굴 및 제작** - 다양한 재료 속도를 가진 도구 기반 채굴
- **블록 배치** - 고스트 미리보기로 블록을 배치하기 위한 오른쪽 클릭
- **인벤토리 시스템** - 핫바가 있는 32슬롯 인벤토리 (9슬롯)
- **전투 시스템** - 공격 애니메이션이 있는 건강/데미지 시스템
- **낮/밤 주기** - 동적 조명과 시간 진행
- **날씨 효과** - 비, 눈 및 폭풍 시스템
- **저장/불러오기 시스템** - 자동 저장이 있는 지속적인 세계 상태

### 🎮 **컨트롤**
- **WASD / 화살표**: 이동
- **스페이스**: 점프 / 공격
- **왼쪽 클릭**: 블록 채굴
- **오른쪽 클릭**: 블록 배치
- **E**: 제작 메뉴 열기
- **Q**: 선택된 아이템 놓기
- **마우스 휠**: 핫바 선택
- **1-9**: 직접 핫바 선택
- **F5**: 수동 저장
- **F9**: 수동 불러오기
- **ESC**: 메뉴 / 메뉴 닫기

## 설치 및 설정

### 사전 요구사항
- **Go 1.19+** - 코어 엔진
- **Git** - 버전 관리

### 빠른 시작
```bash
# 저장소 복제
git clone https://github.com/tesselstudio/TesselBox-game.git
cd TesselBox-game

# 게임 빌드
go build ./cmd/client

# 게임 실행
./client
```

### 개발 설정
```bash
# 종속성 설치
go mod tidy

# 테스트 실행
go test ./...

# 개발용 빌드
go build -tags debug ./cmd/client
```

## 시스템 요구사항

### 최소
- **OS**: Windows 10+, macOS 10.15+, Linux
- **CPU**: 듀얼 코어 프로세서
- **RAM**: 4GB
- **GPU**: OpenGL 3.3+ 호환
- **저장공간**: 500MB 여유 공간

### 권장
- **CPU**: 쿼드 코어 프로세서
- **RAM**: 8GB+
- **GPU**: 전용 그래픽 카드
- **저장공간**: 1GB+ 여유 공간

## 아키텍처

### 코어 기술
- **언어**: Go (Golang)
- **그래픽스**: Ebiten (2D 게임 라이브러리)
- **빌드 시스템**: Go 모듈

### 프로젝트 구조
```
TesselBox/
├── cmd/client/          # 메인 게임 실행 파일
├── pkg/                 # 코어 패키지
│   ├── world/          # 세계 생성 및 관리
│   ├── player/         # 플레이어 메커닉 및 물리
│   ├── blocks/         # 블록 유형 및 속성
│   ├── items/          # 아이템 시스템 및 제작
│   ├── crafting/       # 제작 레시피 및 UI
│   ├── weather/        # 날씨 시뮬레이션
│   ├── gametime/       # 낮/밤 주기
│   ├── save/           # 저장/불러오기 기능
│   └── render/         # 렌더링 및 UI 시스템
├── config/             # 설정 파일
└── assets/             # 게임 에셋 (있는 경우)
```

## 기여

### 개발자를 위해
1. 저장소 포크
2. 기능 브랜치 생성 (`git checkout -b feature/amazing-feature`)
3. 변경사항 커밋 (`git commit -m 'Add amazing feature'`)
4. 브랜치로 푸시 (`git push origin feature/amazing-feature`)
5. Pull Request 열기

### 개발 가이드라인
- Go 코딩 표준 따르기
- 새 기능에 대한 테스트 추가
- 문서 업데이트
- 크로스 플랫폼 호환성 보장

## 라이선스

**CC BY-NC-SA 4.0 라이선스** - 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 크레딧

- **영감**: Terraria 게임 메커닉
- **구축**: Ebiten 게임 엔진 사용
- **기여자**: 오픈 소스 커뮤니티

## 지원

- **이슈**: [GitHub Issues](https://github.com/tesselstudio/TesselBox-game/issues)
- **토론**: [GitHub Discussions](https://github.com/tesselstudio/TesselBox-game/discussions)
- **위키**: [프로젝트 위키](https://github.com/tesselstudio/TesselBox-game/wiki)

---

*TesselBox의 육각형 세계 탐험을 즐기세요!*
