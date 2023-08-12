封裝套件與持續整合(CI)。

將你的雲端 (AMI)、容器 (Docker)、OS (deb, rpm, pkg) 套件的組態設定與腳本放在 /build/package 目錄下。

將你的 CI (Travis CI, CircleCI, Drone CI) 的組態設定與腳本放在 /build/ci 目錄中。請注意：有些 CI 工具 (例如 Travis CI 等)，它們對這些組態設定檔案的位置非常挑剔。如果可能的話，請嘗試將檔案放在 /build/ci 目錄中，並連結 (linking) 這些檔案到 CI 工具期望它們出現的位置。