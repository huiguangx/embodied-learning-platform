package repository

type Dashboard struct {
	ClusterStatus string `json:"clusterStatus"`
	GPUUtilization float64 `json:"gpuUtilization"`
	MemoryUtilization float64 `json:"memoryUtilization"`
	CPUUtilization float64 `json:"cpuUtilization"`
	TotalTasks int `json:"totalTasks"`
	TodayTasks int `json:"todayTasks"`
	WeeklyTasks int `json:"weeklyTasks"`
	SuccessRate float64 `json:"successRate"`
}

func GetDashboard(projectID string) (Dashboard, bool) {
	if projectID == "" { return Dashboard{}, false }
	return Dashboard{ClusterStatus:"healthy", GPUUtilization:79, MemoryUtilization:71.7, CPUUtilization:17.7, TotalTasks:4353, TodayTasks:0, WeeklyTasks:21, SuccessRate:30}, true
}
