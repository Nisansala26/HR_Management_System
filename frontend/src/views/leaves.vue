package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"erp-backend/database"
	"erp-backend/models"
)

func HandleLeaves(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getLeaves(w, r)
	case http.MethodPost:
		applyLeave(w, r)
	case http.MethodPut:
		updateLeave(w, r)
	case http.MethodDelete:
		deleteLeave(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getLeaves(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")

	query := `SELECT l.id, l.employee_id, e.first_name || ' ' || e.last_name, l.leave_type, l.start_date, l.end_date, l.status, l.reason, l.created_at 
	          FROM leaves l 
	          JOIN employees e ON l.employee_id = e.id`

	var args []interface{}
	if statusFilter != "" {
		query += " WHERE l.status = $1"
		args = append(args, statusFilter)
	}

	query += " ORDER BY l.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying leaves: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	leaves := []models.Leave{}
	for rows.Next() {
		var lv models.Leave
		var startVal, endVal time.Time
		err := rows.Scan(&lv.ID, &lv.EmployeeID, &lv.EmployeeName, &lv.LeaveType, &startVal, &endVal, &lv.Status, &lv.Reason, &lv.CreatedAt)
		if err != nil {
			log.Printf("Error scanning leave: %v", err)
			continue
		}
		lv.StartDate = startVal.Format("2006-01-02")
		lv.EndDate = endVal.Format("2006-01-02")
		leaves = append(leaves, lv)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaves)
}

func applyLeave(w http.ResponseWriter, r *http.Request) {
	var lv models.Leave
	err := json.NewDecoder(r.Body).Decode(&lv)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if lv.EmployeeID == 0 || lv.LeaveType == "" || lv.StartDate == "" || lv.EndDate == "" || lv.Reason == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO leaves (employee_id, leave_type, start_date, end_date, status, reason) 
	          VALUES ($1, $2, $3, $4, 'Pending', $5) 
	          RETURNING id, created_at`
	err = database.DB.QueryRow(query, lv.EmployeeID, lv.LeaveType, lv.StartDate, lv.EndDate, lv.Reason).Scan(&lv.ID, &lv.CreatedAt)
	if err != nil {
		log.Printf("Error inserting leave: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	lv.Status = "Pending"
	_ = database.DB.QueryRow("SELECT first_name || ' ' || last_name FROM employees WHERE id = $1", lv.EmployeeID).Scan(&lv.EmployeeName)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lv)
}

func ApproveLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	res, err := database.DB.Exec("UPDATE leaves SET status = 'Approved' WHERE id = $1", id)
	if err != nil {
		log.Printf("Error approving leave: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Leave request not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Leave approved successfully"})
}

func RejectLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	res, err := database.DB.Exec("UPDATE leaves SET status = 'Rejected' WHERE id = $1", id)
	if err != nil {
		log.Printf("Error rejecting leave: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Leave request not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Leave rejected successfully"})
}

func updateLeave(w http.ResponseWriter, r *http.Request) {
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

	var lv models.Leave
	err = json.NewDecoder(r.Body).Decode(&lv)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if lv.EmployeeID == 0 || lv.LeaveType == "" || lv.StartDate == "" || lv.EndDate == "" || lv.Reason == "" || lv.Status == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	query := `UPDATE leaves 
	          SET employee_id = $1, leave_type = $2, start_date = $3, end_date = $4, status = $5, reason = $6 
	          WHERE id = $7`
	res, err := database.DB.Exec(query, lv.EmployeeID, lv.LeaveType, lv.StartDate, lv.EndDate, lv.Status, lv.Reason, id)
	if err != nil {
		log.Printf("Error updating leave: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Leave request not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Leave request updated successfully"})
}

func deleteLeave(w http.ResponseWriter, r *http.Request) {
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

	res, err := database.DB.Exec("DELETE FROM leaves WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting leave: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Leave request not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Leave request deleted successfully"})
}
