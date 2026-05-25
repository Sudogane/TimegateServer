package packets

type ClientPayload = isFromClientToServer_Payload
type ServerPayload = isFromServerToClient_Payload

func NewErrorMessage(code ErrorCode) ServerPayload {
	return &FromServerToClient_ErrorResponse{
		ErrorResponse: &ErrorResponse{
			Code: code,
		},
	}
}

func NewAuthenticationResponse(token string, userData *UserData, dialogueId string, userId string) ServerPayload {
	return &FromServerToClient_AuthenticationResponse{
		AuthenticationResponse: &AuthenticationResponse{
			AccessToken: token,
			UserId:      userId,
			UserData:    userData,
			DialogueTrigger: &DialogueTrigger{
				DialogueId: dialogueId,
			},
		},
	}
}

func NewWebsocketIdResponse(id string) ServerPayload {
	return &FromServerToClient_WebsocketId{
		WebsocketId: &WebsocketIDResponse{
			Id: id,
		},
	}
}

func NewDialogueTrigger(dialogueId string) ServerPayload {
	return &FromServerToClient_DialogueTrigger{
		DialogueTrigger: &DialogueTrigger{
			DialogueId: dialogueId,
		},
	}
}

func NewDigimonTeamViewResponse(digimon *DigimonData) ServerPayload {
	return &FromServerToClient_DigimonTeamViewResponse{
		DigimonTeamViewResponse: &DigimonTeamViewResponse{
			Digimon: digimon,
		},
	}
}

func NewAchievementResponse(achievementName, achievementGroupName string) *AchievementData {
	return &AchievementData{
		AchievementGroupName: achievementGroupName,
		AchievementName:      achievementName,
	}
}
