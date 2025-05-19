package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	//"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	logger *log.Logger
	db     *pgxpool.Pool
}

type ItemType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category_id int    `json:"category_id"`
}

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category_id int    `json:"category_id"`
	Type_id     int    `json:"type_id"`
	Hall_id     int    `json:"hall_id"`
	Status      string `json:"status"`
}

type ReadyItem struct {
	ItemId int    `json:"item_id"`
	Count  int    `json:"count"`
	Name   string `json:"name"`
}

func NewHandler(logger *log.Logger, db *pgxpool.Pool) *Handler {
	return &Handler{
		logger: logger,
		db:     db,
	}
}

func (h *Handler) InitRoutes(router chi.Router) {

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/items", func(r chi.Router) {
		r.Get("/ping", h.ping)

		r.Get("/types/category/{category_id}", h.getItemsTypeByCategory)
		r.Get("/types/hall/{hall_id}", h.getAllItemsTypeByHallId)
		r.Get("/types/hall/{hall_id}/category/{category_id}", h.getAllItemsTypeByCategoryAndHallId)
		r.Get("/category/{category_id}", h.getAllItemsByCategory)
		r.Get("/area/{area_id}", h.getAllItemsByAreaId)
		r.Get("/TimeInterval", h.getAllItemsByTimeInterval)
	})
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "pong"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getItemsTypeByCategory получает типы изделий по ID категории
func (h *Handler) getItemsTypeByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name FROM type_item WHERE category_id = $1`

	rows, err := h.db.Query(context.Background(), query, categoryID)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []ItemType
	for rows.Next() {
		var item ItemType
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsType получает все типы изделий по id цеха
func (h *Handler) getAllItemsTypeByHallId(w http.ResponseWriter, r *http.Request) {
	hallId := chi.URLParam(r, "hall_id")
	query := `SELECT type_item.id, type_item.name, type_item.category_id FROM (SELECT * FROM item WHERE hall_id = $1) AS filtered RIGHT JOIN type_item ON filtered.type_id = type_item.id`

	rows, err := h.db.Query(context.Background(), query, hallId)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []ItemType
	for rows.Next() {
		var item ItemType
		if err := rows.Scan(&item.ID, &item.Name, &item.Category_id); err != nil { // Исправлено item.Category_id на &item.Category_id
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsByCategory получает все изделия по id категории
func (h *Handler) getAllItemsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT item.id, item.name, item.type_id, item.hall_id, item.status FROM 
                item INNER JOIN (
                SELECT type_item.id 
                FROM type_item 
                WHERE type_item.category_id = $1
                ) as filtered on item.type_id = filtered.id`

	rows, err := h.db.Query(context.Background(), query, categoryID)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Type_id, &item.Hall_id, &item.Status); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsTypeByCategoryAndHallId получает все изделия по id категории и id цеха
func (h *Handler) getAllItemsTypeByCategoryAndHallId(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	hallId := chi.URLParam(r, "hall_id")
	if categoryID == "" || hallId == "" {
		http.Error(w, "Category ID and Hall ID are required", http.StatusBadRequest)
		return
	}

	query := `SELECT type_item.id, type_item.name FROM (SELECT * FROM item WHERE hall_id = $1 AND type_id = $2) AS filtered RIGHT JOIN type_item ON filtered.type_id = type_item.id`

	rows, err := h.db.Query(context.Background(), query, hallId, categoryID)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []ItemType
	for rows.Next() {
		var item ItemType
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsByHallId получает все изделия по id цеха
func (h *Handler) getAllItemsByHallId(w http.ResponseWriter, r *http.Request) {
	hallId := chi.URLParam(r, "hall_id")
	query := `SELECT item.id, item.name, item.type_id, item.hall_id, item.status FROM item WHERE item.hall_id = $1`

	rows, err := h.db.Query(context.Background(), query, hallId)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Type_id, &item.Hall_id, &item.Status); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsByAreaId получает все изделия по id участка
func (h *Handler) getAllItemsByAreaId(w http.ResponseWriter, r *http.Request) {
	areaId := chi.URLParam(r, "area_id")
	query := `SELECT item.id, item.name, item.type_id, item.hall_id, item.status
				FROM item
					INNER JOIN (
					SELECT areas_items.item_id
					FROM areas_items
					WHERE areas_items.area_id = $1
				) AS ai ON item.id = ai.item_id`

	rows, err := h.db.Query(context.Background(), query, areaId)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Type_id, &item.Hall_id, &item.Status); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsByTimeInterval получает все изделия по интервалу времени
func (h *Handler) getAllItemsByTimeInterval(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	finishDate := r.URL.Query().Get("finish_date")

	query := `SELECT
				i.id,
				i.name,
				SUM(ri.counter) as item_count
				FROM item i
         			INNER JOIN ready_item ri ON i.id = ri.item_id
				WHERE ri.start_date >= $1 AND ri.completion_date <= $2
				GROUP BY i.id`

	rows, err := h.db.Query(context.Background(), query, startDate, finishDate)

	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []ReadyItem

	for rows.Next() {
		var item ReadyItem
		if err := rows.Scan(&item.ItemId, &item.Name, &item.Count); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
