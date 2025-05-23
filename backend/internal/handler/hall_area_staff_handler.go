package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Модели для представления данных об участках
type ProductionArea struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	HallID int    `json:"hall_id"`
}

type AreaWithBoss struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	HallID    int    `json:"hall_id"`
	HallName  string `json:"hall_name,omitempty"`
	BossID    int    `json:"boss_id"`
	BossName  string `json:"boss_name"`
	BossEmail string `json:"boss_email,omitempty"`
}

type AreasResponse struct {
	TotalCount int            `json:"total_count"`
	Areas      []AreaWithBoss `json:"areas"`
}

// getAreasByHall получает список участков указанного цеха и их начальников
func (h *Handler) getAreasByHall(w http.ResponseWriter, r *http.Request) {
	hallID := chi.URLParam(r, "hall_id")
	if hallID == "" {
		http.Error(w, "Hall ID is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			pa.id, 
			pa.name, 
			pa.hall_id,
			ph.name as hall_name,
			e.id as boss_id,
			e.name as boss_name,
			e.current_position as boss_position
		FROM production_area pa
		JOIN production_halls ph ON pa.hall_id = ph.id
		LEFT JOIN area_boss ab ON pa.id = ab.area_id
		LEFT JOIN engineer eng ON ab.engineer_id = eng.employee_id
		LEFT JOIN employee e ON eng.employee_id = e.id
		WHERE pa.hall_id = $1
		ORDER BY pa.name
	`

	rows, err := h.db.Query(context.Background(), query, hallID)
	if err != nil {
		h.logger.Printf("Error querying database for areas by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var areas []AreaWithBoss
	for rows.Next() {
		var area AreaWithBoss
		var bossID, bossName, bossPosition, hallName sql.NullString

		if err := rows.Scan(
			&area.ID,
			&area.Name,
			&area.HallID,
			&hallName,
			&bossID,
			&bossName,
			&bossPosition,
		); err != nil {
			h.logger.Printf("Error scanning row for areas by hall: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if hallName.Valid {
			area.HallName = hallName.String
		}

		if bossID.Valid {
			area.BossID, _ = strconv.Atoi(bossID.String)
		}

		if bossName.Valid {
			area.BossName = bossName.String
		}

		areas = append(areas, area)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for areas by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := AreasResponse{
		TotalCount: len(areas),
		Areas:      areas,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding json for areas by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllAreas получает список всех участков предприятия и их начальников
func (h *Handler) getAllAreas(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			pa.id, 
			pa.name, 
			pa.hall_id,
			ph.name as hall_name,
			e.id as boss_id,
			e.name as boss_name,
			e.current_position as boss_position
		FROM production_area pa
		JOIN production_halls ph ON pa.hall_id = ph.id
		LEFT JOIN area_boss ab ON pa.id = ab.area_id
		LEFT JOIN engineer eng ON ab.engineer_id = eng.employee_id
		LEFT JOIN employee e ON eng.employee_id = e.id
		ORDER BY pa.hall_id, pa.name
	`

	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		h.logger.Printf("Error querying database for all areas: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var areas []AreaWithBoss
	for rows.Next() {
		var area AreaWithBoss
		var bossID, bossName, bossPosition, hallName sql.NullString

		if err := rows.Scan(
			&area.ID,
			&area.Name,
			&area.HallID,
			&hallName,
			&bossID,
			&bossName,
			&bossPosition,
		); err != nil {
			h.logger.Printf("Error scanning row for all areas: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if hallName.Valid {
			area.HallName = hallName.String
		}

		if bossID.Valid {
			area.BossID, _ = strconv.Atoi(bossID.String)
		}

		if bossName.Valid {
			area.BossName = bossName.String
		}

		areas = append(areas, area)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for all areas: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := AreasResponse{
		TotalCount: len(areas),
		Areas:      areas,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding json for all areas: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAreaBoss получает информацию о начальнике конкретного участка
func (h *Handler) getAreaBoss(w http.ResponseWriter, r *http.Request) {
	areaID := chi.URLParam(r, "area_id")
	if areaID == "" {
		http.Error(w, "Area ID is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			e.id as employee_id,
			e.name,
			e.current_position,
			eng.category_id,
			ce.name as category_name
		FROM area_boss ab
		JOIN engineer eng ON ab.engineer_id = eng.employee_id
		JOIN employee e ON eng.employee_id = e.id
		JOIN category_engineer ce ON eng.category_id = ce.id
		WHERE ab.area_id = $1
	`

	var boss struct {
		EmployeeID      int    `json:"employee_id"`
		Name            string `json:"name"`
		CurrentPosition string `json:"current_position,omitempty"`
		CategoryID      int    `json:"category_id"`
		CategoryName    string `json:"category_name"`
	}

	var currentPosition, categoryName sql.NullString

	err := h.db.QueryRow(context.Background(), query, areaID).Scan(
		&boss.EmployeeID,
		&boss.Name,
		&currentPosition,
		&boss.CategoryID,
		&categoryName,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Area boss not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Printf("Error querying database for area boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if currentPosition.Valid {
		boss.CurrentPosition = currentPosition.String
	}

	if categoryName.Valid {
		boss.CategoryName = categoryName.String
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(boss); err != nil {
		h.logger.Printf("Error encoding json for area boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
