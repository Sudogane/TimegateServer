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

func NewAuthenticationResponse(token string, userData *UserData) ServerPayload {
	return &FromServerToClient_AuthenticationResponse{
		AuthenticationResponse: &AuthenticationResponse{
			AccessToken: token,
			UserData:    userData,
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
