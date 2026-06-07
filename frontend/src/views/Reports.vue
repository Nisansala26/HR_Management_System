<template>
  <div class="reports-container">
    <!-- Breadcrumbs Header -->
    <div class="reports-header">
      <div class="header-left">
        <span class="breadcrumb">Dashboard &gt; Reports</span>
        <h1 class="page-title">
          <span class="icon-wrapper">📊</span>
          Advanced Enterprise Reports Generator
        </h1>
      </div>
    </div>

    <!-- Console Control Console -->
    <div class="control-console-card">
      <div class="control-filters">
        <!-- 1. Report Category Select -->
        <div class="filter-group category-select-group">
          <label class="filter-label">Report Category</label>
          <div class="select-container">
            <span class="icon-addon">{{ currentCategoryIcon }}</span>
            <select v-model="selectedCategory" class="console-select" @change="onCategoryChange">
              <option value="attendance">Attendance Summary (Live DB)</option>
              <option value="payroll">Payroll & Salaries (Live DB)</option>
              <option value="leaves">Leave Management (Live DB)</option>
            </select>
          </div>
        </div>

        <!-- 2. Department Filter -->
        <div class="filter-group">
          <label class="filter-label">Department</label>
          <div class="select-container">
            <span class="icon-addon">🏢</span>
            <select v-model="selectedDept" class="console-select">
              <option value="">All Departments</option>
              <option v-for="dept in departments" :key="dept.id" :value="dept.id">
                {{ dept.name }}
              </option>
            </select>
          </div>
        </div>

        <!-- 3. Status Filter -->
        <div class="filter-group">
          <label class="filter-label">Status</label>
          <div class="select-container">
            <span class="icon-addon">⚡</span>
            <select v-model="selectedStatus" class="console-select">
              <option value="">All Status</option>
              <option v-for="status in availableStatuses" :key="status" :value="status">
                {{ status }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Action Buttons Row -->
      <div class="console-actions">
        <button class="btn btn-compile" @click="compileReport" :disabled="loading">
          <span class="btn-icon">⚡</span> Compile Report
        </button>
        <button class="btn btn-download-csv" @click="downloadCSV" :disabled="loading || records.length === 0">
          <span class="btn-icon">📥</span> Download CSV
        </button>
        <button class="btn btn-print-pdf" @click="printPDF" :disabled="loading || records.length === 0">
          <span class="btn-icon">🖨️</span> Print / PDF
        </button>
        <button class="btn btn-add-entry" @click="openAddModal">
          <span class="btn-icon">+</span> Add Entry
        </button>
      </div>
    </div>

    <!-- Master Ledger Block -->
    <div class="ledger-card" v-if="compiled">
      <div class="ledger-header">
        <h2 class="ledger-title">{{ ledgerTitle }} Master Ledger</h2>
        <div class="ledger-meta">
          Compiled Date: <span class="meta-val">{{ compiledDate }}</span>
          <span class="meta-divider">|</span>
          Total Records: <span class="meta-val">{{ filteredRecords.length }}</span>
        </div>
      </div>

      <div class="table-container-custom">
        <table v-if="!loading && filteredRecords.length > 0">
          <thead>
            <tr v-if="selectedCategory === 'attendance'">
              <th>Employee Name</th>
              <th>Department</th>
              <th>Check In</th>
              <th>Check Out</th>
              <th>Status</th>
              <th class="text-right">Actions</th>
            </tr>
            <tr v-else-if="selectedCategory === 'payroll'">
              <th>Employee Name</th>
              <th>Department</th>
              <th>Period</th>
              <th>Base Salary</th>
              <th>Allowances</th>
              <th>Deductions</th>
              <th>Net Salary</th>
              <th>Status</th>
              <th class="text-right">Actions</th>
            </tr>
            <tr v-else-if="selectedCategory === 'leaves'">
              <th>Employee Name</th>
              <th>Department</th>
              <th>Leave Type</th>
              <th>From</th>
              <th>To</th>
              <th>Days</th>
              <th>Reason</th>
              <th>Status</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in filteredRecords" :key="row.id" class="table-row">
              <!-- Employee Column -->
              <td>
                <div class="emp-meta-wrapper">
                  <span class="emp-name-text">{{ row.employee_name }}</span>
                  <span class="emp-id-badge">ID: EMP-{{ String(row.employee_id).padStart(4, '0') }}</span>
                </div>
              </td>
              
              <!-- Department Column -->
              <td>{{ row.department_name }}</td>

              <!-- Attendance Specific fields -->
              <template v-if="selectedCategory === 'attendance'">
                <td>{{ formatDateTime(row.clock_in) }}</td>
                <td>{{ formatDateTime(row.clock_out) }}</td>
                <td>
                  <span class="badge" :class="'badge-' + row.status.toLowerCase()">
                    {{ row.status }}
                  </span>
                </td>
              </template>

              <!-- Payroll Specific fields -->
              <template v-else-if="selectedCategory === 'payroll'">
                <td class="font-medium">{{ row.month }}</td>
                <td class="currency-val-col">LKR {{ row.basic_salary.toLocaleString('en-US') }}</td>
                <td class="currency-val-col text-success">+LKR {{ row.allowances.toLocaleString('en-US') }}</td>
                <td class="currency-val-col text-danger">-LKR {{ row.deductions.toLocaleString('en-US') }}</td>
                <td class="currency-val-col net-salary-col">LKR {{ row.net_salary.toLocaleString('en-US') }}</td>
                <td>
                  <span class="badge" :class="'badge-' + row.status.toLowerCase()">
                    {{ row.status }}
                  </span>
                </td>
              </template>

              <!-- Leaves Specific fields -->
              <template v-else-if="selectedCategory === 'leaves'">
                <td>
                  <span class="leave-type-pill">{{ row.leave_type }}</span>
                </td>
                <td>{{ row.start_date }}</td>
                <td>{{ row.end_date }}</td>
                <td class="days-col">{{ calculateDays(row.start_date, row.end_date) }}</td>
                <td class="reason-col-text" :title="row.reason">{{ truncateText(row.reason, 25) }}</td>
                <td>
                  <span class="badge" :class="'badge-' + row.status.toLowerCase()">
                    {{ row.status }}
                  </span>
                </td>
              </template>

              <!-- Actions -->
              <td class="text-right">
                <div class="actions-cell-wrapper">
                  <!-- Quick actions for pending leaves -->
                  <template v-if="selectedCategory === 'leaves' && row.status === 'Pending'">
                    <button class="action-quick-btn quick-approve" @click="quickApproveLeave(row.id)" title="Quick Approve">
                      ✓
                    </button>
                    <button class="action-quick-btn quick-reject" @click="quickRejectLeave(row.id)" title="Quick Reject">
                      ✗
                    </button>
                  </template>
                  
                  <!-- Edit Button -->
                  <button class="action-btn-circle edit-btn" @click="openEditModal(row)" title="Edit Record">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                      <path d="M12 20h9M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z" />
                    </svg>
                  </button>

                  <!-- Delete Button -->
                  <button class="action-btn-circle delete-btn" @click="deleteRecord(row)" title="Delete Record">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                      <polyline points="3 6 5 6 21 6" />
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Loading Spinner -->
        <div v-else-if="loading" class="empty-state">
          <div class="spinner"></div>
          <p>Compiling database logs...</p>
        </div>

        <!-- Empty state info -->
        <div v-else class="empty-state">
          <span class="empty-icon">📁</span>
          <h3>No matching logs found</h3>
          <p>Modify your filter values or click "⚡ Compile Report" to fetch the latest updates from the live database.</p>
        </div>
      </div>
    </div>

    <!-- 1. Attendance Modals (Add / Edit) -->
    <div class="modal-overlay" v-if="showAddModal && selectedCategory === 'attendance'" @click.self="closeAddModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">New Attendance Log</h3>
          <button class="close-btn" @click="closeAddModal">&times;</button>
        </div>
        <form @submit.prevent="submitAddAttendance">
          <div class="form-group">
            <label>Employee</label>
            <select class="form-control" v-model="attendanceForm.employee_id" required>
              <option value="" disabled>Select employee...</option>
              <option v-for="emp in employees" :key="emp.id" :value="emp.id">
                {{ emp.first_name }} {{ emp.last_name }} ({{ emp.role }})
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>Date</label>
            <input type="date" class="form-control" v-model="attendanceForm.date" required />
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Clock In (HH:MM)</label>
              <input type="time" class="form-control" v-model="attendanceForm.clock_in" />
            </div>
            <div class="form-group half">
              <label>Clock Out (HH:MM)</label>
              <input type="time" class="form-control" v-model="attendanceForm.clock_out" />
            </div>
          </div>
          <div class="form-group">
            <label>Status</label>
            <select class="form-control" v-model="attendanceForm.status" required>
              <option value="Present">Present</option>
              <option value="Late">Late</option>
              <option value="Absent">Absent</option>
            </select>
          </div>
          <div class="modal-form-actions">
            <button type="button" class="btn btn-secondary" @click="closeAddModal">Cancel</button>
            <button type="submit" class="btn btn-primary">Save Entry</button>
          </div>
        </form>
      </div>
    </div>

    <div class="modal-overlay" v-if="showEditModal && editForm.category === 'attendance'" @click.self="closeEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Edit Attendance Log</h3>
          <button class="close-btn" @click="closeEditModal">&times;</button>
        </div>
        <form @submit.prevent="submitEditAttendance">
          <div class="form-group">
            <label>Employee</label>
            <div class="read-only-box">{{ editForm.employee_name }}</div>
          </div>
          <div class="form-group">
            <label>Date</label>
            <input type="date" class="form-control" v-model="attendanceForm.date" required />
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Clock In</label>
              <input type="time" class="form-control" v-model="attendanceForm.clock_in" />
            </div>
            <div class="form-group half">
              <label>Clock Out</label>
              <input type="time" class="form-control" v-model="attendanceForm.clock_out" />
            </div>
          </div>
          <div class="form-group">
            <label>Status</label>
            <select class="form-control" v-model="attendanceForm.status" required>
              <option value="Present">Present</option>
              <option value="Late">Late</option>
              <option value="Absent">Absent</option>
            </select>
          </div>
          <div class="modal-form-actions">
            <button type="button" class="btn btn-secondary" @click="closeEditModal">Cancel</button>
            <button type="submit" class="btn btn-primary">Save Changes</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 2. Leaves Modals (Add / Edit) -->
    <div class="modal-overlay" v-if="showAddModal && selectedCategory === 'leaves'" @click.self="closeAddModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Request Leave Log</h3>
          <button class="close-btn" @click="closeAddModal">&times;</button>
        </div>
        <form @submit.prevent="submitAddLeave">
          <div class="form-group">
            <label>Employee</label>
            <select class="form-control" v-model="leavesForm.employee_id" required>
              <option value="" disabled>Select employee...</option>
              <option v-for="emp in employees" :key="emp.id" :value="emp.id">
                {{ emp.first_name }} {{ emp.last_name }} ({{ emp.role }})
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>Leave Type</label>
            <select class="form-control" v-model="leavesForm.leave_type" required>
              <option value="Annual">Annual</option>
              <option value="Medical">Medical</option>
              <option value="Personal">Personal</option>
              <option value="Sick">Sick</option>
            </select>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>From Date</label>
              <input type="date" class="form-control" v-model="leavesForm.start_date" required />
            </div>
            <div class="form-group half">
              <label>To Date</label>
              <input type="date" class="form-control" v-model="leavesForm.end_date" required />
            </div>
          </div>
          <div class="form-group">
            <label>Reason / Remarks</label>
            <textarea class="form-control text-area" v-model="leavesForm.reason" rows="3" required></textarea>
          </div>
          <div class="modal-form-actions">
            <button type="button" class="btn btn-secondary" @click="closeAddModal">Cancel</button>
            <button type="submit" class="btn btn-primary">File Request</button>
          </div>
        </form>
      </div>
    </div>

    <div class="modal-overlay" v-if="showEditModal && editForm.category === 'leaves'" @click.self="closeEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Edit Leave Application</h3>
          <button class="close-btn" @click="closeEditModal">&times;</button>
        </div>
        <form @submit.prevent="submitEditLeave">
          <div class="form-group">
            <label>Employee</label>
            <div class="read-only-box">{{ editForm.employee_name }}</div>
          </div>
          <div class="form-group">
            <label>Leave Type</label>
            <select class="form-control" v-model="leavesForm.leave_type" required>
              <option value="Annual">Annual</option>
              <option value="Medical">Medical</option>
              <option value="Personal">Personal</option>
              <option value="Sick">Sick</option>
            </select>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>From Date</label>
              <input type="date" class="form-control" v-model="leavesForm.start_date" required />
            </div>
            <div class="form-group half">
              <label>To Date</label>
              <input type="date" class="form-control" v-model="leavesForm.end_date" required />
            </div>
          </div>
          <div class="form-group">
            <label>Status</label>
            <select class="form-control" v-model="leavesForm.status" required>
              <option value="Pending">Pending</option>
              <option value="Approved">Approved</option>
              <option value="Rejected">Rejected</option>
            </select>
          </div>
          <div class="form-group">
            <label>Reason / Remarks</label>
            <textarea class="form-control text-area" v-model="leavesForm.reason" rows="3" required></textarea>
          </div>
          <div class="modal-form-actions">
            <button type="button" class="btn btn-secondary" @click="closeEditModal">Cancel</button>
            <button type="submit" class="btn btn-primary">Save Changes</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 3. Payroll Modals (Add / Edit) -->
    <div class="modal-overlay" v-if="showAddModal && selectedCategory === 'payroll'" @click.self="closeAddModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">New Payroll Entry</h3>
          <button class="close-btn" @click="closeAddModal">&times;</button>
        </div>
        <form @submit.prevent="submitAddPayroll">
          <div class="form-group">
            <label>Employee</label>
            <select class="form-control" v-model="payrollForm.employee_id" required>
              <option value="" disabled>Select employee...</option>
              <option v-for="emp in employees" :key="emp.id" :value="emp.id">
                {{ emp.first_name }} {{ emp.last_name }} ({{ emp.role }})
              </option>
            </select>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Month (YYYY-MM)</label>
              <input type="month" class="form-control" v-model="payrollForm.month" required />
            </div>
            <div class="form-group half">
              <label>Status</label>
              <select class="form-control" v-model="payrollForm.status" required>
                <option value="Unpaid">Unpaid</option>
                <option value="Processing">Processing</option>
                <option value="Paid">Paid</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Basic Salary (LKR)</label>
              <input type="number" class="form-control" v-model="payrollForm.basic_salary" @input="onBasicSalaryChange" required />
            </div>
            <div class="form-group half">
              <label>Calculated Net Salary</label>
              <div class="calculated-net">
                LKR {{ (payrollForm.basic_salary + payrollForm.allowances - payrollForm.deductions).toLocaleString('en-US') }}
              </div>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Allowances (LKR)</label>
              <input type="number" class="form-control" v-model="payrollForm.allowances" required />
            </div>
            <div class="form-group half">
              <label>Deductions (LKR)</label>
              <input type="number" class="form-control" v-model="payrollForm.deductions" required />
            </div>
          </div>
          <div class="modal-form-actions">
            <button type="button" class="btn btn-secondary" @click="closeAddModal">Cancel</button>
            <button type="submit" class="btn btn-primary">Save Log</button>
          </div>
        </form>
      </div>
    </div>

    <div class="modal-overlay" v-if="showEditModal && editForm.category === 'payroll'" @click.self="closeEditModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Edit Payroll Record</h3>
          <button class="close-btn" @click="closeEditModal">&times;</button>
        </div>
        <form @submit.prevent="submitEditPayroll">
          <div class="form-group">
            <label>Employee</label>
            <div class="read-only-box">{{ editForm.employee_name }}</div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Period</label>
              <div class="read-only-box">{{ payrollForm.month }}</div>
            </div>
            <div class="form-group half">
              <label>Status</label>
              <select class="form-control" v-model="payrollForm.status" required>
                <option value="Unpaid">Unpaid</option>
                <option value="Processing">Processing</option>
                <option value="Paid">Paid</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Basic Salary (LKR)</label>
              <input type="number" class="form-control" v-model="payrollForm.basic_salary" @input="onBasicSalaryChange" required />
            </div>
            <div class="form-group half">
              <label>Calculated Net Salary</label>
              <div class="calculated-net">
                LKR {{ (payrollForm.basic_salary + payrollForm.allowances - payrollForm.deductions).toLocaleString('en-US') }}
              </div>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label>Allowances (LKR)</label>
              <input type="number" class="form-control" v-model="payrollForm.allowances" required />
            </div>
            <div class="form-group half">
              <label>Deductions (LKR)</label>
              <input type="number" class="form-control" v-model="payrollForm.deductions" required />
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
const searchQuery = ref('')
watch(() => route.query.search, (newVal) => {
  searchQuery.value = newVal || ''
}, { immediate: true })

// Console Filters Refs
const selectedCategory = ref('attendance')
const selectedDept = ref('')
const selectedStatus = ref('')

const departments = ref([])
const employees = ref([])
const records = ref([])
const loading = ref(false)
const compiled = ref(false)
const compiledDate = ref('--')

// Add/Edit Modals Visible Refs
const showAddModal = ref(false)
const showEditModal = ref(false)

// Modals Forms Data Models
const attendanceForm = ref({
  employee_id: '',
  date: '',
  clock_in: '',
  clock_out: '',
  status: 'Present'
})

const leavesForm = ref({
  employee_id: '',
  leave_type: 'Annual',
  start_date: '',
  end_date: '',
  reason: '',
  status: 'Pending'
})

const payrollForm = ref({
  employee_id: '',
  month: '',
  basic_salary: 0,
  allowances: 0,
  deductions: 0,
  status: 'Unpaid'
})

const editForm = ref({
  id: null,
  category: '',
  employee_name: ''
})

// Dynamic selectors computed options
const currentCategoryIcon = computed(() => {
  if (selectedCategory.value === 'attendance') return '📅'
  if (selectedCategory.value === 'payroll') return '💰'
  if (selectedCategory.value === 'leaves') return '🌴'
  return '📊'
})

const availableStatuses = computed(() => {
  if (selectedCategory.value === 'attendance') {
    return ['Present', 'Late', 'Absent']
  } else if (selectedCategory.value === 'payroll') {
    return ['Paid', 'Processing', 'Unpaid']
  } else if (selectedCategory.value === 'leaves') {
    return ['Pending', 'Approved', 'Rejected']
  }
  return []
})

const ledgerTitle = computed(() => {
  if (selectedCategory.value === 'attendance') return 'Attendance'
  if (selectedCategory.value === 'payroll') return 'Payroll & Salaries'
  if (selectedCategory.value === 'leaves') return 'Leave Management'
  return ''
})

// Lifecycle Loaders
onMounted(() => {
  fetchCommonData()
  compileReport() // compile automatically on mounted
})

const fetchCommonData = async () => {
  try {
    const [deptRes, empRes] = await Promise.all([
      fetch(`${API_URL}/api/departments`),
      fetch(`${API_URL}/api/employees`)
    ])
    if (deptRes.ok) departments.value = await deptRes.json()
    if (empRes.ok) employees.value = await empRes.json()
  } catch (err) {
    console.error('Error loading metadata:', err)
  }
}

// category change handler
const onCategoryChange = () => {
  selectedStatus.value = ''
  records.value = []
  compiled.value = false
}

// Compile Report
const compileReport = async () => {
  loading.value = true
  compiled.value = false
  try {
    let url = ''
    if (selectedCategory.value === 'attendance') {
      url = `${API_URL}/api/attendance`
    } else if (selectedCategory.value === 'payroll') {
      url = `${API_URL}/api/payroll`
    } else if (selectedCategory.value === 'leaves') {
      url = `${API_URL}/api/leaves`
    }

    const res = await fetch(url)
    if (res.ok) {
      records.value = await res.json()
    }
    
    // Format Compile Date
    const now = new Date()
    compiledDate.value = now.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
    
    compiled.value = true
  } catch (err) {
    console.error('Error compiling report:', err)
    alert('Failed to compile report from live database.')
  } finally {
    loading.value = false
  }
}

// Client-Side computed values filtering
const mappedRecords = computed(() => {
  return records.value.map(row => {
    // lookup employee department
    const emp = employees.value.find(e => e.id === row.employee_id)
    return {
      ...row,
      employee_name: row.employee_name || (emp ? `${emp.first_name} ${emp.last_name}` : `EMP #${row.employee_id}`),
      department_id: row.department_id || (emp ? emp.department_id : null),
      department_name: row.department_name || (emp ? emp.department_name : 'Unassigned')
    }
  })
})

const filteredRecords = computed(() => {
  return mappedRecords.value.filter(row => {
    // 1. Department Filter
    if (selectedDept.value && row.department_id != selectedDept.value) {
      return false
    }
    // 2. Status Filter
    if (selectedStatus.value && row.status.toLowerCase() !== selectedStatus.value.toLowerCase()) {
      return false
    }
    // 3. Search Query Filter
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchesName = row.employee_name && row.employee_name.toLowerCase().includes(q)
      const matchesDept = row.department_name && row.department_name.toLowerCase().includes(q)
      if (!matchesName && !matchesDept) {
        return false
      }
    }
    return true
  })
})

