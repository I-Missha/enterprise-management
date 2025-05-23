package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// WorkTeam представляет информацию о бригаде
type WorkTeam struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	AreaID   int    `json:"area_id"`
	AreaName string `json:"area_name,omitempty"`
	HallID   int    `json:"hall_id"`
	HallName string `json:"hall_name,omitempty"`
}

// WorkTeamMember представляет информацию о члене бригады
type WorkTeamMember struct {
	EmployeeID      int    `json:"employee_id"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	HireDate        string `json:"hire_date,omitempty"`
	CurrentPosition string `json:"current_position,omitempty"`
}

// WorkTeamWithMembers представляет бригаду с её участниками
type WorkTeamWithMembers struct {
	Team    WorkTeam         `json:"team"`
	Members []WorkTeamMember `json:"members"`
}

// getWorkTeamsByArea получает все бригады указанного участка и их состав
func (h *Handler) getWorkTeamsByArea(w http.ResponseWriter, r *http.Request) {
	areaID := chi.URLParam(r, "area_id")
	if areaID == "" {
		http.Error(w, "Area ID is required", http.StatusBadRequest)
		return
	}

	// Сначала получаем все бригады участка
	teamsQuery := `
		SELECT 
			wt.id, 
			wt.name, 
			wt.area_id, 
			pa.name as area_name, 
			wt.hall_id, 
			ph.name as hall_name
		FROM work_team wt
		JOIN production_area pa ON wt.area_id = pa.id
		JOIN production_halls ph ON wt.hall_id = ph.id
		WHERE wt.area_id = $1
		ORDER BY wt.name
	`

	teamRows, err := h.db.Query(context.Background(), teamsQuery, areaID)
	if err != nil {
		h.logger.Printf("Error querying database for work teams by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer teamRows.Close()

	var teams []WorkTeam
	for teamRows.Next() {
		var team WorkTeam
		var areaName, hallName sql.NullString

		if err := teamRows.Scan(
			&team.ID,
			&team.Name,
			&team.AreaID,
			&areaName,
			&team.HallID,
			&hallName,
		); err != nil {
			h.logger.Printf("Error scanning row for work teams by area: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if areaName.Valid {
			team.AreaName = areaName.String
		}

		if hallName.Valid {
			team.HallName = hallName.String
		}

		teams = append(teams, team)
	}

	if err := teamRows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for work teams by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var result []WorkTeamWithMembers

	// Для каждой бригады получаем её членов
	for _, team := range teams {
		membersQuery := `
			SELECT 
				w.employee_id, 
				e.name, 
				w.category, 
				e.hire_date, 
				e.current_position
			FROM worker w
			JOIN employee e ON w.employee_id = e.id
			WHERE w.work_team_id = $1
			ORDER BY e.name
		`

		memberRows, err := h.db.Query(context.Background(), membersQuery, team.ID)
		if err != nil {
			h.logger.Printf("Error querying database for work team members: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var members []WorkTeamMember
		for memberRows.Next() {
			var member WorkTeamMember
			var hireDate, currentPosition sql.NullString

			if err := memberRows.Scan(
				&member.EmployeeID,
				&member.Name,
				&member.Category,
				&hireDate,
				&currentPosition,
			); err != nil {
				memberRows.Close()
				h.logger.Printf("Error scanning row for work team members: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if hireDate.Valid {
				member.HireDate = hireDate.String
			}

			if currentPosition.Valid {
				member.CurrentPosition = currentPosition.String
			}

			members = append(members, member)
		}
		memberRows.Close()

		if err := memberRows.Err(); err != nil {
			h.logger.Printf("Error after iterating rows for work team members: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		result = append(result, WorkTeamWithMembers{
			Team:    team,
			Members: members,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Printf("Error encoding json for work teams by area: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getWorkTeamsByHall получает все бригады указанного цеха и их состав
func (h *Handler) getWorkTeamsByHall(w http.ResponseWriter, r *http.Request) {
	hallID := chi.URLParam(r, "hall_id")
	if hallID == "" {
		http.Error(w, "Hall ID is required", http.StatusBadRequest)
		return
	}

	// Сначала получаем все бригады цеха
	teamsQuery := `
		SELECT 
			wt.id, 
			wt.name, 
			wt.area_id, 
			pa.name as area_name, 
			wt.hall_id, 
			ph.name as hall_name
		FROM work_team wt
		JOIN production_area pa ON wt.area_id = pa.id
		JOIN production_halls ph ON wt.hall_id = ph.id
		WHERE wt.hall_id = $1
		ORDER BY pa.name, wt.name
	`

	teamRows, err := h.db.Query(context.Background(), teamsQuery, hallID)
	if err != nil {
		h.logger.Printf("Error querying database for work teams by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer teamRows.Close()

	var teams []WorkTeam
	for teamRows.Next() {
		var team WorkTeam
		var areaName, hallName sql.NullString

		if err := teamRows.Scan(
			&team.ID,
			&team.Name,
			&team.AreaID,
			&areaName,
			&team.HallID,
			&hallName,
		); err != nil {
			h.logger.Printf("Error scanning row for work teams by hall: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if areaName.Valid {
			team.AreaName = areaName.String
		}

		if hallName.Valid {
			team.HallName = hallName.String
		}

		teams = append(teams, team)
	}

	if err := teamRows.Err(); err != nil {
		h.logger.Printf("Error after iterating rows for work teams by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var result []WorkTeamWithMembers

	// Для каждой бригады получаем её членов
	for _, team := range teams {
		membersQuery := `
			SELECT 
				w.employee_id, 
				e.name, 
				w.category, 
				e.hire_date, 
				e.current_position
			FROM worker w
			JOIN employee e ON w.employee_id = e.id
			WHERE w.work_team_id = $1
			ORDER BY e.name
		`

		memberRows, err := h.db.Query(context.Background(), membersQuery, team.ID)
		if err != nil {
			h.logger.Printf("Error querying database for work team members: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var members []WorkTeamMember
		for memberRows.Next() {
			var member WorkTeamMember
			var hireDate, currentPosition sql.NullString

			if err := memberRows.Scan(
				&member.EmployeeID,
				&member.Name,
				&member.Category,
				&hireDate,
				&currentPosition,
			); err != nil {
				memberRows.Close()
				h.logger.Printf("Error scanning row for work team members: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if hireDate.Valid {
				member.HireDate = hireDate.String
			}

			if currentPosition.Valid {
				member.CurrentPosition = currentPosition.String
			}

			members = append(members, member)
		}
		memberRows.Close()

		if err := memberRows.Err(); err != nil {
			h.logger.Printf("Error after iterating rows for work team members: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		result = append(result, WorkTeamWithMembers{
			Team:    team,
			Members: members,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Printf("Error encoding json for work teams by hall: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
