package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loganalyzer",
	Short: "Loganalyzer est un outil pour vérifier des logs.",
	Long:  `Un outil CLI en Go pour centraliser l'analyse de multiples logs en parallèle et d'en extraire des informations clés, tout en gérant les erreurs potentielles de manière robuste.`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "fichier de configuration (par défaut est $HOME/.gowatcher.yaml)")
}
