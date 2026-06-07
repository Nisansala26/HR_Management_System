<template>
  <div class="attendance-view">
    <!-- Header Section -->
    <div class="attendance-header">
      <h2 class="page-title">Attendance</h2>
      <p class="breadcrumbs-path">Dashboard > <span class="active-crumb">Attendance</span></p>
    </div>

    <!-- 4 KPI Metrics Row (Exactly replicating screenshot 2 cards with heart icons and dots) -->
    <div class="metrics-cards-row">
      <!-- Card 1: Present Today -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper circle-green">
          <!-- Heart Icon -->
          <svg width="20" height="20" viewBox="0 0 24 24" fill="#10B981" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
        </div>
        <div class="metric-text-group">
          <span class="label">PRESENT TODAY</span>
          <h3 class="value">{{ attendanceStats.present_today }}</h3>
          <span class="pct-dot-label text-success">
            <span class="dot-marker">•</span> {{ attendanceStats.present_pct.toFixed(0) }}%
          </span>
        </div>
      </div>

      <!-- Card 2: Absent Today -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper circle-red">
          <!-- Heart Icon -->
          <svg width="20" height="20" viewBox="0 0 24 24" fill="#EF4444" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
        </div>
        <div class="metric-text-group">
          <span class="label">ABSENT TODAY</span>
          <h3 class="value">{{ attendanceStats.absent_today }}</h3>
          <span class="pct-dot-label text-danger">
            <span class="dot-marker">•</span> {{ attendanceStats.absent_pct.toFixed(0) }}%
          </span>
        </div>
      </div>

      <!-- Card 3: Late Comers -->
      <div class="card metric-card">
        <div class="metric-icon-wrapper circle-yellow">
          <!-- Heart Icon -->
          <svg width="20" height="20" viewBox="0 0 24 24" fill="#F59E0B" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
        </div>
        <div class="metric-text-group">
          <span class="label">LATE COMERS</span>
          <h3 class="value">{{ attendanceStats.late_comers }}</h3>
          <span class="pct-dot-label text-warning">
            <span class="dot-marker">•</span> {{ attendanceStats.late_pct.toFixed(0) }}%
          </span>
        </div>
      </div>

    </div>

    <!-- Attendance split layout section (65/35 grid columns) -->
    <div class="attendance-content-split">
      
      <!-- Left Column: Log Registry -->
      <div class="left-split">
        <div class="card logs-console-card">
          
          <!-- Filters Console (Exactly matching screenshot 2 filters layout) -->
          <div class="filters-console-row">
            <div class="filters-left-group">
              <!-- Date button input -->
              <div class="date-selector-btn-wrapper">
                <span class="calendar-icon-prefix">🗓️</span>
                <input type="date" class="date-picker-element" v-model="filterDate" @change="onDateChange" />
              </div>

              <!-- Department Selector -->
              <select class="custom-console-select" v-model="selectedDept" @change="fetchAttendance">
                <option value="">All Departments</option>
                <option v-for="dept in departments" :key="dept.id" :value="dept.id">
                  {{ dept.name }}
                </option>
              </select>

              <!-- Status Selector -->
              <select class="custom-console-select" v-model="selectedStatus" @change="fetchAttendance">
                <option value="">All Status</option>
                <option value="Present">Present</option>
                <option value="Late">Late</option>
                <option value="Absent">Absent</option>
              </select>
            </div>

            <div class="filters-right-group">
              <button class="btn btn-primary add-attendance-btn" @click="openClockModal">
                + Add Attendance
              </button>
            </div>
          </div>

          <!-- Tabular logs -->
          <div class="custom-table-container">
            <table class="attendance-table">
              <thead>
                <tr class="header-strip">
                  <th>ID</th>
                  <th>Employee</th>
                  <th>Department</th>
                  <th>Check In</th>
                  <th>Check Out</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="log in filteredLogs" :key="log.id" class="table-row">
                  <td class="id-cell">#{{ String(log.id).padStart(4, '0') }}</td>
                  <td>
                    <span class="full-name">{{ log.employee_name }}</span>
                  </td>
                  <td>{{ log.department_name || 'Unassigned' }}</td>
                  <td class="time-cell">{{ formatTime(log.clock_in) }}</td>
                  <td class="time-cell">{{ formatTime(log.clock_out) }}</td>
                  <td>
                    <span class="badge" :class="'badge-' + log.status.toLowerCase()">
                      {{ log.status }}
                    </span>
                  </td>
                </tr>
                <tr v-if="filteredLogs.length === 0">
                  <td colspan="6" class="no-records-cell">
                    No results found for current filters.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Entries Footer / pagination -->
          <div class="logs-footer-row">
            <span class="results-info">Showing 1 to {{ filteredLogs.length }} of {{ filteredLogs.length }} results</span>
            <div class="simple-pagination-circles">
              <button class="circle-pag-btn active">1</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Visual Charts Overview panels -->
      <div class="right-split">
        <!-- Panel 1: Donut Representation overview -->
        <div class="card visual-metrics-card">
          <h3 class="card-title">Attendance Overview</h3>
          <div class="donut-chart-container">
            <!-- Dynamic SVG Donut Chart representing active stats -->
            <svg class="donut-svg" width="130" height="130" viewBox="0 0 40 40">
              <circle class="donut-hole" cx="20" cy="20" r="15.91549430918954" fill="transparent"></circle>
              <circle class="donut-ring" cx="20" cy="20" r="15.91549430918954" fill="transparent" stroke="#e2e8f0" stroke-width="4.5"></circle>
              
              <!-- Present (On Time) segment -->
              <circle 
                class="donut-segment" 
                cx="20" 
                cy="20" 
                r="15.91549430918954" 
                fill="transparent" 
                stroke="#10B981" 
                stroke-width="4.5"
                :stroke-dasharray="`${Math.max(0, attendanceStats.present_pct - attendanceStats.late_pct)} ${100 - Math.max(0, attendanceStats.present_pct - attendanceStats.late_pct)}`"
                stroke-dashoffset="25"
              ></circle>
              
              <!-- Late segment offset -->
              <circle 
                class="donut-segment" 
                cx="20" 
                cy="20" 
                r="15.91549430918954" 
                fill="transparent" 
                stroke="#F59E0B" 
                stroke-width="4.5"
                :stroke-dasharray="`${attendanceStats.late_pct} ${100 - attendanceStats.late_pct}`"
                :stroke-dashoffset="125 - Math.max(0, attendanceStats.present_pct - attendanceStats.late_pct)"
              ></circle>

              <!-- Absent segment offset -->
              <circle 
                class="donut-segment" 
                cx="20" 
                cy="20" 
                r="15.91549430918954" 
                fill="transparent" 
                stroke="#EF4444" 
                stroke-width="4.5"
                :stroke-dasharray="`${attendanceStats.absent_pct} ${100 - attendanceStats.absent_pct}`"
                :stroke-dashoffset="125 - attendanceStats.present_pct"
              ></circle>
            </svg>

            <!-- Legend labels -->
            <div class="donut-legend-group">
              <div class="legend-row">
                <span class="legend-color-box present"></span>
                <span class="legend-text">Present</span>
              </div>
              <div class="legend-row">
                <span class="legend-color-box absent"></span>
                <span class="legend-text">Absent</span>
              </div>
              <div class="legend-row">
                <span class="legend-color-box late"></span>
                <span class="legend-text">Late</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Panel 2: Department wise progress bars -->
        <div class="card department-progress-card">
          <h3 class="card-title">Department Wise (Present %)</h3>
          <div class="dept-progress-list">
            <!-- Dynamic vertical/horizontal progress lists -->
            <div v-for="dept in deptProgressList" :key="dept.name" class="progress-bar-row">
              <div class="progress-bar-meta">
                <span class="dept-title-label">{{ dept.name }}</span>
                <span class="dept-pct-label">{{ dept.pct }}%</span>
              </div>
              <div class="progress-outer-frame">
                <div class="progress-inner-bar" :style="{ width: dept.pct + '%', backgroundColor: dept.color }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>

    <!-- Clock / Add Attendance Modal -->
    <div class="modal-overlay" v-if="showClockModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Record Attendance Log</h3>
          <button class="close-btn" @click="closeClockModal">&times;</button>
        </div>
        <form @submit.prevent="submitAttendanceLog">
          <div class="form-group">
            <label for="emp-select">Select Employee</label>
            <select id="emp-select" class="form-control" v-model="clockForm.employee_id" required>
              <option v-for="emp in employees" :key="emp.id" :value="emp.id">
                {{ emp.first_name }} {{ emp.last_name }} ({{ emp.role }})
              </option>
            </select>
          </div>
          <div class="form-row">
            <div class="form-group half">
              <label for="clock-in-time">Check In Time</label>
              <input type="time" id="clock-in-time" class="form-control" v-model="clockForm.clock_in" required />
            </div>
            <div class="form-group half">
              <label for="clock-out-time">Check Out Time</label>
              <input type="time" id="clock-out-time" class="form-control" v-model="clockForm.clock_out" />
            </div>
          </div>
          <div class="form-group">
            <label for="status-select">Status</label>
            <select id="status-select" class="form-control" v-model="clockForm.status">
              <option value="Present">Present (On Time)</option>
              <option value="Late">Late Check-in</option>
              <option value="Absent">Absent</option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="closeClockModal">Cancel</button>
            <button type="submit" class="btn btn-primary">Save Log</button>
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
const searchQuery = ref('')
watch(() => route.query.search, (newVal) => {
  searchQuery.value = newVal || ''
}, { immediate: true })

