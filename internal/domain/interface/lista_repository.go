package interfaces

import "comparei-servico-listas/internal/domain/listas"

type ListaRepository interface {
	HasOpenList(userID string) (bool, error)
	Create(lista *listas.Lista) (int64, error)

	GetByID(id int64, userID string) (*listas.Lista, error)
	FinalizaLista(listaID int64, userID string) error
	GetAll(userID string) ([]*listas.Lista, error)
	Update(lista *listas.Lista) error

	AddItem(item *listas.ItemLista) error
	RemoveItem(itemID int64) error
	UpdateItem(item *listas.ItemLista) error
	GetItem(itemID int64) (*listas.ItemLista, error)

	UpdatePriceInOpenLists(produtoID int64, mercadoID int64, novoPreco float64) error
}
