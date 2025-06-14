# AlfKoreanSearch : Korean Search Workflow for Alfred 
![Test](../../actions/workflows/test-korean-ac.yml/badge.svg) ![Release](../../actions/workflows/release.yml/badge.svg)  
![GitHub stars](https://img.shields.io/github/stars/inchans/alfkoreansearch?style=flat&logo=apachespark)
![GitHub all releases](https://img.shields.io/github/downloads/inchanS/alfkoreansearch/total?logo=github) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/inchanS/alfkoreansearch?logo=rocket)  ![GitHub](https://img.shields.io/github/license/inchanS/alfkoreansearch)



Korean Search Workflow for Alfred
---------------------------------

Alfred에서 국립국어원 [우리말샘 웹사이트](https://opendict.korean.go.kr) 검색어 자동완성 워크플로우
<br>

> 이 프로젝트는 [@Kuniz](https://github.com/Kuniz)님의 [alfnaversearch 워크플로우](https://github.com/Kuniz/alfnaversearch)를 기반으로 구현하였습니다.

<br>  

Preview
--------

<img src="images/alfKoreanSearch.gif" width="600">

<br>  

Install workflow
--------------

- [releases](../../releases/latest) 페이지의 `AlfKoreanSearch.alfredworkflow`를 다운로드 받아서 실행한다.

- MacOS 12.3 이상의 경우
  - python3 설치
    - `brew install python`
    - `xcode-select --install`

- Alfred 4.0 이상 지원
- Python 2 사용 불가


General Usage
--------------
* `k ...`  : 검색어 입력 (연관검색어가 나열됨)  
* **Cmd + Y** : 검색결과를 미리 보기(웹브라우져 출력)

트리거가 되는 키워드`k`는 Alfred-workflow-AlfKoreanSearch 에서 개인에 맞게 직접 수정할 수 있습니다. 


Externel Module
--------------
 This workflow used alfred-workflow more than v0.0.2. Alfred-workflow can find there(https://github.com/deanishe/alfred-workflow).
 This workflow used alp(A Python Module for Alfred Workflows) module at v0.0.1. It created by Daniel Shannon. 
 Certifi : using ssl with default urllib

LICENSE
--------------
 - MIT
