<template>
  <div class="dashboard-view">
    <!-- Top Row: 4 Stat Cards Grid -->
    <div class="grid-cols-4 stat-cards-grid">
      <!-- 1. Total Employees -->
      <div class="card stat-card">
        <div class="stat-icon-container user-group">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">TOTAL EMPLOYEES</span>
          <h3 class="stat-value">{{ stats.total_employees }}</h3>
          <span class="stat-subtext trend-up">▲ {{ stats.total_employees }} this month</span>
        </div>
      </div>

      <!-- 2. Departments -->
      <div class="card stat-card">
        <div class="stat-icon-container building">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="4" y="2" width="16" height="20" rx="2" ry="2"></rect>
            <line x1="9" y1="22" x2="9" y2="16"></line>
            <line x1="15" y1="22" x2="15" y2="16"></line>
            <line x1="9" y1="16" x2="15" y2="16"></line>
            <path d="M8 6h.01M16 6h.01M8 10h.01M16 10h.01M12 6h.01M12 10h.01"></path>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">DEPARTMENTS</span>
          <h3 class="stat-value">{{ stats.total_departments }}</h3>
          <span class="stat-subtext trend-none">No change</span>
        </div>
      </div>

      <!-- 3. Present Today -->
      <div class="card stat-card">
        <div class="stat-icon-container user-single">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
            <circle cx="12" cy="7" r="4"></circle>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">PRESENT TODAY</span>
          <h3 class="stat-value">{{ presentTodayCount }}</h3>
          <span class="stat-subtext trend-present">● {{ presentTodayPct }}%</span>
        </div>
      </div>

      <!-- 4. Pending Leaves -->
      <div class="card stat-card">
        <div class="stat-icon-container calendar">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
            <line x1="16" y1="2" x2="16" y2="6"></line>
            <line x1="8" y1="2" x2="8" y2="6"></line>
            <line x1="3" y1="10" x2="21" y2="10"></line>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">PENDING LEAVES</span>
          <h3 class="stat-value">{{ stats.pending_leaves }}</h3>
          <router-link to="/leaves" class="stat-subtext trend-leaves">View all</router-link>
        </div>
      </div>
    </div>

    <!-- Middle Row: Attendance Overview Bar Chart & Recent Leave Requests -->
    <div class="dashboard-middle-row">
      <!-- Attendance Overview Vertical Bar Chart Card -->
      <div class="card chart-card">
        <h3 class="card-header-title">Attendance Overview</h3>
        <div class="chart-content-wrapper">
          <!-- Y-Axis Labels -->
          <div class="y-axis-labels">
            <span>100</span>
            <span>90</span>
            <span>80</span>
            <span>70</span>
            <span>60</span>
            <span>50</span>
            <span>40</span>
            <span>30</span>
          </div>
          
          <!-- Grid Area & Columns -->
          <div class="chart-grid-area">
            <!-- Grid Horizontal Lines -->
            <div class="grid-lines">
              <div class="grid-line"></div>
              <div class="grid-line"></div>
              <div class="grid-line"></div>
              <div class="grid-line"></div>
              <div class="grid-line"></div>
              <div class="grid-line"></div>
              <div class="grid-line"></div>
              <div class="grid-line"></div>
            </div>

            <!-- Vertical Bars Columns -->
            <div class="bars-columns">
              <div v-for="day in weeklyAttendance" :key="day.label" class="bar-column">
                <div class="bar-wrapper">
                  <div 
                    class="bar-fill" 
                    :style="{ height: getBarHeight(day.pct) + '%' }"
                    :title="day.pct + '% Attendance'"
                  >
                    <!-- Tooltip percentage on hover -->
                    <span class="bar-tooltip">{{ day.pct }}%</span>
                  </div>
                </div>
                <span class="day-label">{{ day.label }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Leave Requests Card -->
      <div class="card recent-leaves-card">
        <h3 class="card-header-title">Recent Leave Requests</h3>
        <div class="leaves-list">
          <div v-for="leave in recentLeaves" :key="leave.id" class="leave-row-item">
            <span class="leave-emp-name">{{ leave.employee_name }}</span>
            <span class="badge" :class="'badge-' + leave.status.toLowerCase()">
              {{ leave.status }}
            </span>
          </div>
          <div v-if="recentLeaves.length === 0" class="leaves-empty-msg">
            No recent leave applications.
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Row: Recent Employees Table & Department Distribution Chart -->
    <div class="dashboard-bottom-row">
      <!-- Recent Employees Table -->
      <div class="card employees-table-card">
        <h3 class="card-header-title">Recent Employees</h3>
        <div class="table-scroll-container">
          <table class="recent-emps-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Department</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="emp in recentEmployees" :key="emp.id">
                <td class="emp-id-text">EMP{{ String(emp.id).padStart(3, '0') }}</td>
                <td class="emp-name-text">{{ emp.first_name }} {{ emp.last_name }}</td>
                <td>{{ emp.department_name || 'Unassigned' }}</td>
                <td>
                  <span class="badge" :class="'badge-' + (emp.status ? emp.status.toLowerCase() : 'active')">
                    {{ emp.status || 'Active' }}
                  </span>
                </td>
              </tr>
              <tr v-if="recentEmployees.length === 0">
                <td colspan="4" class="table-empty-row">No employees registered yet.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Department Distribution Donut Chart -->
      <div class="card distribution-chart-card">
        <h3 class="card-header-title">Department Distribution</h3>
        <div class="distribution-body">
          <!-- Donut SVG Segment -->
          <div class="donut-chart-container">
            <svg width="130" height="130" viewBox="0 0 36 36" class="donut-chart-element">
              <!-- Background Ring -->
              <circle cx="18" cy="18" r="15.9155" fill="transparent" stroke="#f1f5f9" stroke-width="4" />
              <!-- Segments -->
              <circle 
                v-for="item in deptDistribution" 
                :key="item.name" 
                cx="18" 
                cy="18" 
                r="15.9155" 
                fill="transparent" 
                :stroke="item.color" 
                stroke-width="4.2" 
                :stroke-dasharray="`${item.pct} ${100 - item.pct}`" 
                :stroke-dashoffset="item.dashoffset" 
                transform="rotate(-90 18 18)" 
                class="donut-segment"
              />
            </svg>
          </div>

          <!-- Color Legends -->
          <div class="chart-legend-grid">
            <div v-for="item in deptDistribution" :key="item.name" class="legend-item-row">
              <span class="color-bullet" :style="{ backgroundColor: item.color }"></span>
              <span class="legend-lbl">{{ item.name }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const stats = ref({
  total_employees: 0,
  active_employees: 0,
  total_departments: 0,
  pending_leaves: 0,
  monthly_payroll_sum: 0,
  recent_activities: []
})

const employees = ref([])
const leaves = ref([])

const presentTodayCount = ref(0)
const presentTodayPct = ref(81.67)

// Attendance Weekly Overview mock data matching the screenshot values
const weeklyAttendance = ref([
  { label: 'Mon', pct: 85 },
  { label: 'Tue', pct: 98 },
  { label: 'Wed', pct: 92 },
  { label: 'Thu', pct: 90 },
  { label: 'Fri', pct: 96 },
  { label: 'Sat', pct: 35 },
  { label: 'Sun', pct: 70 }
])

const getBarHeight = (pct) => {
  const min = 30
  const max = 100
  if (pct <= min) return 0
  if (pct >= max) return 100
  return ((pct - min) / (max - min)) * 100
}

// Select recent 3 leave requests
const recentLeaves = computed(() => {
  if (leaves.value.length > 0) {
    return leaves.value.slice(0, 3)
  }
  // Fallback Mock values matching the screenshot
  return [
    { id: 1, employee_name: 'Kavindu Perera', status: 'Pending' },
    { id: 2, employee_name: 'Nimesha Fernando', status: 'Pending' },
    { id: 3, employee_name: 'Tharushi Silva', status: 'Approved' }
  ]
})

// Select recent 3 employees
const recentEmployees = computed(() => {
  if (employees.value.length > 0) {
    // Sort by id desc
    const sorted = [...employees.value].sort((a, b) => b.id - a.id)
    return sorted.slice(0, 3)
  }
  // Fallback Mock values matching screenshot
  return [
    { id: 1, first_name: 'Dhanushka', last_name: 'Wickramasinghe', department_name: 'IT', status: 'Active' },
    { id: 2, first_name: 'Hasini', last_name: 'Fernando', department_name: 'HR', status: 'Active' },
    { id: 3, first_name: 'Kasun', last_name: 'Perera', department_name: 'Finance', status: 'Active' }
  ]
})

// Department Distribution Donut Segment percentages
const deptDistribution = computed(() => {
  const counts = { IT: 0, HR: 0, Finance: 0, Marketing: 0, Other: 0 }
  
  if (employees.value.length > 0) {
    employees.value.forEach(emp => {
      const dept = emp.department_name ? emp.department_name.toLowerCase() : ''
      if (dept.includes('engineering') || dept.includes('it') || dept.includes('information')) {
        counts.IT++
      } else if (dept.includes('human') || dept.includes('hr')) {
        counts.HR++
      } else if (dept.includes('finance')) {
        counts.Finance++
      } else if (dept.includes('marketing')) {
        counts.Marketing++
      } else {
        counts.Other++
      }
    })
  } else {
    // Mock values matching donut segments proportion
    counts.IT = 3
    counts.HR = 2
    counts.Finance = 1
    counts.Marketing = 1
    counts.Other = 1
  }
  
  const total = (counts.IT + counts.HR + counts.Finance + counts.Marketing + counts.Other) || 1
  const list = [
    { name: 'IT', count: counts.IT, color: '#3b82f6' },
    { name: 'HR', count: counts.HR, color: '#10b981' },
    { name: 'Finance', count: counts.Finance, color: '#f59e0b' },
    { name: 'Marketing', count: counts.Marketing, color: '#8b5cf6' },
    { name: 'Other', count: counts.Other, color: '#64748b' }
  ]
  
  let offset = 0
  return list.map(item => {
    const pct = (item.count / total) * 100
    const dashoffset = offset
    offset += pct
    return {
      ...item,
      pct,
      dashoffset: -dashoffset
    }
  })
})

const fetchStats = async () => {
  try {
    const res = await fetch(`${API_URL}/api/dashboard/stats`)
    if (res.ok) {
      stats.value = await res.json()
    }
  } catch (err) {
    console.error('Error fetching dashboard stats: ', err)
    stats.value = {
      total_employees: 8,
      active_employees: 7,
      total_departments: 4,
      pending_leaves: 2,
      monthly_payroll_sum: 47970.00
    }
  }
}

const fetchEmployees = async () => {
  try {
    const res = await fetch(`${API_URL}/api/employees`)
    if (res.ok) {
      employees.value = await res.json()
    }
  } catch (err) {
    console.error('Error loading employees:', err)
  }
}

const fetchLeaves = async () => {
  try {
    const res = await fetch(`${API_URL}/api/leaves`)
    if (res.ok) {
      leaves.value = await res.json()
    }
  } catch (err) {
    console.error('Error loading leaves:', err)
  }
}

const fetchPresentToday = async () => {
  try {
    const todayStr = new Date().toISOString().split('T')[0]
    const res = await fetch(`${API_URL}/api/attendance?date=${todayStr}`)
    if (res.ok) {
      const logs = await res.json()
      const presentLogs = logs.filter(l => l.status === 'Present' || l.status === 'Late')
      presentTodayCount.value = presentLogs.length
      
      const denominator = stats.value.total_employees || employees.value.length || 8
      presentTodayPct.value = Math.round((presentLogs.length / denominator) * 10000) / 100
    }
  } catch (err) {
    console.error('Error fetching present stats:', err)
  }
}

onMounted(async () => {
  await fetchStats()
  await fetchEmployees()
  await fetchLeaves()
  await fetchPresentToday()
})
</script>

<style scoped>
.dashboard-view {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.stat-cards-grid {
  gap: 1.5rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 1.25rem 1.5rem;
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.stat-icon-container {
  width: 46px;
  height: 46px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-group { background-color: #eff6ff; color: #3b82f6; }
.building { background-color: #ecfdf5; color: #10b981; }
.user-single { background-color: #f5f3ff; color: #8b5cf6; }
.calendar { background-color: #fff1f2; color: #f43f5e; }

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 0.72rem;
  color: var(--text-muted);
  font-weight: 700;
  letter-spacing: 0.05em;
}

.stat-value {
  font-family: var(--font-display);
  font-size: 1.6rem;
  font-weight: 700;
  margin: 0.1rem 0;
  color: var(--text-primary);
  line-height: 1.1;
}

.stat-subtext {
  font-size: 0.72rem;
  font-weight: 600;
  text-decoration: none;
}

.trend-up { color: #10b981; }
.trend-none { color: #94a3b8; }
.trend-present { color: #f59e0b; }
.trend-leaves { color: #ef4444; }

/* Middle Row: Chart & Leaves split panel */
.dashboard-middle-row {
  display: grid;
  grid-template-columns: 2.2fr 1fr;
  gap: 1.5rem;
}

.chart-card, .recent-leaves-card, .employees-table-card, .distribution-chart-card {
  background: white;
  border-radius: var(--radius-md);
  padding: 1.5rem;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.chart-card:hover, .recent-leaves-card:hover, .employees-table-card:hover, .distribution-chart-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.card-header-title {
  font-family: var(--font-display);
  font-size: 1.05rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 1.5rem;
}

/* Vertical Bar Chart Custom Layout */
.chart-content-wrapper {
  display: flex;
  gap: 1rem;
  height: 200px;
}

.y-axis-labels {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  font-size: 0.75rem;
  color: #94a3b8;
  text-align: right;
  width: 20px;
  font-weight: 500;
}

.chart-grid-area {
  flex: 1;
  position: relative;
  height: 100%;
}

.grid-lines {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  pointer-events: none;
}

.grid-line {
  border-top: 1px dashed #f1f5f9;
  width: 100%;
}

.bars-columns {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  z-index: 2;
}

.bar-column {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 10%;
  height: 100%;
}

.bar-wrapper {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  position: relative;
}

.bar-fill {
  width: 24px;
  background: linear-gradient(180deg, #10b981, #059669);
  border-radius: 6px 6px 0 0;
  transition: height 0.6s cubic-bezier(0.16, 1, 0.3, 1);
  position: relative;
  cursor: pointer;
}

.bar-fill:hover {
  background: linear-gradient(180deg, #059669, #047857);
}

.bar-tooltip {
  visibility: hidden;
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translate(-50%, -6px);
  background: #0f172a;
  color: white;
  font-size: 0.65rem;
  font-weight: 600;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  white-space: nowrap;
  z-index: 10;
}

.bar-fill:hover .bar-tooltip {
  visibility: visible;
}

.day-label {
  font-size: 0.75rem;
  color: #94a3b8;
  font-weight: 600;
  margin-top: 0.5rem;
}

/* Recent Leaves list */
.leaves-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.leave-row-item {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: var(--radius-sm);
  padding: 0.85rem 1.1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.leave-emp-name {
  font-size: 0.85rem;
  font-weight: 600;
  color: #1e293b;
}

.leaves-empty-msg {
  text-align: center;
  color: var(--text-muted);
  padding: 2rem 0;
  font-size: 0.85rem;
}

/* Bottom Row: Employees & Distribution donut split */
.dashboard-bottom-row {
  display: grid;
  grid-template-columns: 2.2fr 1fr;
  gap: 1.5rem;
}

.table-scroll-container {
  overflow-x: auto;
}

.recent-emps-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.85rem;
}

.recent-emps-table th {
  background-color: #f8fafc;
  padding: 0.75rem 1rem;
  color: #64748b;
  font-weight: 600;
  border-bottom: 1px solid #e2e8f0;
  text-transform: uppercase;
  font-size: 0.72rem;
  letter-spacing: 0.05em;
}

.recent-emps-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid #e2e8f0;
  color: #334155;
  font-weight: 500;
}

.recent-emps-table tr:last-child td {
  border-bottom: none;
}

.emp-id-text {
  font-weight: 600;
  color: #64748b;
}

.emp-name-text {
  font-weight: 600;
  color: #1e293b;
}

.status-active-badge {
  color: #10b981;
  font-weight: 600;
}

.table-empty-row {
  text-align: center;
  color: var(--text-muted);
  padding: 2rem 0;
}

/* Donut chart styles */
.distribution-body {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 2rem;
  width: 100%;
  padding: 0.5rem 0;
}

.donut-chart-container {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.donut-chart-element {
  display: block;
}

.donut-segment {
  transition: stroke-dasharray 0.3s ease, stroke-dashoffset 0.3s ease;
}

.chart-legend-grid {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  flex-shrink: 0;
}

@media (max-width: 480px) {
  .distribution-body {
    flex-direction: column;
    gap: 1rem;
  }
}

.legend-item-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.color-bullet {
  width: 12px;
  height: 12px;
  border-radius: 3px;
  flex-shrink: 0;
}

.legend-lbl {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 600;
}

@media (max-width: 1024px) {
  .dashboard-middle-row, .dashboard-bottom-row {
    grid-template-columns: 1fr;
  }
}
</style>