const filterDate = ref(new Date().toISOString().split('T')[0])
const selectedDept = ref('')
const selectedStatus = ref('')

const logs = ref([])
const departments = ref([])
const employees = ref([])

const showClockModal = ref(false)
const clockForm = ref({
  employee_id: 1,
  clock_in: '09:00',
  clock_out: '17:00',
  status: 'Present'
})

// Dynamic calculated attendance stats computed property
const attendanceStats = computed(() => {
  const activeEmps = employees.value.filter(e => {
    const matchesDept = !selectedDept.value || e.department_id === parseInt(selectedDept.value)
    return e.status === 'Active' && matchesDept
  })

  const deptLogs = logs.value.filter(log => {
    return !selectedDept.value || log.department_id === parseInt(selectedDept.value)
  })

  const checkedInNonActive = deptLogs.filter(log => {
    return !activeEmps.some(e => e.id === log.employee_id)
  })

  const total = activeEmps.length + checkedInNonActive.length

  const present = deptLogs.filter(l => l.status === 'Present').length
  const late = deptLogs.filter(l => l.status === 'Late').length
  
  const presentToday = present + late
  const absent = Math.max(0, total - presentToday)

  return {
    present_today: presentToday,
    absent_today: absent,
    late_comers: late,
    total_employees: total,
    present_pct: total > 0 ? (presentToday / total) * 100 : 0,
    absent_pct: total > 0 ? (absent / total) * 100 : 0,
    late_pct: total > 0 ? (late / total) * 100 : 0
  }
})

