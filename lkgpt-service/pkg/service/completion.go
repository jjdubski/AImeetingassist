package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	lksdk "github.com/livekit/server-sdk-go"
	"github.com/sashabaranov/go-openai"
)

// A sentence in the conversation (Used for the history)
type SpeechEvent struct {
	ParticipantName string
	IsBot           bool
	Text            string
}

type JoinLeaveEvent struct {
	Leave           bool
	ParticipantName string
	Time            time.Time
}

type MeetingEvent struct {
	Speech *SpeechEvent
	Join   *JoinLeaveEvent
}

type ChatCompletion struct {
	client *openai.Client
}

func NewChatCompletion(client *openai.Client) *ChatCompletion {
	return &ChatCompletion{
		client: client,
	}
}

type ResponseData struct {
	SummaryText string `json:"summary_text"`
}

func (c *ChatCompletion) Complete(ctx context.Context, events []*MeetingEvent, prompt *SpeechEvent,
	participant *lksdk.RemoteParticipant, room *lksdk.Room, language *Language) (string, error) {

	fmt.Println("Complete called")
	var sb strings.Builder
	sb.WriteString("The participants in the meeting are:")
	participants := room.GetParticipants()
	for i, participant := range participants {
		sb.WriteString(participant.Identity())
		if i != len(participants)-1 {
			sb.WriteString(", ")
		}
	}
	participantNames := sb.String()
	sb.Reset()

	// curl https://api-inference.huggingface.co/models/knkarthick/MEETING_SUMMARY \
	// -X POST \
	// -d '{"inputs": "The tower is 324 metres (1,063 ft) tall, about the same height as an 81-storey building, and the tallest structure in Paris. Its base is square, measuring 125 metres (410 ft) on each side. During its construction, the Eiffel Tower surpassed the Washington Monument to become the tallest man-made structure in the world, a title it held for 41 years until the Chrysler Building in New York City was finished in 1930. It was the first structure to reach a height of 300 metres. Due to the addition of a broadcasting aerial at the top of the tower in 1957, it is now taller than the Chrysler Building by 5.2 metres (17 ft). Excluding transmitters, the Eiffel Tower is the second tallest free-standing structure in France after the Millau Viaduct."}' \
	// -H 'Content-Type: application/json' \
	// -H "Authorization: Bearer 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

	url := "https://api-inference.huggingface.co/models/knkarthick/MEETING_SUMMARY"
	payload := []byte(fmt.Sprintf(`{"inputs": "%s"}`, participantNames+"\nThe transcription of the meeting is:\n"+prompt.Text))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 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")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	var responseData []ResponseData
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", err
	}

	if len(responseData) > 0 {
		summaryText := responseData[0].SummaryText
		fmt.Println("Summary Text:", summaryText)
		return summaryText, nil
	} else {
		fmt.Println("No data in the response array.")
		return "No data in the response array.", nil
	}

	// messages := make([]openai.ChatCompletionMessage, 0, len(events)+3)
	// messages = append(messages, openai.ChatCompletionMessage{
	// 	Role: openai.ChatMessageRoleSystem,
	// 	Content: "You are KITT, a voice assistant in a meeting created by LiveKit. " +
	// 		"Keep your responses concise while still being friendly and personable. " +
	// 		"If your response is a question, please append a question mark symbol to the end of it. " + // Used for auto-trigger
	// 		fmt.Sprintf("There are actually %d participants in the meeting: %s. ", len(participants), participantNames) +
	// 		fmt.Sprintf("Current language: %s Current date: %s", language.Label, time.Now().Format("January 2, 2006 3:04pm")),
	// })

	// for _, e := range events {
	// 	if e.Speech != nil {
	// 		if e.Speech.IsBot {
	// 			messages = append(messages, openai.ChatCompletionMessage{
	// 				Role:    openai.ChatMessageRoleAssistant,
	// 				Content: e.Speech.Text,
	// 				Name:    BotIdentity,
	// 			})
	// 		} else {
	// 			messages = append(messages, openai.ChatCompletionMessage{
	// 				Role:    openai.ChatMessageRoleUser,
	// 				Content: fmt.Sprintf("%s said %s", e.Speech.ParticipantName, e.Speech.Text),
	// 				Name:    e.Speech.ParticipantName,
	// 			})
	// 		}
	// 	}

	// 	if e.Join != nil {
	// 		if e.Join.Leave {
	// 			messages = append(messages, openai.ChatCompletionMessage{
	// 				Role:    openai.ChatMessageRoleSystem,
	// 				Content: fmt.Sprintf("%s left the meeting at %s", e.Join.ParticipantName, e.Join.Time.Format("3:04pm")),
	// 			})
	// 		} else {
	// 			messages = append(messages, openai.ChatCompletionMessage{
	// 				Role:    openai.ChatMessageRoleSystem,
	// 				Content: fmt.Sprintf("%s joined the meeting at %s", e.Join.ParticipantName, e.Join.Time.Format("3:04pm")),
	// 			})
	// 		}
	// 	}
	// }

	// messages = append(messages, openai.ChatCompletionMessage{
	// 	Role:    openai.ChatMessageRoleSystem,
	// 	Content: fmt.Sprintf("You are currently talking to %s", participant.Identity()),
	// })

	// prompt
	// messages = append(messages, openai.ChatCompletionMessage{
	// 	Role:    openai.ChatMessageRoleUser,
	// 	Content: prompt.Text,
	// 	Name:    prompt.ParticipantName,
	// })

	// stream, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
	// 	Model:    openai.GPT3Dot5Turbo,
	// 	Messages: messages,
	// 	Stream:   true,
	// })
}

// Wrapper around openai.ChatCompletionStream to return only complete sentences
// type ChatStream struct {
// 	stream string
// }

// func (c *ChatStream) Recv() (string, error) {
// 	sb := strings.Builder{}
// 	for {
// 		response, err := c.stream.Recv()
// 		if err != nil {
// 			content := sb.String()
// 			if err == io.EOF && len(strings.TrimSpace(content)) != 0 {
// 				return content, nil
// 			}
// 			return "", err
// 		}

// 		if len(response.Choices) == 0 {
// 			continue
// 		}

// 		delta := response.Choices[0].Delta.Content
// 		sb.WriteString(delta)

// 		if strings.HasSuffix(strings.TrimSpace(delta), ".") {
// 			return sb.String(), nil
// 		}
// 	}
// }

// func (c *ChatStream) Close() {
// 	c.stream.Close()
// }
