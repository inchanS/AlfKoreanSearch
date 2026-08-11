# AlfKoreanSearch : Korean Search Workflow for Alfred
![Test](../../actions/workflows/test-go.yml/badge.svg) ![Release](../../actions/workflows/release.yml/badge.svg)  
![GitHub stars](https://img.shields.io/github/stars/inchans/alfkoreansearch?style=flat&logo=apachespark)
![GitHub all releases](https://img.shields.io/github/downloads/inchanS/alfkoreansearch/total?logo=github) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/inchanS/alfkoreansearch?logo=rocket)  ![GitHub](https://img.shields.io/github/license/inchanS/alfkoreansearch)

Korean Search Workflow for Alfred
---------------------------------

Alfred에서 국립국어원 [우리말샘 웹사이트](https://opendict.korean.go.kr) 검색어가 자동완성 되는 워크플로우
<br>  
<br>  

### **Acknowledgments**
[@Kuniz](https://github.com/Kuniz)님의 [alfnaversearch 워크플로우](https://github.com/Kuniz/alfnaversearch)를 기반으로 우리말샘 검색 용도에 맞게 구현한 워크플로우입니다.  
이 프로젝트는 [AlfNaverSearchPlus](https://github.com/inchanS/AlfNaverSearchPlus)가 Go로 전환하기 이전의 Python 코드를 포크하여 이어받았습니다.

변경부분
- 우리말샘(opendict.korean.go.kr) 자동완성 검색 기능
- **자동 업데이트 기능 추가**
- **Go 단일 바이너리로 전면 재작성**
  - **이전 버전은 Python이 필요했지만, 재작성 이후로는 Python이 필요 없습니다.**
  - Python 및 alfred-pyworkflow 의존을 완전히 제거하여 별도 런타임 설치가 불필요
  - Apple Silicon / Intel 유니버설 바이너리, 검색 기동 속도 향상
<br>  

Preview
--------

<img src="images/alfKoreanSearch.gif" width="600">

<br>  

Install workflow
--------------

- [releases](../../releases/latest) 페이지의 `AlfKoreanSearch.alfredworkflow`를 다운로드 받아서 실행한다.

- **별도의 런타임 설치가 필요 없습니다.** (Python 불필요)
  - 워크플로우에 포함된 유니버설 바이너리(Apple Silicon / Intel)로 동작합니다.
  - 다운로드 격리(Gatekeeper)는 워크플로우 내부의 `run` 스크립트가 최초 실행 시 자동으로 해제하므로, 별도 조치 없이 바로 사용할 수 있습니다.
  - 참고: 이전 Python 버전은 python3 설치가 필요했습니다. (`brew install python`, `xcode-select --install`)

- Alfred 4.0 이상 지원

Auto Update
--------------

자동 업데이트를 지원합니다.

- 워크플로우 사용 시 **주 1회** 백그라운드에서 새 릴리스를 확인하며, 검색 속도에는 영향이 없습니다.
- 새 버전이 있으면 검색 결과 맨 위에 `New version of AlfKoreanSearch is available!` 항목이 표시되고,
  선택하면 새 버전을 내려받아 설치합니다.
- 수동 명령어: 검색 keyword 뒤에 `workflow:update`를 입력하면 즉시 최신 버전을 내려받아 설치합니다. (예: `k workflow:update`)

General Usage
--------------
* `k ...`  : 검색어 입력 (연관검색어가 자동완성되어 나열됨)
* **Enter** : 선택한 검색어를 우리말샘에서 검색 (웹브라우저 출력)

### 단축키 관련 기능
* **Cmd + C** : 선택한 검색어를 클립보드에 복사
* **Cmd + L** : 선택한 검색어를 라지 타입(large type)으로 표시
* **Cmd + Y** : 검색결과 미리 보기(Quick Look)

트리거가 되는 키워드 `k`는 Alfred의 워크플로우 설정에서 개인에 맞게 직접 수정할 수 있습니다.

Build from source
--------------
워크플로우 로직은 순수 Go로 작성되어 있습니다. (빌드 시 `golang.org/x/text`만 사용하며 바이너리에 정적 링크됩니다.)

```sh
# 유니버설 바이너리 빌드 + ad-hoc 서명 (workflow/koreansearch 생성)
sh ./build.sh

# .alfredworkflow 패키지 생성
sh ./make.sh

# 테스트
go test ./...
```

- `cmd/koreansearch` : 진입점 (서브커맨드 디스패치, NFC 정규화)
- `internal/` : `alfred`(피드백 JSON), `httpx`(HTTP), `cache`(캐시), `update`(자동 업데이트), `urlx`(URL 인코딩), `handlers`(우리말샘 검색)
- `workflow/run` : 다운로드 격리 해제 후 바이너리를 실행하는 Script Filter 진입 스크립트

이 워크플로우는 [AlfNaverSearchPlus](https://github.com/inchanS/AlfNaverSearchPlus)의 Go 아키텍처를 계승했습니다. 이전 Python 버전은 [alfred-pyworkflow](https://github.com/harrtho/alfred-pyworkflow)에 의존했으나, 런타임 의존성 문제로 Go로 전면 대체되었습니다.

LICENSE
--------------
 - MIT