// Dynamic calculated department progress percentages
const deptProgressList = computed(() => {
  return departments.value.map((dept, index) => {
    const activeDeptEmps = employees.value.filter(e => e.status === 'Active' && e.department_id === dept.id)
    const deptLogs = logs.value.filter(l => l.department_id === dept.id)
    const nonActiveCheckedIn = deptLogs.filter(log => !activeDeptEmps.some(e => e.id === log.employee_id))
    
    const total = activeDeptEmps.length + nonActiveCheckedIn.length
    const present = deptLogs.filter(l => l.status === 'Present' || l.status === 'Late').length
    
    const pct = total > 0 ? Math.round((present / total) * 100) : 0
    
    const colors = ['#3b82f6', '#8b5cf6', '#f59e0b', '#ec4899', '#10b981']
    const color = colors[index % colors.length]
    
    return {
      name: dept.name,
      pct,
      color
    }
  })
})

const formatTime = (isoStr) => {
  if (!isoStr) return '--:--'
  const date = new Date(isoStr)
  return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false })
}



const fetchDepartments = async () => {
  try {
    const res = await fetch(`${API_URL}/api/departments`)
    if (res.ok) {
      departments.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching depts: ', err)
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
  }
}

const fetchAttendance = async () => {
  try {
    const res = await fetch(`${API_URL}/api/attendance?date=${filterDate.value}`)
    if (res.ok) {
      logs.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching logs: ', err)
    // Fallback Mock
    logs.value = [
      { id: 1, employee_id: 1, employee_name: 'Admin HR Manager', clock_in: new Date(filterDate.value + 'T08:45:00').toISOString(), clock_out: new Date(filterDate.value + 'T17:30:00').toISOString(), status: 'Present', department_name: 'Human Resources', department_id: 1 },
      { id: 2, employee_id: 2, employee_name: 'John Doe', clock_in: new Date(filterDate.value + 'T09:05:00').toISOString(), clock_out: new Date(filterDate.value + 'T18:00:00').toISOString(), status: 'Present', department_name: 'Engineering', department_id: 2 },
      { id: 3, employee_id: 3, employee_name: 'Alice Smith', clock_in: new Date(filterDate.value + 'T09:35:00').toISOString(), clock_out: null, status: 'Late', department_name: 'Engineering', department_id: 2 },
      { id: 4, employee_id: 4, employee_name: 'Michael Brown', clock_in: new Date(filterDate.value + 'T08:50:00').toISOString(), clock_out: new Date(filterDate.value + 'T17:00:00').toISOString(), status: 'Present', department_name: 'Marketing', department_id: 3 }
    ]
  }
}

const onDateChange = () => {
  fetchAttendance()
}

// Client filtering
const filteredLogs = computed(() => {
  const activeEmps = employees.value.filter(e => {
    const matchesDept = !selectedDept.value || e.department_id === parseInt(selectedDept.value)
    return e.status === 'Active' && matchesDept
  })

  const actualLogsFiltered = logs.value.filter(log => {
    return !selectedDept.value || log.department_id === parseInt(selectedDept.value)
  })

  const loggedActiveEmpIds = actualLogsFiltered.map(l => l.employee_id)
  const unloggedActiveEmps = activeEmps.filter(e => !loggedActiveEmpIds.includes(e.id))

  const mockAbsentLogs = unloggedActiveEmps.map(e => {
    return {
      id: 9999 + e.id,
      employee_id: e.id,
      employee_name: `${e.first_name} ${e.last_name}`,
      department_name: e.department_name || 'Unassigned',
      department_id: e.department_id,
      date: filterDate.value,
      clock_in: null,
      clock_out: null,
      status: 'Absent'
    }
  })

  const allLogs = [...actualLogsFiltered, ...mockAbsentLogs]

  const searchedLogs = allLogs.filter(log => {
    if (!searchQuery.value) return true
    const q = searchQuery.value.toLowerCase()
    return log.employee_name.toLowerCase().includes(q) || log.department_name.toLowerCase().includes(q)
  })

  searchedLogs.sort((a, b) => {
    const statusOrder = { 'present': 1, 'late': 2, 'absent': 3 }
    const orderA = statusOrder[a.status.toLowerCase()] || 4
    const orderB = statusOrder[b.status.toLowerCase()] || 4
    if (orderA !== orderB) return orderA - orderB
    return a.employee_name.localeCompare(b.employee_name)
  })

  return searchedLogs.filter(log => {
    const matchesStatus = selectedStatus.value ? log.status.toLowerCase() === selectedStatus.value.toLowerCase() : true
    return matchesStatus
  })
})

const openClockModal = () => {
  clockForm.value = {
    employee_id: employees.value.length > 0 ? employees.value[0].id : 1,
    clock_in: '09:00',
    clock_out: '17:00',
    status: 'Present'
  }
  showClockModal.value = true
}

const closeClockModal = () => {
  showClockModal.value = false
}

const submitAttendanceLog = async () => {
  try {
    const res = await fetch(`${API_URL}/api/attendance`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        employee_id: parseInt(clockForm.value.employee_id),
        date: filterDate.value,
        clock_in: clockForm.value.clock_in,
        clock_out: clockForm.value.clock_out,
        status: clockForm.value.status
      })
    })

    if (res.ok) {
      closeClockModal()
      await fetchAttendance()
    } else {
      const errMsg = await res.text()
      alert(`Failed to save manual attendance: ${errMsg || res.statusText}`)
    }
  } catch (err) {
    console.error('Error recording attendance: ', err)
    // Local insertion fallback
    const targetEmp = employees.value.find(e => e.id == clockForm.value.employee_id)
    const existingLogIdx = logs.value.findIndex(l => l.employee_id == clockForm.value.employee_id && l.date == filterDate.value)
    
    const logData = {
      id: existingLogIdx !== -1 ? logs.value[existingLogIdx].id : Math.max(...logs.value.map(l => l.id), 0) + 1,
      employee_id: clockForm.value.employee_id,
      employee_name: targetEmp ? `${targetEmp.first_name} ${targetEmp.last_name}` : 'Employee',
      department_name: targetEmp ? targetEmp.department_name : 'Engineering',
      date: filterDate.value,
      clock_in: clockForm.value.clock_in ? new Date(`${filterDate.value}T${clockForm.value.clock_in}:00`).toISOString() : null,
      clock_out: clockForm.value.clock_out ? new Date(`${filterDate.value}T${clockForm.value.clock_out}:00`).toISOString() : null,
      status: clockForm.value.status
    }

    if (existingLogIdx !== -1) {
      logs.value[existingLogIdx] = logData
    } else {
      logs.value.push(logData)
    }
    closeClockModal()
  }
}

