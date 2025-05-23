package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Модели для представления данных о персонале
type Worker struct {
	EmployeeID      int    `json:"employee_id"`
	Name            string `json:"name"`
	HallID          int    `json:"hall_id"`
	AreaID          int    `json:"area_id"`
	WorkTeamID      int    `json:"work_team_id"`
	Category        string `json:"category"`
	HireDate        string `json:"hire_date,omitempty"`
	CurrentPosition string `json:"current_position,omitempty"`
}

type Engineer struct {
	EmployeeID      int    `json:"employee_id"`
	Name            string `json:"name"`
	HallID          int    `json:"hall_id"`
	AreaID          int    `json:"area_id"`
	CategoryID      int    `json:"category_id"`
	CategoryName    string `json:"category_name,omitempty"`
	HireDate        string `json:"hire_date,omitempty"`
	CurrentPosition string `json:"current_position,omitempty"`
}

type StaffMember struct {
	EmployeeID      int    `json:"employee_id"`
	Name            string `json:"name"`
	HallID          int    `json:"hall_id"`
	AreaID          int    `json:"area_id,omitempty"`
	Category        string `json:"category,omitempty"`
	CategoryID      int    `json:"category_id,omitempty"`
	CategoryName    string `json:"category_name,omitempty"`
	WorkTeamID      int    `json:"work_team_id,omitempty"`
	HireDate        string `json:"hire_date,omitempty"`
	CurrentPosition string `json:"current_position,omitempty"`
	Type            string `json:"type"` // "worker" или "engineer"
}