// Formatting and Helpers
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return '--'
  const date = new Date(dateTimeStr)
  const yyyy = date.getFullYear()
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const dd = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd} ${hh}:${min}`
}

const getHHMM = (dateTimeStr) => {
  if (!dateTimeStr) return ''
  const date = new Date(dateTimeStr)
  const hh = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  return `${hh}:${min}`
}

const calculateDays = (start, end) => {
  if (!start || !end) return 0
  const s = new Date(start)
  const e = new Date(end)
  const diffTime = Math.abs(e - s)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)) + 1
  return diffDays || 0
}

const truncateText = (text, limit) => {
  if (!text) return ''
  return text.length > limit ? text.substring(0, limit) + '...' : text
}

const onBasicSalaryChange = () => {
  payrollForm.value.allowances = Math.round(payrollForm.value.basic_salary * 0.05)
  payrollForm.value.deductions = Math.round(payrollForm.value.basic_salary * 0.02)
}

// Modals Controls
const openAddModal = () => {
  if (employees.value.length === 0) {
    alert("Please seed or add employees first.")
    return
  }

  const firstEmpId = employees.value[0].id

  if (selectedCategory.value === 'attendance') {
    attendanceForm.value = {
      employee_id: firstEmpId,
      date: new Date().toISOString().slice(0, 10),
      clock_in: '09:00',
      clock_out: '17:00',
      status: 'Present'
    }
  } else if (selectedCategory.value === 'leaves') {
    leavesForm.value = {
      employee_id: firstEmpId,
      leave_type: 'Annual',
      start_date: new Date().toISOString().slice(0, 10),
      end_date: new Date().toISOString().slice(0, 10),
      reason: '',
      status: 'Pending'
    }
  } else if (selectedCategory.value === 'payroll') {
    payrollForm.value = {
      employee_id: firstEmpId,
      month: new Date().toISOString().slice(0, 7),
      basic_salary: 60000,
      allowances: 3000,
      deductions: 1200,
      status: 'Unpaid'
    }
  }

  showAddModal.value = true
}

const closeAddModal = () => {
  showAddModal.value = false
}

const openEditModal = (row) => {
  editForm.value.id = row.id
  editForm.value.category = selectedCategory.value
  editForm.value.employee_name = row.employee_name

  if (selectedCategory.value === 'attendance') {
    attendanceForm.value = {
      employee_id: row.employee_id,
      date: row.date,
      clock_in: getHHMM(row.clock_in),
      clock_out: getHHMM(row.clock_out),
      status: row.status
    }
  } else if (selectedCategory.value === 'leaves') {
    leavesForm.value = {
      employee_id: row.employee_id,
      leave_type: row.leave_type,
      start_date: row.start_date,
      end_date: row.end_date,
      reason: row.reason,
      status: row.status
    }
  } else if (selectedCategory.value === 'payroll') {
    payrollForm.value = {
      employee_id: row.employee_id,
      month: row.month,
      basic_salary: row.basic_salary,
      allowances: row.allowances,
      deductions: row.deductions,
      status: row.status
    }
  }

  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
}

// Add CRUD Submits
const submitAddAttendance = async () => {
  try {
    const res = await fetch(`${API_URL}/api/attendance`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(attendanceForm.value)
    })
    if (res.ok) {
      showAddModal.value = false
      compileReport()
    } else {
      const text = await res.text()
      alert(`Error saving: ${text}`)
    }
  } catch (err) {
    console.error('Error adding attendance:', err)
  }
}

const submitAddLeave = async () => {
  try {
    const res = await fetch(`${API_URL}/api/leaves`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(leavesForm.value)
    })
    if (res.ok) {
      showAddModal.value = false
      compileReport()
    } else {
      const text = await res.text()
      alert(`Error saving: ${text}`)
    }
  } catch (err) {
    console.error('Error adding leave:', err)
  }
}

const submitAddPayroll = async () => {
  try {
    const res = await fetch(`${API_URL}/api/payroll`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payrollForm.value)
    })
    if (res.ok) {
      showAddModal.value = false
      compileReport()
    } else {
      const text = await res.text()
      alert(`Error saving: ${text}`)
    }
  } catch (err) {
    console.error('Error adding payroll:', err)
  }
}

// Edit CRUD Submits
const submitEditAttendance = async () => {
  try {
    const res = await fetch(`${API_URL}/api/attendance?id=${editForm.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(attendanceForm.value)
    })
    if (res.ok) {
      showEditModal.value = false
      compileReport()
    } else {
      const text = await res.text()
      alert(`Error updating: ${text}`)
    }
  } catch (err) {
    console.error('Error editing attendance:', err)
  }
}

