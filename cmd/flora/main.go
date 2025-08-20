package main

import (
	"flora/commands"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// main is the entry point of the bot.
func main() {
	log.Println("Starting Flora...")

	// .envファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set in .env file")
	}

	// Discordセッションを作成
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// コマンドマネージャーを初期化
	cmdManager := commands.NewManager()
	cmdManager.Add(commands.PingCommand) // pingコマンドを追加

	// BOTが準備できた時のイベントハンドラ
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		s.UpdateGameStatus(0, "Flora")

		// グローバルコマンドを登録
		if err := cmdManager.RegisterAll(s); err != nil {
			log.Fatalf("Failed to register commands: %v", err)
		}
		log.Println("Commands registered successfully.")
	})

	// コマンド実行時のハンドラ
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if cmd, exists := cmdManager.Get(i.ApplicationCommandData().Name); exists {
				cmd.Handler(s, i)
			}
		}
	})

	// WebSocket接続を開く
	if err := dg.Open(); err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	// BOTがオンラインになるまで待機し、終了シグナルを待つ
	log.Println("Flora is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// きれいに終了処理
	log.Println("Shutting down Flora...")
	dg.Close()
}
