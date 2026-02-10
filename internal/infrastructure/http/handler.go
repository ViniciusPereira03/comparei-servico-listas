package http

import (
	"comparei-servico-listas/internal/app"
	"comparei-servico-listas/internal/domain/listas"
	"comparei-servico-listas/internal/infrastructure/http/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type ListaHandler struct {
	Service *app.ListaService
}

func NewListaHandler(service *app.ListaService) *ListaHandler {
	return &ListaHandler{Service: service}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":    err.Error(),
		"mensagem": message,
	})
}

func validaToken(w http.ResponseWriter, r *http.Request) (string, error) {
	secret := os.Getenv("USER_JWT_SECRET")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("Missing token")
	}

	// Remover o prefixo "Bearer " se existir
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Verificar o token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	// Acessar os dados (claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id := claims["id"]
		return fmt.Sprintf("%v", id), nil
	}

	return "", fmt.Errorf("Erro ao decodificar token")
}

func (h *ListaHandler) CreateLista(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CRIAR LISTA")
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	fmt.Println("USER ID: ", userID)

	var req dto.CreateListaDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro no payload JSON", http.StatusBadRequest)
		return
	}

	novaLista := &listas.Lista{
		UserID: userID,
		Nome:   req.Nome,
	}

	fmt.Println("NOVA LISTA: ", novaLista)

	id, err := h.Service.CreateLista(novaLista)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lista, err := h.Service.GetByID(userID, id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lista)
}

func (h *ListaHandler) GetListas(w http.ResponseWriter, r *http.Request) {
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	listas, err := h.Service.GetListasUsuario(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(listas)
}

func (h *ListaHandler) GetListaByID(w http.ResponseWriter, r *http.Request) {
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	lista, err := h.Service.GetByID(userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)

}

func (h *ListaHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	vars := mux.Vars(r)
	listaID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var req dto.AddItemDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro no payload JSON", http.StatusBadRequest)
		return
	}

	item := &listas.ItemLista{
		ListaID:       listaID,
		ProdutoID:     req.ProdutoID,
		MercadoID:     req.MercadoID,
		Quantidade:    req.Quantidade,
		PrecoUnitario: req.PrecoUnitario,
		Checked:       false,
	}

	if err := h.Service.AddItem(userID, item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ListaHandler) DelItem(w http.ResponseWriter, r *http.Request) {
	_, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	vars := mux.Vars(r)
	itemListaID, _ := strconv.ParseInt(vars["id"], 10, 64)

	err := h.Service.RemoveItem(itemListaID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ListaHandler) CheckItem(w http.ResponseWriter, r *http.Request) {
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	vars := mux.Vars(r)
	itemID, _ := strconv.ParseInt(vars["item_id"], 10, 64)

	var req dto.ToggleItemDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro no payload JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.ToggleItemCheck(userID, itemID, req.Checked)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
