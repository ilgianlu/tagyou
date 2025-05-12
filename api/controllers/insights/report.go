package insights

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/ilgianlu/tagyou/conf"
)

type ReportDTO struct {
	ClientID          string
	Model             string
}

func (a *ReportDTO) Validate() bool {
	if a.ClientID == "" {
		return false
	}
	if a.Model == "" {
		return false
	}
	return true
}

func (dc InsightsController) Report(w http.ResponseWriter, r *http.Request) {
	reportParams := ReportDTO{}
	if err := json.NewDecoder(r.Body).Decode(&reportParams); err != nil {
		slog.Error("error decoding json input", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !reportParams.Validate() {
		slog.Warn("data passed is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	slog.Debug("report request", "client-id", reportParams.ClientID)

	client := openai.NewClient(
    	option.WithBaseURL(conf.AI_URL),
		option.WithAPIKey(conf.AI_API_KEY),
	)

	curPrompt, err := prompt(reportParams.Model, reportParams.ClientID)
	if err != nil {
		slog.Warn("could not creae meaningful prompt", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Debug("sending prompt", "prompt", curPrompt)
	stream := client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(curPrompt),
		},
		Model: conf.AI_MODEL,
	})

	w.WriteHeader(http.StatusOK)
	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)
	
		if content, ok := acc.JustFinishedContent(); ok {
			slog.Debug("Content stream finished", "content", content)
		}
	
		// if using tool calls
		if tool, ok := acc.JustFinishedToolCall(); ok {
			slog.Debug("Tool call stream finished", "index", tool.Index, "name", tool.Name, "arguments", tool.Arguments)
		}
	
		if refusal, ok := acc.JustFinishedRefusal(); ok {
			slog.Debug("Refusal stream finished", "refusal", refusal)
		}
	
		// it's best to use chunks after handling JustFinished events
		if len(chunk.Choices) > 0 {
	        w.Write([]byte(chunk.Choices[0].Delta.Content))
			slog.Debug("chunk", "data", chunk.Choices[0].Delta.Content)
		}
	}
	
	if stream.Err() != nil {
		slog.Error("error contacting ai", "err", stream.Err().Error)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
