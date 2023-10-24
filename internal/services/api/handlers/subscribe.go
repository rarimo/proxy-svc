package handlers

import (
	"bytes"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/proxy-svc/internal/types"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"io"
	"net/http"
	"net/url"
	"path"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func newSubscribeRequest(r *http.Request) (*types.SubscribeRequest, error) {
	var req types.SubscribeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	return &req, validation.Errors{
		"email": validation.Validate(req.Email, validation.Required),
	}.Filter()
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	request, err := newSubscribeRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	body := types.SubscribeRequestBody{
		Email:  request.Email,
		Status: "subscribed",
		Tags:   []string{"Landing subscribers"},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		Log(r).WithError(err).Error("failed to marshal body")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	reqUrl, err := url.Parse(Proxy(r).Url)

	if err != nil {
		Log(r).WithError(err).Error("failed to parse url")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	reqUrl.Path = path.Join(reqUrl.Path, "/3.0/lists/", Proxy(r).ListId, "/members")
	req, err := http.NewRequestWithContext(r.Context(), "POST", reqUrl.String(), bytes.NewReader(jsonBody))

	if err != nil {
		Log(r).WithError(err).Error("failed to create request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey "+Proxy(r).ApiKey)

	resp, err := Proxy(r).Client.Do(req)
	if err != nil {
		Log(r).WithError(err).Error("failed to do request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		Log(r).WithError(err).Error("failed to read response body")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}