onMounted(() => {
  fetchDepartments().then(() => {
    fetchEmployees().then(() => {
      fetchAttendance()
    })
  })
})
</script>

<style scoped>
.attendance-view {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.attendance-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.25rem 0;
}

.page-title {
  font-family: var(--font-display);
  font-size: 1.8rem;
  font-weight: 800;
  color: #1e293b;
  letter-spacing: -0.02em;
}

.breadcrumbs-path {
  font-size: 0.8rem;
  font-weight: 500;
  color: var(--text-muted);
}

.active-crumb {
  color: #3b82f6;
}

/* 4 Metrics Row containing heart icons */
.metrics-cards-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
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
  gap: 1.25rem;
}

.metric-icon-wrapper {
  width: 44px;
  height: 44px;
  border-radius: 50%; /* Circle backgrounds */
  display: flex;
  align-items: center;
  justify-content: center;
}

.metric-icon-wrapper.circle-green { background-color: #ecfdf5; }
.metric-icon-wrapper.circle-red { background-color: #fff5f5; }
.metric-icon-wrapper.circle-yellow { background-color: #fef3c7; }
.metric-icon-wrapper.circle-blue { background-color: #eff6ff; }

.metric-text-group {
  display: flex;
  flex-direction: column;
}

.metric-text-group .label {
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--text-muted);
  letter-spacing: 0.03em;
}

.metric-text-group .value {
  font-family: var(--font-display);
  font-size: 1.7rem;
  font-weight: 800;
  color: var(--text-primary);
  line-height: 1.1;
  margin: 0.1rem 0;
}

.pct-dot-label {
  font-size: 0.78rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.dot-marker {
  font-size: 1.1rem;
  line-height: 1;
}

.text-success { color: #10b981; }
.text-danger { color: #ef4444; }
.text-warning { color: #f59e0b; }
.text-info { color: #3b82f6; }

/* 65/35 Split grids */
.attendance-content-split {
  display: grid;
  grid-template-columns: 2fr 1.1fr;
  gap: 1.5rem;
}

/* Left split styling */
.logs-console-card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.filters-console-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filters-left-group {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.date-selector-btn-wrapper {
  background-color: #eff6ff;
  border: 1px solid rgba(59, 130, 246, 0.1);
  border-radius: var(--radius-sm);
  padding: 0.5rem 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.calendar-icon-prefix {
  font-size: 0.85rem;
}

.date-picker-element {
  background: none;
  border: none;
  font-family: var(--font-body);
  font-size: 0.82rem;
  font-weight: 700;
  color: #1e40af;
  outline: none;
  cursor: pointer;
  width: 120px;
}

.custom-console-select {
  background-color: white;
  border: 1px solid #e2e8f0;
  padding: 0.5rem 0.75rem;
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-size: 0.82rem;
  color: var(--text-secondary);
  font-weight: 500;
  cursor: pointer;
}

.custom-console-select:focus {
  outline: none;
  border-color: #cbd5e1;
}

.filters-right-group {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.add-attendance-btn {
  background-color: #3b82f6;
  border: none;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: var(--radius-sm);
  font-weight: 700;
  font-size: 0.82rem;
  cursor: pointer;
}

.add-attendance-btn:hover {
  background-color: #2563eb;
}

.filter-icon-btn {
  background-color: #6366f1; /* Purple filter button */
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: var(--radius-sm);
  font-weight: 700;
  font-size: 0.82rem;
  cursor: pointer;
}

.filter-icon-btn:hover {
  background-color: #4f46e5;
}

/* Custom Table Overrides */
.custom-table-container {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  overflow: hidden;
}

.attendance-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.85rem;
}

.header-strip {
  background-color: #f8fafc;
}

.attendance-table th {
  padding: 0.85rem 1rem;
  font-weight: 700;
  color: var(--text-secondary);
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  border-bottom: 1px solid var(--border-color);
}

.attendance-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}

.attendance-table tr:last-child td {
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

.full-name {
  font-weight: 600;
  color: var(--text-primary);
}

.time-cell {
  font-family: var(--font-display);
  font-weight: 500;
}

.no-records-cell {
  text-align: center;
  padding: 3rem;
  color: var(--text-muted);
}

/* Footer entries list */
.logs-footer-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.5rem;
}

.results-info {
  font-size: 0.8rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.simple-pagination-circles {
  display: flex;
  align-items: center;
}

.circle-pag-btn {
  background-color: #3b82f6;
  border: none;
  color: white;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-weight: 700;
  font-size: 0.8rem;
  cursor: pointer;
}

/* Right splits visual charts */
.right-split {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.visual-metrics-card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
}

.donut-chart-container {
  display: flex;
  align-items: center;
  justify-content: space-around;
  margin-top: 1rem;
}

.donut-svg {
  transform: rotate(-90deg);
}

.donut-segment {
  stroke-linecap: round;
  transition: stroke-dasharray 0.5s ease;
}

.donut-legend-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.legend-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.legend-color-box {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.legend-color-box.present { background-color: #10B981; }
.legend-color-box.absent { background-color: #EF4444; }
.legend-color-box.late { background-color: #F59E0B; }

.legend-text {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text-secondary);
}

/* Department Progress bar wise */
.department-progress-card {
  padding: 1.5rem;
}

.dept-progress-list {
  display: flex;
  flex-direction: column;
  gap: 1.1rem;
  margin-top: 1rem;
}

.progress-bar-row {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.progress-bar-meta {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--text-secondary);
}

.progress-outer-frame {
  height: 6px;
  background-color: #f1f5f9;
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-inner-bar {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width 0.8s ease;
}

/* Modal form halfs */
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
  .attendance-content-split {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .metrics-cards-row {
    grid-template-columns: 1fr;
  }
  .filters-console-row {
    flex-direction: column;
    align-items: stretch;
    gap: 1rem;
  }
  .filters-left-group {
    flex-direction: column;
    align-items: stretch;
  }
  .filters-right-group {
    justify-content: flex-end;
  }
}
</style>
