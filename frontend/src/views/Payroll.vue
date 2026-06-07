<template>
  <div class="payroll-view">
    <!-- Header Section -->
    <div class="payroll-header">
      <div class="header-left">
        <h2 class="page-title">
          <!-- Small graphic indicator matching header icon -->
          <span class="payroll-title-icon">💵</span>
          Payroll Management
        </h2>
        <p class="breadcrumbs-path">Dashboard > <span class="active-crumb">Payroll</span></p>
      </div>
      <div class="header-right">
        <!-- Green pill Compile button -->
        <button class="btn btn-success run-payroll-btn" @click="generateMonthlyPayroll">
          + Run Payroll
        </button>
      </div>
    </div>

    <!-- 3 Summary Metrics Cards Row -->
    <div class="metrics-cards-row grid-cols-3">
      <!-- Card 1: Total Payroll (Month) -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper bg-green-light">
          <!-- Money Bag SVG Icon -->
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2a4 4 0 0 0-4 4v2H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V10a2 2 0 0 0-2-2h-3V6a4 4 0 0 0-4-4z" />
            <circle cx="12" cy="14" r="2" fill="#10b981" />
            <path d="M12 11v1" />
          </svg>
        </div>
        <div class="metric-info">
          <span class="metric-label">Total Payroll (Month)</span>
          <h3 class="metric-value">LKR {{ totalNetPayroll.toLocaleString('en-US', {maximumFractionDigits: 0}) }}</h3>
          <span class="trend-label text-success">
            ↑ 4.2% vs last month
          </span>
        </div>
      </div>

      <!-- Card 2: Processed (Paid) -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper bg-green-light">
          <!-- Checkbox SVG Icon -->
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5" xmlns="http://www.w3.org/2000/svg">
            <rect x="3" y="3" width="18" height="18" rx="3" fill="#10b981" fill-opacity="0.1" />
            <polyline points="9 11 12 14 16 8" />
          </svg>
        </div>
        <div class="metric-info">
          <span class="metric-label">Processed</span>
          <h3 class="metric-value">{{ paidCount }} Employees</h3>
          <span class="trend-label text-success">
            Fully Paid
          </span>
        </div>
      </div>

      <!-- Card 3: Pending Approvals -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper bg-orange-light">
          <!-- Hourglass SVG Icon -->
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" xmlns="http://www.w3.org/2000/svg">
            <path d="M5 2h14M5 22h14M19 2v4a7 7 0 0 1-7 7 7 7 0 0 1-7-7V2M5 22v-4a7 7 0 0 1 7-7 7 7 0 0 1 7 7v4" />
          </svg>
        </div>
        <div class="metric-info">
          <span class="metric-label">Pending Approvals</span>
          <h3 class="metric-value">{{ pendingCount }} Employees</h3>
          <span class="trend-label text-warning font-semibold">
            Requires Action
          </span>
        </div>
      </div>
    </div>

    <!-- Filters Control Bar Card -->
    <div class="table-control-bar card">
      <!-- Search input on left -->
      <div class="search-input-container">
        <svg class="magnifier-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input 
          type="text" 
          placeholder="Search employee payroll..." 
          class="custom-search-input" 
          v-model="searchQuery"
        />
      </div>

      <!-- Filters on right -->
      <div class="control-actions-group">
        <!-- Period Month Input -->
        <div class="month-selector-wrapper">
          <span class="calendar-icon">🗓️</span>
          <input 
            type="month" 
            class="period-month-input" 
            v-model="selectedMonth" 
            @change="fetchPayroll"
          />
        </div>

        <!-- Departments dropdown select -->
        <select class="custom-select" v-model="selectedDept">
          <option value="">All Departments</option>
          <option v-for="dept in departments" :key="dept.id" :value="dept.id">
            {{ dept.name }}
          </option>
        </select>

        <!-- Export Report Button -->
        <button class="btn btn-export" @click="exportReport">
          Export Report
        </button>

        <!-- Add Custom Record Button -->
        <button class="btn btn-primary add-record-btn" @click="openAddModal">
          + Add Record
        </button>
      </div>
    </div>

    <!-- Payroll Ledger Sheet Data Grid Table -->
    <div class="table-container">
      <table v-if="!loading && filteredPayrollRecords.length > 0">
        <thead>
          <tr class="header-strip">
            <th>Employee</th>
            <th>Period</th>
            <th>Base Salary</th>
            <th>Allowances (5%)</th>
            <th>Deductions (2%)</th>
            <th>Net Pay</th>
            <th>Status</th>
            <th>Payment Details</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rec in filteredPayrollRecords" :key="rec.id" class="table-row">
            <td>
              <div class="emp-meta">
                <span class="emp-name">{{ rec.employee_name }}</span>
                <span class="emp-id">ID: EMP-{{ String(rec.employee_id).padStart(4, '0') }}</span>
              </div>
            </td>
            <td class="font-medium">{{ rec.month }}</td>
            <td class="currency-cell">LKR {{ rec.basic_salary.toLocaleString('en-US') }}</td>
            <td class="currency-cell text-success">+LKR {{ rec.allowances.toLocaleString('en-US') }}</td>
            <td class="currency-cell text-danger">-LKR {{ rec.deductions.toLocaleString('en-US') }}</td>
            <td class="currency-cell net-pay">LKR {{ rec.net_salary.toLocaleString('en-US') }}</td>
            <td>
              <span class="badge" :class="'badge-' + rec.status.toLowerCase()">
                {{ rec.status }}
              </span>
            </td>
            <td>
              <span class="payment-date" v-if="rec.payment_date">
                Paid: {{ formatDate(rec.payment_date) }}
              </span>
              <span class="payment-date text-muted" v-else>
                --/--/----
              </span>
            </td>
            <td class="text-right">
              <div class="actions-wrapper">
                <button 
                  v-if="rec.status !== 'Paid'" 
                  class="btn btn-success btn-xs" 
                  @click="processPayment(rec.id)"
                  title="Mark as Paid"
                >
                  Pay
                </button>
                <span v-else class="status-locked text-success">
                  Paid ✓
                </span>

                <!-- Edit Button -->
                <button class="action-circle-btn edit-btn" @click="openEditModal(rec)" title="Edit Record">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 20h9M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z" />
                  </svg>
                </button>

                <!-- Delete Button -->
                <button class="action-circle-btn delete-btn" @click="deletePayrollRecord(rec.id)" title="Delete Record">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" xmlns="http://www.w3.org/2000/svg">
                    <polyline points="3 6 5 6 21 6" />
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Spinner Loader -->
      <div v-else-if="loading" class="empty-state">
        <div class="spinner"></div>
        <p>Fetching monthly payroll records...</p>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-state">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
        <p v-if="searchQuery || selectedDept">No matches found for current search/filters.</p>
        <p v-else>No payroll compiled for {{ selectedMonth }}. Click "+ Run Payroll" on the header to generate details.</p>
      </div>
    </div>

    <!-- Add Payroll Record Modal -->
    <div class="modal-overlay" v-if="showAddModal" @click.self="closeAddModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Add Payroll Record</h3>
          <button class="close-btn" @click="closeAddModal">&times;</button>
        </div>
        <form @submit.prevent="submitAddPayroll">
          <div class="form-group">
            <label for="add-emp">Employee</label>
            <select id="add-emp" class="form-control" v-model="payrollForm.employee_id" required>
              <option value="" disabled>Select employee...</option>
              <option v-for="emp in employees" :key="emp.id" :value="emp.id">
                {{ emp.first_name }} {{ emp.last_name }} ({{ emp.role }})
              </option>
            </select>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label for="add-month">Month</label>
              <input type="month" id="add-month" class="form-control" v-model="payrollForm.month" required />
            </div>
            <div class="form-group half">
              <label for="add-status">Status</label>
              <select id="add-status" class="form-control" v-model="payrollForm.status" required>
                <option value="Unpaid">Unpaid</option>
                <option value="Processing">Processing</option>
                <option value="Paid">Paid</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label for="add-salary">Base Salary (LKR)</label>
              <input type="number" id="add-salary" class="form-control" v-model="payrollForm.basic_salary" @input="onBasicSalaryChange" required />
            </div>
            <div class="form-group half">
              <label>Net Salary</label>
              <div class="calculated-net-box">
                LKR {{ (payrollForm.basic_salary + payrollForm.allowances - payrollForm.deductions).toLocaleString('en-US') }}
              </div>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label for="add-allowances">Allowances (LKR)</label>
              <input type="number" id="add-allowances" class="form-control" v-model="payrollForm.allowances" required />
            </div>
            <div class="form-group half">
              <label for="add-deductions">Deductions (LKR)</label>
              <input type="number" id="add-deductions" class="form-control" v-model="payrollForm.deductions" required />
            </div>
          </div>
          <div class="modal-form-actions">
            <button type="button" class="btn btn-secondary" @click="closeAddModal">Cancel</button>
            <button type="submit" class="btn btn-primary">Save</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Payroll Record Modal -->
    <div class="modal-overlay" v-if="showEditModal" @click.self="closeEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Edit Payroll Record</h3>
          <button class="close-btn" @click="closeEditModal">&times;</button>
        </div>
        <form @submit.prevent="submitEditPayroll">
          <div class="form-group">
            <label>Employee Name</label>
            <div class="read-only-field">{{ editForm.employee_name }}</div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Period</label>
              <div class="read-only-field">{{ editForm.month }}</div>
            </div>
            <div class="form-group half">
              <label for="edit-status">Status</label>
              <select id="edit-status" class="form-control" v-model="editForm.status" required>
                <option value="Unpaid">Unpaid</option>
                <option value="Processing">Processing</option>
                <option value="Paid">Paid</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label for="edit-salary">Base Salary (LKR)</label>
              <input type="number" id="edit-salary" class="form-control" v-model="editForm.basic_salary" @input="onEditBasicSalaryChange" required />
            </div>
            <div class="form-group half">
              <label>Net Salary</label>
              <div class="calculated-net-box">
                LKR {{ (editForm.basic_salary + editForm.allowances - editForm.deductions).toLocaleString('en-US') }}
              </div>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label for="edit-allowances">Allowances (LKR)</label>
              <input type="number" id="edit-allowances" class="form-control" v-model="editForm.allowances" required />
            </div>
            <div class="form-group half">
              <label for="edit-deductions">Deductions (LKR)</label>
              <input type="number" id="edit-deductions" class="form-control" v-model="editForm.deductions" required />
            </div>
          </div>
          <div class="modal-form-actions">
            <button type="button" class="btn btn-secondary" @click="closeEditModal">Cancel</button>
            <button type="submit" class="btn btn-primary">Save Changes</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { jsPDF } from 'jspdf'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const route = useRoute()

