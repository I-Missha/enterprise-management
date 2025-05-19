package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Worker struct {
	EmployeeID int    `json:"employee_id"`
	Category   string `json:"category"`
	Name       string `json:"name,omitempty"`
	IsBoss     bool   `json:"is_boss,omitempty"`
}

type WorkerRepository struct {
	db *sql.DB
}

func NewWorkerRepository(db *sql.DB) *WorkerRepository {
	return &WorkerRepository{
		db: db,
	}
}

var ErrNotFound = errors.New("record not found")

func (r *WorkerRepository) GetAll(ctx context.Context) ([]Worker, error) {
	query := `
		SELECT w.employee_id, w.category, e.name, 
		       CASE WHEN wb.worker_id IS NOT NULL THEN true ELSE false END as is_boss
		FROM worker w
		JOIN employee e ON w.employee_id = e.id
		LEFT JOIN worker_boss wb ON w.employee_id = wb.worker_id
		ORDER BY w.employee_id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка рабочих: %w", err)
	}
	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var w Worker
		if err := rows.Scan(&w.EmployeeID, &w.Category, &w.Name, &w.IsBoss); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки рабочего: %w", err)
		}
		workers = append(workers, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return workers, nil
}

func (r *WorkerRepository) GetByID(ctx context.Context, id int) (Worker, error) {
	query := `
		SELECT w.employee_id, w.category, e.name, 
		       CASE WHEN wb.worker_id IS NOT NULL THEN true ELSE false END as is_boss
		FROM worker w
		JOIN employee e ON w.employee_id = e.id
		LEFT JOIN worker_boss wb ON w.employee_id = wb.worker_id
		WHERE w.employee_id = $1
	`

	var w Worker
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.EmployeeID, &w.Category, &w.Name, &w.IsBoss,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Worker{}, ErrNotFound
		}
		return Worker{}, fmt.Errorf("ошибка при получении рабочего: %w", err)
	}

	return w, nil
}

func (r *WorkerRepository) GetByCategory(ctx context.Context, category string) ([]Worker, error) {
	query := `
		SELECT w.employee_id, w.category, e.name, 
		       CASE WHEN wb.worker_id IS NOT NULL THEN true ELSE false END as is_boss
		FROM worker w
		JOIN employee e ON w.employee_id = e.id
		LEFT JOIN worker_boss wb ON w.employee_id = wb.worker_id
		WHERE w.category = $1
		ORDER BY w.employee_id
	`

	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рабочих по категории: %w", err)
	}
	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var w Worker
		if err := rows.Scan(&w.EmployeeID, &w.Category, &w.Name, &w.IsBoss); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки рабочего: %w", err)
		}
		workers = append(workers, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return workers, nil
}

func (r *WorkerRepository) GetByHall(ctx context.Context, hallID int) ([]Worker, error) {
	query := `
		SELECT DISTINCT w.employee_id, w.category, e.name, 
		       CASE WHEN wb.worker_id IS NOT NULL THEN true ELSE false END as is_boss
		FROM worker w
		JOIN employee e ON w.employee_id = e.id
		LEFT JOIN worker_boss wb ON w.employee_id = wb.worker_id
		JOIN work_team_member wtm ON w.employee_id = wtm.worker_id
		JOIN work_team wt ON wtm.work_team_id = wt.id
		JOIN production_area pa ON wt.area_id = pa.id
		WHERE pa.hall_id = $1
		ORDER BY w.employee_id
	`

	rows, err := r.db.QueryContext(ctx, query, hallID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рабочих цеха: %w", err)
	}
	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var w Worker
		if err := rows.Scan(&w.EmployeeID, &w.Category, &w.Name, &w.IsBoss); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки рабочего: %w", err)
		}
		workers = append(workers, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return workers, nil
}

func (r *WorkerRepository) GetByArea(ctx context.Context, areaID int) ([]Worker, error) {
	query := `
		SELECT DISTINCT w.employee_id, w.category, e.name, 
		       CASE WHEN wb.worker_id IS NOT NULL THEN true ELSE false END as is_boss
		FROM worker w
		JOIN employee e ON w.employee_id = e.id
		LEFT JOIN worker_boss wb ON w.employee_id = wb.worker_id
		JOIN work_team_member wtm ON w.employee_id = wtm.worker_id
		JOIN work_team wt ON wtm.work_team_id = wt.id
		WHERE wt.area_id = $1
		ORDER BY w.employee_id
	`

	rows, err := r.db.QueryContext(ctx, query, areaID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рабочих участка: %w", err)
	}
	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var w Worker
		if err := rows.Scan(&w.EmployeeID, &w.Category, &w.Name, &w.IsBoss); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки рабочего: %w", err)
		}
		workers = append(workers, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return workers, nil
}

func (r *WorkerRepository) GetByTeam(ctx context.Context, teamID int) ([]Worker, error) {
	query := `
		SELECT w.employee_id, w.category, e.name, 
		       CASE WHEN wb.worker_id IS NOT NULL THEN true ELSE false END as is_boss
		FROM worker w
		JOIN employee e ON w.employee_id = e.id
		LEFT JOIN worker_boss wb ON w.employee_id = wb.worker_id
		JOIN work_team_member wtm ON w.employee_id = wtm.worker_id
		WHERE wtm.work_team_id = $1
		ORDER BY w.employee_id
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рабочих бригады: %w", err)
	}
	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var w Worker
		if err := rows.Scan(&w.EmployeeID, &w.Category, &w.Name, &w.IsBoss); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки рабочего: %w", err)
		}
		workers = append(workers, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return workers, nil
}
