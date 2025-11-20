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

func NewLoginResponse(access_token, username string, level, exp int32, bits, yen int64, staminaCurrent, staminaMax int32) ServerPayload {
	return &FromServerToClient_LoginResponse{
		LoginResponse: &LoginResponse{
			AccessToken: access_token,
			UserData: &UserData{
				Username:       username,
				Level:          level,
				Exp:            exp,
				Bits:           bits,
				Yen:            yen,
				StaminaCurrent: staminaCurrent,
				StaminaMax:     staminaMax,
			},
		},
	}
}

func NewRegisterResponse(access_token string) ServerPayload {
	return &FromServerToClient_RegisterResponse{
		RegisterResponse: &RegisterResponse{
			AccessToken: access_token,
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