const submitEditLeave = async () => {
  try {
    const res = await fetch(`${API_URL}/api/leaves?id=${editForm.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(leavesForm.value)
    })
    if (res.ok) {
      showEditModal.value = false
      compileReport()
    } else {
      const text = await res.text()
      alert(`Error updating: ${text}`)
    }
  } catch (err) {
    console.error('Error editing leave:', err)
  }
}

const submitEditPayroll = async () => {
  try {
    const res = await fetch(`${API_URL}/api/payroll?id=${editForm.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payrollForm.value)
    })
    if (res.ok) {
      showEditModal.value = false
      compileReport()
    } else {
      const text = await res.text()
      alert(`Error updating: ${text}`)
    }
  } catch (err) {
    console.error('Error editing payroll:', err)
  }
}

// Quick Actions for Leaves
const quickApproveLeave = async (id) => {
  try {
    const res = await fetch(`${API_URL}/api/leaves/approve?id=${id}`, {
      method: 'PUT'
    })
    if (res.ok) {
      compileReport()
    }
  } catch (err) {
    console.error('Error approving leave:', err)
  }
}

const quickRejectLeave = async (id) => {
  try {
    const res = await fetch(`${API_URL}/api/leaves/reject?id=${id}`, {
      method: 'PUT'
    })
    if (res.ok) {
      compileReport()
    }
  } catch (err) {
    console.error('Error rejecting leave:', err)
  }
}

