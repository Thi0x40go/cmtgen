package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	OneMB       = 1_000_000
	GeminiModel = "gemini-2.5-flash"
)

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("erro ao obter git diff: %w", err)
	}

	return out.String(), nil
}

func truncate(text string, max int) string {
	if len(text) > max {
		return text[:max]
	}
	return text
}

func getSubjectFromUser() string {
	fmt.Print("Digite o subject do commit (ou deixe vazio): ")
	reader := bufio.NewReader(os.Stdin)
	subject, _ := reader.ReadString('\n')
	return strings.TrimSpace(subject)
}

func loadPrompt() (string, error) {
	data, err := os.ReadFile("commit_message_prompt.txt")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func buildPrompt(basePrompt, subject, diff string) string {
	if subject == "" {
		return fmt.Sprintf(`%s

Here are the changes in this commit:
%s
`, basePrompt, diff)
	}

	return fmt.Sprintf(`%s

Here is the user's subject line:
%s

Here are the changes in this commit:
%s
`, basePrompt, subject, diff)
}

type geminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func generateCommitMessage(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY não definido")
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		GeminiModel,
		apiKey,
	)

	reqBody := geminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed geminiResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Candidates) == 0 {
		return "", fmt.Errorf("resposta vazia da API")
	}

	return strings.TrimSpace(parsed.Candidates[0].Content.Parts[0].Text), nil
}

func main() {
	diff, err := getGitDiff()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	diff = truncate(diff, OneMB)

	subject := getSubjectFromUser()

	basePrompt, err := loadPrompt()
	if err != nil {
		fmt.Println("Erro ao carregar prompt:", err)
		os.Exit(1)
	}

	prompt := buildPrompt(basePrompt, subject, diff)

	fmt.Println("\nGerando mensagem de commit com Gemini...\n")

	commitMessage, err := generateCommitMessage(prompt)
	if err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}

	fmt.Println("Mensagem sugerida:\n")
	fmt.Println(commitMessage)

	fmt.Print("\nDeseja fazer o commit com essa mensagem? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')

	if strings.TrimSpace(strings.ToLower(confirm)) == "y" {
		cmd := exec.Command("git", "commit", "-m", commitMessage)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
