package main

import (
	"comparei-servico-listas/internal/app"
	"comparei-servico-listas/internal/infrastructure/http"
	"comparei-servico-listas/internal/infrastructure/messaging/subscriber"
	"comparei-servico-listas/internal/infrastructure/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	httpNet "net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Carregar variáveis de ambiente
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	// 2. Conexão com MySQL
	dsn := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ")/" + os.Getenv("MYSQL_DB") + "?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao abrir conexão MySQL:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Erro ao conectar no MySQL:", err)
	}
	log.Println("✅ Conexão com MySQL estabelecida com sucesso!")

	// 3. Conexão com Redis (Mensageria)
	redisHost := os.Getenv("REDIS_MESSAGING_HOST")
	redisPort := os.Getenv("REDIS_MESSAGING_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	// Testar conexão Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("Erro ao conectar no Redis de mensageria:", err)
	}
	log.Println("✅ Conexão com Redis estabelecida com sucesso!")

	// 4. Inicialização de Dependências (Injeção de Dependência)

	// Repository
	listaRepo := repository.NewMySQLRepository(db)

	// Service
	listaService := app.NewListaService(listaRepo)

	// Handler
	listaHandler := http.NewListaHandler(listaService)

	// 5. Configurar Subscriber (Mensageria)
	// Injeta o service no subscriber para que ele possa chamar a lógica de negócio
	subscriber.SetListaService(listaService)

	// Inicia o subscriber em uma Goroutine (background) para não bloquear o servidor HTTP
	go func() {
		log.Println("📡 Iniciando Subscriber...")
		subscriber.SubPriceUpdates()
	}()

	// 6. Configurar Roteamento e Servidor HTTP
	router := http.NewRouter(listaHandler)

	// Middleware de Autenticação (Sugestão Simplificada)
	// Aqui você deve garantir que o handler consiga extrair o ID do usuário.
	// Se tiver um middleware criado, adicione aqui: router.Use(middleware.AuthMiddleware)

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8083" // Porta padrão sugerida para o serviço de listas
	}

	log.Println("🚀 Servidor rodando na porta " + serverPort)
	if err := httpNet.ListenAndServe(":"+serverPort, router); err != nil {
		log.Fatal("Erro fatal no servidor HTTP:", err)
	}
}