// getHallStaff получает данные о кадровом составе указанного цеха
func (h *Handler) getHallStaff(w http.ResponseWriter, r *http.Request) {
	hallID := chi.URLParam(r, "hall_id")
	if hallID == "" {
		http.Error(w, "Hall ID is required", http.StatusBadRequest)
		return
	}

	// Объединенный запрос, получающий и рабочих, и инженеров цеха
	query := `
		SELECT 
			'worker' as type, 
			w.employee_id, 
			e.name, 
			w.hall_id, 
			w.area_id, 
			w.work_team_id,
			w.category,
			NULL as category_id,
			NULL as category_name,
			e.hire_date,
			e.current_position
		FROM worker w
		JOIN employee e ON w.employee_id = e.id
		WHERE w.hall_id = $1
		UNION ALL
		SELECT 
			'engineer' as type, 
			eng.employee_id, 
			e.name, 
			eng.hall_id, 
			eng.area_id,
			NULL as work_team_id,
			NULL as category,
			eng.category_id,
			ce.name as category_name,
			e.hire_date,
			e.current_position
		FROM engineer eng
		JOIN employee e ON eng.employee_id = e.id
		LEFT JOIN category_engineer ce ON eng.category_id = ce.id
		WHERE eng.hall_id = $1
	`

	rows, err := h.db.Query(context.Background(), query, hallID)
	if err != nil {
		h.logger.Printf("Error querying database for hall staff: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var staff []StaffMember
	for rows.Next() {
		var member StaffMember
		var category, categoryName, hireDate, currentPosition sql.NullString
		var workTeamID, categoryID sql.NullInt32

		if err := rows.Scan(
			&member.Type,
			&member.EmployeeID,
			&member.Name,
			&member.HallID,
			&member.AreaID,
			&workTeamID,
			&category,
			&categoryID,
			&categoryName,
			&hireDate,
			&currentPosition,
		); err != nil {
			h.logger.Printf("Error scanning row for hall staff: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if workTeamID.Valid {
			member.WorkTeamID = int(workTeamID.Int32)
		}

		if category.Valid {
			member.Category = category.String
		}

		if categoryID.Valid {
			member.CategoryID = int(categoryID.Int32)
		}

		if categoryName.Valid {
			member.CategoryName = categoryName.String
		}

		if hireDate.Valid {
			member.HireDate = hireDate.String
		}

		if currentPosition.Valid {
			member.CurrentPosition = currentPosition.String
		}

		staff = append(staff, member)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for hall staff: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(staff); err != nil {
		h.logger.Printf("Error encoding json for hall staff: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllStaff получает данные о всем кадровом составе предприятия
func (h *Handler) getAllStaff(w http.ResponseWriter, r *http.Request) {
	// Объединенный запрос, получающий всех рабочих и инженеров
	query := `
		SELECT 
			'worker' as type, 
			w.employee_id, 
			e.name, 
			w.hall_id, 
			w.area_id, 
			w.work_team_id,
			w.category,
			NULL as category_id,
			NULL as category_name,
			e.hire_date,
			e.current_position
		FROM worker w
		JOIN employee e ON w.employee_id = e.id
		UNION ALL
		SELECT 
			'engineer' as type, 
			eng.employee_id, 
			e.name, 
			eng.hall_id, 
			eng.area_id,
			NULL as work_team_id,
			NULL as category,
			eng.category_id,
			ce.name as category_name,
			e.hire_date,
			e.current_position
		FROM engineer eng
		JOIN employee e ON eng.employee_id = e.id
		LEFT JOIN category_engineer ce ON eng.category_id = ce.id
	`

	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		h.logger.Printf("Error querying database for all staff: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var staff []StaffMember
	for rows.Next() {
		var member StaffMember
		var category, categoryName, hireDate, currentPosition sql.NullString
		var workTeamID, categoryID sql.NullInt32

		if err := rows.Scan(
			&member.Type,
			&member.EmployeeID,
			&member.Name,
			&member.HallID,
			&member.AreaID,
			&workTeamID,
			&category,
			&categoryID,
			&categoryName,
			&hireDate,
			&currentPosition,
		); err != nil {
			h.logger.Printf("Error scanning row for all staff: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if workTeamID.Valid {
			member.WorkTeamID = int(workTeamID.Int32)
		}

		if category.Valid {
			member.Category = category.String
		}

		if categoryID.Valid {
			member.CategoryID = int(categoryID.Int32)
		}

		if categoryName.Valid {
			member.CategoryName = categoryName.String
		}

		if hireDate.Valid {
			member.HireDate = hireDate.String
		}

		if currentPosition.Valid {
			member.CurrentPosition = currentPosition.String
		}

		staff = append(staff, member)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for all staff: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(staff); err != nil {
		h.logger.Printf("Error encoding json for all staff: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getStaffByCategory получает данные о кадровом составе по указанной категории
func (h *Handler) getStaffByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	staffType := chi.URLParam(r, "type") // "worker" или "engineer"

	if categoryID == "" || staffType == "" {
		http.Error(w, "Category ID and staff type are required", http.StatusBadRequest)
		return
	}

	var query string
	var args []interface{}

	switch staffType {
	case "worker":
		// Для рабочих используем category
		query = `
			SELECT 
				'worker' as type, 
				w.employee_id, 
				e.name, 
				w.hall_id, 
				w.area_id, 
				w.work_team_id,
				w.category,
				NULL as category_id,
				NULL as category_name,
				e.hire_date,
				e.current_position
			FROM worker w
			JOIN employee e ON w.employee_id = e.id
			WHERE w.category = $1
		`
		args = append(args, categoryID)
	case "engineer":
		// Для инженеров используем category_id
		query = `
			SELECT 
				'engineer' as type, 
				eng.employee_id, 
				e.name, 
				eng.hall_id, 
				eng.area_id,
				NULL as work_team_id,
				NULL as category,
				eng.category_id,
				ce.name as category_name,
				e.hire_date,
				e.current_position
			FROM engineer eng
			JOIN employee e ON eng.employee_id = e.id
			LEFT JOIN category_engineer ce ON eng.category_id = ce.id
			WHERE eng.category_id = $1
		`
		args = append(args, categoryID)
	default:
		http.Error(w, "Invalid staff type. Use 'worker' or 'engineer'", http.StatusBadRequest)
		return
	}

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		h.logger.Printf("Error querying database for staff by category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var staff []StaffMember
	for rows.Next() {
		var member StaffMember
		var category, categoryName, hireDate, currentPosition sql.NullString
		var workTeamID, categoryID sql.NullInt32

		if err := rows.Scan(
			&member.Type,
			&member.EmployeeID,
			&member.Name,
			&member.HallID,
			&member.AreaID,
			&workTeamID,
			&category,
			&categoryID,
			&categoryName,
			&hireDate,
			&currentPosition,
		); err != nil {
			h.logger.Printf("Error scanning row for staff by category: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if workTeamID.Valid {
			member.WorkTeamID = int(workTeamID.Int32)
		}

		if category.Valid {
			member.Category = category.String
		}

		if categoryID.Valid {
			member.CategoryID = int(categoryID.Int32)
		}

		if categoryName.Valid {
			member.CategoryName = categoryName.String
		}

		if hireDate.Valid {
			member.HireDate = hireDate.String
		}

		if currentPosition.Valid {
			member.CurrentPosition = currentPosition.String
		}

		staff = append(staff, member)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for staff by category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(staff); err != nil {
		h.logger.Printf("Error encoding json for staff by category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
