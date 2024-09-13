# go-layout

[Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## Software Architecture

![project_layout](./doc/project_layout.png)

首先要分成兩個視角
1. 實際物理目錄存放位置
2. 抽象架構意義

思路
- 上述架構圖以抽象架構的角度進行繪圖, 圖片內容與物理存放位置不一定正相關  
  比如 `utility/`, `inject/` 雖然物理位置在 `pkg/` 底下, 但只是為了集中 go 程式才放在一起, 抽象意義並不會跟 `pkg/` 一起思考
- `pkg/` 是所有 go 程式的物理存放位置  
  唯一例外: config.go 抽象意義應該放在 `pkg/`, 但為了方便查找 config.yml, 所以 config.go 物理存放位置是 `configs/`
- `pkg/` 物理位置有四個目錄 `adapters/`, `app/`, `utility/`, `inject/`, 用抽象角度思考可以分為三個類型
    1. `pkg/`, `app/`, `adapters/`: 和商業邏輯相關的程式碼 及 global variable 依賴 (defaultErrorRegistry, defaultLogger, defaultShutdown, defaultVersion ...)
    2. `utility/`: 提供技術類型的元件, 期望未來可以成為公司專用 package, 避免每個 project 的做法不同
    3. `inject/`: 組裝 project code, 讓不同的進入端點(main.go, xxx_test.go)使用
- 箭頭表示套件依賴方向, 而不是資料流方向. 大多數的情況子層目錄依賴父層目錄
- project 商業邏輯的依賴終點是 `app/`, 有人也稱為 service layer
- Dependency Injection 用來處理和商業邏輯相關的依賴, 讓程式有可測試性. 而與商業邏輯無關的依賴則放在全域變數, 因為不會影響到邏輯的變更
- 不是 high level 去依賴 low level. 應該是 high level 提出它需要的功能, low level 去實作出這些功能, 達成 high level 的目標
