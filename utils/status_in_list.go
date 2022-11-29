package utils

// Checks if status is in list
func StatusInList(status int, statusList []int) bool {
	for _, value := range statusList {
		if value == status {
			return true
		}
	}
	return false
}