// Delete CRUD Action
const deleteRecord = async (row) => {
  if (!confirm(`Are you sure you want to delete this log?`)) return
  
  try {
    let url = ''
    if (selectedCategory.value === 'attendance') {
      url = `${API_URL}/api/attendance?id=${row.id}`
    } else if (selectedCategory.value === 'payroll') {
      url = `${API_URL}/api/payroll?id=${row.id}`
    } else if (selectedCategory.value === 'leaves') {
      url = `${API_URL}/api/leaves?id=${row.id}`
    }

    const res = await fetch(url, { method: 'DELETE' })
    if (res.ok) {
      compileReport()
    } else {
      const text = await res.text()
      alert(`Error deleting: ${text}`)
    }
  } catch (err) {
    console.error('Error deleting record:', err)
  }
}

// CSV Exporter
const downloadCSV = () => {
  if (filteredRecords.value.length === 0) return
  let csvContent = '\uFEFF' // BOM for Excel encoding issues
  
  if (selectedCategory.value === 'attendance') {
    csvContent += 'Employee Name,Department,Check In,Check Out,Status\n'
    filteredRecords.value.forEach(row => {
      csvContent += `"${row.employee_name}","${row.department_name}","${formatDateTime(row.clock_in)}","${formatDateTime(row.clock_out)}","${row.status}"\n`
    })
  } else if (selectedCategory.value === 'payroll') {
    csvContent += 'Employee Name,Department,Period,Base Salary,Allowances,Deductions,Net Salary,Status\n'
    filteredRecords.value.forEach(row => {
      csvContent += `"${row.employee_name}","${row.department_name}","${row.month}",${row.basic_salary},${row.allowances},${row.deductions},${row.net_salary},"${row.status}"\n`
    })
  } else if (selectedCategory.value === 'leaves') {
    csvContent += 'Employee Name,Department,Leave Type,From,To,Days,Reason,Status\n'
    filteredRecords.value.forEach(row => {
      const days = calculateDays(row.start_date, row.end_date)
      const cleanReason = row.reason.replace(/"/g, '""')
      csvContent += `"${row.employee_name}","${row.department_name}","${row.leave_type}","${row.start_date}","${row.end_date}",${days},"${cleanReason}","${row.status}"\n`
    })
  }
  
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement("a")
  link.setAttribute("href", url)
  link.setAttribute("download", `${ledgerTitle.value.toLowerCase().replace(/ /g, '_')}_report.csv`)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// Print / PDF Exporter
const printPDF = () => {
  if (filteredRecords.value.length === 0) {
    alert("No records to export.")
    return
  }
  const doc = new jsPDF()
  
  // Header Design
  doc.setFont("helvetica", "bold")
  doc.setFontSize(18)
  doc.setTextColor(31, 66, 54) // Forest Green branding color (#1f4236)
  doc.text(`${ledgerTitle.value} Master Ledger`, 14, 20)
  
  doc.setFont("helvetica", "normal")
  doc.setFontSize(10)
  doc.setTextColor(100, 116, 139) // slate-500
  doc.text(`Compiled Date: ${compiledDate.value} | Total Records: ${filteredRecords.value.length}`, 14, 28)

  // Draw table header line
  doc.setFillColor(31, 66, 54)
  doc.rect(14, 34, 182, 8, "F")
  doc.setTextColor(255, 255, 255)
  doc.setFont("helvetica", "bold")
  doc.setFontSize(9)

  let y = 48
  if (selectedCategory.value === 'attendance') {
    doc.text("Employee", 16, 39)
    doc.text("Department", 65, 39)
    doc.text("Check In", 100, 39)
    doc.text("Check Out", 145, 39)
    doc.text("Status", 180, 39)

    doc.setFont("helvetica", "normal")
    doc.setTextColor(15, 23, 42)
    filteredRecords.value.forEach(row => {
      if (y > 275) {
        doc.addPage()
        y = 20
      }
      doc.text(truncateText(row.employee_name, 22), 16, y)
      doc.text(truncateText(row.department_name, 16), 65, y)
      doc.text(formatDateTime(row.clock_in), 100, y)
      doc.text(formatDateTime(row.clock_out), 145, y)
      doc.text(row.status, 180, y)
      y += 8
    })
  } else if (selectedCategory.value === 'payroll') {
    doc.text("Employee", 16, 39)
    doc.text("Department", 60, 39)
    doc.text("Period", 95, 39)
    doc.text("Net Salary", 125, 39)
    doc.text("Status", 170, 39)

    doc.setFont("helvetica", "normal")
    doc.setTextColor(15, 23, 42)
    filteredRecords.value.forEach(row => {
      if (y > 275) {
        doc.addPage()
        y = 20
      }
      doc.text(truncateText(row.employee_name, 22), 16, y)
      doc.text(truncateText(row.department_name, 16), 60, y)
      doc.text(row.month, 95, y)
      doc.text(`LKR ${row.net_salary.toLocaleString()}`, 125, y)
      doc.text(row.status, 170, y)
      y += 8
    })
  } else if (selectedCategory.value === 'leaves') {
    doc.text("Employee", 16, 39)
    doc.text("Type", 58, 39)
    doc.text("From", 80, 39)
    doc.text("To", 110, 39)
    doc.text("Days", 140, 39)
    doc.text("Status", 165, 39)

    doc.setFont("helvetica", "normal")
    doc.setTextColor(15, 23, 42)
    filteredRecords.value.forEach(row => {
      if (y > 275) {
        doc.addPage()
        y = 20
      }
      doc.text(truncateText(row.employee_name, 22), 16, y)
      doc.text(row.leave_type, 58, y)
      doc.text(row.start_date, 80, y)
      doc.text(row.end_date, 110, y)
      doc.text(String(calculateDays(row.start_date, row.end_date)), 140, y)
      doc.text(row.status, 165, y)
      y += 8
    })
  }
  
  doc.save(`${ledgerTitle.value.toLowerCase().replace(/ /g, '_')}_report.pdf`)
}
</script>

<style scoped>
.reports-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  font-family: var(--font-body);
}

.reports-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.breadcrumb {
  font-size: 0.8rem;
  color: var(--text-muted);
  display: block;
  margin-bottom: 0.25rem;
  font-weight: 500;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.6rem;
  color: var(--primary);
}

.icon-wrapper {
  background: var(--accent-light);
  padding: 0.5rem;
  border-radius: var(--radius-md);
  font-size: 1.3rem;
}

/* Console Control Card */
.control-console-card {
  background: white;
  border-radius: var(--radius-md);
  padding: 1.5rem;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.control-filters {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.filter-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-secondary);
}

.select-container {
  position: relative;
  display: flex;
  align-items: center;
}

.icon-addon {
  position: absolute;
  left: 0.75rem;
  font-size: 1rem;
  pointer-events: none;
}

.console-select {
  width: 100%;
  padding: 0.625rem 0.75rem 0.625rem 2.2rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
  background: #f8fafc;
  font-family: var(--font-body);
  font-size: 0.875rem;
  color: var(--text-primary);
  font-weight: 500;
  cursor: pointer;
  outline: none;
  transition: all 0.2s ease;
}

.console-select:focus {
  border-color: var(--primary-light);
  box-shadow: 0 0 0 3px rgba(45, 161, 121, 0.15);
  background: white;
}

/* Action buttons styling */
.console-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  padding-top: 0.5rem;
  border-top: 1px solid #f1f5f9;
}

