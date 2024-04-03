package authapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/apiclient"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/utils"
)

type API interface {
	Admin() (*UserDTO, error)
}

type authAPI struct {
	ctx    context.Context
	host   string
	client apiclient.BaseClient
}

func NewAuthAPI(ctx context.Context, host string, client apiclient.BaseClient) API {
	return &authAPI{
		ctx:    ctx,
		host:   host,
		client: client,
	}
}

func (a *authAPI) Admin() (*UserDTO, error) {
	var url = fmt.Sprintf("%s/api/internal/admin", a.host)

	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Send(req)
	if err != nil {
		return nil, err
	}

	defer utils.CloseReadCloser(resp.Body)

	var user = UserDTO{}

	err = apiclient.ScanResponse(resp, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
