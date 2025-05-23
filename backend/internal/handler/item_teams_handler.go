package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ItemTeamResponse представляет ответ с данными о бригадах, участвующих в сборке изделия
type ItemTeamResponse struct {
	ItemID   int                   `json:"item_id"`
	ItemName string                `json:"item_name"`
	Teams    []WorkTeamWithMembers `json:"teams"`
}

// getItemTeams получает состав бригад, участвующих в сборке указанного изделия
func (h *Handler) getItemTeams(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "item_id")
	if itemID == "" {
		http.Error(w, "Item ID is required", http.StatusBadRequest)
		return
	}

	// Получаем название изделия
	var itemName string
	itemQuery := `SELECT name FROM item WHERE id = $1`
	err := h.db.QueryRow(context.Background(), itemQuery, itemID).Scan(&itemName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		h.logger.Printf("Error querying item name: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Получаем список бригад, работающих над изделием через соответствующие работы
	teamsQuery := `
		SELECT DISTINCT 
			wt.id, 
			wt.name, 
			wt.area_id, 
			pa.name as area_name, 
			wt.hall_id, 
			ph.name as hall_name
		FROM item_work_type iwt
		JOIN work_type wtp ON iwt.work_type_id = wtp.id
		JOIN work_team wt ON wtp.work_team_id = wt.id
		JOIN production_area pa ON wt.area_id = pa.id
		JOIN production_halls ph ON wt.hall_id = ph.id
		WHERE iwt.item_id = $1
		ORDER BY wt.area_id, wt.name
	`

	teamRows, err := h.db.Query(context.Background(), teamsQuery, itemID)
	if err != nil {
		h.logger.Printf("Error querying database for teams working on item: %v", err)
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
			h.logger.Printf("Error scanning row for teams working on item: %v", err)
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
		h.logger.Printf("Error after iterating rows for teams working on item: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var teamsWithMembers []WorkTeamWithMembers

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

		teamsWithMembers = append(teamsWithMembers, WorkTeamWithMembers{
			Team:    team,
			Members: members,
		})
	}

	response := ItemTeamResponse{
		ItemID:   int(streamAtoi(itemID)),
		ItemName: itemName,
		Teams:    teamsWithMembers,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding json response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Вспомогательная функция для конвертации строки в число
func streamAtoi(s string) int64 {
	var n int64
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int64(c-'0')
	}
	return n
}
