package main

import (
"fmt"
"os"

"aoroa/internal/server"
)

func main() {
	// 명령행 인자 확인
	if len(os.Args) > 1 && os.Args[1] == "server" {
		// 서버 모드
		runServer()
	} else {
		// 기본값: 사용법 출력
		printUsage()
	}
}

func runServer() {
	fmt.Println("=== 이슈 관리 API 서버 시작 ===")
	srv := server.New()
	srv.Run()
}

func printUsage() {
	fmt.Println("=== 이슈 관리 API ===")
	fmt.Println("사용법:")
	fmt.Println("  go run main.go server    # 서버 시작")
	fmt.Println("  go test ./... -v         # 테스트 실행")
	fmt.Println("\n서버 시작 후 다음 엔드포인트를 사용할 수 있습니다:")
	fmt.Println("  POST   /issue           # 이슈 생성")
	fmt.Println("  GET    /issues          # 이슈 목록 조회")
	fmt.Println("  GET    /issue/:id       # 특정 이슈 조회")
	fmt.Println("  PUT    /issue/:id       # 이슈 수정")
	fmt.Println("  GET    /health          # 헬스 체크")
}
