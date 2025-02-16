package api

import (
	"context"

	"merchshop/internal/service"
)

type APIServer struct {
	merchService *service.MerchService
}

func NewAPIServer(ms *service.MerchService) *APIServer {
	return &APIServer{merchService: ms}
}

func (s *APIServer) PostApiAuth(ctx context.Context, req PostApiAuthRequestObject) (PostApiAuthResponseObject, error) {
	if req.Body == nil {
		return PostApiAuth400JSONResponse(ErrorResponse{Errors: ptr("Invalid request body")}), nil
	}
	username := req.Body.Username
	password := req.Body.Password

	token, err := s.merchService.Authenticate(ctx, username, password)
	if err != nil {
		return PostApiAuth500JSONResponse(ErrorResponse{Errors: ptr(err.Error())}), nil
	}
	resp := PostApiAuth200JSONResponse(AuthResponse{Token: &token})
	return resp, nil
}

func (s *APIServer) GetApiBuyItem(ctx context.Context, req GetApiBuyItemRequestObject) (GetApiBuyItemResponseObject, error) {
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return GetApiBuyItem400JSONResponse(ErrorResponse{Errors: ptr("missing or invalid user")}), nil
	}
	if err := s.merchService.BuyItem(ctx, username, req.Item); err != nil {
		return GetApiBuyItem500JSONResponse(ErrorResponse{Errors: ptr(err.Error())}), nil
	}
	return GetApiBuyItem200Response{}, nil
}

func (s *APIServer) GetApiInfo(ctx context.Context, req GetApiInfoRequestObject) (GetApiInfoResponseObject, error) {
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return GetApiInfo400JSONResponse(ErrorResponse{Errors: ptr("missing or invalid user")}), nil
	}
	info, err := s.merchService.GetInfo(ctx, username)
	if err != nil {
		return GetApiInfo500JSONResponse(ErrorResponse{Errors: ptr(err.Error())}), nil
	}

	var invAPI []struct {
		Quantity *int    `json:"quantity,omitempty"`
		Type     *string `json:"type,omitempty"`
	}
	for _, item := range info.Inventory {
		qty := int(item.Amount)
		typ := item.Item
		invAPI = append(invAPI, struct {
			Quantity *int    `json:"quantity,omitempty"`
			Type     *string `json:"type,omitempty"`
		}{
			Quantity: &qty,
			Type:     &typ,
		})
	}

	var sentAPI []struct {
		Amount *int    `json:"amount,omitempty"`
		ToUser *string `json:"toUser,omitempty"`
	}
	for _, t := range info.CoinHistory.Sent {
		amt := int(t.Amount)
		to := t.ToUsername
		sentAPI = append(sentAPI, struct {
			Amount *int    `json:"amount,omitempty"`
			ToUser *string `json:"toUser,omitempty"`
		}{
			Amount: &amt,
			ToUser: &to,
		})
	}

	var receivedAPI []struct {
		Amount   *int    `json:"amount,omitempty"`
		FromUser *string `json:"fromUser,omitempty"`
	}
	for _, t := range info.CoinHistory.Received {
		amt := int(t.Amount)
		from := t.FromUsername
		receivedAPI = append(receivedAPI, struct {
			Amount   *int    `json:"amount,omitempty"`
			FromUser *string `json:"fromUser,omitempty"`
		}{
			Amount:   &amt,
			FromUser: &from,
		})
	}

	var coinHistory struct {
		Received *[]struct {
			Amount   *int    `json:"amount,omitempty"`
			FromUser *string `json:"fromUser,omitempty"`
		} `json:"received,omitempty"`
		Sent *[]struct {
			Amount *int    `json:"amount,omitempty"`
			ToUser *string `json:"toUser,omitempty"`
		} `json:"sent,omitempty"`
	}
	if len(receivedAPI) > 0 {
		coinHistory.Received = &receivedAPI
	}
	if len(sentAPI) > 0 {
		coinHistory.Sent = &sentAPI
	}

	coinsVal := int(info.Coins)

	resp := GetApiInfo200JSONResponse(InfoResponse{
		Coins:       &coinsVal,
		Inventory:   &invAPI,
		CoinHistory: &coinHistory,
	})
	return resp, nil
}

func (s *APIServer) PostApiSendCoin(ctx context.Context, req PostApiSendCoinRequestObject) (PostApiSendCoinResponseObject, error) {
	fromUsername, ok := ctx.Value("username").(string)
	if !ok || fromUsername == "" {
		return PostApiSendCoin400JSONResponse(ErrorResponse{Errors: ptr("missing or invalid user")}), nil
	}
	if req.Body == nil {
		return PostApiSendCoin400JSONResponse(ErrorResponse{Errors: ptr("invalid request body")}), nil
	}
	if err := s.merchService.SendCoin(ctx, fromUsername, req.Body.ToUser, int32(req.Body.Amount)); err != nil {
		return PostApiSendCoin500JSONResponse(ErrorResponse{Errors: ptr(err.Error())}), nil
	}
	return PostApiSendCoin200Response{}, nil
}

func ptr(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}
