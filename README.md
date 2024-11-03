# 주식 거래 계산기 (Stock Trade Calculator)

이 프로젝트는 주식 거래 시 필요한 다양한 계산 기능을 제공하는 데스크톱 애플리케이션입니다. [Cursor AI](https://cursor.sh/)의 도움을 받아 개발되었습니다.

## 주요 기능

- 매수/매도 계산
- 물타기(평단가) 계산
- 복리 수익률 계산
- 거래 수수료 설정
- 다크/라이트 모드 지원

## 설치 방법

### 바이너리 실행
1. [Releases](https://github.com/yourusername/trade-calculator/releases) 페이지에서 최신 버전 다운로드
2. 압축 해제 후 실행 파일 실행

### 소스코드 실행
1. 저장소 클론
   git clone https://github.com/yourusername/trade-calculator.git
2. 의존성 설치
   go mod download
3. 실행
   go run .

### 패키지로 설치
go install github.com/yourusername/trade-calculator@latest

## 사용 방법

### 매수 탭
- 매수가격과 수량을 입력하여 매수 정보 등록
- 자동으로 수수료와 목표가 계산
- 매수 내역 저장 및 관리

### 매도 탭
- 저장된 매수 내역에서 선택하여 매도 계산
- 예상 수익률과 수수료 자동 계산

### 물타기 탭
- 기존 매수 내역에 추가 매수 시 평균단가 계산
- 수정된 매수 정보 자동 업데이트

### 복리 계산 탭
- 목표 수익률 달성을 위한 일일 수익률 계산
- 투자 기간별 복리 수익 계산

### 설정 탭
- 거래 수수료율 설정
- 손절 기준 설정
- 목표 수익률 설정
- 테마 설정

## 기여하기

이슈와 풀 리퀘스트를 환영합니다. 주요 변경사항의 경우 먼저 이슈를 열어 논의해주세요.

## 라이센스

[MIT 라이센스](LICENSE)로 배포됩니다. Cursor AI를 활용한 개발이 허용됩니다.
