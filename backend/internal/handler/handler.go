package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	//"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	logger *log.Logger
	db     *pgxpool.Pool
}

type ItemType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category_id int    `json:"category_id"`
}

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category_id int    `json:"category_id"`
	Type_id     int    `json:"type_id"`
	Hall_id     int    `json:"hall_id"`
	Status      string `json:"status"`
}

type ReadyItem struct {
	ItemId int    `json:"item_id"`
	Count  int    `json:"count"`
	Name   string `json:"name"`
}

func NewHandler(logger *log.Logger, db *pgxpool.Pool) *Handler {
	return &Handler{
		logger: logger,
		db:     db,
	}
}

func (h *Handler) InitRoutes(router chi.Router) {

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/items", func(r chi.Router) {
		r.Get("/ping", h.ping)

		r.Get("/types/category/{category_id}", h.getItemsTypeByCategory)
		r.Get("/types/hall/{hall_id}", h.getAllItemsTypeByHallId)
		r.Get("/types/hall/{hall_id}/category/{category_id}", h.getAllItemsTypeByCategoryAndHallId)
		r.Get("/category/{category_id}", h.getAllItemsByCategory)
		r.Get("/area/{area_id}", h.getAllItemsByAreaId)
		r.Get("/TimeInterval", h.getAllItemsByTimeInterval)
		r.Get("/{item_id}/works", h.getItemWorks)      // Новый маршрут для получения работ по изделию
		r.Get("/{item_id}/teams", h.getItemTeams)      // Новый маршрут для получения бригад, участвующих в сборке изделия
		r.Get("/{item_id}/labs", h.getItemTestingLabs) // Получение лабораторий по изделию
	})

	router.Route("/staff", func(r chi.Router) {
		r.Get("/ping", h.ping)

		// Новые маршруты для работы с персоналом
		r.Get("/hall/{hall_id}", h.getHallStaff)                      // Кадровый состав цеха
		r.Get("/all", h.getAllStaff)                                  // Кадровый состав всего предприятия
		r.Get("/{type}/category/{category_id}", h.getStaffByCategory) // Кадры по категориям
	})

	router.Route("/areas", func(r chi.Router) {
		r.Get("/ping", h.ping)

		// Новые маршруты для работы с участками
		r.Get("/hall/{hall_id}", h.getAreasByHall) // Участки указанного цеха
		r.Get("/all", h.getAllAreas)               // Все участки предприятия
		r.Get("/{area_id}/boss", h.getAreaBoss)    // Начальник конкретного участка
	})

	router.Route("/work-teams", func(r chi.Router) {
		r.Get("/ping", h.ping)

		// Новые маршруты для работы с бригадами
		r.Get("/area/{area_id}", h.getWorkTeamsByArea) // Бригады указанного участка
		r.Get("/hall/{hall_id}", h.getWorkTeamsByHall) // Бригады указанного цеха
	})

	router.Route("/masters", func(r chi.Router) {
		r.Get("/ping", h.ping)

		// Новые маршруты для работы с мастерами
		r.Get("/area/{area_id}", h.getMastersByArea) // Мастера указанного участка
		r.Get("/hall/{hall_id}", h.getMastersByHall) // Мастера указанного цеха
	})

	router.Route("/current-items", func(r chi.Router) {
		r.Get("/ping", h.ping)

		// Маршруты для получения текущих изделий по категории
		r.Get("/category/{category_id}/area/{area_id}", h.getCurrentItemsByCategoryAndArea) // По категории и участку
		r.Get("/category/{category_id}/hall/{hall_id}", h.getCurrentItemsByCategoryAndHall) // По категории и цеху
		r.Get("/category/{category_id}", h.getCurrentItemsByCategory)                       // По категории на всем предприятии

		// Маршруты для получения всех текущих изделий (без фильтрации по категории)
		r.Get("/area/{area_id}", h.getCurrentItemsByArea) // По участку
		r.Get("/hall/{hall_id}", h.getCurrentItemsByHall) // По цеху
		r.Get("/all", h.getAllCurrentItems)               // По всему предприятию
	})

	router.Route("/tested-items", func(r chi.Router) {
		r.Get("/ping", h.ping)
		// Маршруты для получения изделий, прошедших испытания
		r.Get("/lab/{lab_id}/category/{category_id}", h.getTestedItemsByLabAndCategoryInPeriod) // По лаборатории, категории и периоду
		r.Get("/lab/{lab_id}", h.getTestedItemsByLabInPeriod)                                   // По лаборатории и периоду
	})

	router.Route("/lab-workers", func(r chi.Router) {
		r.Get("/ping", h.ping)
		// Маршруты для получения информации об испытателях
		// GET /lab-workers/lab/{lab_id}?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&item_id=X&category_id=Y
		r.Get("/lab/{lab_id}", h.getLabWorkersByItemOrCategoryInLabInPeriod) // Фильтры item_id и category_id - опциональные query параметры
		// GET /lab-workers/all/lab/{lab_id}?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
		r.Get("/all/lab/{lab_id}", h.getLabWorkersInLabInPeriod) // Все испытатели в лаборатории за период
	})

	router.Route("/lab-equipment", func(r chi.Router) {
		r.Get("/ping", h.ping)
		// Маршруты для получения информации об оборудовании лабораторий
		// GET /lab-equipment/lab/{lab_id}?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&item_id=X&category_id=Y
		r.Get("/lab/{lab_id}", h.getLabEquipmentByItemOrCategoryInLabInPeriod)
		// GET /lab-equipment/all/lab/{lab_id}?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
		r.Get("/all/lab/{lab_id}", h.getLabEquipmentInLabInPeriod)
	})

	// Новый маршрут для создания сущностей
	router.Route("/create", func(r chi.Router) {
		r.Get("/ping", h.ping)

		// POST /create/hall создаёт новый цех
		r.Post("/hall", h.createProductionHall)
		// POST /create/area создаёт новый участок
		r.Post("/area", h.createProductionArea)
		// POST /create/item-category создаёт новую категорию изделия
		r.Post("/item-category", h.createCategoryItem)
		// POST /create/item-type создаёт новый тип изделия
		r.Post("/item-type", h.createTypeItem)
		// POST /create/item создаёт новое изделие
		r.Post("/item", h.createItem)
		// POST /create/employee создаёт нового сотрудника
		r.Post("/employee", h.createEmployee)
		// POST /create/lab создаёт новую лабораторию
		r.Post("/lab", h.createTestingLaboratory)
		// POST /create/work-team создаёт новую рабочую бригаду
		r.Post("/work-team", h.createWorkTeam)
		// POST /create/area-boss назначает начальника участка
		r.Post("/area-boss", h.createAreaBoss)
		// POST /create/hall-boss назначает начальника цеха
		r.Post("/hall-boss", h.createHallBoss)
		// POST /create/master создаёт нового мастера
		r.Post("/master", h.createMaster)
		// POST /create/worker-boss назначает бригадира
		r.Post("/worker-boss", h.createWorkerBoss)

	})
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "pong"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getItemsTypeByCategory получает типы изделий по ID категории
func (h *Handler) getItemsTypeByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, name FROM type_item WHERE category_id = $1`

	rows, err := h.db.Query(context.Background(), query, categoryID)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []ItemType
	for rows.Next() {
		var item ItemType
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsType получает все типы изделий по id цеха
func (h *Handler) getAllItemsTypeByHallId(w http.ResponseWriter, r *http.Request) {
	hallId := chi.URLParam(r, "hall_id")
	query := `SELECT type_item.id, type_item.name, type_item.category_id FROM (SELECT * FROM item WHERE hall_id = $1) AS filtered RIGHT JOIN type_item ON filtered.type_id = type_item.id`

	rows, err := h.db.Query(context.Background(), query, hallId)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []ItemType
	for rows.Next() {
		var item ItemType
		if err := rows.Scan(&item.ID, &item.Name, &item.Category_id); err != nil { // Исправлено item.Category_id на &item.Category_id
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsByCategory получает все изделия по id категории
func (h *Handler) getAllItemsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	if categoryID == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	query := `SELECT item.id, item.name, item.type_id, item.hall_id, item.status FROM 
                item INNER JOIN (
                SELECT type_item.id 
                FROM type_item 
                WHERE type_item.category_id = $1
                ) as filtered on item.type_id = filtered.id`

	rows, err := h.db.Query(context.Background(), query, categoryID)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Type_id, &item.Hall_id, &item.Status); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsTypeByCategoryAndHallId получает все изделия по id категории и id цеха
func (h *Handler) getAllItemsTypeByCategoryAndHallId(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	hallId := chi.URLParam(r, "hall_id")
	if categoryID == "" || hallId == "" {
		http.Error(w, "Category ID and Hall ID are required", http.StatusBadRequest)
		return
	}

	query := `SELECT type_item.id, type_item.name FROM (SELECT * FROM item WHERE hall_id = $1 AND type_id = $2) AS filtered RIGHT JOIN type_item ON filtered.type_id = type_item.id`

	rows, err := h.db.Query(context.Background(), query, hallId, categoryID)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []ItemType
	for rows.Next() {
		var item ItemType
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsByHallId получает все изделия по id цеха
func (h *Handler) getAllItemsByHallId(w http.ResponseWriter, r *http.Request) {
	hallId := chi.URLParam(r, "hall_id")
	query := `SELECT item.id, item.name, item.type_id, item.hall_id, item.status FROM item WHERE item.hall_id = $1`

	rows, err := h.db.Query(context.Background(), query, hallId)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Type_id, &item.Hall_id, &item.Status); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsByAreaId получает все изделия по id участка
func (h *Handler) getAllItemsByAreaId(w http.ResponseWriter, r *http.Request) {
	areaId := chi.URLParam(r, "area_id")
	query := `SELECT item.id, item.name, item.type_id, item.hall_id, item.status
				FROM item
					INNER JOIN (
					SELECT areas_items.item_id
					FROM areas_items
					WHERE areas_items.area_id = $1
				) AS ai ON item.id = ai.item_id`

	rows, err := h.db.Query(context.Background(), query, areaId)
	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Type_id, &item.Hall_id, &item.Status); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getAllItemsByTimeInterval получает все изделия по интервалу времени
func (h *Handler) getAllItemsByTimeInterval(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	finishDate := r.URL.Query().Get("finish_date")

	query := `SELECT
				i.id,
				i.name,
				SUM(ri.counter) as item_count
				FROM item i
         			INNER JOIN ready_item ri ON i.id = ri.item_id
				WHERE ri.start_date >= $1 AND ri.completion_date <= $2
				GROUP BY i.id`

	rows, err := h.db.Query(context.Background(), query, startDate, finishDate)

	if err != nil {
		h.logger.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []ReadyItem

	for rows.Next() {
		var item ReadyItem
		if err := rows.Scan(&item.ItemId, &item.Name, &item.Count); err != nil {
			h.logger.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		h.logger.Printf("Error encoding json: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
