package repository

import (
	"comparei-servico-listas/internal/domain/listas"
	"database/sql"
	"log"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// --- Listas ---

func (r *MySQLRepository) HasOpenList(userID string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM listas WHERE user_id = ? AND status = 'ABERTA' AND deleted_at IS NULL"
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MySQLRepository) Create(lista *listas.Lista) (int64, error) {
	query := "INSERT INTO listas (user_id, nome, status, total_previsto, total_final) VALUES (?, ?, ?, ?, ?)"
	res, err := r.db.Exec(query, lista.UserID, lista.Nome, lista.Status, lista.TotalPrevisto, lista.TotalFinal)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *MySQLRepository) GetByID(id int64, userID string) (*listas.Lista, error) {
	query := "SELECT id, user_id, nome, status, total_previsto, total_final, created_at, updated_at FROM listas WHERE id = ? AND user_id = ? AND deleted_at IS NULL"

	lista := &listas.Lista{}
	err := r.db.QueryRow(query, id, userID).Scan(
		&lista.ID, &lista.UserID, &lista.Nome, &lista.Status,
		&lista.TotalPrevisto, &lista.TotalFinal, &lista.CreatedAt, &lista.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Buscar itens da lista
	itens, err := r.getItemsByListaID(lista.ID)
	if err != nil {
		return nil, err
	}
	lista.Itens = itens

	return lista, nil
}

func (r *MySQLRepository) FinalizaLista(listaID int64, userID string) error {
	query := "UPDATE listas SET status=? WHERE id=? AND user_id=? AND deleted_at IS NULL"
	_, err := r.db.Exec(query, listas.StatusFechada, listaID, userID)
	return err
}

func (r *MySQLRepository) GetAll(userID string) ([]*listas.Lista, error) {
	query := "SELECT id, user_id, nome, status, total_previsto, total_final, created_at, updated_at FROM listas WHERE user_id = ? AND deleted_at IS NULL ORDER BY created_at DESC"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listasArr []*listas.Lista
	for rows.Next() {
		l := &listas.Lista{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.Nome, &l.Status, &l.TotalPrevisto, &l.TotalFinal, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		listasArr = append(listasArr, l)
	}
	return listasArr, nil
}

func (r *MySQLRepository) Update(lista *listas.Lista) error {
	query := "UPDATE listas SET nome=?, status=?, total_previsto=?, total_final=? WHERE id=? AND user_id=? AND deleted_at IS NULL"
	_, err := r.db.Exec(query, lista.Nome, lista.Status, lista.TotalPrevisto, lista.TotalFinal, lista.ID, lista.UserID)
	return err
}

// --- Itens ---

func (r *MySQLRepository) AddItem(item *listas.ItemLista) error {
	query := "INSERT INTO itens_lista (lista_id, produto_id, mercado_id, quantidade, preco_unitario, checked) VALUES (?, ?, ?, ?, ?, ?)"
	res, err := r.db.Exec(query, item.ListaID, item.ProdutoID, item.MercadoID, item.Quantidade, item.PrecoUnitario, item.Checked)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	item.ID = id
	return nil
}

func (r *MySQLRepository) RemoveItem(itemID int64) error {
	query := "UPDATE itens_lista SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.Exec(query, itemID)
	return err
}

func (r *MySQLRepository) UpdateItem(item *listas.ItemLista) error {
	query := "UPDATE itens_lista SET quantidade=?, preco_unitario=?, checked=?, mercado_id=? WHERE id=? AND deleted_at IS NULL"
	_, err := r.db.Exec(query, item.Quantidade, item.PrecoUnitario, item.Checked, item.MercadoID, item.ID)
	return err
}

func (r *MySQLRepository) GetItem(itemID int64) (*listas.ItemLista, error) {
	query := "SELECT id, lista_id, produto_id, mercado_id, quantidade, preco_unitario, checked FROM itens_lista WHERE id = ? AND deleted_at IS NULL"
	item := &listas.ItemLista{}
	err := r.db.QueryRow(query, itemID).Scan(&item.ID, &item.ListaID, &item.ProdutoID, &item.MercadoID, &item.Quantidade, &item.PrecoUnitario, &item.Checked)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Auxiliar privado para buscar itens
func (r *MySQLRepository) getItemsByListaID(listaID int64) ([]listas.ItemLista, error) {
	query := "SELECT id, lista_id, produto_id, mercado_id, quantidade, preco_unitario, checked FROM itens_lista WHERE lista_id = ? AND deleted_at IS NULL"
	rows, err := r.db.Query(query, listaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itens []listas.ItemLista
	for rows.Next() {
		var i listas.ItemLista
		if err := rows.Scan(&i.ID, &i.ListaID, &i.ProdutoID, &i.MercadoID, &i.Quantidade, &i.PrecoUnitario, &i.Checked); err != nil {
			return nil, err
		}
		itens = append(itens, i)
	}
	return itens, nil
}

// --- Atualização em Massa (RF4) ---

func (r *MySQLRepository) UpdatePriceInOpenLists(produtoID int64, mercadoID int64, novoPreco float64) error {
	// Atualiza o preço unitário de itens que estão em listas ABERTAS e correspondem ao produto/mercado
	query := `
		UPDATE itens_lista il
		JOIN listas l ON il.lista_id = l.id
		SET il.preco_unitario = ?
		WHERE il.produto_id = ? 
		  AND il.mercado_id = ? 
		  AND l.status = 'ABERTA'
		  AND il.checked = FALSE -- não mudar preço se já comprou
		  AND il.deleted_at IS NULL
	`
	result, err := r.db.Exec(query, novoPreco, produtoID, mercadoID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Preço atualizado em %d itens de listas abertas.", rowsAffected)

	// Nota: Idealmente, dispararíamos o recálculo dos totais das listas afetadas aqui.
	return nil
}