// Default to current month (e.g. YYYY-MM)
const selectedMonth = ref(new Date().toISOString().slice(0, 7))
const payrollRecords = ref([])
const employees = ref([])
const departments = ref([])
const loading = ref(true)

// Filters console refs
const searchQuery = ref('')

watch(() => route.query.search, (newVal) => {
  searchQuery.value = newVal || ''
}, { immediate: true })

const selectedDept = ref('')

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

// Fetch base data
const fetchPayroll = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL}/api/payroll?month=${selectedMonth.value}`)
    if (res.ok) {
      payrollRecords.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching payroll: ', err)
    // Fallback Mock data
    payrollRecords.value = [
      { id: 1, employee_id: 1, employee_name: 'Admin HR Manager', month: selectedMonth.value, basic_salary: 9500, allowances: 475, deductions: 190, net_salary: 9785, status: 'Paid', payment_date: new Date().toISOString() },
      { id: 2, employee_id: 2, employee_name: 'John Doe', month: selectedMonth.value, basic_salary: 8500, allowances: 425, deductions: 170, net_salary: 8755, status: 'Paid', payment_date: new Date().toISOString() },
      { id: 3, employee_id: 3, employee_name: 'Alice Smith', month: selectedMonth.value, basic_salary: 5500, allowances: 275, deductions: 110, net_salary: 5665, status: 'Paid', payment_date: new Date().toISOString() },
      { id: 4, employee_id: 4, employee_name: 'Michael Brown', month: selectedMonth.value, basic_salary: 4800, allowances: 240, deductions: 96, net_salary: 4944, status: 'Paid', payment_date: new Date().toISOString() },
      { id: 5, employee_id: 5, employee_name: 'Jane Wilson', month: selectedMonth.value, basic_salary: 6000, allowances: 300, deductions: 120, net_salary: 6180, status: 'Paid', payment_date: new Date().toISOString() },
      { id: 6, employee_id: 6, employee_name: 'David Lee', month: selectedMonth.value, basic_salary: 6200, allowances: 310, deductions: 124, net_salary: 6386, status: 'Processing', payment_date: null },
      { id: 7, employee_id: 7, employee_name: 'Sophia Davis', month: selectedMonth.value, basic_salary: 4000, allowances: 200, deductions: 80, net_salary: 4120, status: 'Unpaid', payment_date: null }
    ]
  } finally {
    loading.value = false
  }
}

const fetchEmployees = async () => {
  try {
    const res = await fetch(`${API_URL}/api/employees`)
    if (res.ok) {
      employees.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching employees list:', err)
  }
}

const fetchDepartments = async () => {
  try {
    const res = await fetch(`${API_URL}/api/departments`)
    if (res.ok) {
      departments.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching departments:', err)
    departments.value = [
      { id: 1, name: 'Human Resources' },
      { id: 2, name: 'Engineering' },
      { id: 3, name: 'Marketing' },
      { id: 4, name: 'Finance' }
    ]
  }
}

// Client-side computed mapping to resolve departments
const mappedPayrollRecords = computed(() => {
  return payrollRecords.value.map(rec => {
    const emp = employees.value.find(e => e.id === rec.employee_id)
    return {
      ...rec,
      department_id: emp ? emp.department_id : null,
      department_name: emp ? emp.department_name : 'Unassigned'
    }
  })
})

// Filtered payroll sheet list
const filteredPayrollRecords = computed(() => {
  return mappedPayrollRecords.value.filter(rec => {
    // 1. Search Query filter (Employee Name or ID)
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchName = rec.employee_name && rec.employee_name.toLowerCase().includes(q)
      const matchId = String(rec.employee_id).includes(q)
      if (!matchName && !matchId) return false
    }

    // 2. Department filter
    if (selectedDept.value && rec.department_id != selectedDept.value) {
      return false
    }

    return true
  })
})

// KPI metrics calculations (recomputes reactively based on filtered results)
const totalNetPayroll = computed(() => {
  return filteredPayrollRecords.value.reduce((acc, curr) => acc + curr.net_salary, 0)
})

const paidCount = computed(() => {
  return filteredPayrollRecords.value.filter(r => r.status === 'Paid').length
})

const pendingCount = computed(() => {
  return filteredPayrollRecords.value.filter(r => r.status !== 'Paid').length
})

// Compile payroll sheet for month
const generateMonthlyPayroll = async () => {
  try {
    const res = await fetch(`${API_URL}/api/payroll/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ month: selectedMonth.value })
    })
    if (res.ok) {
      fetchPayroll()
    } else {
      const text = await res.text()
      alert(text || 'Payroll sheet for this month has already been compiled.')
    }
  } catch (err) {
    console.error('Error generating payroll: ', err)
  }
}

