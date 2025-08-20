package main

import (
	"flora/configs"
	"flora/database"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting Flora...")

	// 設定を読み込む
	cfg := configs.NewConfig()

	// DIコンテナを初期化してアプリケーションを生成
	app, err := InitializeApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// データベースのマイグレーションを実行
	if err := database.Migrate(app.DB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// WebSocket接続を開く
	dg := app.Session
	if err := dg.Open(); err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()

	log.Println("Flora is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down Flora...")
}
