package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Модель для представления работ над изделием
type ItemWork struct {
	ID           int            `json:"id"`
	SequenceNum  int            `json:"sequence_num"`
	ItemID       int            `json:"item_id"`
	WorkTypeID   int            `json:"work_type_id"`
	WorkName     string         `json:"work_name"`
	AreaID       int            `json:"area_id"`
	AreaName     string         `json:"area_name,omitempty"`
	WorkTeamID   int            `json:"work_team_id"`
	WorkTeamName string         `json:"work_team_name,omitempty"`
	StartDate    sql.NullString `json:"start_date,omitempty"`
	EndDate      sql.NullString `json:"end_date,omitempty"`
}

// getItemWorks получает перечень работ, которые проходит указанное изделие
func (h *Handler) getItemWorks(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "item_id")
	if itemID == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			iwt.id, 
			iwt.seq_number, 
			iwt.item_id, 
			iwt.work_type_id,
			wt.work_name,
			wt.area_id,
			pa.name as area_name,
			wt.work_team_id,
			wteam.name as work_team_name,
			iwt.start_date,
			iwt.end_date
		FROM item_work_type iwt
		JOIN work_type wt ON iwt.work_type_id = wt.id
		JOIN production_area pa ON wt.area_id = pa.id
		JOIN work_team wteam ON wt.work_team_id = wteam.id
		WHERE iwt.item_id = $1
		ORDER BY iwt.seq_number
	`

	rows, err := h.db.Query(context.Background(), query, itemID)
	if err != nil {
		h.logger.Printf("Error querying database for item works: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var works []ItemWork
	for rows.Next() {
		var work ItemWork
		var areaName, workTeamName sql.NullString

		if err := rows.Scan(
			&work.ID,
			&work.SequenceNum,
			&work.ItemID,
			&work.WorkTypeID,
			&work.WorkName,
			&work.AreaID,
			&areaName,
			&work.WorkTeamID,
			&workTeamName,
			&work.StartDate,
			&work.EndDate,
		); err != nil {
			h.logger.Printf("Error scanning row for item works: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if areaName.Valid {
			work.AreaName = areaName.String
		}

		if workTeamName.Valid {
			work.WorkTeamName = workTeamName.String
		}

		works = append(works, work)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for item works: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(works); err != nil {
		h.logger.Printf("Error encoding json for item works: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
