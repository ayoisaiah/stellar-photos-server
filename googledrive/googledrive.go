package googledrive

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/ayoisaiah/stellar-photos-server/config"
	"github.com/ayoisaiah/stellar-photos-server/unsplash"
	"github.com/ayoisaiah/stellar-photos-server/utils"
)

const (
	saveToDriveTimeout = 180
)

// Google drive application key.
type key struct {
	GoogleDriveKey string `json:"googledrive_key"`
}

type googleDriveAuth struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// SendGoogleDriveKey sends the application key to the client on request to avoid
// exposing it in the extension code.
func SendGoogleDriveKey(w http.ResponseWriter, r *http.Request) error {
	d := key{
		GoogleDriveKey: config.Conf.GoogleDrive.Key,
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return utils.JSONResponse(w, b)
}

// AuthorizeGoogleDrive redeems the authorization code received from the client for
// an access token.
func AuthorizeGoogleDrive(w http.ResponseWriter, r *http.Request) error {
	values, err := utils.GetURLQueryParams(r.URL.String())
	if err != nil {
		return err
	}

	code := values.Get("code")

	id := config.Conf.GoogleDrive.Key
	secret := config.Conf.GoogleDrive.Secret

	formValues := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     id,
		"client_secret": secret,
		"code":          code,
		"redirect_uri":  config.Conf.RedirectURL,
	}

	endpoint := "https://oauth2.googleapis.com/token"

	body, err := utils.SendPOSTRequest(endpoint, formValues, &googleDriveAuth{})
	if err != nil {
		return err
	}

	return utils.JSONResponse(w, body)
}

// RefreshGoogleDriveToken generates additional access tokens after the initial
// token has expired.
func RefreshGoogleDriveToken(w http.ResponseWriter, r *http.Request) error {
	values, err := utils.GetURLQueryParams(r.URL.String())
	if err != nil {
		return err
	}

	refreshToken := values.Get("refresh_token")
	if refreshToken == "" {
		return errors.New("refresh token not specified")
	}

	id := config.Conf.GoogleDrive.Key
	secret := config.Conf.GoogleDrive.Secret

	formValues := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     id,
		"client_secret": secret,
		"refresh_token": refreshToken,
	}

	endpoint := "https://oauth2.googleapis.com/token"

	body, err := utils.SendPOSTRequest(endpoint, formValues, &googleDriveAuth{})
	if err != nil {
		return err
	}

	return utils.JSONResponse(w, body)
}

// SaveToGoogleDrive saves the requested photo to the current user's
// Google Drive account.
func SaveToGoogleDrive(w http.ResponseWriter, r *http.Request) error {
	values, err := utils.GetURLQueryParams(r.URL.String())
	if err != nil {
		return err
	}

	token := values.Get("token")
	id := values.Get("id")
	url := values.Get("url")

	_, err = unsplash.TrackPhotoDownload(id)
	if err != nil {
		return err
	}

	v := fmt.Sprintf("Bearer %s", token)

	ctx, cncl := context.WithTimeout(
		context.Background(),
		time.Second*saveToDriveTimeout,
	)
	defer cncl()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		http.NoBody,
	)
	if err != nil {
		return err
	}

	resp, err := utils.Client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := utils.CheckForErrors(resp)
	if err != nil {
		return err
	}

	endpoint := "https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart"

	// Metadata content.
	// New multipart writer.
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	// Metadata part.
	metadata := fmt.Sprintf(`{"name": "photo-%s.jpeg"}`, id)

	metadataHeader := textproto.MIMEHeader{}

	metadataHeader.Set("Content-Type", "application/json; charset=UTF-8")

	part, _ := writer.CreatePart(metadataHeader)

	_, err = part.Write([]byte(metadata))
	if err != nil {
		return err
	}

	mediaHeader := textproto.MIMEHeader{}
	mediaHeader.Set("Content-Type", "image/jpeg")

	mediaPart, _ := writer.CreatePart(mediaHeader)

	_, err = io.Copy(mediaPart, bytes.NewReader(respBody))
	if err != nil {
		return err
	}

	// Close multipart writer.
	writer.Close()

	request, err = http.NewRequest(
		"POST",
		endpoint,
		bytes.NewReader(body.Bytes()),
	)
	if err != nil {
		return err
	}

	contentType := fmt.Sprintf(
		"multipart/related; boundary=%s",
		writer.Boundary(),
	)

	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Authorization", v)
	request.Header.Set("Content-Length", fmt.Sprintf("%d", body.Len()))

	response, err := utils.Client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	_, err = utils.CheckForErrors(response)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
