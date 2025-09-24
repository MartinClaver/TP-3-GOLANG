package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Log struct {
	Id   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

func LoadLogsFromFile(filePath string) ([]Log, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}

	var targets []Log
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}
	return targets, nil
}

// SaveTargetsToFile écrit une liste d'InputTarget dans un fichier JSON.
func SaveTargetsToFile(filePath string, targets []Log) error {
	data, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}
	return nil
}
