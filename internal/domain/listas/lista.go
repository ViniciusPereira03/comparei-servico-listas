package listas

import "time"

type StatusLista string

const (
	StatusAberta    StatusLista = "ABERTA"
	StatusFechada   StatusLista = "FECHADA"
	StatusCancelada StatusLista = "CANCELADA"
)

type Lista struct {
	ID            int64       `json:"id"`
	UserID        string      `json:"user_id"`
	Nome          string      `json:"nome"`
	Status        StatusLista `json:"status"`
	TotalPrevisto float64     `json:"total_previsto"`
	TotalFinal    float64     `json:"total_final"`
	Itens         []ItemLista `json:"itens"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type ItemLista struct {
	ID            int64   `json:"id"`
	ListaID       int64   `json:"lista_id"`
	ProdutoID     int64   `json:"produto_id"`
	MercadoID     *int64  `json:"mercado_id"`
	Quantidade    float64 `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Checked       bool    `json:"checked"`
}
