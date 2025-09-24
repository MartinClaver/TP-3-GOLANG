package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/MartinClaver/TP-3-GOLANG/internal/analyzer"
	"github.com/MartinClaver/TP-3-GOLANG/internal/config"
	"github.com/MartinClaver/TP-3-GOLANG/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	inputFilePath  string
	outputFilePath string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse un fichier de configuration JSON de logs..",
	Long:  `La commande prendra un chemin vers un fichier de configuration JSON via un drapeau --config <path> (raccourci -c). Ce fichier contiendra la liste des logs à analyser.`,
	Run: func(cmd *cobra.Command, args []string) {

		if inputFilePath == "" {
			fmt.Println("Erreur: le chemin du fichier d'entrée (--config) est obligatoire.")
			return
		}

		// Charger les cibles depuis le fichier JSON d'entrée
		targets, err := config.LoadLogsFromFile(inputFilePath)
		if err != nil {
			fmt.Printf("Erreur lors du chargement des Paths: %v\n", err)
			return
		}

		if len(targets) == 0 {
			fmt.Println("Aucune Path à vérifier trouvée dans le fichier d'entrée.")
			return
		}

		var wg sync.WaitGroup
		resultsChan := make(chan analyzer.CheckResult, len(targets)) // Canal pour collecter les résultats

		wg.Add(len(targets))
		for _, target := range targets {
			go func(t config.Log) {
				defer wg.Done()
				result, err := analyzer.CheckPathSync(t)
				if err != nil {
					fmt.Println("Le fichier n'est pas trouvé :", err)
				}
				resultsChan <- result // Envoyer le resultat au canal
			}(target)
		}

		wg.Wait()          // Attendre que toutes les goroutines aient fini
		close(resultsChan) // Fermer le canal après que tous les résultats ont été envoyés

		var finalReport []analyzer.ReportEntry
		for res := range resultsChan { // Récupérer tous les résultats du canal
			reportEntry := analyzer.ConvertToReportEntry(res)
			finalReport = append(finalReport, reportEntry)

			// Affichage immédiat comme avant
			if res.Err != nil {
				var unreachable *analyzer.UnreachablePathError
				if errors.As(res.Err, &unreachable) {
					fmt.Printf("🚫 %s (%s) est inaccessible : %v\n", res.Log.Id, unreachable.Path, unreachable.Err)
				} else {
					fmt.Printf("❌ %s (%s) : erreur - %v\n", res.Log.Id, res.Log.Path, res.Err)
				}
			} else {
				fmt.Printf("✅ %s (%s) : OK - %s\n", res.Log.Id, res.Log.Path, res.Status)
			}
		}

		// Exporter les résultats si outputFilePath est spécifié
		if outputFilePath != "" {
			err := reporter.ExportResultsToJsonFile(outputFilePath, finalReport)
			if err != nil {
				fmt.Printf("Erreur lors de l'exportation des résultats: %v\n", err)
			} else {
				fmt.Printf("✅ Résultats exportés vers %s\n", outputFilePath)
			}
		}
	},
}

// init() est une fonction spéciale de Go, exécutée lors de l'initialisation du package.
func init() {
	// Cette ligne est cruciale : elle "ajoute" la sous-commande `checkCmd` à la commande racine `rootCmd`.
	// C'est ainsi que Cobra sait que 'analyze' est une commande valide sous 'gowatcher'.
	rootCmd.AddCommand(analyzeCmd)

	// Ici, vous pouvez ajouter des drapeaux (flags) spécifiques à la commande 'check'.
	// Ces drapeaux ne seront disponibles que lorsque la commande 'check' est utilisée.
	// Exemple (commenté) : analyzeCmd.Flags().StringVarP(&sourceFile, "source", "s", "", "Fichier contenant les Paths à vérifier")

	// Ajout des drapeaux spécifiques à la commande 'check'
	analyzeCmd.Flags().StringVarP(&inputFilePath, "config", "c", "", "Chemin vers le fichier JSON d'entrée contenant les Paths")
	analyzeCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Chemin vers le fichier JSON de sortie pour les résultats (optionnel)")

	// Marquer le drapeau "input" comme obligatoire
	analyzeCmd.MarkFlagRequired("config")
}
