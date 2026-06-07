package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"erp-backend/database"
	"erp-backend/models"
)

func HandleDepartments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getDepartments(w, r)
	case http.MethodPost:
		createDepartment(w, r)
	case http.MethodPut:
		updateDepartment(w, r)
	case http.MethodDelete:
		deleteDepartment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getDepartments(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}
		var dept models.Department
		query := `SELECT d.id, d.name, d.manager_name, d.budget, d.created_at, COUNT(e.id) 
		          FROM departments d 
		          LEFT JOIN employees e ON d.id = e.department_id 
		          WHERE d.id = $1 
		          GROUP BY d.id`
		err = database.DB.QueryRow(query, id).Scan(&dept.ID, &dept.Name, &dept.ManagerName, &dept.Budget, &dept.CreatedAt, &dept.EmployeeCount)
		if err != nil {
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dept)
		return
	}

	query := `SELECT d.id, d.name, d.manager_name, d.budget, d.created_at, COUNT(e.id) 
	          FROM departments d 
	          LEFT JOIN employees e ON d.id = e.department_id 
	          GROUP BY d.id 
	          ORDER BY d.id ASC`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error querying departments: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	depts := []models.Department{}
	for rows.Next() {
		var dept models.Department
		err := rows.Scan(&dept.ID, &dept.Name, &dept.ManagerName, &dept.Budget, &dept.CreatedAt, &dept.EmployeeCount)
		if err != nil {
			log.Printf("Error scanning department: %v", err)
			continue
		}
		depts = append(depts, dept)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(depts)
}

func createDepartment(w http.ResponseWriter, r *http.Request) {
	var dept models.Department
	err := json.NewDecoder(r.Body).Decode(&dept)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if dept.Name == "" || dept.ManagerName == "" || dept.Budget < 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO departments (name, manager_name, budget) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, created_at`
	err = database.DB.QueryRow(query, dept.Name, dept.ManagerName, dept.Budget).Scan(&dept.ID, &dept.CreatedAt)
	if err != nil {
		log.Printf("Error inserting department: %v", err)
		http.Error(w, "Could not create department (check duplicate name)", http.StatusBadRequest)
		return
	}

	dept.EmployeeCount = 0
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func updateDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	var dept models.Department
	err = json.NewDecoder(r.Body).Decode(&dept)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `UPDATE departments SET name = $1, manager_name = $2, budget = $3 WHERE id = $4`
	res, err := database.DB.Exec(query, dept.Name, dept.ManagerName, dept.Budget, id)
	if err != nil {
		log.Printf("Error updating department: %v", err)
		http.Error(w, "Could not update department", http.StatusBadRequest)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	dept.ID = id
	// Fetch actual employee count
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM employees WHERE department_id = $1", id).Scan(&dept.EmployeeCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

func deleteDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec("DELETE FROM departments WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting department: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetDepartmentStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var stats struct {
		TotalDepartments int `json:"total_departments"`
		ActiveSectors    int `json:"active_sectors"`
		TotalEmployees   int `json:"total_employees"`
		DepartmentHeads  int `json:"department_heads"`
	}

	err := database.DB.QueryRow("SELECT COUNT(*) FROM departments").Scan(&stats.TotalDepartments)
	if err != nil {
		log.Printf("Error querying department total stats: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = database.DB.QueryRow("SELECT COUNT(*) FROM departments WHERE budget > 0").Scan(&stats.ActiveSectors)
	if err != nil {
		log.Printf("Error querying active sectors stats: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = database.DB.QueryRow("SELECT COUNT(*) FROM employees WHERE status = 'Active'").Scan(&stats.TotalEmployees)
	if err != nil {
		log.Printf("Error querying total employees stats: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = database.DB.QueryRow("SELECT COUNT(DISTINCT manager_name) FROM departments").Scan(&stats.DepartmentHeads)
	if err != nil {
		log.Printf("Error querying department heads stats: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

