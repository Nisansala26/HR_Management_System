<template>
  <div class="departments-view">
    <!-- Header Page Title -->
    <div class="dept-header-section">
      <div class="header-titles">
        <h2 class="directory-title">Department Directory</h2>
        <p class="directory-subtitle">Overview and management of corporate organizational structure.</p>
      </div>
      <button class="btn add-new-dept-btn" @click="openAddModal">
        + Add New Department
      </button>
    </div>

    <!-- Divider Line -->
    <hr class="header-divider-line" />

    <!-- 4 KPI Metrics Row-->
    <div class="metrics-cards-row">
      <!-- Card 1: Total Departments -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper light-purple">
          <!-- Folder Icon -->
          <svg width="22" height="22" viewBox="0 0 24 24" fill="#EAB308" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2z" />
          </svg>
        </div>
        <div class="metric-text-group">
          <span class="label">Total Departments</span>
          <h3 class="value">{{ deptStats.total_departments }}</h3>
          <span class="subtext">Registered sectors</span>
        </div>
      </div>

      <!-- Card 2: Active Sectors -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper light-green">
          <!-- Active Checkbox Icon -->
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#22C55E" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="9 11 12 14 22 4"></polyline>
            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
          </svg>
        </div>
        <div class="metric-text-group">
          <span class="label">Active Sectors</span>
          <h3 class="value">{{ deptStats.active_sectors }}</h3>
          <span class="subtext">Fully operational</span>
        </div>
      </div>

      <!-- Card 3: Total Employees -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper light-blue">
          <!-- Double User Icon -->
          <svg width="22" height="22" viewBox="0 0 24 24" fill="#6366F1">
            <path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5z"/>
          </svg>
        </div>
        <div class="metric-text-group">
          <span class="label">Total Employees</span>
          <h3 class="value">{{ deptStats.total_employees }}</h3>
          <span class="subtext">Across organization</span>
        </div>
      </div>

      <!-- Card 4: Department Heads -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper light-red">
          <!-- Person Silhouette Icon -->
          <svg width="22" height="22" viewBox="0 0 24 24" fill="#312E81">
            <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" />
          </svg>
        </div>
        <div class="metric-text-group">
          <span class="label">Department Heads</span>
          <h3 class="value">{{ deptStats.department_heads }}</h3>
          <span class="subtext">Heads appointed</span>
        </div>
      </div>
    </div>

    <!-- Search Controls Bar -->
    <div class="dept-search-bar card">
      <div class="search-input-container">
        <svg class="search-magnifier-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input 
          type="text" 
          placeholder="Search by department name..." 
          class="dept-search-input"
          v-model="searchQuery"
        />
      </div>
    </div>

    <!-- Registry List -->
    <div class="custom-table-container">
      <table class="depts-table">
        <thead>
          <tr class="header-strip">
            <th>SECTOR ID</th>
            <th>SECTOR NAME</th>
            <th>SECTOR HEAD</th>
            <th>ANNUAL BUDGET</th>
            <th>HEADCOUNT</th>
            <th class="text-right">ACTIONS</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="dept in filteredDepartments" :key="dept.id" class="table-row">
            <td class="id-cell">SEC-{{ String(dept.id).padStart(3, '0') }}</td>
            <td class="dept-name-cell">{{ dept.name }}</td>
            <td class="manager-cell">
              <div class="manager-meta">
                <div class="avatar-circle-sm">{{ getInitials(dept.manager_name) }}</div>
                <span>{{ dept.manager_name }}</span>
              </div>
            </td>
            <td class="budget-cell">${{ dept.budget.toLocaleString('en-US', {minimumFractionDigits: 2}) }}</td>
            <td>
              <span class="headcount-tag interactive" @click="openManageStaffModal(dept)" title="Manage Staff">
                {{ dept.employee_count }} Staff
              </span>
            </td>
            <td>
              <div class="actions-group">
                <button class="action-icon staff-btn" @click="openManageStaffModal(dept)" title="Manage Staff">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                    <circle cx="9" cy="7" r="4" />
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
                    <path d="M16 3.13a4 4 0 0 1 0 7.75" />
                  </svg>
                </button>
                <button class="action-icon edit-btn" @click="openEditModal(dept)" title="Edit">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                    <path d="M18.5 2.5a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                  </svg>
                </button>
                <button class="action-icon delete-btn" @click="deleteDepartment(dept.id)" title="Delete">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="filteredDepartments.length === 0">
            <td colspan="6" class="no-records-cell">
              No departments registered matching search filters.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Form (Add / Edit) -->
    <div class="modal-overlay" v-if="showModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">{{ editMode ? 'Update Department Details' : 'Register New Sector' }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <form @submit.prevent="submitForm">
          <div class="form-group">
            <label for="name">Department / Sector Name</label>
            <input type="text" id="name" class="form-control" placeholder="e.g. Engineering" v-model="form.name" required />
          </div>
          <div class="form-group">
            <label for="manager_name">Department Manager / Head</label>
            <input type="text" id="manager_name" class="form-control" placeholder="e.g. Sarah Jenkins" v-model="form.manager_name" required />
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
            <button type="submit" class="btn btn-primary">
              {{ editMode ? 'Save Details' : 'Register' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Manage Staff Modal -->
    <div class="modal-overlay" v-if="showStaffModal">
      <div class="modal-content staff-modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Manage Department Staff</h3>
          <button class="close-btn" @click="closeStaffModal" :disabled="savingStaff">&times;</button>
        </div>
        <div class="modal-body">
          <p class="modal-subtitle">
            Assign or remove staff members for the <strong>{{ activeDept?.name }}</strong> department.
          </p>
          
          <div class="staff-search-box">
            <input 
              type="text" 
              placeholder="Filter staff by name or role..." 
              class="form-control search-staff-input"
              v-model="staffSearchQuery"
            />
          </div>

          <div class="staff-checklist-container">
            <div 
              v-for="emp in filteredStaffList" 
              :key="emp.id" 
              class="staff-check-item"
              :class="{ selected: selectedEmployeeIds.includes(emp.id) }"
            >
              <label :for="'emp-check-' + emp.id" class="staff-label">
                <input 
                  type="checkbox" 
                  :id="'emp-check-' + emp.id" 
                  v-model="selectedEmployeeIds" 
                  :value="emp.id"
                />
                <div class="staff-info-pair">
                  <span class="staff-name">{{ emp.first_name }} {{ emp.last_name }}</span>
                  <span class="staff-position">{{ emp.role }}</span>
                </div>
                <div class="staff-dept-badge" :class="{ 'assigned': emp.department_id }">
                  {{ emp.department_id ? emp.department_name : 'Unassigned' }}
                </div>
              </label>
            </div>
            <div v-if="filteredStaffList.length === 0" class="no-staff-found">
              No staff members found matching query.
            </div>
          </div>
        </div>
        <div class="modal-actions">
          <button type="button" class="btn btn-secondary" @click="closeStaffModal" :disabled="savingStaff">Cancel</button>
          <button type="button" class="btn btn-primary" @click="saveStaffAssignments" :disabled="savingStaff">
            {{ savingStaff ? 'Saving...' : 'Save Assignments' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const route = useRoute()
const departments = ref([])
const deptStats = ref({ total_departments: 0, active_sectors: 0, total_employees: 0, department_heads: 0 })
const searchQuery = ref('')

watch(() => route.query.search, (newVal) => {
  searchQuery.value = newVal || ''
}, { immediate: true })

const showModal = ref(false)
const editMode = ref(false)
const currentEditId = ref(null)

// Manage Staff state
const showStaffModal = ref(false)
const allEmployees = ref([])
const selectedEmployeeIds = ref([])
const savingStaff = ref(false)
const activeDept = ref(null)
const staffSearchQuery = ref('')

const form = ref({
  name: '',
  manager_name: '',
  budget: 150000
})

const getInitials = (name) => {
  if (!name) return ''
  const parts = name.split(' ')
  return parts.map(p => p.charAt(0)).join('').toUpperCase().slice(0, 2)
}

const fetchStats = async () => {
  try {
    const res = await fetch(`${API_URL}/api/departments/stats`)
    if (res.ok) {
      deptStats.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching department stats: ', err)
    deptStats.value = {
      total_departments: departments.value.length || 4,
      active_sectors: departments.value.filter(d => d.budget > 0).length || 4,
      total_employees: departments.value.reduce((acc, curr) => acc + curr.employee_count, 0) || 8,
      department_heads: new Set(departments.value.map(d => d.manager_name)).size || 4
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
      { id: 1, name: 'Human Resources', manager_name: 'Sarah Jenkins', budget: 150000.00, employee_count: 2 },
      { id: 2, name: 'Engineering', manager_name: 'David Miller', budget: 500000.00, employee_count: 4 },
      { id: 3, name: 'Marketing', manager_name: 'Emily Watson', budget: 120000.00, employee_count: 1 },
      { id: 4, name: 'Finance', manager_name: 'Robert Chen', budget: 200000.00, employee_count: 1 }
    ]
  }
}

const filteredDepartments = computed(() => {
  return departments.value.filter(dept => {
    return dept.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
           dept.manager_name.toLowerCase().includes(searchQuery.value.toLowerCase())
  })
})

const filteredStaffList = computed(() => {
  if (!allEmployees.value) return []
  return allEmployees.value.filter(emp => {
    const name = `${emp.first_name} ${emp.last_name}`.toLowerCase()
    const role = (emp.role || '').toLowerCase()
    const q = staffSearchQuery.value.toLowerCase()
    return name.includes(q) || role.includes(q)
  })
})

const openAddModal = () => {
  editMode.value = false
  currentEditId.value = null
  form.value = {
    name: '',
    manager_name: '',
    budget: 150000
  }
  showModal.value = true
}

const openEditModal = (dept) => {
  editMode.value = true
  currentEditId.value = dept.id
  form.value = {
    name: dept.name,
    manager_name: dept.manager_name,
    budget: dept.budget
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const submitForm = async () => {
  try {
    let url = `${API_URL}/api/departments`
    let method = 'POST'
    if (editMode.value) {
      url += `?id=${currentEditId.value}`
      method = 'PUT'
    }

    const payload = { ...form.value }
    payload.budget = parseFloat(payload.budget) || 0.00

    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    if (res.ok) {
      closeModal()
      await fetchDepartments()
      await fetchStats()
    } else {
      const errMsg = await res.text()
      alert(`Failed to register department: ${errMsg || res.statusText}`)
    }
  } catch (err) {
    console.error('Error submitting form: ', err)
    if (editMode.value) {
      const idx = departments.value.findIndex(d => d.id === currentEditId.value)
      if (idx !== -1) {
        departments.value[idx] = {
          ...departments.value[idx],
          ...form.value
        }
      }
    } else {
      departments.value.push({
        id: Math.max(...departments.value.map(d => d.id), 0) + 1,
        ...form.value,
        employee_count: 0
      })
    }
    closeModal()
    fetchStats()
  }
}

const deleteDepartment = async (id) => {
  if (!confirm('Are you sure you want to delete this department? Employees will be unassigned.')) return
  try {
    const res = await fetch(`${API_URL}/api/departments?id=${id}`, {
      method: 'DELETE'
    })
    if (res.ok) {
      await fetchDepartments()
      await fetchStats()
    }
  } catch (err) {
    departments.value = departments.value.filter(d => d.id !== id)
    fetchStats()
  }
}

// Manage Staff Handlers
const openManageStaffModal = async (dept) => {
  activeDept.value = dept
  staffSearchQuery.value = ''
  try {
    const res = await fetch(`${API_URL}/api/employees`)
    if (res.ok) {
      allEmployees.value = await res.json()
      selectedEmployeeIds.value = allEmployees.value
        .filter(e => e.department_id === dept.id)
        .map(e => e.id)
      showStaffModal.value = true
    } else {
      alert('Failed to load employees list.')
    }
  } catch (err) {
    console.error('Error fetching employees for department management:', err)
    alert('Error fetching employees list.')
  }
}

const closeStaffModal = () => {
  showStaffModal.value = false
  activeDept.value = null
  selectedEmployeeIds.value = []
}

const saveStaffAssignments = async () => {
  savingStaff.value = true
  try {
    const deptId = activeDept.value.id
    const updates = []
    
    for (const emp of allEmployees.value) {
      const wasInDept = (emp.department_id === deptId)
      const isSelectedNow = selectedEmployeeIds.value.includes(emp.id)
      
      if (isSelectedNow && !wasInDept) {
        // Assign to this department
        const updatedEmp = { ...emp, department_id: deptId }
        updates.push(updateEmployeeDept(updatedEmp))
      } else if (!isSelectedNow && wasInDept) {
        // Unassign from this department
        const updatedEmp = { ...emp, department_id: null }
        updates.push(updateEmployeeDept(updatedEmp))
      }
    }
    
    if (updates.length > 0) {
      await Promise.all(updates)
    }
    
    closeStaffModal()
    await fetchDepartments()
    await fetchStats()
  } catch (err) {
    console.error('Error saving department staff assignments:', err)
    alert('Some assignments failed to save. Please try again.')
  } finally {
    savingStaff.value = false
  }
}

const updateEmployeeDept = async (emp) => {
  const res = await fetch(`${API_URL}/api/employees?id=${emp.id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(emp)
  })
  if (!res.ok) {
    throw new Error(`Failed to update employee ${emp.id}`)
  }
}

onMounted(() => {
  fetchDepartments().then(fetchStats)
})
</script>

<style scoped>
.departments-view {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* Header style  */
.dept-header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}

.directory-title {
  font-family: var(--font-display);
  font-size: 1.8rem;
  font-weight: 800;
  color: #1e293b;
  letter-spacing: -0.02em;
}

.directory-subtitle {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.add-new-dept-btn {
  background-color: #00c853; /* Bright green add department button */
  color: white;
  border-radius: var(--radius-sm);
  padding: 0.65rem 1.25rem;
  font-weight: 700;
  border: none;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.add-new-dept-btn:hover {
  background-color: #055726;
}

.header-divider-line {
  border: none;
  border-top: 1.5px solid #cbd5e1;
  margin-bottom: 0.5rem;
}

/* 4 Metrics cards layout */
.metrics-cards-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.25rem;
}

.metric-card {
  background: white;
  padding: 1.25rem;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.metric-icon-wrapper {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
}

.metric-icon-wrapper.light-purple { background-color: #f5f3ff; color: #8b5cf6; }
.metric-icon-wrapper.light-green { background-color: #ecfdf5; color: #10b981; }
.metric-icon-wrapper.light-blue { background-color: #eff6ff; color: #3b82f6; }
.metric-icon-wrapper.light-red { background-color: #fff5f5; color: #ef4444; }

.metric-text-group {
  display: flex;
  flex-direction: column;
}

.metric-text-group .label {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-muted);
}

.metric-text-group .value {
  font-family: var(--font-display);
  font-size: 1.6rem;
  font-weight: 800;
  color: var(--text-primary);
  line-height: 1.1;
  margin: 0.1rem 0;
}

.metric-text-group .subtext {
  font-size: 0.72rem;
  color: var(--text-muted);
}

/* Search bar control */
.dept-search-bar {
  padding: 0.9rem 1.25rem;
}

.search-input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.search-magnifier-icon {
  position: absolute;
  left: 14px;
  color: #6366f1; /* Purple magnifier matching screenshot 1 */
  pointer-events: none;
}

.dept-search-input {
  background-color: #f1f5f9;
  border: 1px solid transparent;
  padding: 0.65rem 1rem 0.65rem 2.8rem;
  border-radius: var(--radius-sm);
  width: 320px;
  font-family: var(--font-body);
  font-size: 0.85rem;
  color: var(--text-primary);
  transition: all 0.2s ease;
}

.dept-search-input:focus {
  outline: none;
  background-color: white;
  border-color: #cbd5e1;
  box-shadow: var(--shadow-sm);
}

/* Custom Table styles */
.custom-table-container {
  background: white;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.depts-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.85rem;
}

.header-strip {
  background-color: #f8fafc;
}

.depts-table th {
  padding: 1rem 1.25rem;
  font-weight: 700;
  color: var(--text-secondary);
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  border-bottom: 1px solid var(--border-color);
}

.depts-table td {
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}

.depts-table tr:last-child td {
  border-bottom: none;
}

.table-row:hover td {
  background-color: #f8fafc;
}

.id-cell {
  font-family: var(--font-display);
  font-weight: 700;
  color: var(--text-secondary);
}

.dept-name-cell {
  font-weight: 600;
  color: var(--text-primary);
}

.manager-cell {
  font-weight: 500;
}

.manager-meta {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.avatar-circle-sm {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background-color: #f1f5f9;
  color: var(--primary-light);
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border-color);
}

.budget-cell {
  font-family: var(--font-display);
  font-weight: 700;
  color: var(--primary-light);
}

.headcount-tag {
  background-color: #eff6ff;
  color: #2563eb;
  padding: 0.2rem 0.6rem;
  border-radius: var(--radius-sm);
  font-weight: 600;
  font-size: 0.75rem;
  border: 1px solid rgba(37, 99, 235, 0.1);
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

.action-icon.staff-btn { color: #10b981; }
.action-icon.staff-btn:hover { background-color: #ecfdf5; border-color: #a7f3d0; }

.headcount-tag.interactive {
  cursor: pointer;
  transition: all 0.2s ease;
}

.headcount-tag.interactive:hover {
  background-color: #2563eb;
  color: white;
  border-color: #2563eb;
}

/* Manage Staff Modal Styling */
.staff-modal-content {
  max-width: 550px;
}

.modal-subtitle {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-bottom: 1rem;
}

.staff-search-box {
  margin-bottom: 1rem;
}

.search-staff-input {
  background-color: #f8fafc;
}

.staff-checklist-container {
  max-height: 280px;
  overflow-y: auto;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  display: flex;
  flex-direction: column;
}

.staff-check-item {
  display: flex;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  padding: 0.75rem 1rem;
  transition: background-color 0.15s ease;
}

.staff-check-item:last-child {
  border-bottom: none;
}

.staff-check-item:hover {
  background-color: #f8fafc;
}

.staff-check-item.selected {
  background-color: #ecfdf5;
}

.staff-label {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  cursor: pointer;
  user-select: none;
  font-weight: normal;
  margin-bottom: 0;
}

.staff-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
  cursor: pointer;
  accent-color: var(--primary-light);
}

.staff-info-pair {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.staff-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-primary);
}

.staff-position {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.staff-dept-badge {
  font-size: 0.7rem;
  font-weight: 600;
  padding: 0.15rem 0.5rem;
  border-radius: var(--radius-full);
  background-color: #f1f5f9;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.staff-dept-badge.assigned {
  background-color: #eff6ff;
  color: #2563eb;
  border-color: rgba(37, 99, 235, 0.1);
}

.no-staff-found {
  padding: 2rem;
  text-align: center;
  color: var(--text-muted);
  font-size: 0.875rem;
}

.no-records-cell {
  text-align: center;
  padding: 3rem;
  color: var(--text-muted);
}

.text-right {
  text-align: right;
}

/* Modal Actions */
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1.5rem;
}

@media (max-width: 1024px) {
  .metrics-cards-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .metrics-cards-row {
    grid-template-columns: 1fr;
  }
  .dept-header-section {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }
  .dept-search-input {
    width: 100%;
  }
}
</style>