// Confirm wire payment
const processPayment = async (id) => {
  try {
    const res = await fetch(`${API_URL}/api/payroll/pay?id=${id}`, {
      method: 'PUT'
    })
    if (res.ok) {
      fetchPayroll()
    }
  } catch (err) {
    const record = payrollRecords.value.find(r => r.id === id)
    if (record) {
      record.status = 'Paid'
      record.payment_date = new Date().toISOString()
    }
  }
}

// Add Modal refs
const showAddModal = ref(false)
const payrollForm = ref({
  employee_id: '',
  month: selectedMonth.value,
  basic_salary: 5000,
  allowances: 250,
  deductions: 100,
  status: 'Unpaid'
})

// Edit Modal refs
const showEditModal = ref(false)
const editForm = ref({
  id: null,
  employee_name: '',
  month: '',
  basic_salary: 0,
  allowances: 0,
  deductions: 0,
  status: 'Unpaid'
})

// Open/Close Modals
const openAddModal = () => {
  if (employees.value.length > 0) {
    payrollForm.value.employee_id = employees.value[0].id
  }
  payrollForm.value.month = selectedMonth.value
  payrollForm.value.basic_salary = 5000
  payrollForm.value.allowances = 250
  payrollForm.value.deductions = 100
  payrollForm.value.status = 'Unpaid'
  showAddModal.value = true
}

