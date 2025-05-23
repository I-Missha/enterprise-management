package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Master представляет информацию о мастере
type Master struct {
	ID              int    `json:"id"`
	AreaID          int    `json:"area_id"`
	AreaName        string `json:"area_name,omitempty"`
	EngineerID      int    `json:"engineer_id"`
	Name            string `json:"name"`
	CategoryID      int    `json:"category_id"`
	CategoryName    string `json:"category_name,omitempty"`
	HireDate        string `json:"hire_date,omitempty"`
	CurrentPosition string `json:"current_position,omitempty"`
	HallID          int    `json:"hall_id"`
	HallName        string `json:"hall_name,omitempty"`
}

// getMastersByArea получает список мастеров указанного участка
func (h *Handler) getMastersByArea(w http.ResponseWriter, r *http.Request) {
	areaID := chi.URLParam(r, "area_id")
	if areaID == "" {
		http.Error(w, "Area ID is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			m.id, 
			m.area_id, 
			pa.name as area_name,
			m.engineer_id, 
			e.name, 
			eng.category_id,
			ce.name as category_name,
			e.hire_date,
			e.current_position,
			pa.hall_id,
			ph.name as hall_name
		FROM masters m
		JOIN engineer eng ON m.engineer_id = eng.employee_id
		JOIN employee e ON eng.employee_id = e.id
		JOIN category_engineer ce ON eng.category_id = ce.id
		JOIN production_area pa ON m.area_id = pa.id
		JOIN production_halls ph ON pa.hall_id = ph.id
		WHERE m.area_id = $1
		ORDER BY e.name
	`

	rows, err := h.db.Query(context.Background(), query, areaID)
	if err != nil {
		h.logger.Printf("Error querying database for masters by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var masters []Master
	for rows.Next() {
		var master Master
		var areaName, categoryName, hireDate, currentPosition, hallName sql.NullString

		if err := rows.Scan(
			&master.ID,
			&master.AreaID,
			&areaName,
			&master.EngineerID,
			&master.Name,
			&master.CategoryID,
			&categoryName,
			&hireDate,
			&currentPosition,
			&master.HallID,
			&hallName,
		); err != nil {
			h.logger.Printf("Error scanning row for masters by area: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if areaName.Valid {
			master.AreaName = areaName.String
		}

		if categoryName.Valid {
			master.CategoryName = categoryName.String
		}

		if hireDate.Valid {
			master.HireDate = hireDate.String
		}

		if currentPosition.Valid {
			master.CurrentPosition = currentPosition.String
		}

		if hallName.Valid {
			master.HallName = hallName.String
		}

		masters = append(masters, master)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for masters by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(masters); err != nil {
		h.logger.Printf("Error encoding json for masters by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getMastersByHall получает список мастеров указанного цеха
func (h *Handler) getMastersByHall(w http.ResponseWriter, r *http.Request) {
	hallID := chi.URLParam(r, "hall_id")
	if hallID == "" {
		http.Error(w, "Hall ID is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			m.id, 
			m.area_id, 
			pa.name as area_name,
			m.engineer_id, 
			e.name, 
			eng.category_id,
			ce.name as category_name,
			e.hire_date,
			e.current_position,
			pa.hall_id,
			ph.name as hall_name
		FROM masters m
		JOIN engineer eng ON m.engineer_id = eng.employee_id
		JOIN employee e ON eng.employee_id = e.id
		JOIN category_engineer ce ON eng.category_id = ce.id
		JOIN production_area pa ON m.area_id = pa.id
		JOIN production_halls ph ON pa.hall_id = ph.id
		WHERE pa.hall_id = $1
		ORDER BY pa.name, e.name
	`

	rows, err := h.db.Query(context.Background(), query, hallID)
	if err != nil {
		h.logger.Printf("Error querying database for masters by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var masters []Master
	for rows.Next() {
		var master Master
		var areaName, categoryName, hireDate, currentPosition, hallName sql.NullString

		if err := rows.Scan(
			&master.ID,
			&master.AreaID,
			&areaName,
			&master.EngineerID,
			&master.Name,
			&master.CategoryID,
			&categoryName,
			&hireDate,
			&currentPosition,
			&master.HallID,
			&hallName,
		); err != nil {
			h.logger.Printf("Error scanning row for masters by hall: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if areaName.Valid {
			master.AreaName = areaName.String
		}

		if categoryName.Valid {
			master.CategoryName = categoryName.String
		}

		if hireDate.Valid {
			master.HireDate = hireDate.String
		}

		if currentPosition.Valid {
			master.CurrentPosition = currentPosition.String
		}

		if hallName.Valid {
			master.HallName = hallName.String
		}

		masters = append(masters, master)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for masters by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(masters); err != nil {
		h.logger.Printf("Error encoding json for masters by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
