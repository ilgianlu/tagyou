package insights

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/ilgianlu/tagyou/conf"
)

type BlueprintDTO struct {
	ClientID             string
	Model                string
}

func (a *BlueprintDTO) Validate() bool {
	if a.ClientID == "" {
		return false
	}
	if a.Model == "" {
		return false
	}
	return true
}

func (dc InsightsController) Blueprint(w http.ResponseWriter, r *http.Request) {
	blueprintParams := BlueprintDTO{}
	if err := json.NewDecoder(r.Body).Decode(&blueprintParams); err != nil {
		slog.Error("error decoding json input", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !blueprintParams.Validate() {
		slog.Warn("data passed is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	slog.Debug("blueprint request", "client-id", blueprintParams.ClientID, "model", blueprintParams.Model)
	input, err := os.OpenFile(conf.DebugDataFilepath(blueprintParams.ClientID), os.O_RDONLY, 0644)
	if (err != nil) {
		slog.Error("Error opening debug data file", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer input.Close()

	output, err := os.OpenFile(conf.BlueprintFilepath(blueprintParams.Model), os.O_CREATE|os.O_WRONLY, 0644)
	if (err != nil) {
		slog.Error("Error creating blueprint file", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer output.Close()

	// blueprint csv header
	output.Write([]byte("\"timestamp\",\"senderid\",\"topic\",\"message\"\n"))
	if _, err := io.Copy(output, input); err != nil {
		slog.Error("Error creating blueprint file", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