const closeAddModal = () => {
  showAddModal.value = false
}

const openEditModal = (rec) => {
  editForm.value = {
    id: rec.id,
    employee_name: rec.employee_name,
    month: rec.month,
    basic_salary: rec.basic_salary,
    allowances: rec.allowances,
    deductions: rec.deductions,
    status: rec.status
  }
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
}

// Watchers / handlers
const onBasicSalaryChange = () => {
  payrollForm.value.allowances = Math.round(payrollForm.value.basic_salary * 0.05)
  payrollForm.value.deductions = Math.round(payrollForm.value.basic_salary * 0.02)
}

const onEditBasicSalaryChange = () => {
  editForm.value.allowances = Math.round(editForm.value.basic_salary * 0.05)
  editForm.value.deductions = Math.round(editForm.value.basic_salary * 0.02)
}

// CRUD actions
const submitAddPayroll = async () => {
  try {
    const res = await fetch(`${API_URL}/api/payroll`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        employee_id: Number(payrollForm.value.employee_id),
        month: payrollForm.value.month,
        basic_salary: Number(payrollForm.value.basic_salary),
        allowances: Number(payrollForm.value.allowances),
        deductions: Number(payrollForm.value.deductions),
        status: payrollForm.value.status
      })
    })
    if (res.ok) {
      showAddModal.value = false
      fetchPayroll()
    } else {
      const text = await res.text()
      alert(`Error: ${text}`)
    }
  } catch (err) {
    console.error('Error adding payroll:', err)
  }
}