.btn-icon {
  font-size: 1.1rem;
}

.btn-compile {
  background: #0f172a;
  color: white;
}
.btn-compile:hover {
  background: #1e293b;
}

.btn-download-csv {
  background: #10b981;
  color: white;
}
.btn-download-csv:hover {
  background: #059669;
}

.btn-print-pdf {
  background: #3b82f6;
  color: white;
}
.btn-print-pdf:hover {
  background: #2563eb;
}

.btn-add-entry {
  background: white;
  color: var(--primary);
  border: 1.5px solid var(--border-color);
}
.btn-add-entry:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
}

/* Ledger details section styling */
.ledger-card {
  background: white;
  border-radius: var(--radius-md);
  padding: 1.5rem;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
  animation: fadeIn 0.3s ease;
}

.ledger-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.25rem;
  padding-bottom: 0.75rem;
  border-bottom: 2.5px solid #f1f5f9;
}

.ledger-title {
  font-size: 1.2rem;
  color: var(--primary);
  font-weight: 600;
}

.ledger-meta {
  font-size: 0.8rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.meta-val {
  font-weight: 600;
  color: var(--text-primary);
}

.meta-divider {
  margin: 0 0.5rem;
  color: var(--text-muted);
}

/* Custom table styling to ensure premium feel */
.table-container-custom {
  overflow-x: auto;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
}

.table-row {
  transition: all 0.2s ease;
}

.emp-meta-wrapper {
  display: flex;
  flex-direction: column;
}

.emp-name-text {
  font-weight: 600;
  color: var(--text-primary);
}

.emp-id-badge {
  font-size: 0.7rem;
  color: var(--text-muted);
  font-weight: 500;
}

.currency-val-col {
  font-family: var(--font-display);
  font-weight: 600;
}

.net-salary-col {
  color: var(--primary);
  font-weight: 700;
}

.leave-type-pill {
  background: #eff6ff;
  color: #1e40af;
  padding: 0.2rem 0.5rem;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: 600;
}

.days-col {
  font-weight: 600;
  color: var(--text-secondary);
}

.reason-col-text {
  color: var(--text-secondary);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.actions-cell-wrapper {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  align-items: center;
}

/* Quick Actions buttons */
.action-quick-btn {
  width: 22px;
  height: 22px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.quick-approve {
  background: #ecfdf5;
  color: #059669;
}
.quick-approve:hover {
  background: #059669;
  color: white;
}

.quick-reject {
  background: #fef2f2;
  color: #dc2626;
}
.quick-reject:hover {
  background: #dc2626;
  color: white;
}

/* Action Circle buttons */
.action-btn-circle {
  width: 26px;
  height: 26px;
  border-radius: var(--radius-full);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-color);
  cursor: pointer;
  background: white;
  transition: all 0.2s ease;
  color: var(--text-secondary);
}

.edit-btn:hover {
  background: var(--info-light);
  color: var(--info);
  border-color: var(--info);
}

.delete-btn:hover {
  background: var(--danger-light);
  color: var(--danger);
  border-color: var(--danger);
}

/* Form row elements */
.form-row {
  display: flex;
  gap: 1rem;
}

.form-group.half {
  flex: 1;
}

.read-only-box {
  background: #f1f5f9;
  padding: 0.625rem 0.875rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 600;
}

.calculated-net {
  background: #ecfdf5;
  padding: 0.625rem 0.875rem;
  border-radius: var(--radius-sm);
  border: 1.5px dashed var(--success);
  font-size: 0.9rem;
  color: var(--success);
  font-weight: 700;
}

.text-area {
  resize: vertical;
}

.modal-form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
}

/* Spinner Loader */
.spinner {
  border: 3px solid #f3f3f3;
  border-top: 3px solid var(--primary);
  border-radius: 50%;
  width: 32px;
  height: 32px;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 2.5rem;
  display: block;
  margin-bottom: 0.5rem;
}

/* Media Query Responsiveness */
@media (max-width: 768px) {
  .control-filters {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
}
</style>
