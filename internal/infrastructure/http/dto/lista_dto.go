package dto

type CreateListaDTO struct {
	Nome string `json:"nome"`
}

type AddItemDTO struct {
	ProdutoID     int64   `json:"produto_id"`
	MercadoID     *int64  `json:"mercado_id"`
	Quantidade    float64 `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
}

type ToggleItemDTO struct {
	Checked bool `json:"checked"`
}