const submitEditPayroll = async () => {
  try {
    const res = await fetch(`${API_URL}/api/payroll?id=${editForm.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        basic_salary: Number(editForm.value.basic_salary),
        allowances: Number(editForm.value.allowances),
        deductions: Number(editForm.value.deductions),
        status: editForm.value.status
      })
    })
    if (res.ok) {
      showEditModal.value = false
      fetchPayroll()
    } else {
      const text = await res.text()
      alert(`Error updating: ${text}`)
    }
  } catch (err) {
    console.error('Error editing payroll:', err)
  }
}

const deletePayrollRecord = async (id) => {
  if (!confirm("Are you sure you want to delete this payroll record?")) return
  try {
    const res = await fetch(`${API_URL}/api/payroll?id=${id}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      fetchPayroll()
    } else {
      const text = await res.text()
      alert(`Error deleting: ${text}`)
    }
  } catch (err) {
    console.error('Error deleting payroll:', err)
  }
}

// Working Export to PDF Report Feature
const exportReport = () => {
  if (filteredPayrollRecords.value.length === 0) {
    alert("No payroll records to export for this period/filter.")
    return
  }

  const doc = new jsPDF()

  // 1. Title Header
  doc.setFont("helvetica", "bold")
  doc.setFontSize(18)
  doc.setTextColor(31, 66, 54) // Forest Green branding color (#1f4236)
  doc.text("Payroll Management Report", 14, 20)

  // 2. Metadata / Period details
  doc.setFont("helvetica", "normal")
  doc.setFontSize(10)
  doc.setTextColor(100)
  doc.text(`Period: ${selectedMonth.value}`, 14, 28)
  doc.text(`Generated On: ${new Date().toLocaleDateString('en-US')}`, 14, 34)

  // Separator Line
  doc.setDrawColor(226, 232, 240) // Slate-200 border color
  doc.setLineWidth(0.5)
  doc.line(14, 38, 196, 38)

  // 3. Summary metrics card info
  doc.setFont("helvetica", "bold")
  doc.setFontSize(11)
  doc.setTextColor(15, 23, 42) // Slate-900 text color
  doc.text("Summary Metrics:", 14, 46)

  doc.setFont("helvetica", "normal")
  doc.setFontSize(9.5)
  doc.text(`Total Payroll Liabilities: LKR ${totalNetPayroll.value.toLocaleString('en-US')}`, 14, 53)
  doc.text(`Processed (Paid): ${paidCount.value} Employees`, 14, 59)
  doc.text(`Pending Approvals: ${pendingCount.value} Employees`, 14, 65)

  // Separator Line
  doc.line(14, 71, 196, 71)

  // 4. Employee Detailed table
  doc.setFont("helvetica", "bold")
  doc.setFontSize(11)
  doc.text("Employee Breakdown:", 14, 80)

  // Table Headers
  doc.setFont("helvetica", "bold")
  doc.setFontSize(8.5)
  doc.setTextColor(100)
  let y = 88
  doc.text("Employee", 14, y)
  doc.text("Base Salary", 68, y)
  doc.text("Allowances", 98, y)
  doc.text("Deductions", 128, y)
  doc.text("Net Pay", 155, y)
  doc.text("Status", 182, y)

  doc.line(14, y + 2, 196, y + 2)
  doc.setFont("helvetica", "normal")
  doc.setTextColor(0)
  y += 9

  filteredPayrollRecords.value.forEach((rec) => {
    // Page break handling
    if (y > 275) {
      doc.addPage()
      y = 20
      
      // Reprint Table Headers on new page
      doc.setFont("helvetica", "bold")
      doc.setFontSize(8.5)
      doc.setTextColor(100)
      doc.text("Employee", 14, y)
      doc.text("Base Salary", 68, y)
      doc.text("Allowances", 98, y)
      doc.text("Deductions", 128, y)
      doc.text("Net Pay", 155, y)
      doc.text("Status", 182, y)
      doc.line(14, y + 2, 196, y + 2)
      doc.setFont("helvetica", "normal")
      doc.setTextColor(0)
      y += 9
    }

    // Row Data
    doc.setFont("helvetica", "bold")
    doc.setFontSize(9)
    doc.setTextColor(15, 23, 42)
    doc.text(rec.employee_name, 14, y)
    
    doc.setFont("helvetica", "normal")
    doc.setFontSize(7.5)
    doc.setTextColor(148, 163, 184) // Slate-400
    doc.text(`ID: EMP-${String(rec.employee_id).padStart(4, '0')}`, 14, y + 4.5)

    doc.setFontSize(8.5)
    doc.setTextColor(15, 23, 42)
    doc.text(`LKR ${rec.basic_salary.toLocaleString('en-US')}`, 68, y)
    
    doc.setTextColor(16, 185, 129) // Green for positive allowances
    doc.text(`+LKR ${rec.allowances.toLocaleString('en-US')}`, 98, y)
    
    doc.setTextColor(239, 68, 68) // Red for negative deductions
    doc.text(`-LKR ${rec.deductions.toLocaleString('en-US')}`, 128, y)
    
    doc.setTextColor(15, 23, 42)
    doc.setFont("helvetica", "bold")
    doc.text(`LKR ${rec.net_salary.toLocaleString('en-US')}`, 155, y)
    doc.setFont("helvetica", "normal")

    // Dynamic colored status
    if (rec.status === 'Paid') {
      doc.setTextColor(16, 185, 129) // Green
    } else {
      doc.setTextColor(245, 158, 11) // Orange
    }
    doc.text(rec.status, 182, y)

    // Reset drawing settings for row separator line
    doc.setDrawColor(241, 245, 249) // Slate-100 line
    doc.setLineWidth(0.3)
    doc.line(14, y + 6.5, 196, y + 6.5)
    doc.setDrawColor(226, 232, 240) // Restore default draw color

    y += 12.5
  })

  // Download trigger
  doc.save(`payroll_report_${selectedMonth.value}.pdf`)
}

