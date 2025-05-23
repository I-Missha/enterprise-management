package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// LabEquipmentInfo представляет информацию об оборудовании лаборатории
type LabEquipmentInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LabID    int    `json:"lab_id"`
	LabName  string `json:"lab_name,omitempty"`
	ItemID   int    `json:"item_id,omitempty"`   // ID изделия, для которого использовалось оборудование
	ItemName string `json:"item_name,omitempty"` // Название изделия
	TestDate string `json:"test_date,omitempty"` // Дата использования оборудования для теста
}

// getLabEquipmentByItemOrCategoryInLabInPeriod получает список оборудования
func (h *Handler) getLabEquipmentByItemOrCategoryInLabInPeriod(w http.ResponseWriter, r *http.Request) {
	labID := chi.URLParam(r, "lab_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	itemID := r.URL.Query().Get("item_id")         // Необязательный параметр
	categoryID := r.URL.Query().Get("category_id") // Необязательный параметр

	if labID == "" || startDate == "" || endDate == "" {
		http.Error(w, "Lab ID, start date, and end date are required", http.StatusBadRequest)
		return
	}

	var queryArgs []interface{}
	queryBuilder := strings.Builder{}

	queryBuilder.WriteString(`
		SELECT DISTINCT
			le.id,
			le.name,
			le.lab_id,
			tl.name AS lab_name,
			i.id AS item_id,
			i.name AS item_name,
			leit.test_date
		FROM lab_equip le
		JOIN testing_laboratory tl ON le.lab_id = tl.id
		JOIN lab_equip_item_test leit ON le.id = leit.lab_equip_id
		JOIN item i ON leit.item_id = i.id
		JOIN type_item ti ON i.type_id = ti.id
		WHERE le.lab_id = $1
		  AND leit.test_date >= $2
		  AND leit.test_date <= $3
	`)
	queryArgs = append(queryArgs, labID, startDate, endDate)
	argCounter := 3

	if itemID != "" {
		argCounter++
		queryBuilder.WriteString(" AND i.id = $" + strconv.Itoa(argCounter))
		queryArgs = append(queryArgs, itemID)
	}

	if categoryID != "" {
		argCounter++
		queryBuilder.WriteString(" AND ti.category_id = $" + strconv.Itoa(argCounter))
		queryArgs = append(queryArgs, categoryID)
	}

	queryBuilder.WriteString(" ORDER BY le.name, leit.test_date")

	rows, err := h.db.Query(context.Background(), queryBuilder.String(), queryArgs...)
	if err != nil {
		h.logger.Printf("Error querying database for lab equipment: %v, query: %s, args: %v", err, queryBuilder.String(), queryArgs)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var equipmentList []LabEquipmentInfo
	for rows.Next() {
		var equip LabEquipmentInfo
		var labName, itemName sql.NullString
		var testDate pgtype.Date

		if err := rows.Scan(
			&equip.ID,
			&equip.Name,
			&equip.LabID,
			&labName,
			&equip.ItemID,
			&itemName,
			&testDate,
		); err != nil {
			h.logger.Printf("Error scanning row for lab equipment: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if labName.Valid {
			equip.LabName = labName.String
		}
		if itemName.Valid {
			equip.ItemName = itemName.String
		}
		if testDate.Valid {
			equip.TestDate = testDate.Time.Format("2006-01-02")
		}
		equipmentList = append(equipmentList, equip)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for lab equipment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(equipmentList); err != nil {
		h.logger.Printf("Error encoding json for lab equipment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getLabEquipmentInLabInPeriod получает список всего оборудования в лаборатории за период
func (h *Handler) getLabEquipmentInLabInPeriod(w http.ResponseWriter, r *http.Request) {
	labID := chi.URLParam(r, "lab_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if labID == "" || startDate == "" || endDate == "" {
		http.Error(w, "Lab ID, start date, and end date are required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT DISTINCT
			le.id,
			le.name,
			le.lab_id,
			tl.name AS lab_name,
			i.id AS item_id,
			i.name AS item_name,
			leit.test_date
		FROM lab_equip le
		JOIN testing_laboratory tl ON le.lab_id = tl.id
		JOIN lab_equip_item_test leit ON le.id = leit.lab_equip_id
		JOIN item i ON leit.item_id = i.id
		WHERE le.lab_id = $1
		  AND leit.test_date >= $2
		  AND leit.test_date <= $3
		ORDER BY le.name, leit.test_date
	`

	rows, err := h.db.Query(context.Background(), query, labID, startDate, endDate)
	if err != nil {
		h.logger.Printf("Error querying database for lab equipment in period: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var equipmentList []LabEquipmentInfo
	for rows.Next() {
		var equip LabEquipmentInfo
		var labName, itemName sql.NullString
		var testDate pgtype.Date

		if err := rows.Scan(
			&equip.ID,
			&equip.Name,
			&equip.LabID,
			&labName,
			&equip.ItemID,
			&itemName,
			&testDate,
		); err != nil {
			h.logger.Printf("Error scanning row for lab equipment in period: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if labName.Valid {
			equip.LabName = labName.String
		}
		if itemName.Valid {
			equip.ItemName = itemName.String
		}
		if testDate.Valid {
			equip.TestDate = testDate.Time.Format("2006-01-02")
		}
		equipmentList = append(equipmentList, equip)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for lab equipment in period: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(equipmentList); err != nil {
		h.logger.Printf("Error encoding json for lab equipment in period: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
