package utils

import "github.com/google/uuid"

// Функция contains проверяет, содержится ли элемент в слайсе
func ContainsUUID(slice []uuid.UUID, element uuid.UUID) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}
