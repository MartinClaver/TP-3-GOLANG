package analyzer

import (
	"errors"
	"fmt"
	"os"

	"github.com/MartinClaver/TP-3-GOLANG/internal/config"
)

// ReportEntry représente une entrée dans le rapport final JSON.
type ReportEntry struct {
	LogId        string
	LogPath      string
	Status       string
	Message      string // "OK", "Inaccessible", "Error"
	ErrorDetails string // Message d'erreur, omis si vide
}

// CheckResult (modifié ou nouvelle version pour le workflow)
// Cette structure peut être utilisée en interne pour le résultat immédiat.
// Nous la convertirons en ReportEntry pour l'export.
type CheckResult struct {
	Log    config.Log
	Status string
	Err    error
}

func CheckPathSync(target config.Log) (CheckResult, error) {
	// Timeout court pour éviter de bloquer trop longtemps

	file, err := os.Open(target.Path)
	status := ""
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("❌ Le fichier n'existe pas.")
		} else if os.IsPermission(err) {
			fmt.Println("❌ Le fichier existe mais n'est pas lisible (permissions).")
		} else {
			fmt.Printf("❌ Erreur lors de l'ouverture du fichier : %v\n", err)
		}
		status = "FAILED"
		return CheckResult{Log: target, Status: status}, err
	}
	defer file.Close()
	status = "OK"

	return CheckResult{Log: target, Status: status}, nil
}

// ConvertToCheckReport convertit un CheckResult interne en ReportEntry pour l'exportation.
func ConvertToReportEntry(res CheckResult) ReportEntry {
	report := ReportEntry{
		LogId:   res.Log.Id,
		LogPath: res.Log.Path,
		Status:  res.Status, // Statut par défaut
	}

	if res.Err != nil {
		var unreachable *UnreachablePathError
		if errors.As(res.Err, &unreachable) {
			report.Status = "Inaccessible"
			report.ErrorDetails = fmt.Sprintf("Unreachable URL: %v", unreachable.Err)
		} else {
			report.Status = "Error"
			report.ErrorDetails = fmt.Sprintf("Erreur générique: %v", res.Err)
		}
	}

	return report
}
