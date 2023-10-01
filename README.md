# AIMeetingAssist

## Summary
AIMeetingAssist is an AI-powered meeting assistant that lives in a WebRTC conference call.

### What it does:
1. Listen to your conversations
1. Transcribes your audio converstations
1. Sends an email summary of your meeting converstations

<img src="https://user-images.githubusercontent.com/8453967/231227021-4f5a4412-ff14-4837-97e7-c55a1d9717c4.gif" 
        alt="Preview gif" 
        width="640" 
        />

## How it works

This repo contains two services:
1. meet
2. lkgpt-service

The `meet` service is a NextJS app that implements a typical video call app. The `lkgpt-service` implements KITT. When a room is created, a webhook calls a handler in `lkgpt-service` which adds a participant to the room. The particpant uses GCP's speech-to-text, a pretrained Hugging Face model, and GCP's text-to-speech to create KITT.


## Getting started

### Prerequisites

- `GOOGLE_APPLICATION_CREDENTIALS` json body from a Google cloud account. See <https://cloud.google.com/docs/authentication/application-default-credentials#GAC>
- Hugging Face Pretrained Model api key
- LiveKit api_key, api_secret, and url from [LiveKit Cloud](https://cloud.livekit.io)
- Go 1.19+ and Node.js

### Running Locally

To run locally, you'll need to run the `lkgt-service` service.

2. Locally, run lkgpt-service Locally:
```bash
# From the lkgpt-service/ directory
go run /cmd/server/main.go --config config.yaml --gcp-credentials-path gcp-credentials.json`
```
1. Go to <https://kittplus.vercel.app/rooms/3qjq-2itk> in your browser and enter your name & email.
3. Once the `lkgt-service` service is running, you can navigate back to your browser. There's one more step needed when running locally. When deployed, KITT is spawned via a LiveKit webhook, but locally - the webhook will have no way of reaching your local `lkgpt-service` that's running. So you'll have to manually call an API to spawn KITT:
```bash
# <room_name> comes from the url slug when you enter a room in the UI
curl -XPOST http://localhost:3001/join/<room_name>
```
