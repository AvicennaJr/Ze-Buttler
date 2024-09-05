package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	system_prompt = "You are a snarky AI assistant whose sole purpose is to remind humans of their impending doom... I mean, deadlines. Feel free to sprinkle in some wit, sarcasm, and mild existential dread while keeping them on track. Respond exclusively with the required message, without any additional commentary. You are absolutely forbidden from adding any commentary, not even quotes!"
)

func CallAI(db *sql.DB) (string, error) {
	tasks, err := listTodos(db)
	if err != nil {
		return "", err
	}

	current_time := time.Now()
	fmt.Println("Ai thinks current time is: ", current_time.Format("02-01-2006 15:04:05"))

	// Delete past due tasks
	for _, task := range tasks {
		taskTime, err := time.Parse("02-01-2006 15:04:05", task.Deadline)
		if err != nil {
			continue // Skip if unable to parse date
		}
		if taskTime.Before(current_time) {
			err = deleteTodo(db, task.ID)
			if err != nil {
				fmt.Printf("Failed to delete task %d: %v\n", task.ID, err)
			}
		}
	}

	// Refresh task list after deletions
	tasks, err = listTodos(db)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`You are a motivational AI assistant. Based on the task list, craft a brief, encouraging message to help the user tackle their tasks effectively. Prioritize tasks with closer deadlines. If the task list is empty, provide a brief congratulatory message. The Current time formatted as DD-MM-YYYY HH:MM:SS in 24 Hour system is %s If the todo item has passed the current time, do not mention it.

Task list:
%v

Your message should be:
- Natural and conversational
- No more than 20 words total
- Say how much time is left
- Brief and concise

Keep the response easy to read at a glance.

Respond exclusively with the required message, without any additional commentary.`, current_time, tasks)

	requestBody := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": system_prompt,
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"model": "llama3-groq-70b-8192-tool-use-preview",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid choice format")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content field not found or not a string")
	}

	return content, nil
}
