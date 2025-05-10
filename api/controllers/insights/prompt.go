package insights

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"text/template"
)

type TemplateData struct {
	Blueprint string
	DebugData string
}

func prompt(model string, clientID string) (string, error) {
	blueprint, err := readBlueprint(model)
	if err != nil {
		return blueprint, err
	}
	debugData, err := readDebugData(clientID)
	if err != nil {
		return debugData, err
	}
	tmpl, err := template.ParseFiles("./api/controllers/insights/prompt.tmpl")
	if err != nil {
		return "no template to load!", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, TemplateData{ Blueprint: blueprint, DebugData: debugData})
	if err != nil {
		return "error running template", err
	}
	return buf.String(), nil
}

func readBlueprint(model string) (string, error) {
    modelBluePrint, err := os.ReadFile(blueprintFilepath(model))
	if err != nil {
		slog.Debug("Failed to read file", "err", err)
		return fmt.Sprintf("No data for model %s", model), err
	}
	return string(modelBluePrint), nil
}

func readDebugData(clientID string) (string, error) {
    debugData, err := os.ReadFile(debugDataFilepath(clientID))
	if err != nil {
		slog.Debug("Failed to read file", "err", err)
		return fmt.Sprintf("no debug data for client %s", clientID), err
	}
	return string(debugData), nil
}
