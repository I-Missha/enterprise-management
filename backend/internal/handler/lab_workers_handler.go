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

// LabWorkerInfo представляет информацию об испытателе
type LabWorkerInfo struct {
	EmployeeID      int    `json:"employee_id"`
	Name            string `json:"name"`
	HireDate        string `json:"hire_date,omitempty"`
	CurrentPosition string `json:"current_position,omitempty"`
	LabID           int    `json:"lab_id"`
	LabName         string `json:"lab_name,omitempty"`
	TestDate        string `json:"test_date,omitempty"` // Дата участия в конкретном тесте
}

// getLabWorkersByItemOrCategoryInLabInPeriod получает список испытателей
func (h *Handler) getLabWorkersByItemOrCategoryInLabInPeriod(w http.ResponseWriter, r *http.Request) {
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
			e.id AS employee_id,
			e.name,
			e.hire_date,
			e.current_position,
			lw.lab_id,
			tl.name AS lab_name,
			lwtri.test_date
		FROM employee e
		JOIN lab_worker lw ON e.id = lw.employee_id
		JOIN testing_laboratory tl ON lw.lab_id = tl.id
		JOIN lab_worker_test_ready_item lwtri ON lw.employee_id = lwtri.lab_worker_id
		JOIN ready_item ri ON lwtri.ready_item_id = ri.id
		JOIN item i ON ri.item_id = i.id
		JOIN type_item ti ON i.type_id = ti.id
		WHERE lw.lab_id = $1
		  AND lwtri.test_date >= $2
		  AND lwtri.test_date <= $3
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

	queryBuilder.WriteString(" ORDER BY e.name, lwtri.test_date")

	rows, err := h.db.Query(context.Background(), queryBuilder.String(), queryArgs...)
	if err != nil {
		h.logger.Printf("Error querying database for lab workers: %v, query: %s, args: %v", err, queryBuilder.String(), queryArgs)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workers []LabWorkerInfo
	for rows.Next() {
		var worker LabWorkerInfo
		var hireDate, testDate pgtype.Date
		var currentPosition, labName sql.NullString

		if err := rows.Scan(
			&worker.EmployeeID,
			&worker.Name,
			&hireDate,
			&currentPosition,
			&worker.LabID,
			&labName,
			&testDate,
		); err != nil {
			h.logger.Printf("Error scanning row for lab worker: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if hireDate.Valid {
			worker.HireDate = hireDate.Time.Format("2006-01-02")
		}
		if currentPosition.Valid {
			worker.CurrentPosition = currentPosition.String
		}
		if labName.Valid {
			worker.LabName = labName.String
		}
		if testDate.Valid {
			worker.TestDate = testDate.Time.Format("2006-01-02")
		}
		workers = append(workers, worker)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for lab workers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workers); err != nil {
		h.logger.Printf("Error encoding json for lab workers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getLabWorkersInLabInPeriod получает список всех испытателей в лаборатории за период
func (h *Handler) getLabWorkersInLabInPeriod(w http.ResponseWriter, r *http.Request) {
	labID := chi.URLParam(r, "lab_id")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if labID == "" || startDate == "" || endDate == "" {
		http.Error(w, "Lab ID, start date, and end date are required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT DISTINCT
			e.id AS employee_id,
			e.name,
			e.hire_date,
			e.current_position,
			lw.lab_id,
			tl.name AS lab_name,
			lwtri.test_date
		FROM employee e
		JOIN lab_worker lw ON e.id = lw.employee_id
		JOIN testing_laboratory tl ON lw.lab_id = tl.id
		JOIN lab_worker_test_ready_item lwtri ON lw.employee_id = lwtri.lab_worker_id
		WHERE lw.lab_id = $1
		  AND lwtri.test_date >= $2
		  AND lwtri.test_date <= $3
		ORDER BY e.name, lwtri.test_date
	`

	rows, err := h.db.Query(context.Background(), query, labID, startDate, endDate)
	if err != nil {
		h.logger.Printf("Error querying database for lab workers in period: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workers []LabWorkerInfo
	for rows.Next() {
		var worker LabWorkerInfo
		var hireDate, testDate pgtype.Date
		var currentPosition, labName sql.NullString

		if err := rows.Scan(
			&worker.EmployeeID,
			&worker.Name,
			&hireDate,
			&currentPosition,
			&worker.LabID,
			&labName,
			&testDate,
		); err != nil {
			h.logger.Printf("Error scanning row for lab worker in period: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if hireDate.Valid {
			worker.HireDate = hireDate.Time.Format("2006-01-02")
		}
		if currentPosition.Valid {
			worker.CurrentPosition = currentPosition.String
		}
		if labName.Valid {
			worker.LabName = labName.String
		}
		if testDate.Valid {
			worker.TestDate = testDate.Time.Format("2006-01-02")
		}
		workers = append(workers, worker)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for lab workers in period: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workers); err != nil {
		h.logger.Printf("Error encoding json for lab workers in period: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
