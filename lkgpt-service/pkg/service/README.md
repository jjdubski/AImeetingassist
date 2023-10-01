# gptparticipant contains the listen to the coversation function
## In func (p *GPTParticipant) activateParticipant(rp *lksdk.RemoteParticipant) set the activation timeout
## func (p *GPTParticipant) onTranscriptionReceived(result RecognizeResult, rp *lksdk.RemoteParticipant, transcriber *Transcriber) contains the recieved transcrition from Google
### In here it also checks for activation word if there is more than one person, should probably remove this feature.
### More formally, the function accepts transcribed data returned from google and uses the data to respond to the client
## It looks like in func (p *GPTParticipant) trackSubscribed(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) the participants speech is tracked and send to Google
## The actual answering function func (p *GPTParticipant) answer(events []*MeetingEvent, prompt *SpeechEvent, rp *lksdk.RemoteParticipant, language *Language) (string, error)
## In func (t *Transcriber) newStream() (sttpb.Speech_StreamingRecognizeClient, error) the STT is set to listen for a certain phrase. We do not need this.
