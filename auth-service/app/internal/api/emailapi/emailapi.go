package emailapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/e"
	"net/http"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/parser"
)

type API interface {
	Write(emailContent *EmailDTO) error
}

type emailApi struct {
	ctx    context.Context
	host   string
	client http.Client
}

func NewEmailAPI(ctx context.Context, host string) API {
	const serviceTimeout = 7 * time.Second

	return &emailApi{
		ctx:  ctx,
		host: host,
		client: http.Client{
			Timeout: serviceTimeout,
		},
	}
}

func (eApi *emailApi) Write(emailContent *EmailDTO) error {
	var url = fmt.Sprintf("%s/api/email", eApi.host)

	requestJSON, err := json.Marshal(*emailContent)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(eApi.ctx, http.MethodPost, url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return err
	}

	resp, err := eApi.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 400 {
		return nil
	}

	apiError, err := parser.Decode[e.ResponseError](resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return apiError
}
