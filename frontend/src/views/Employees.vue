<template>
  <div class="employees-view">
    <!-- Top Stats Cards Grid (Exactly matching the user screenshot) -->
    <div class="stats-cards-row">
      <!-- Card 1: All Employees -->
      <div class="card metric-card blue-border">
        <div class="metric-icon-bg blue">
          <svg class="metric-icon" width="22" height="22" viewBox="0 0 24 24" fill="currentColor">
            <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
          </svg>
        </div>
        <div class="metric-info">
          <span class="metric-label">ALL EMPLOYEES</span>
          <h3 class="metric-value">{{ employeeStats.total_employees }}</h3>
        </div>
      </div>

      <!-- Card 2: Active Staff -->
      <div class="card metric-card green-border">
        <div class="metric-icon-bg green">
          <svg class="metric-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
        </div>
        <div class="metric-info">
          <span class="metric-label">ACTIVE STAFF</span>
          <h3 class="metric-value">{{ employeeStats.active_staff }}</h3>
        </div>
      </div>

      <!-- Card 3: On Leave -->
      <div class="card metric-card yellow-border">
        <div class="metric-icon-bg yellow">
          <svg class="metric-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M5 3h14M5 21h14M17 3v4a5 5 0 0 1-10 0V3M17 21v-4a5 5 0 0 0-10 0v4M12 12v.01"></path>
          </svg>
        </div>
        <div class="metric-info">
          <span class="metric-label">ON LEAVE</span>
          <h3 class="metric-value">{{ employeeStats.on_leave }}</h3>
        </div>
      </div>
    </div>

    <!-- Filters & Action Bar (Exactly matching user screenshot layout) -->
    <div class="table-control-bar card">
      <!-- Search Input with magnifying glass -->
      <div class="search-input-container">
        <svg class="magnifier-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input 
          type="text" 
          placeholder="Search employees by name, role or ID..." 
          class="custom-search-input" 
          v-model="searchQuery"
          @input="resetToFirstPage"
        />
      </div>

      <!-- Right Side Filters & Button -->
      <div class="control-actions-group">
        <select class="custom-select" v-model="selectedDept" @change="resetToFirstPage">
          <option value="">All Departments</option>
          <option v-for="dept in departments" :key="dept.id" :value="dept.id">
            {{ dept.name }}
          </option>
        </select>
        <button class="btn btn-primary add-employee-btn" @click="openAddModal">
          + Add Employee
        </button>
      </div>
    </div>

    <!-- Table Section (Exactly matching screenshot styling) -->
    <div class="custom-table-container">
      <table class="employees-table">
        <thead>
          <tr class="header-strip">
            <th>EMPLOYEE ID</th>
            <th>FULL NAME</th>
            <th>DEPARTMENT</th>
            <th>POSITION</th>
            <th>STATUS</th>
            <th class="text-right">ACTIONS</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="emp in paginatedEmployees" :key="emp.id" class="table-row">
            <td class="id-cell">EMP-{{ String(emp.id).padStart(4, '0') }}</td>
            <td>
              <div class="name-cell-wrapper">
                <div class="avatar-circle">{{ getInitials(emp.first_name, emp.last_name) }}</div>
                <span class="full-name">{{ emp.first_name }} {{ emp.last_name }}</span>
              </div>
            </td>
            <td>
              <span class="dept-bubble">{{ emp.department_name }}</span>
            </td>
            <td class="position-cell">{{ emp.role }}</td>
            <td>
              <span class="badge" :class="'badge-' + emp.status.toLowerCase()">
                {{ emp.status }}
              </span>
            </td>
            <td>
              <div class="actions-group">
                <button class="action-icon edit-btn" @click="openEditModal(emp)" title="Edit Employee">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                    <path d="M18.5 2.5a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                  </svg>
                </button>
                <button class="action-icon delete-btn" @click="deleteEmployee(emp.id)" title="Delete Employee">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="paginatedEmployees.length === 0">
            <td colspan="6" class="no-records-cell">
              No employees matched the query.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Table Footer / Pagination (Exactly matching screenshot layout) -->
    <div class="table-footer-pagination" v-if="employees.length > 0">
      <div class="footer-entries-info">
        Showing {{ entriesStart }} to {{ entriesEnd }} of {{ totalEntries }} entries
      </div>
      <div class="pagination-controls-group">
        <button 
          class="pag-btn prev" 
          :disabled="currentPage === 1" 
          @click="goToPage(currentPage - 1)"
        >
          Previous
        </button>
        <button 
          v-for="page in totalPages" 
          :key="page" 
          class="pag-btn page-num" 
          :class="{ active: page === currentPage }"
          @click="goToPage(page)"
        >
          {{ page }}
        </button>
        <button 
          class="pag-btn next" 
          :disabled="currentPage === totalPages" 
          @click="goToPage(currentPage + 1)"
        >
          Next
        </button>
      </div>
    </div>

    <!-- Modal Form (Add / Edit) -->
    <div class="modal-overlay" v-if="showModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">{{ editMode ? 'Modify Employee Profile' : 'Add New Employee' }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <form @submit.prevent="submitForm">
          <div class="form-row">
            <div class="form-group half">
              <label for="first_name">First Name</label>
              <input type="text" id="first_name" class="form-control" v-model="form.first_name" required />
            </div>
            <div class="form-group half">
              <label for="last_name">Last Name</label>
              <input type="text" id="last_name" class="form-control" v-model="form.last_name" required />
            </div>
          </div>
          <div class="form-group">
            <label for="email">Email</label>
            <input type="email" id="email" class="form-control" v-model="form.email" required />
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label for="role">Position / Role</label>
              <input type="text" id="role" class="form-control" placeholder="e.g. Lead Engineer" v-model="form.role" required />
            </div>
            <div class="form-group half">
              <label for="salary">Base Salary ($)</label>
              <input type="number" id="salary" class="form-control" v-model="form.salary" required min="1"/>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label for="department">Department</label>
              <select id="department" class="form-control" v-model="form.department_id" required>
                <option :value="null">Unassigned</option>
                <option v-for="dept in departments" :key="dept.id" :value="dept.id">
                  {{ dept.name }}
                </option>
              </select>
            </div>
            <div class="form-group half">
              <label for="join_date">Join Date</label>
              <input type="date" id="join_date" class="form-control" v-model="form.join_date" required />
            </div>
          </div>
          <div class="form-group" v-if="editMode">
            <label for="status">Work Status</label>
            <select id="status" class="form-control" v-model="form.status">
              <option value="Active">Active</option>
              <option value="Suspended">Suspended</option>
              <option value="Resigned">Resigned</option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
            <button type="submit" class="btn btn-primary">
              {{ editMode ? 'Save Changes' : 'Register' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const route = useRoute()
const employees = ref([])
const departments = ref([])
const employeeStats = ref({ total_employees: 0, active_staff: 0, on_leave: 0 })
const searchQuery = ref('')

watch(() => route.query.search, (newVal) => {
  searchQuery.value = newVal || ''
}, { immediate: true })

const selectedDept = ref('')

// Pagination
const currentPage = ref(1)
const itemsPerPage = ref(5)

const showModal = ref(false)
const editMode = ref(false)
const currentEditId = ref(null)

const form = ref({
  first_name: '',
  last_name: '',
  email: '',
  role: '',
  salary: 5000,
  department_id: null,
  join_date: new Date().toISOString().split('T')[0],
  status: 'Active'
})

const getInitials = (first, last) => {
  return `${first?.charAt(0) || ''}${last?.charAt(0) || ''}`.toUpperCase()
}

const fetchStats = async () => {
  try {
    const res = await fetch(`${API_URL}/api/employees/stats`)
    if (res.ok) {
      employeeStats.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching employee stats: ', err)
    // Compute locally from mock data
    const total = employees.value.length
    const active = employees.value.filter(e => e.status === 'Active').length
    employeeStats.value = {
      total_employees: total || 8,
      active_staff: active || 7,
      on_leave: 1
    }
  }
}

const fetchDepartments = async () => {
  try {
    const res = await fetch(`${API_URL}/api/departments`)
    if (res.ok) {
      departments.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching depts: ', err)
    departments.value = [
      { id: 1, name: 'Human Resources' },
      { id: 2, name: 'Engineering' },
      { id: 3, name: 'Marketing' },
      { id: 4, name: 'Finance' }
    ]
  }
}

const fetchEmployees = async () => {
  try {
    const res = await fetch(`${API_URL}/api/employees`)
    if (res.ok) {
      employees.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching employees list: ', err)
    // Fallback Mock
    employees.value = [
      { id: 1, first_name: 'Admin', last_name: 'HR Manager', email: 'admin@erp.com', role: 'HR Manager', salary: 9500, status: 'Active', department_id: 1, department_name: 'Human Resources', join_date: '2024-01-15' },
      { id: 2, first_name: 'John', last_name: 'Doe', email: 'john.doe@erp.com', role: 'Senior Software Engineer', salary: 8500, status: 'Active', department_id: 2, department_name: 'Engineering', join_date: '2024-03-01' },
      { id: 3, first_name: 'Alice', last_name: 'Smith', email: 'alice.smith@erp.com', role: 'QA Engineer', salary: 5500, status: 'Active', department_id: 2, department_name: 'Engineering', join_date: '2024-06-10' },
      { id: 4, first_name: 'Michael', last_name: 'Brown', email: 'michael.brown@erp.com', role: 'Marketing Specialist', salary: 4800, status: 'Active', department_id: 3, department_name: 'Marketing', join_date: '2024-11-01' },
      { id: 5, first_name: 'Jane', last_name: 'Wilson', email: 'jane.wilson@erp.com', role: 'Financial Analyst', salary: 6000, status: 'Active', department_id: 4, department_name: 'Finance', join_date: '2025-02-20' },
      { id: 6, first_name: 'David', last_name: 'Lee', email: 'david.lee@erp.com', role: 'Frontend Developer', salary: 6200, status: 'Active', department_id: 2, department_name: 'Engineering', join_date: '2025-04-01' },
      { id: 7, first_name: 'Sophia', last_name: 'Davis', email: 'sophia.davis@erp.com', role: 'HR Assistant', salary: 4000, status: 'Active', department_id: 1, department_name: 'Human Resources', join_date: '2025-05-15' },
      { id: 8, first_name: 'James', last_name: 'Taylor', email: 'james.taylor@erp.com', role: 'DevOps Specialist', salary: 8000, status: 'Suspended', department_id: 2, department_name: 'Engineering', join_date: '2024-08-01' }
    ]
  }
}

// Client-side Searching and Filtering
const filteredEmployees = computed(() => {
  return employees.value.filter(emp => {
    const fullName = `${emp.first_name} ${emp.last_name}`.toLowerCase()
    const matchesSearch = searchQuery.value ? (
      fullName.includes(searchQuery.value.toLowerCase()) ||
      emp.role.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      `emp-${String(emp.id).padStart(4, '0')}`.includes(searchQuery.value.toLowerCase())
    ) : true

    const matchesDept = selectedDept.value ? emp.department_id === parseInt(selectedDept.value) : true

    return matchesSearch && matchesDept
  })
})

// Pagination Calculations
const totalEntries = computed(() => filteredEmployees.value.length)
const totalPages = computed(() => Math.ceil(totalEntries.value / itemsPerPage.value))
const entriesStart = computed(() => {
  if (totalEntries.value === 0) return 0
  return (currentPage.value - 1) * itemsPerPage.value + 1
})
const entriesEnd = computed(() => {
  const end = currentPage.value * itemsPerPage.value
  return end > totalEntries.value ? totalEntries.value : end
})

const paginatedEmployees = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredEmployees.value.slice(start, end)
})

const resetToFirstPage = () => {
  currentPage.value = 1
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

// Modal open/close actions
const openAddModal = () => {
  editMode.value = false
  currentEditId.value = null
  form.value = {
    first_name: '',
    last_name: '',
    email: '',
    role: '',
    salary: 5000,
    department_id: departments.value.length > 0 ? departments.value[0].id : null,
    join_date: new Date().toISOString().split('T')[0],
    status: 'Active'
  }
  showModal.value = true
}

const openEditModal = (emp) => {
  editMode.value = true
  currentEditId.value = emp.id
  form.value = {
    first_name: emp.first_name,
    last_name: emp.last_name,
    email: emp.email,
    role: emp.role,
    salary: emp.salary,
    department_id: emp.department_id,
    join_date: emp.join_date ? new Date(emp.join_date).toISOString().split('T')[0] : '',
    status: emp.status
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const submitForm = async () => {
  try {
    let url = `${API_URL}/api/employees`
    let method = 'POST'
    if (editMode.value) {
      url += `?id=${currentEditId.value}`
      method = 'PUT'
    }

    const payload = { ...form.value }
    if (payload.department_id) {
      payload.department_id = parseInt(payload.department_id)
    }

    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    if (res.ok) {
      closeModal()
      await fetchEmployees()
      await fetchStats()
    } else {
      const errMsg = await res.text()
      alert(`Failed to register employee: ${errMsg || res.statusText}`)
    }
  } catch (err) {
    console.error('Error submitting form: ', err)
    // Fallback simulation
    if (editMode.value) {
      const idx = employees.value.findIndex(e => e.id === currentEditId.value)
      if (idx !== -1) {
        employees.value[idx] = {
          ...employees.value[idx],
          ...form.value,
          department_name: departments.value.find(d => d.id == form.value.department_id)?.name || 'Unassigned'
        }
      }
    } else {
      employees.value.push({
        id: Math.max(...employees.value.map(e => e.id), 0) + 1,
        ...form.value,
        department_name: departments.value.find(d => d.id == form.value.department_id)?.name || 'Unassigned'
      })
    }
    closeModal()
    fetchStats()
  }
}

const deleteEmployee = async (id) => {
  if (!confirm('Are you sure you want to remove this employee from payroll files?')) return
  try {
    const res = await fetch(`${API_URL}/api/employees?id=${id}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      await fetchEmployees()
      await fetchStats()
    }
  } catch (err) {
    employees.value = employees.value.filter(e => e.id !== id)
    fetchStats()
  }
}

onMounted(() => {
  fetchDepartments().then(() => {
    fetchEmployees().then(fetchStats)
  })
})
</script>

<style scoped>
.employees-view {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Stats Row Styling matching screenshots */
.stats-cards-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
}

.metric-card {
  background: white;
  padding: 1.25rem 1.5rem;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  display: flex;
  align-items: center;
  gap: 1.25rem;
  border: 1px solid var(--border-color);
  border-left-width: 4px; /* Left colored border accent */
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.metric-card.blue-border { border-left-color: #3b82f6; }
.metric-card.green-border { border-left-color: #10b981; }
.metric-card.yellow-border { border-left-color: #f59e0b; }

.metric-icon-bg {
  width: 46px;
  height: 46px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
}

.metric-icon-bg.blue { background-color: #eff6ff; color: #3b82f6; }
.metric-icon-bg.green { background-color: #ecfdf5; color: #10b981; }
.metric-icon-bg.yellow { background-color: #fef3c7; color: #f59e0b; }

.metric-info {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.metric-label {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-muted);
  letter-spacing: 0.05em;
}

.metric-value {
  font-family: var(--font-display);
  font-size: 1.7rem;
  font-weight: 800;
  color: var(--text-primary);
  line-height: 1;
}

/* Filters Control Bar (matching layout screenshot) */
.table-control-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  background-color: white;
}

.search-input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.magnifier-icon {
  position: absolute;
  left: 14px;
  color: #3b82f6; /* Accent color matching search icon */
  pointer-events: none;
}

.custom-search-input {
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 0.65rem 1rem 0.65rem 2.8rem;
  border-radius: var(--radius-sm);
  width: 320px;
  font-family: var(--font-body);
  font-size: 0.85rem;
  color: var(--text-primary);
  transition: all 0.2s ease;
}

.custom-search-input:focus {
  outline: none;
  background-color: white;
  border-color: #cbd5e1;
  box-shadow: var(--shadow-sm);
}

.control-actions-group {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.custom-select {
  background-color: white;
  border: 1px solid #e2e8f0;
  padding: 0.65rem 1rem;
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-size: 0.85rem;
  color: var(--text-secondary);
  width: 180px;
  cursor: pointer;
}

.custom-select:focus {
  outline: none;
  border-color: #cbd5e1;
}

.add-employee-btn {
  background-color: #3b82f6; /* Blue button color */
  color: white;
  border-radius: var(--radius-sm);
  padding: 0.65rem 1.25rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.add-employee-btn:hover {
  background-color: #2563eb;
}

/* Rounded strip headers table */
.custom-table-container {
  background: white;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.employees-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.85rem;
}

.header-strip {
  background-color: #f8fafc; /* Rounded strip bg color */
}

.employees-table th {
  padding: 1rem 1.25rem;
  font-weight: 700;
  color: var(--text-secondary);
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  border-bottom: 1px solid var(--border-color);
}

.employees-table td {
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}

.employees-table tr:last-child td {
  border-bottom: none;
}

.table-row:hover td {
  background-color: #f8fafc;
}

/* Inner elements styling */
.id-cell {
  font-family: var(--font-display);
  font-weight: 700;
  color: var(--text-secondary);
}

.name-cell-wrapper {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.avatar-circle {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background-color: #f1f5f9;
  color: #3b82f6;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-color);
}

.full-name {
  font-weight: 600;
  color: var(--text-primary);
}

.dept-bubble {
  background-color: #f1f5f9;
  border: 1px solid #e2e8f0;
  padding: 0.2rem 0.6rem;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-secondary);
}

.position-cell {
  font-weight: 500;
  color: var(--text-primary);
}

.actions-group {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.action-icon {
  background: none;
  border: 1px solid #e2e8f0;
  width: 30px;
  height: 30px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-icon.edit-btn { color: #3b82f6; }
.action-icon.edit-btn:hover { background-color: #eff6ff; border-color: #93c5fd; }

.action-icon.delete-btn { color: #ef4444; }
.action-icon.delete-btn:hover { background-color: #fef2f2; border-color: #fca5a5; }

.no-records-cell {
  text-align: center;
  padding: 3rem;
  color: var(--text-muted);
  font-weight: 500;
}

.text-right {
  text-align: right;
}

/* Pagination Row Styling matching screenshot */
.table-footer-pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
}

.footer-entries-info {
  font-size: 0.8rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.pagination-controls-group {
  display: flex;
  gap: 0.35rem;
  align-items: center;
}

.pag-btn {
  background-color: white;
  border: 1px solid #e2e8f0;
  color: var(--text-secondary);
  font-family: var(--font-body);
  font-size: 0.8rem;
  font-weight: 500;
  padding: 0.45rem 0.85rem;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pag-btn:hover:not(:disabled) {
  background-color: #f8fafc;
  color: var(--text-primary);
  border-color: #cbd5e1;
}

.pag-btn:disabled {
  color: var(--text-muted);
  background-color: #f8fafc;
  cursor: not-allowed;
  opacity: 0.7;
}

.pag-btn.page-num.active {
  background-color: #3b82f6; /* Blue page background */
  color: white;
  border-color: #3b82f6;
  font-weight: 700;
}

/* Modals */
.form-row {
  display: flex;
  gap: 1rem;
}

.form-group.half {
  flex: 1;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
}

@media (max-width: 768px) {
  .stats-cards-row {
    grid-template-columns: 1fr;
  }
  .table-control-bar {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }
  .custom-search-input {
    width: 100%;
  }
  .control-actions-group {
    flex-direction: column;
    align-items: stretch;
  }
  .custom-select {
    width: 100%;
  }
  .table-footer-pagination {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>
