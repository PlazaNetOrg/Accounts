package i18n

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type MessageSet map[string]string

var (
	translations = make(map[string]MessageSet)
	mu            sync.RWMutex
)

func Init() {
	mu.Lock()
	defer mu.Unlock()

	supportedLangs := []string{"en", "pl"}

	for _, lang := range supportedLangs {
		filePath := "translations/" + lang + ".json"
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read translation file %s: %v", filePath, err)
		}

		var rawTranslations map[string]interface{}
		if err := json.Unmarshal(data, &rawTranslations); err != nil {
			log.Fatalf("Failed to parse translation file %s: %v", filePath, err)
		}

		messages := make(MessageSet)
		for key, value := range rawTranslations {
			if msgMap, ok := value.(map[string]interface{}); ok {
				if other, ok := msgMap["other"].(string); ok {
					messages[key] = other
				}
			}
		}

		translations[lang] = messages
		log.Printf("Loaded %d messages for language: %s", len(messages), lang)
	}
}

func Get(lang, messageID string) string {
	mu.RLock()
	defer mu.RUnlock()

	if _, exists := translations[lang]; !exists {
		lang = "en"
	}

	if msg, exists := translations[lang][messageID]; exists {
		return msg
	}

	if lang != "en" {
		if msg, exists := translations["en"][messageID]; exists {
			return msg
		}
	}

	return messageID
}

func GetAll(lang string) MessageSet {
	mu.RLock()
	defer mu.RUnlock()

	if _, exists := translations[lang]; !exists {
		lang = "en"
	}

	return translations[lang]
}
