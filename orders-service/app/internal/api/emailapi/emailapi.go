package emailapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/apiclient"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/utils"
	"net/http"
	"time"
)

type API interface {
	Write(emailContent *EmailDTO) error
}

type emailApi struct {
	ctx    context.Context
	host   string
	client apiclient.BaseClient
}

func NewEmailAPI(ctx context.Context, host string, client apiclient.BaseClient) API {
	const serviceTimeout = 7 * time.Second

	return &emailApi{
		ctx:    ctx,
		host:   host,
		client: client,
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

	resp, err := eApi.client.Send(req)
	if err != nil {
		return err
	}
	defer utils.CloseReadCloser(resp.Body)

	err = apiclient.ScanResponse[any](resp, nil)

	return err
}
