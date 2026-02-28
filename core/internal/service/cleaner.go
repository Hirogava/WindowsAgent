package service

import "strings"

// удаляем префиксы нейронки из ответа, чтобы остался чистый JSON для выполнения команды, по типу таких ```json
func JsonCleaner(jsonStr string) string {
	cleaned := strings.Split(jsonStr, "```json")
	if len(cleaned) > 1 {
		cleaned = strings.Split(cleaned[1], "```")
	}
	return cleaned[0]
}