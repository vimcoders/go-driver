package googleplay

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

var client *Client

func Connect(jsonKey []byte) error {
	if len(jsonKey) <= 0 {
		return nil
	}
	c, err := New(jsonKey)
	if err != nil {
		return err
	}
	client = c
	return nil
}

func VerifyProduct(ctx context.Context, packageName, productID, token string) (*androidpublisher.ProductPurchase, error) {
	if client == nil {
		return nil, errors.New("client == nil")
	}
	return client.VerifyProduct(ctx, packageName, productID, token)
}

func ConsumeProduct(ctx context.Context, packageName, productID, token string) error {
	if client == nil {
		return errors.New("client == nil")
	}
	return client.ConsumeProduct(ctx, packageName, productID, token)
}

// The Client type implements VerifySubscription method
type Client struct {
	service *androidpublisher.Service
}

// New returns http client which includes the credentials to access androidpublisher API.
// You should create a service account for your project at
// https://console.developers.google.com and download a JSON key file to set this argument.
func New(jsonKey []byte) (*Client, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)
	conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}
	val := conf.Client(ctx).Transport.(*oauth2.Transport)
	_, err = val.Source.Token()
	if err != nil {
		return nil, err
	}
	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}
	return &Client{service}, err
}

func (x *Client) VerifyProduct(ctx context.Context, packageName, productID, token string) (*androidpublisher.ProductPurchase, error) {
	ps := androidpublisher.NewPurchasesProductsService(x.service)
	return ps.Get(packageName, productID, token).Context(ctx).Do()
}

func (x *Client) ConsumeProduct(ctx context.Context, packageName, productID, token string) error {
	ps := androidpublisher.NewPurchasesProductsService(x.service)
	return ps.Consume(packageName, productID, token).Context(ctx).Do()
}

// VerifySignature verifies in app billing signature.
// You need to prepare a public key for your Android app's in app billing
// at https://play.google.com/apps/publish/
func VerifySignature(base64EncodedPublicKey string, receipt []byte, signature string) (err error) {
	// prepare public key
	decodedPublicKey, err := base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key")
	}
	publicKey, _ := publicKeyInterface.(*rsa.PublicKey)
	// generate hash value from receipt
	hasher := sha1.New()
	hasher.Write(receipt)
	hashedReceipt := hasher.Sum(nil)
	// decode signature
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature")
	}
	// verify
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashedReceipt, decodedSignature); err != nil {
		return err
	}
	return nil
}
