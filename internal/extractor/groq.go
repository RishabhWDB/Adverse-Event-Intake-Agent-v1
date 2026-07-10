package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CaseData struct {
	PatientAge       int    `json:"patient_age"`
	PatientSex       string `json:"patient_sex"`
	SuspectDrug      string `json:"suspect_drug"`
	Dose             string `json:"dose"`
	EventDescription string `json:"event_description"`
	OnsetDate        string `json:"onset_date"`
	Outcome          string `json:"outcome"`
	Reporter         string `json:"reporter"`
}

const systemPrompt = `You are a pharmacovigilance data extraction assistant. Extract adverse event case details from the narrative and return ONLY valid JSON matching this exact structure, no other text:

{
  "patient_age": <number>,
  "patient_sex": "<M/F/Unknown>",
  "suspect_drug": "<drug name>",
  "dose": "<dose if mentioned>",
  "event_description": "<the adverse event>",
  "onset_date": "<date or relative time if mentioned>",
  "outcome": "<recovered/recovering/not recovered/fatal/unknown>",
  "reporter": "<physician/patient/pharmacist/unknown>"
}

If a field is not mentioned, use "Unknown" for strings or 0 for numbers.`

func ExtractCase(narrative string) (*CaseData, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY not set")
	}

	payload := map[string]interface{}{
		"model": "openai/gpt-oss-20b",
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": narrative},
		},
		"response_format": map[string]string{"type": "json_object"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("groq api error: %s", string(respBody))
	}

	var groqResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &groqResp); err != nil {
		return nil, err
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from groq")
	}

	var caseData CaseData
	if err := json.Unmarshal([]byte(groqResp.Choices[0].Message.Content), &caseData); err != nil {
		return nil, fmt.Errorf("failed to parse extracted JSON: %w", err)
	}

	return &caseData, nil
}
