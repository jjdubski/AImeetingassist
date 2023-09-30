package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	sb.WriteString("The participants in the meeting are: ")
	participants := room.GetParticipants()
	for i, participant := range participants {
		sb.WriteString(participant.Identity())
		if i != len(participants)-1 {
			sb.WriteString(", ")
		}
	}
	// participantNames := sb.String()
	// comp := participantNames + "\nThe transcription of the meeting is:\n" + prompt.Text
	comp := sb.String()
	comp = prompt.Text
	fmt.Println("comp:", comp)
	sb.Reset()

	url := "https://api-inference.huggingface.co/models/knkarthick/MEETING_SUMMARY"
	payload := []byte(fmt.Sprintf(`{"inputs": "%s"}`, comp))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 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")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(string(body))

	var respData []ResponseData

	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Access the parsed data
	for _, info := range respData {
		fmt.Println("Summary Text:", info.SummaryText)
	}

	return respData[0].SummaryText, nil

}
