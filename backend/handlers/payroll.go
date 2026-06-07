package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"erp-backend/database"
	"erp-backend/models"
)

func HandlePayroll(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPayroll(w, r)
	case http.MethodPost:
		createPayroll(w, r)
	case http.MethodPut:
		updatePayroll(w, r)
	case http.MethodDelete:
		deletePayroll(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPayroll(w http.ResponseWriter, r *http.Request) {
	monthFilter := r.URL.Query().Get("month")
	empIdStr := r.URL.Query().Get("employee_id")

	query := `SELECT p.id, p.employee_id, e.first_name || ' ' || e.last_name, p.month, p.basic_salary, p.allowances, p.deductions, p.net_salary, p.status, p.payment_date 
	          FROM payroll p 
	          JOIN employees e ON p.employee_id = e.id`

	var args []interface{}
	argCount := 1

	if monthFilter != "" {
		query += " WHERE p.month = $1"
		args = append(args, monthFilter)
		argCount++
	}

	if empIdStr != "" {
		empId, err := strconv.Atoi(empIdStr)
		if err == nil {
			if argCount == 1 {
				query += " WHERE p.employee_id = $1"
			} else {
				query += " AND p.employee_id = $2"
			}
			args = append(args, empId)
		}
	}

	query += " ORDER BY p.month DESC, p.id ASC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying payroll: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	records := []models.Payroll{}
	for rows.Next() {
		var rec models.Payroll
		err := rows.Scan(
			&rec.ID, &rec.EmployeeID, &rec.EmployeeName, &rec.Month,
			&rec.BasicSalary, &rec.Allowances, &rec.Deductions, &rec.NetSalary,
			&rec.Status, &rec.PaymentDate,
		)
		if err != nil {
			log.Printf("Error scanning payroll: %v", err)
			continue
		}
		records = append(records, rec)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func GeneratePayroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Month string `json:"month"` // YYYY-MM
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Month) != 7 {
		http.Error(w, "Invalid month. Use YYYY-MM format", http.StatusBadRequest)
		return
	}

	// Fetch all active employees
	rows, err := database.DB.Query("SELECT id, salary FROM employees WHERE status = 'Active'")
	if err != nil {
		log.Printf("Error fetching active employees for payroll: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	insertedCount := 0

	for rows.Next() {
		var empID int
		var salary float64
		if err := rows.Scan(&empID, &salary); err != nil {
			continue
		}

		// Compute basic allowances (e.g. 5%) and deductions (e.g. 2%)
		allowances := salary * 0.05
		deductions := salary * 0.02
		netSalary := salary + allowances - deductions

		// Insert payroll record if it doesn't already exist for this month
		query := `INSERT INTO payroll (employee_id, month, basic_salary, allowances, deductions, net_salary, status) 
		          VALUES ($1, $2, $3, $4, $5, $6, 'Unpaid') 
		          ON CONFLICT (employee_id, month) DO NOTHING`
		
		res, err := tx.Exec(query, empID, req.Month, salary, allowances, deductions, netSalary)
		if err != nil {
			log.Printf("Error inserting payroll for emp %d: %v", empID, err)
			continue
		}
		
		affected, _ := res.RowsAffected()
		if affected > 0 {
			insertedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing payroll: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Payroll compilation completed",
		"records_created": insertedCount,
		"month":          req.Month,
	})
}

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
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

	res, err := database.DB.Exec("UPDATE payroll SET status = 'Paid', payment_date = $1 WHERE id = $2", time.Now(), id)
	if err != nil {
		log.Printf("Error processing payment: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Payroll record not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Payroll payment processed successfully"})
}

func createPayroll(w http.ResponseWriter, r *http.Request) {
	var rec models.Payroll
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if rec.EmployeeID == 0 || rec.Month == "" || rec.BasicSalary == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Calculate net salary if not provided
	if rec.NetSalary == 0 {
		rec.NetSalary = rec.BasicSalary + rec.Allowances - rec.Deductions
	}

	query := `INSERT INTO payroll (employee_id, month, basic_salary, allowances, deductions, net_salary, status, payment_date) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	          RETURNING id`
	
	var paymentDate *time.Time
	if rec.Status == "Paid" {
		t := time.Now()
		paymentDate = &t
	}

	err = database.DB.QueryRow(query, rec.EmployeeID, rec.Month, rec.BasicSalary, rec.Allowances, rec.Deductions, rec.NetSalary, rec.Status, paymentDate).Scan(&rec.ID)
	if err != nil {
		log.Printf("Error inserting custom payroll: %v", err)
		http.Error(w, "Internal Server Error or Conflict (record already exists for this month)", http.StatusInternalServerError)
		return
	}

	// Fetch employee name
	_ = database.DB.QueryRow("SELECT first_name || ' ' || last_name FROM employees WHERE id = $1", rec.EmployeeID).Scan(&rec.EmployeeName)
	if paymentDate != nil {
		rec.PaymentDate = paymentDate
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rec)
}

func updatePayroll(w http.ResponseWriter, r *http.Request) {
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

	var rec models.Payroll
	err = json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Re-calculate Net
	rec.NetSalary = rec.BasicSalary + rec.Allowances - rec.Deductions

	var res sql.Result
	if rec.Status == "Paid" {
		t := time.Now()
		res, err = database.DB.Exec(
			"UPDATE payroll SET basic_salary = $1, allowances = $2, deductions = $3, net_salary = $4, status = $5, payment_date = $6 WHERE id = $7",
			rec.BasicSalary, rec.Allowances, rec.Deductions, rec.NetSalary, rec.Status, t, id,
		)
	} else {
		res, err = database.DB.Exec(
			"UPDATE payroll SET basic_salary = $1, allowances = $2, deductions = $3, net_salary = $4, status = $5, payment_date = NULL WHERE id = $6",
			rec.BasicSalary, rec.Allowances, rec.Deductions, rec.NetSalary, rec.Status, id,
		)
	}

	if err != nil {
		log.Printf("Error updating payroll: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Payroll record not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Payroll updated successfully"})
}

func deletePayroll(w http.ResponseWriter, r *http.Request) {
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

	res, err := database.DB.Exec("DELETE FROM payroll WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting payroll: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Payroll record not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Payroll deleted successfully"})
}
