package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// ProductionHallRequest представляет данные для создания нового цеха
type ProductionHallRequest struct {
	ID   int    `json:"id"` // ID will be populated by the database
	Name string `json:"name"`
}

// createProductionHall обрабатывает запрос на создание нового цеха
func (h *Handler) createProductionHall(w http.ResponseWriter, r *http.Request) {
	var hallReq ProductionHallRequest

	if err := json.NewDecoder(r.Body).Decode(&hallReq); err != nil {
		h.logger.Printf("Error decoding request body for new hall: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(hallReq.Name) == "" {
		http.Error(w, "Name is required and cannot be empty for hall", http.StatusBadRequest)
		return
	}

	// Проверяем, не существует ли уже цех с таким именем
	var nameExists bool
	err := h.db.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM production_halls WHERE name = $1)",
		hallReq.Name).Scan(&nameExists)
	if err != nil {
		h.logger.Printf("Error checking if hall name exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if nameExists {
		http.Error(w, "Hall with name '"+hallReq.Name+"' already exists", http.StatusConflict)
		return
	}

	// Вставляем новый цех в базу данных и получаем сгенерированный ID
	err = h.db.QueryRow(context.Background(),
		"INSERT INTO production_halls (name) VALUES ($1) RETURNING id",
		hallReq.Name).Scan(&hallReq.ID)
	if err != nil {
		h.logger.Printf("Error inserting new hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(hallReq); err != nil {
		h.logger.Printf("Error encoding response for new hall: %v", err)
	}
}

// ProductionAreaRequest представляет данные для создания нового участка
type ProductionAreaRequest struct {
	ID     int    `json:"id"` // ID will be populated by the database
	Name   string `json:"name"`
	HallID int    `json:"hall_id"`
}

// createProductionArea обрабатывает запрос на создание нового производственного участка
func (h *Handler) createProductionArea(w http.ResponseWriter, r *http.Request) {
	var areaReq ProductionAreaRequest
	if err := json.NewDecoder(r.Body).Decode(&areaReq); err != nil {
		h.logger.Printf("Error decoding request body for new area: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(areaReq.Name) == "" {
		http.Error(w, "Name is required for area", http.StatusBadRequest)
		return
	}
	if areaReq.HallID <= 0 {
		http.Error(w, "Valid HallID is required for area", http.StatusBadRequest)
		return
	}

	var hallExists bool
	err := h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM production_halls WHERE id = $1)", areaReq.HallID).Scan(&hallExists)
	if err != nil {
		h.logger.Printf("Error checking hall existence for new area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !hallExists {
		http.Error(w, "Hall with ID "+strconv.Itoa(areaReq.HallID)+" not found", http.StatusNotFound)
		return
	}

	var nameExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM production_area WHERE name = $1 AND hall_id = $2)", areaReq.Name, areaReq.HallID).Scan(&nameExists)
	if err != nil {
		h.logger.Printf("Error checking area name existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if nameExists {
		http.Error(w, "Area with name '"+areaReq.Name+"' already exists in this hall", http.StatusConflict)
		return
	}

	err = h.db.QueryRow(context.Background(), "INSERT INTO production_area (name, hall_id) VALUES ($1, $2) RETURNING id", areaReq.Name, areaReq.HallID).Scan(&areaReq.ID)
	if err != nil {
		h.logger.Printf("Error inserting new area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(areaReq)
}

// CategoryItemRequest представляет данные для создания категории изделия
type CategoryItemRequest struct {
	ID        int    `json:"id"` // ID will be populated by the database
	Name      string `json:"name"`
	Attribute string `json:"attribute,omitempty"`
}

// createCategoryItem обрабатывает запрос на создание новой категории изделия
func (h *Handler) createCategoryItem(w http.ResponseWriter, r *http.Request) {
	var catReq CategoryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&catReq); err != nil {
		h.logger.Printf("Error decoding request body for new item category: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(catReq.Name) == "" {
		http.Error(w, "Name is required for item category", http.StatusBadRequest)
		return
	}

	var nameExists bool
	err := h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM category_item WHERE name = $1)", catReq.Name).Scan(&nameExists)
	if err != nil {
		h.logger.Printf("Error checking item category name existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if nameExists {
		http.Error(w, "Item category with name '"+catReq.Name+"' already exists", http.StatusConflict)
		return
	}
	attribute := sql.NullString{String: catReq.Attribute, Valid: catReq.Attribute != ""}

	err = h.db.QueryRow(context.Background(), "INSERT INTO category_item (name, attribute) VALUES ($1, $2) RETURNING id", catReq.Name, attribute).Scan(&catReq.ID)
	if err != nil {
		h.logger.Printf("Error inserting new item category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(catReq)
}

// TypeItemRequest представляет данные для создания типа изделия
type TypeItemRequest struct {
	ID         int    `json:"id"` // ID will be populated by the database
	Name       string `json:"name"`
	CategoryID int    `json:"category_id"`
}

// createTypeItem обрабатывает запрос на создание нового типа изделия
func (h *Handler) createTypeItem(w http.ResponseWriter, r *http.Request) {
	var typeReq TypeItemRequest
	if err := json.NewDecoder(r.Body).Decode(&typeReq); err != nil {
		h.logger.Printf("Error decoding request body for new item type: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(typeReq.Name) == "" {
		http.Error(w, "Name is required for item type", http.StatusBadRequest)
		return
	}
	if typeReq.CategoryID <= 0 {
		http.Error(w, "Valid CategoryID is required for item type", http.StatusBadRequest)
		return
	}

	var categoryExists bool
	err := h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM category_item WHERE id = $1)", typeReq.CategoryID).Scan(&categoryExists)
	if err != nil {
		h.logger.Printf("Error checking category existence for new item type: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !categoryExists {
		http.Error(w, "Category with ID "+strconv.Itoa(typeReq.CategoryID)+" not found", http.StatusNotFound)
		return
	}

	var nameExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM type_item WHERE name = $1)", typeReq.Name).Scan(&nameExists)
	if err != nil {
		h.logger.Printf("Error checking item type name existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if nameExists {
		http.Error(w, "Item type with name '"+typeReq.Name+"' already exists", http.StatusConflict)
		return
	}

	err = h.db.QueryRow(context.Background(), "INSERT INTO type_item (name, category_id) VALUES ($1, $2) RETURNING id", typeReq.Name, typeReq.CategoryID).Scan(&typeReq.ID)
	if err != nil {
		h.logger.Printf("Error inserting new item type: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(typeReq)
}

// ItemRequest представляет данные для создания изделия
type ItemRequest struct {
	ID      int    `json:"id"` // ID will be populated by the database
	Name    string `json:"name"`
	TypeID  int    `json:"type_id"`
	HallID  int    `json:"hall_id"`
	Status  string `json:"status"` // 'in_progress', 'testing', 'completed'
	AreaIDs []int  `json:"area_ids,omitempty"`
}

// createItem обрабатывает запрос на создание нового изделия
func (h *Handler) createItem(w http.ResponseWriter, r *http.Request) {
	var itemReq ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&itemReq); err != nil {
		h.logger.Printf("Error decoding request body for new item: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(itemReq.Name) == "" {
		http.Error(w, "Name is required for item", http.StatusBadRequest)
		return
	}
	if itemReq.TypeID <= 0 {
		http.Error(w, "Valid TypeID is required for item", http.StatusBadRequest)
		return
	}
	if itemReq.HallID <= 0 {
		http.Error(w, "Valid HallID is required for item", http.StatusBadRequest)
		return
	}
	validStatuses := map[string]bool{"in_progress": true, "testing": true, "completed": true}
	if _, ok := validStatuses[itemReq.Status]; !ok {
		http.Error(w, "Invalid status. Must be one of 'in_progress', 'testing', 'completed'", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin(context.Background())
	if err != nil {
		h.logger.Printf("Error beginning transaction for new item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	var typeExists bool
	err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM type_item WHERE id = $1)", itemReq.TypeID).Scan(&typeExists)
	if err != nil {
		h.logger.Printf("Error checking type existence for new item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !typeExists {
		http.Error(w, "Type with ID "+strconv.Itoa(itemReq.TypeID)+" not found", http.StatusNotFound)
		return
	}

	var hallExists bool
	err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM production_halls WHERE id = $1)", itemReq.HallID).Scan(&hallExists)
	if err != nil {
		h.logger.Printf("Error checking hall existence for new item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !hallExists {
		http.Error(w, "Hall with ID "+strconv.Itoa(itemReq.HallID)+" not found", http.StatusNotFound)
		return
	}

	err = tx.QueryRow(context.Background(), "INSERT INTO item (name, type_id, hall_id, status) VALUES ($1, $2, $3, $4) RETURNING id",
		itemReq.Name, itemReq.TypeID, itemReq.HallID, itemReq.Status).Scan(&itemReq.ID)
	if err != nil {
		h.logger.Printf("Error inserting new item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(itemReq.AreaIDs) > 0 {
		for _, areaID := range itemReq.AreaIDs {
			if areaID <= 0 {
				http.Error(w, "Invalid AreaID "+strconv.Itoa(areaID)+" in area_ids list", http.StatusBadRequest)
				return
			}
			var areaHallID int
			err := tx.QueryRow(context.Background(), "SELECT hall_id FROM production_area WHERE id = $1", areaID).Scan(&areaHallID)
			if err == sql.ErrNoRows {
				http.Error(w, "Area with ID "+strconv.Itoa(areaID)+" not found", http.StatusNotFound)
				return
			} else if err != nil {
				h.logger.Printf("Error checking area existence for item association: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if areaHallID != itemReq.HallID {
				http.Error(w, "Area with ID "+strconv.Itoa(areaID)+" belongs to hall "+strconv.Itoa(areaHallID)+", not hall "+strconv.Itoa(itemReq.HallID), http.StatusBadRequest)
				return
			}

			_, err = tx.Exec(context.Background(), "INSERT INTO areas_items (item_id, area_id) VALUES ($1, $2)", itemReq.ID, areaID)
			if err != nil {
				if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
					http.Error(w, "Item with ID "+strconv.Itoa(itemReq.ID)+" is already associated with area ID "+strconv.Itoa(areaID), http.StatusConflict)
					return
				}
				h.logger.Printf("Error inserting into areas_items for item %d, area %d: %v", itemReq.ID, areaID, err)
				http.Error(w, "Internal server error during area association", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		h.logger.Printf("Error committing transaction for new item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(itemReq)
}

// EmployeeRequest представляет данные для создания сотрудника
type EmployeeRequest struct {
	ID              int    `json:"id"` // ID will be populated by the database
	Name            string `json:"name"`
	HireDate        string `json:"hire_date"` // Expected format YYYY-MM-DD
	CurrentPosition string `json:"current_position,omitempty"`
}

// createEmployee обрабатывает запрос на создание нового сотрудника
func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	var empReq EmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&empReq); err != nil {
		h.logger.Printf("Error decoding request body for new employee: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(empReq.Name) == "" {
		http.Error(w, "Name is required for employee", http.StatusBadRequest)
		return
	}
	parsedHireDate, err := time.Parse("2006-01-02", empReq.HireDate)
	if err != nil {
		http.Error(w, "Invalid hire_date format. Expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	pgDate := pgtype.Date{Time: parsedHireDate, Valid: true}

	err = h.db.QueryRow(context.Background(), "INSERT INTO employee (name, hire_date, current_position) VALUES ($1, $2, $3) RETURNING id",
		empReq.Name, pgDate, currentPosition).Scan(&empReq.ID)
	if err != nil {
		h.logger.Printf("Error inserting new employee: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(empReq)
}

// TestingLaboratoryRequest представляет данные для создания лаборатории
type TestingLaboratoryRequest struct {
	ID      int    `json:"id"` // ID will be populated by the database
	Name    string `json:"name"`
	HallIDs []int  `json:"hall_ids,omitempty"`
}

// createTestingLaboratory обрабатывает запрос на создание новой лаборатории
func (h *Handler) createTestingLaboratory(w http.ResponseWriter, r *http.Request) {
	var labReq TestingLaboratoryRequest
	if err := json.NewDecoder(r.Body).Decode(&labReq); err != nil {
		h.logger.Printf("Error decoding request body for new lab: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(labReq.Name) == "" {
		http.Error(w, "Name is required for lab", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin(context.Background())
	if err != nil {
		h.logger.Printf("Error beginning transaction for new lab: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	var nameExists bool
	err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM testing_laboratory WHERE name = $1)", labReq.Name).Scan(&nameExists)
	if err != nil {
		h.logger.Printf("Error checking lab name existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if nameExists {
		http.Error(w, "Lab with name '"+labReq.Name+"' already exists", http.StatusConflict)
		return
	}

	err = tx.QueryRow(context.Background(), "INSERT INTO testing_laboratory (name) VALUES ($1) RETURNING id", labReq.Name).Scan(&labReq.ID)
	if err != nil {
		h.logger.Printf("Error inserting new lab: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(labReq.HallIDs) > 0 {
		for _, hallID := range labReq.HallIDs {
			if hallID <= 0 {
				http.Error(w, "Invalid HallID "+strconv.Itoa(hallID)+" in hall_ids list", http.StatusBadRequest)
				return
			}
			var hallExists bool
			err := tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM production_halls WHERE id = $1)", hallID).Scan(&hallExists)
			if err != nil {
				h.logger.Printf("Error checking hall existence for lab association: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if !hallExists {
				http.Error(w, "Hall with ID "+strconv.Itoa(hallID)+" not found", http.StatusNotFound)
				return
			}

			_, err = tx.Exec(context.Background(), "INSERT INTO lab_hall (lab_id, hall_id) VALUES ($1, $2)", labReq.ID, hallID)
			if err != nil {
				if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
					http.Error(w, "Lab with ID "+strconv.Itoa(labReq.ID)+" is already associated with hall ID "+strconv.Itoa(hallID), http.StatusConflict)
					return
				}
				h.logger.Printf("Error inserting into lab_hall for lab %d, hall %d: %v", labReq.ID, hallID, err)
				http.Error(w, "Internal server error during hall association", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		h.logger.Printf("Error committing transaction for new lab: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(labReq)
}

// WorkTeamRequest представляет данные для создания рабочей бригады
type WorkTeamRequest struct {
	ID     int    `json:"id"` // ID will be populated by the database
	Name   string `json:"name"`
	AreaID int    `json:"area_id"`
	HallID int    `json:"hall_id"`
}

// createWorkTeam обрабатывает запрос на создание новой рабочей бригады
func (h *Handler) createWorkTeam(w http.ResponseWriter, r *http.Request) {
	var teamReq WorkTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&teamReq); err != nil {
		h.logger.Printf("Error decoding request body for new work team: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(teamReq.Name) == "" {
		http.Error(w, "Name is required for work team", http.StatusBadRequest)
		return
	}
	if teamReq.AreaID <= 0 {
		http.Error(w, "Valid AreaID is required for work team", http.StatusBadRequest)
		return
	}
	if teamReq.HallID <= 0 {
		http.Error(w, "Valid HallID is required for work team", http.StatusBadRequest)
		return
	}

	var areaHallID int
	err := h.db.QueryRow(context.Background(), "SELECT hall_id FROM production_area WHERE id = $1", teamReq.AreaID).Scan(&areaHallID)
	if err == sql.ErrNoRows {
		http.Error(w, "Area with ID "+strconv.Itoa(teamReq.AreaID)+" not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Printf("Error checking area existence for new work team: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if areaHallID != teamReq.HallID {
		http.Error(w, "HallID "+strconv.Itoa(teamReq.HallID)+" does not match area's hall ID "+strconv.Itoa(areaHallID), http.StatusBadRequest)
		return
	}

	var nameExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM work_team WHERE name = $1 AND area_id = $2)", teamReq.Name, teamReq.AreaID).Scan(&nameExists)
	if err != nil {
		h.logger.Printf("Error checking work team name existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if nameExists {
		http.Error(w, "Work team with name '"+teamReq.Name+"' already exists in this area", http.StatusConflict)
		return
	}

	err = h.db.QueryRow(context.Background(), "INSERT INTO work_team (name, area_id, hall_id) VALUES ($1, $2, $3) RETURNING id",
		teamReq.Name, teamReq.AreaID, teamReq.HallID).Scan(&teamReq.ID)
	if err != nil {
		h.logger.Printf("Error inserting new work team: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(teamReq)
}

// AreaBossRequest представляет данные для назначения начальника участка
type AreaBossRequest struct {
	AreaID     int `json:"area_id"`
	EngineerID int `json:"engineer_id"`
}

// createAreaBoss обрабатывает запрос на назначение начальника участка
func (h *Handler) createAreaBoss(w http.ResponseWriter, r *http.Request) {
	var req AreaBossRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding request body for new area boss: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AreaID <= 0 {
		http.Error(w, "Valid AreaID is required", http.StatusBadRequest)
		return
	}
	if req.EngineerID <= 0 {
		http.Error(w, "Valid EngineerID is required", http.StatusBadRequest)
		return
	}

	// Проверяем существование участка
	var areaExists bool
	err := h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM production_area WHERE id = $1)", req.AreaID).Scan(&areaExists)
	if err != nil {
		h.logger.Printf("Error checking area existence for new area boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !areaExists {
		http.Error(w, "Area with ID "+strconv.Itoa(req.AreaID)+" not found", http.StatusNotFound)
		return
	}

	// Проверяем существование инженера
	var engineerExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM engineer WHERE employee_id = $1)", req.EngineerID).Scan(&engineerExists)
	if err != nil {
		h.logger.Printf("Error checking engineer existence for new area boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !engineerExists {
		http.Error(w, "Engineer with ID "+strconv.Itoa(req.EngineerID)+" not found", http.StatusNotFound)
		return
	}

	// Проверяем, не назначен ли уже начальник этому участку
	var bossExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM area_boss WHERE area_id = $1)", req.AreaID).Scan(&bossExists)
	if err != nil {
		h.logger.Printf("Error checking if area boss already exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if bossExists {
		http.Error(w, "Area with ID "+strconv.Itoa(req.AreaID)+" already has a boss", http.StatusConflict)
		return
	}

	_, err = h.db.Exec(context.Background(), "INSERT INTO area_boss (area_id, engineer_id) VALUES ($1, $2)", req.AreaID, req.EngineerID)
	if err != nil {
		h.logger.Printf("Error inserting new area boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// HallBossRequest представляет данные для назначения начальника цеха
type HallBossRequest struct {
	HallID     int `json:"hall_id"`
	EngineerID int `json:"engineer_id"`
}

// createHallBoss обрабатывает запрос на назначение начальника цеха
func (h *Handler) createHallBoss(w http.ResponseWriter, r *http.Request) {
	var req HallBossRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding request body for new hall boss: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.HallID <= 0 {
		http.Error(w, "Valid HallID is required", http.StatusBadRequest)
		return
	}
	if req.EngineerID <= 0 {
		http.Error(w, "Valid EngineerID is required", http.StatusBadRequest)
		return
	}

	// Проверяем существование цеха
	var hallExists bool
	err := h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM production_halls WHERE id = $1)", req.HallID).Scan(&hallExists)
	if err != nil {
		h.logger.Printf("Error checking hall existence for new hall boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !hallExists {
		http.Error(w, "Hall with ID "+strconv.Itoa(req.HallID)+" not found", http.StatusNotFound)
		return
	}

	// Проверяем существование инженера
	var engineerExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM engineer WHERE employee_id = $1)", req.EngineerID).Scan(&engineerExists)
	if err != nil {
		h.logger.Printf("Error checking engineer existence for new hall boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !engineerExists {
		http.Error(w, "Engineer with ID "+strconv.Itoa(req.EngineerID)+" not found", http.StatusNotFound)
		return
	}

	// Проверяем, не назначен ли уже начальник этому цеху
	var bossExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM hall_bosses WHERE hall_id = $1)", req.HallID).Scan(&bossExists)
	if err != nil {
		h.logger.Printf("Error checking if hall boss already exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if bossExists {
		http.Error(w, "Hall with ID "+strconv.Itoa(req.HallID)+" already has a boss", http.StatusConflict)
		return
	}

	_, err = h.db.Exec(context.Background(), "INSERT INTO hall_bosses (hall_id, engineer_id) VALUES ($1, $2)", req.HallID, req.EngineerID)
	if err != nil {
		h.logger.Printf("Error inserting new hall boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// MasterCreateRequest представляет данные для создания мастера
type MasterCreateRequest struct {
	ID         int `json:"id"` // ID will be populated by the database
	AreaID     int `json:"area_id"`
	EngineerID int `json:"engineer_id"`
}

// createMaster обрабатывает запрос на создание нового мастера
func (h *Handler) createMaster(w http.ResponseWriter, r *http.Request) {
	var req MasterCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding request body for new master: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AreaID <= 0 {
		http.Error(w, "Valid AreaID is required", http.StatusBadRequest)
		return
	}
	if req.EngineerID <= 0 {
		http.Error(w, "Valid EngineerID is required", http.StatusBadRequest)
		return
	}

	var areaExists bool
	err := h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM production_area WHERE id = $1)", req.AreaID).Scan(&areaExists)
	if err != nil {
		h.logger.Printf("Error checking area existence for new master: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !areaExists {
		http.Error(w, "Area with ID "+strconv.Itoa(req.AreaID)+" not found", http.StatusNotFound)
		return
	}

	var engineerExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM engineer WHERE employee_id = $1)", req.EngineerID).Scan(&engineerExists)
	if err != nil {
		h.logger.Printf("Error checking engineer existence for new master: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !engineerExists {
		http.Error(w, "Engineer with ID "+strconv.Itoa(req.EngineerID)+" not found", http.StatusNotFound)
		return
	}

	var uniqueCombinationExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM masters WHERE area_id = $1 AND engineer_id = $2)", req.AreaID, req.EngineerID).Scan(&uniqueCombinationExists)
	if err != nil {
		h.logger.Printf("Error checking master unique combination: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if uniqueCombinationExists {
		http.Error(w, "Master with AreaID "+strconv.Itoa(req.AreaID)+" and EngineerID "+strconv.Itoa(req.EngineerID)+" already exists", http.StatusConflict)
		return
	}

	err = h.db.QueryRow(context.Background(), "INSERT INTO masters (area_id, engineer_id) VALUES ($1, $2) RETURNING id", req.AreaID, req.EngineerID).Scan(&req.ID)
	if err != nil {
		h.logger.Printf("Error inserting new master: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// WorkerBossRequest представляет данные для назначения бригадира
type WorkerBossRequest struct {
	WorkerID int `json:"worker_id"`
}

// createWorkerBoss обрабатывает запрос на назначение бригадира
func (h *Handler) createWorkerBoss(w http.ResponseWriter, r *http.Request) {
	var req WorkerBossRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error decoding request body for new worker boss: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.WorkerID <= 0 {
		http.Error(w, "Valid WorkerID is required", http.StatusBadRequest)
		return
	}

	// Проверяем существование рабочего
	var workerExists bool
	err := h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM worker WHERE employee_id = $1)", req.WorkerID).Scan(&workerExists)
	if err != nil {
		h.logger.Printf("Error checking worker existence for new worker boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !workerExists {
		http.Error(w, "Worker with ID "+strconv.Itoa(req.WorkerID)+" not found", http.StatusNotFound)
		return
	}

	// Проверяем, не назначен ли уже этот рабочий бригадиром
	var bossExists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM worker_boss WHERE worker_id = $1)", req.WorkerID).Scan(&bossExists)
	if err != nil {
		h.logger.Printf("Error checking if worker boss already exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if bossExists {
		http.Error(w, "Worker with ID "+strconv.Itoa(req.WorkerID)+" is already a boss", http.StatusConflict)
		return
	}

	_, err = h.db.Exec(context.Background(), "INSERT INTO worker_boss (worker_id) VALUES ($1)", req.WorkerID)
	if err != nil {
		h.logger.Printf("Error inserting new worker boss: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}
