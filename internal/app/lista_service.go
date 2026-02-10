package app

import (
	interfaces "comparei-servico-listas/internal/domain/interface"
	"comparei-servico-listas/internal/domain/listas"
	"errors"
)

type ListaService struct {
	repo interfaces.ListaRepository
}

func NewListaService(repo interfaces.ListaRepository) *ListaService {
	return &ListaService{repo: repo}
}

func (s *ListaService) CreateLista(lista *listas.Lista) (int64, error) {
	hasOpen, err := s.repo.HasOpenList(lista.UserID)
	if err != nil {
		return 0, err
	}
	if hasOpen {
		return 0, errors.New("usuário já possui uma lista em aberto!")
	}

	lista.Status = listas.StatusAberta
	lista.TotalPrevisto = 0
	lista.TotalFinal = 0

	return s.repo.Create(lista)
}

func (s *ListaService) GetByID(userID string, listaID int64) (*listas.Lista, error) {
	return s.repo.GetByID(listaID, userID)
}

func (s *ListaService) GetListasUsuario(userID string) ([]*listas.Lista, error) {
	return s.repo.GetAll(userID)
}

func (s *ListaService) AddItem(userID string, item *listas.ItemLista) error {
	// 1. Validar se a lista pertence ao usuário
	lista, err := s.repo.GetByID(item.ListaID, userID)
	if err != nil {
		return err
	}
	if lista == nil {
		return errors.New("lista não encontrada ou acesso negado")
	}
	if lista.Status != listas.StatusAberta {
		return errors.New("não é possível editar uma lista fechada")
	}

	// 2. Adicionar Item
	err = s.repo.AddItem(item)
	if err != nil {
		return err
	}

	return s.recalculateTotals(item.ListaID)
}

func (s *ListaService) RemoveItem(itemID int64) error {
	return s.repo.RemoveItem(itemID)
}

func (s *ListaService) ToggleItemCheck(userID string, itemID int64, checked bool) error {
	item, err := s.repo.GetItem(itemID)
	if err != nil {
		return err
	}

	lista, err := s.repo.GetByID(item.ListaID, userID)
	if err != nil || lista == nil {
		return errors.New("lista não encontrada")
	}

	item.Checked = checked
	err = s.repo.UpdateItem(item)
	if err != nil {
		return err
	}

	// if checked {
	// 	// Publicar evento para confirmar preço (Serviço Produtos)
	// 	go publisher.PubConfirmarPreco(item.ProdutoID, item.MercadoID, item.PrecoUnitario)

	// 	// Publicar log (Serviço Logs)
	// 	go publisher.PubLogEvento(userID, "ITEM_COMPRADO", fmt.Sprintf("Produto %d comprado na lista %d", item.ProdutoID, item.ListaID))
	// }

	return s.recalculateTotals(item.ListaID)
}

// RF3: Método auxiliar de cálculo
func (s *ListaService) recalculateTotals(listaID int64) error {
	// Busca lista completa com itens
	// Itera somando (Quantidade * Preco)
	// Atualiza tabela 'listas'
	// *Nota: Implementação detalhada depende do Repository retornando itens*
	return nil
}

func (s *ListaService) UpdatePricesFromEvent(produtoID int64, mercadoID int64, novoPreco float64) error {
	err := s.repo.UpdatePriceInOpenLists(produtoID, mercadoID, novoPreco)
	if err != nil {
		return err
	}

	return nil
}
