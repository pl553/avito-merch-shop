package api

import (
	"context"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (*Server) PostApiAuth(ctx context.Context, request PostApiAuthRequestObject) (PostApiAuthResponseObject, error) {
	return nil, nil
}

func (*Server) GetApiBuyItem(ctx context.Context, request GetApiBuyItemRequestObject) (GetApiBuyItemResponseObject, error) {
	return nil, nil
}

func (*Server) GetApiInfo(ctx context.Context, request GetApiInfoRequestObject) (GetApiInfoResponseObject, error) {
	return nil, nil
}

func (*Server) PostApiSendCoin(ctx context.Context, request PostApiSendCoinRequestObject) (PostApiSendCoinResponseObject, error) {
	return nil, nil
}
