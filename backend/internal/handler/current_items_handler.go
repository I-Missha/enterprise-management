package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CurrentItem представляет информацию о текущем изделии
type CurrentItem struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	TypeID       int    `json:"type_id"`
	TypeName     string `json:"type_name,omitempty"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name,omitempty"`
	HallID       int    `json:"hall_id"`
	HallName     string `json:"hall_name,omitempty"`
	AreaID       int    `json:"area_id,omitempty"`
	AreaName     string `json:"area_name,omitempty"`
	Status       string `json:"status"`
}

// getCurrentItemsByCategoryAndArea получает текущие изделия определенной категории собираемые указанным участком
func (h *Handler) getCurrentItemsByCategoryAndArea(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	areaID := chi.URLParam(r, "area_id")

	if categoryID == "" || areaID == "" {
		http.Error(w, "Category ID and Area ID are required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			i.id, 
			i.name, 
			i.type_id, 
			ti.name as type_name,
			ti.category_id,
			ci.name as category_name,
			i.hall_id,
			ph.name as hall_name,
			ai.area_id,
			pa.name as area_name,
			i.status
		FROM item i
		JOIN type_item ti ON i.type_id = ti.id
		JOIN category_item ci ON ti.category_id = ci.id
		JOIN production_halls ph ON i.hall_id = ph.id
		JOIN areas_items ai ON i.id = ai.item_id
		JOIN production_area pa ON ai.area_id = pa.id
		WHERE ti.category_id = $1 
		  AND ai.area_id = $2
		  AND i.status = 'in_progress'
		ORDER BY i.name
	`

	rows, err := h.db.Query(context.Background(), query, categoryID, areaID)
	if err != nil {
		h.logger.Printf("Error querying database for current items by category and area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CurrentItem
	for rows.Next() {
		var item CurrentItem
		var typeName, categoryName, hallName, areaName sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.TypeID,
			&typeName,
			&item.CategoryID,
			&categoryName,
			&item.HallID,
			&hallName,
			&item.AreaID,
			&areaName,
			&item.Status,
		); err != nil {
			h.logger.Printf("Error scanning row for current items by category and area: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if typeName.Valid {
			item.TypeName = typeName.String
		}

		if categoryName.Valid {
			item.CategoryName = categoryName.String
		}

		if hallName.Valid {
			item.HallName = hallName.String
		}

		if areaName.Valid {
			item.AreaName = areaName.String
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for current items by category and area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json for current items by category and area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getCurrentItemsByCategoryAndHall получает текущие изделия определенной категории собираемые указанным цехом
func (h *Handler) getCurrentItemsByCategoryAndHall(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	hallID := chi.URLParam(r, "hall_id")

	if categoryID == "" || hallID == "" {
		http.Error(w, "Category ID and Hall ID are required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			i.id, 
			i.name, 
			i.type_id, 
			ti.name as type_name,
			ti.category_id,
			ci.name as category_name,
			i.hall_id,
			ph.name as hall_name,
			i.status
		FROM item i
		JOIN type_item ti ON i.type_id = ti.id
		JOIN category_item ci ON ti.category_id = ci.id
		JOIN production_halls ph ON i.hall_id = ph.id
		WHERE ti.category_id = $1 
		  AND i.hall_id = $2
		  AND i.status = 'in_progress'
		ORDER BY i.name
	`

	rows, err := h.db.Query(context.Background(), query, categoryID, hallID)
	if err != nil {
		h.logger.Printf("Error querying database for current items by category and hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CurrentItem
	for rows.Next() {
		var item CurrentItem
		var typeName, categoryName, hallName sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.TypeID,
			&typeName,
			&item.CategoryID,
			&categoryName,
			&item.HallID,
			&hallName,
			&item.Status,
		); err != nil {
			h.logger.Printf("Error scanning row for current items by category and hall: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if typeName.Valid {
			item.TypeName = typeName.String
		}

		if categoryName.Valid {
			item.CategoryName = categoryName.String
		}

		if hallName.Valid {
			item.HallName = hallName.String
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for current items by category and hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json for current items by category and hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getCurrentItemsByCategory получает все текущие изделия определенной категории собираемые на предприятии
func (h *Handler) getCurrentItemsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")

	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			i.id, 
			i.name, 
			i.type_id, 
			ti.name as type_name,
			ti.category_id,
			ci.name as category_name,
			i.hall_id,
			ph.name as hall_name,
			i.status
		FROM item i
		JOIN type_item ti ON i.type_id = ti.id
		JOIN category_item ci ON ti.category_id = ci.id
		JOIN production_halls ph ON i.hall_id = ph.id
		WHERE ti.category_id = $1 
		  AND i.status = 'in_progress'
		ORDER BY i.hall_id, i.name
	`

	rows, err := h.db.Query(context.Background(), query, categoryID)
	if err != nil {
		h.logger.Printf("Error querying database for current items by category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CurrentItem
	for rows.Next() {
		var item CurrentItem
		var typeName, categoryName, hallName sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.TypeID,
			&typeName,
			&item.CategoryID,
			&categoryName,
			&item.HallID,
			&hallName,
			&item.Status,
		); err != nil {
			h.logger.Printf("Error scanning row for current items by category: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if typeName.Valid {
			item.TypeName = typeName.String
		}

		if categoryName.Valid {
			item.CategoryName = categoryName.String
		}

		if hallName.Valid {
			item.HallName = hallName.String
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for current items by category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json for current items by category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getCurrentItemsByArea получает все текущие изделия собираемые указанным участком
func (h *Handler) getCurrentItemsByArea(w http.ResponseWriter, r *http.Request) {
	areaID := chi.URLParam(r, "area_id")

	if areaID == "" {
		http.Error(w, "Area ID is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			i.id, 
			i.name, 
			i.type_id, 
			ti.name as type_name,
			ti.category_id,
			ci.name as category_name,
			i.hall_id,
			ph.name as hall_name,
			ai.area_id,
			pa.name as area_name,
			i.status
		FROM item i
		JOIN type_item ti ON i.type_id = ti.id
		JOIN category_item ci ON ti.category_id = ci.id
		JOIN production_halls ph ON i.hall_id = ph.id
		JOIN areas_items ai ON i.id = ai.item_id
		JOIN production_area pa ON ai.area_id = pa.id
		WHERE ai.area_id = $1
		  AND i.status = 'in_progress'
		ORDER BY i.name
	`

	rows, err := h.db.Query(context.Background(), query, areaID)
	if err != nil {
		h.logger.Printf("Error querying database for current items by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CurrentItem
	for rows.Next() {
		var item CurrentItem
		var typeName, categoryName, hallName, areaName sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.TypeID,
			&typeName,
			&item.CategoryID,
			&categoryName,
			&item.HallID,
			&hallName,
			&item.AreaID,
			&areaName,
			&item.Status,
		); err != nil {
			h.logger.Printf("Error scanning row for current items by area: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if typeName.Valid {
			item.TypeName = typeName.String
		}

		if categoryName.Valid {
			item.CategoryName = categoryName.String
		}

		if hallName.Valid {
			item.HallName = hallName.String
		}

		if areaName.Valid {
			item.AreaName = areaName.String
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for current items by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json for current items by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getCurrentItemsByHall получает все текущие изделия собираемые указанным цехом
func (h *Handler) getCurrentItemsByHall(w http.ResponseWriter, r *http.Request) {
	hallID := chi.URLParam(r, "hall_id")

	if hallID == "" {
		http.Error(w, "Hall ID is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			i.id, 
			i.name, 
			i.type_id, 
			ti.name as type_name,
			ti.category_id,
			ci.name as category_name,
			i.hall_id,
			ph.name as hall_name,
			i.status
		FROM item i
		JOIN type_item ti ON i.type_id = ti.id
		JOIN category_item ci ON ti.category_id = ci.id
		JOIN production_halls ph ON i.hall_id = ph.id
		WHERE i.hall_id = $1
		  AND i.status = 'in_progress'
		ORDER BY i.name
	`

	rows, err := h.db.Query(context.Background(), query, hallID)
	if err != nil {
		h.logger.Printf("Error querying database for current items by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CurrentItem
	for rows.Next() {
		var item CurrentItem
		var typeName, categoryName, hallName sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.TypeID,
			&typeName,
			&item.CategoryID,
			&categoryName,
			&item.HallID,
			&hallName,
			&item.Status,
		); err != nil {
			h.logger.Printf("Error scanning row for current items by hall: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if typeName.Valid {
			item.TypeName = typeName.String
		}

		if categoryName.Valid {
			item.CategoryName = categoryName.String
		}

		if hallName.Valid {
			item.HallName = hallName.String
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for current items by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json for current items by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllCurrentItems получает все текущие изделия собираемые на предприятии
func (h *Handler) getAllCurrentItems(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			i.id, 
			i.name, 
			i.type_id, 
			ti.name as type_name,
			ti.category_id,
			ci.name as category_name,
			i.hall_id,
			ph.name as hall_name,
			i.status
		FROM item i
		JOIN type_item ti ON i.type_id = ti.id
		JOIN category_item ci ON ti.category_id = ci.id
		JOIN production_halls ph ON i.hall_id = ph.id
		WHERE i.status = 'in_progress'
		ORDER BY i.hall_id, i.name
	`

	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		h.logger.Printf("Error querying database for all current items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []CurrentItem
	for rows.Next() {
		var item CurrentItem
		var typeName, categoryName, hallName sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.TypeID,
			&typeName,
			&item.CategoryID,
			&categoryName,
			&item.HallID,
			&hallName,
			&item.Status,
		); err != nil {
			h.logger.Printf("Error scanning row for all current items: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if typeName.Valid {
			item.TypeName = typeName.String
		}

		if categoryName.Valid {
			item.CategoryName = categoryName.String
		}

		if hallName.Valid {
			item.HallName = hallName.String
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for all current items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json for all current items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