onMounted(() => {
  fetchPayroll()
  fetchEmployees()
  fetchDepartments()
})
</script>

<style scoped>
.payroll-view {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.payroll-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.page-title {
  font-size: 1.75rem;
  color: var(--text-primary);
  margin-bottom: 0.25rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.payroll-title-icon {
  font-size: 1.6rem;
}

.breadcrumbs-path {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.active-crumb {
  color: var(--primary-light);
  font-weight: 500;
}

.run-payroll-btn {
  background-color: #10b981;
  color: white;
  font-weight: 600;
  border-radius: var(--radius-full);
  padding: 0.65rem 1.5rem;
}

.run-payroll-btn:hover {
  background-color: #059669;
}

/* 3 Metrics Cards Row */
.metrics-cards-row {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1.5rem;
}

.metric-card {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 1.5rem;
  background: white;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  transition: transform 0.2s, box-shadow 0.2s;
}

.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.metric-icon-wrapper {
  width: 46px;
  height: 46px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.bg-green-light {
  background-color: #ecfdf5;
}

.bg-orange-light {
  background-color: #fffbeb;
}

.metric-info {
  display: flex;
  flex-direction: column;
}

.metric-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
}

.metric-value {
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0.2rem 0;
  line-height: 1.2;
}

.trend-label {
  font-size: 0.78rem;
  font-weight: 600;
}

.text-success { color: #10b981; }
.text-warning { color: #f59e0b; }
.text-danger { color: #ef4444; }

/* Control Bar Card styling */
.table-control-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  background-color: white;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
}

.search-input-container {
  position: relative;
  display: flex;
  align-items: center;
  flex: 1;
  max-width: 380px;
}

.magnifier-icon {
  position: absolute;
  left: 0.85rem;
  color: var(--text-muted);
  pointer-events: none;
}

.custom-search-input {
  width: 100%;
  padding: 0.55rem 0.85rem 0.55rem 2.3rem;
  border-radius: var(--radius-full);
  border: 1px solid var(--border-color);
  font-size: 0.88rem;
  color: var(--text-primary);
  background-color: #f8fafc;
  transition: border-color 0.2s, background-color 0.2s;
}

.custom-search-input:focus {
  outline: none;
  border-color: #10b981;
  background-color: white;
}

.control-actions-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

/* Month Input period */
.month-selector-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.calendar-icon {
  position: absolute;
  left: 0.65rem;
  font-size: 0.95rem;
  pointer-events: none;
}

.period-month-input {
  padding: 0.5rem 0.65rem 0.5rem 2rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
  font-size: 0.85rem;
  color: var(--text-secondary);
  background-color: white;
  cursor: pointer;
}

.custom-select {
  padding: 0.5rem 1.75rem 0.5rem 0.75rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
  font-size: 0.85rem;
  color: var(--text-secondary);
  background-color: white;
  cursor: pointer;
}

.custom-select:focus {
  outline: none;
  border-color: #10b981;
}

.btn-export {
  background-color: white;
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  padding: 0.5rem 1.25rem;
  border-radius: var(--radius-full);
  font-size: 0.88rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s, border-color 0.2s;
}

.btn-export:hover {
  background-color: #f8fafc;
  border-color: #cbd5e1;
}

/* Table overrides */
.emp-meta {
  display: flex;
  flex-direction: column;
}

.emp-name {
  font-weight: 600;
}

.emp-id {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.currency-cell {
  font-family: var(--font-display);
  font-weight: 600;
}

.currency-cell.net-pay {
  color: var(--primary-light);
  font-size: 0.95rem;
}

.payment-date {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.btn-xs {
  padding: 0.25rem 0.65rem;
  font-size: 0.75rem;
  border-radius: var(--radius-sm);
  font-weight: 600;
}

.status-locked {
  font-size: 0.8rem;
  font-weight: 600;
}

.text-right {
  text-align: right;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 4rem 0;
  color: var(--text-muted);
  gap: 1rem;
}

.spinner {
  border: 3px solid #e2e8f0;
  border-top: 3px solid #10b981;
  border-radius: 50%;
  width: 28px;
  height: 28px;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (max-width: 1024px) {
  .metrics-cards-row {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  .table-control-bar {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }
  .search-input-container {
    max-width: 100%;
  }
  .control-actions-group {
    flex-wrap: wrap;
  }
}

/* CRUD styles */
.actions-wrapper {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.45rem;
}

.action-circle-btn {
  width: 25px;
  height: 25px;
  border-radius: 50%;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s;
}

.edit-btn {
  background-color: #eff6ff;
  color: #1d4ed8;
}
.edit-btn:hover {
  background-color: #2563eb;
  color: white;
}

.delete-btn {
  background-color: #fef2f2;
  color: #b91c1c;
}
.delete-btn:hover {
  background-color: #ef4444;
  color: white;
}

.read-only-field {
  background-color: #f1f5f9;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  padding: 0.55rem 0.75rem;
  font-size: 0.88rem;
  color: var(--text-secondary);
  font-weight: 500;
  height: 38px;
  display: flex;
  align-items: center;
}

.calculated-net-box {
  background-color: #f8fafc;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  padding: 0.55rem 0.75rem;
  font-size: 0.88rem;
  font-weight: 700;
  color: #10b981;
  display: flex;
  align-items: center;
  height: 38px;
}

.add-record-btn {
  background-color: white;
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  padding: 0.5rem 1.25rem;
  border-radius: var(--radius-full);
  font-size: 0.88rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s, border-color 0.2s;
}
.add-record-btn:hover {
  background-color: #f8fafc;
  border-color: #cbd5e1;
}

/* Modals overlays override */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(15, 23, 42, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: var(--radius-lg);
  padding: 1.75rem;
  width: 100%;
  max-width: 480px;
  box-shadow: var(--shadow-lg);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.5rem;
}

.modal-title {
  font-size: 1.15rem;
  color: var(--text-primary);
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.4rem;
  color: var(--text-muted);
  cursor: pointer;
}
.close-btn:hover {
  color: var(--text-primary);
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group.half {
  flex: 1;
}

.modal-form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 1.25rem;
  border-top: 1px solid var(--border-color);
  padding-top: 1rem;
}
</style>
