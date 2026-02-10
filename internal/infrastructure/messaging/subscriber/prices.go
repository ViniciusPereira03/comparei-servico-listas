package subscriber

import (
	"comparei-servico-listas/internal/app"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var listaService *app.ListaService

func SetListaService(service *app.ListaService) {
	listaService = service
}

// Struct do evento vindo do Promer
type MercadoProdutos struct {
	ID             int64      `json:"id"`
	ProdutoID      int64      `json:"id_produto"`
	MercadoID      int64      `json:"id_mercado"`
	PrecoUnitario  float32    `json:"preco_unitario"`
	NivelConfianca int32      `json:"nivel_confianca"`
	CreatedAt      time.Time  `json:"created_at"`
	ModifiedAt     time.Time  `json:"modified_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type PriceUpdateEvent struct {
	Id             int64           `json:"id"`
	UserID         string          `json:"user_id"`
	MercadoProduto MercadoProdutos `json:"mercado_produto"`
}

func SubPriceUpdates() {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_MESSAGING_HOST") + ":" + os.Getenv("REDIS_MESSAGING_PORT"),
	})

	ctx := context.Background()
	pubsub := rdb.Subscribe(ctx, "update_product") // Tópico definido no Promer
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var event PriceUpdateEvent
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Println("Erro ao decodificar evento de preço:", err)
			continue
		}

		log.Printf("Atualizando preço do produto %d nas listas...", event.MercadoProduto.ProdutoID)
		err := listaService.UpdatePricesFromEvent(event.MercadoProduto.ProdutoID, event.MercadoProduto.MercadoID, float64(event.MercadoProduto.PrecoUnitario))
		if err != nil {
			log.Println("Erro ao atualizar preços nas listas:", err)
		}
	}
}
