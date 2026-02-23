package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Provider interface {
	Generate(prompt string) (string, error)
}

type GeminiProvider struct {
	APIKey string
	Model  string
}

func (p *GeminiProvider) Generate(prompt string) (string, error) {
	if p.APIKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY não configurada")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.Model, p.APIKey)

	reqData := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	}
	jsonData, _ := json.Marshal(reqData)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta da API: %v", err)
	}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(result.Candidates[0].Content.Parts[0].Text), nil
	}

	return "", fmt.Errorf("API retornou resposta vazia")
}
