package main

// contains checks if a string is present in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to remove an item from a slice
func removeItem(slice []string, item string) []string {
	var newSlice []string = make([]string, 0, len(slice))
	for _, s := range slice {
		if s != item {
			newSlice = append(newSlice, s)
		}
	}
	return newSlice
}
