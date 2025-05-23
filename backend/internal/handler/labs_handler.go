package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TestingLab представляет информацию об испытательной лаборатории
type TestingLab struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// getItemTestingLabs получает перечень испытательных лабораторий, участвующих в испытаниях указанного изделия
func (h *Handler) getItemTestingLabs(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "item_id")
	if itemID == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	// Получаем список лабораторий, участвующих в испытаниях готового изделия
	labsQuery := `
		SELECT DISTINCT 
			tl.id, 
			tl.name
		FROM testing_laboratory tl
		JOIN test_ready_item tri ON tl.id = tri.lab_id
		JOIN ready_item ri ON tri.ready_item_id = ri.id
		WHERE ri.item_id = $1
		ORDER BY tl.name
	`

	rows, err := h.db.Query(context.Background(), labsQuery, itemID)
	if err != nil {
		h.logger.Printf("Error querying database for testing labs by item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var labs []TestingLab
	for rows.Next() {
		var lab TestingLab

		if err := rows.Scan(
			&lab.ID,
			&lab.Name,
		); err != nil {
			h.logger.Printf("Error scanning row for testing labs by item: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		labs = append(labs, lab)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for testing labs by item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(labs); err != nil {
		h.logger.Printf("Error encoding json response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
