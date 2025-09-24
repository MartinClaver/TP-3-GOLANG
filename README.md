# loganalyzer

Outil CLI en Go pour **analyser des fichiers de logs en parallèle**, avec gestion d'erreurs personnalisées et import/export JSON.

## ✨ Fonctionnalités

- Lecture d'un fichier **config JSON** listant les logs à analyser.
- Lancement **concurrent** d'analyses (une goroutine par log).
- **Gestion des erreurs** robuste avec `errors.Is` / `errors.As` :
  - `FileAccessError` (inexistant, permission, non régulier)
  - `ParsingError` (JSON invalide)
- **Export JSON** du rapport (`--output`).
- **CLI Cobra** avec sous-commandes et flags (`--config`, `--output`, `--timeout-ms`).
- **Code modulaire** (`internal/config`, `internal/analyzer`, `internal/reporter`).

## Utilisation
go run . analyze --config sample/config.json --output report.json
## ou
go build -o loganalyzer
./loganalyzer analyze -c sample/config.json -o report.json
