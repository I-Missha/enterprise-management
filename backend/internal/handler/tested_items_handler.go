package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// TestedItem представляет информацию об изделии, прошедшем испытание
type TestedItem struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	TypeName       string `json:"type_name,omitempty"`
	CategoryName   string `json:"category_name,omitempty"`
	TestDateStart  string `json:"test_date_start"`
	TestDateFinish string `json:"test_date_finish"`
	TestResult     string `json:"test_result,omitempty"`
	LaboratoryName string `json:"laboratory_name,omitempty"`
	CompletionDate string `json:"completion_date"`
}

// getTestedItemsByLabAndCategoryInPeriod получает изделия категории, испытанные в лаборатории за период
func (h *Handler) getTestedItemsByLabAndCategoryInPeriod(w http.ResponseWriter, r *http.Request) {
	labID := chi.URLParam(r, "lab_id")
	categoryID := chi.URLParam(r, "category_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if labID == "" || categoryID == "" || startDate == "" || endDate == "" {
		http.Error(w, "Lab ID, Category ID, start date, and end date are required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			i.id,
			i.name,
			ti.name AS type_name,
			ci.name AS category_name,
			tri.test_date_start,
			tri.test_date_finish,
			tri.result AS test_result,
			tl.name AS laboratory_name,
			ri.completion_date
		FROM item i
		JOIN type_item ti ON i.type_id = ti.id
		JOIN category_item ci ON ti.category_id = ci.id
		JOIN ready_item ri ON i.id = ri.item_id
		JOIN test_ready_item tri ON ri.id = tri.ready_item_id
		JOIN testing_laboratory tl ON tri.lab_id = tl.id
		WHERE tl.id = $1
		  AND ci.id = $2
		  AND tri.test_date_start >= $3
		  AND tri.test_date_finish <= $4
		ORDER BY i.name, tri.test_date_start
	`

	rows, err := h.db.Query(context.Background(), query, labID, categoryID, startDate, endDate)
	if err != nil {
		h.logger.Printf("Error querying database for tested items by lab, category, and period: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []TestedItem
	for rows.Next() {
		var item TestedItem
		var typeName, categoryName, testResult, labName sql.NullString
		var scanTestDateStart, scanTestDateFinish, scanCompletionDate pgtype.Date

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&typeName,
			&categoryName,
			&scanTestDateStart,
			&scanTestDateFinish,
			&testResult,
			&labName,
			&scanCompletionDate,
		); err != nil {
			h.logger.Printf("Error scanning row for tested items: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if typeName.Valid {
			item.TypeName = typeName.String
		}
		if categoryName.Valid {
			item.CategoryName = categoryName.String
		}
		if scanTestDateStart.Valid {
			item.TestDateStart = scanTestDateStart.Time.Format("2006-01-02")
		}
		if scanTestDateFinish.Valid {
			item.TestDateFinish = scanTestDateFinish.Time.Format("2006-01-02")
		}
		if testResult.Valid {
			item.TestResult = testResult.String
		}
		if labName.Valid {
			item.LaboratoryName = labName.String
		}
		if scanCompletionDate.Valid {
			item.CompletionDate = scanCompletionDate.Time.Format("2006-01-02")
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for tested items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json for tested items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getTestedItemsByLabInPeriod получает изделия, испытанные в лаборатории за период
func (h *Handler) getTestedItemsByLabInPeriod(w http.ResponseWriter, r *http.Request) {
	labID := chi.URLParam(r, "lab_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if labID == "" || startDate == "" || endDate == "" {
		http.Error(w, "Lab ID, start date, and end date are required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			i.id,
			i.name,
			ti.name AS type_name,
			ci.name AS category_name,
			tri.test_date_start,
			tri.test_date_finish,
			tri.result AS test_result,
			tl.name AS laboratory_name,
			ri.completion_date
		FROM item i
		JOIN type_item ti ON i.type_id = ti.id
		JOIN category_item ci ON ti.category_id = ci.id
		JOIN ready_item ri ON i.id = ri.item_id
		JOIN test_ready_item tri ON ri.id = tri.ready_item_id
		JOIN testing_laboratory tl ON tri.lab_id = tl.id
		WHERE tl.id = $1
		  AND tri.test_date_start >= $2
		  AND tri.test_date_finish <= $3
		ORDER BY i.name, tri.test_date_start
	`

	rows, err := h.db.Query(context.Background(), query, labID, startDate, endDate)
	if err != nil {
		h.logger.Printf("Error querying database for tested items by lab and period: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []TestedItem
	for rows.Next() {
		var item TestedItem
		var typeName, categoryName, testResult, labName sql.NullString
		var scanTestDateStart, scanTestDateFinish, scanCompletionDate pgtype.Date

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&typeName,
			&categoryName,
			&scanTestDateStart,
			&scanTestDateFinish,
			&testResult,
			&labName,
			&scanCompletionDate,
		); err != nil {
			h.logger.Printf("Error scanning row for tested items: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if typeName.Valid {
			item.TypeName = typeName.String
		}
		if categoryName.Valid {
			item.CategoryName = categoryName.String
		}
		if scanTestDateStart.Valid {
			item.TestDateStart = scanTestDateStart.Time.Format("2006-01-02")
		}
		if scanTestDateFinish.Valid {
			item.TestDateFinish = scanTestDateFinish.Time.Format("2006-01-02")
		}
		if testResult.Valid {
			item.TestResult = testResult.String
		}
		if labName.Valid {
			item.LaboratoryName = labName.String
		}
		if scanCompletionDate.Valid {
			item.CompletionDate = scanCompletionDate.Time.Format("2006-01-02")
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for tested items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json for tested items: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
