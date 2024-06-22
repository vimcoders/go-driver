package appstore

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// rootPEM is from `openssl x509 -inform der -in AppleRootCA-G3.cer -out apple_root.pem`
const rootPEM = `
-----BEGIN CERTIFICATE-----
MIICQzCCAcmgAwIBAgIILcX8iNLFS5UwCgYIKoZIzj0EAwMwZzEbMBkGA1UEAwwS
QXBwbGUgUm9vdCBDQSAtIEczMSYwJAYDVQQLDB1BcHBsZSBDZXJ0aWZpY2F0aW9u
IEF1dGhvcml0eTETMBEGA1UECgwKQXBwbGUgSW5jLjELMAkGA1UEBhMCVVMwHhcN
MTQwNDMwMTgxOTA2WhcNMzkwNDMwMTgxOTA2WjBnMRswGQYDVQQDDBJBcHBsZSBS
b290IENBIC0gRzMxJjAkBgNVBAsMHUFwcGxlIENlcnRpZmljYXRpb24gQXV0aG9y
aXR5MRMwEQYDVQQKDApBcHBsZSBJbmMuMQswCQYDVQQGEwJVUzB2MBAGByqGSM49
AgEGBSuBBAAiA2IABJjpLz1AcqTtkyJygRMc3RCV8cWjTnHcFBbZDuWmBSp3ZHtf
TjjTuxxEtX/1H7YyYl3J6YRbTzBPEVoA/VhYDKX1DyxNB0cTddqXl5dvMVztK517
IDvYuVTZXpmkOlEKMaNCMEAwHQYDVR0OBBYEFLuw3qFYM4iapIqZ3r6966/ayySr
MA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgEGMAoGCCqGSM49BAMDA2gA
MGUCMQCD6cHEFl4aXTQY2e3v9GwOAEZLuN+yRhHFD/3meoyhpmvOwgPUnPWTxnS4
at+qIxUCMG1mihDK1A3UT82NQz60imOlM27jbdoXt2QfyFMm+YhidDkLF1vLUagM
6BgD56KyKA==
-----END CERTIFICATE-----
`

// type AppStoreServerNotification struct {
// 	appleRootCert   string
// 	Payload         *NotificationPayload
// 	TransactionInfo *TransactionInfo
// 	RenewalInfo     *RenewalInfo
// 	IsValid         bool
// }

type AppStoreServerRequest struct {
	SignedPayload string `json:"signedPayload"`
}

// type NotificationHeader struct {
// 	Alg string   `json:"alg"`
// 	X5c []string `json:"x5c"`
// }

// type NotificationPayload struct {
// 	jwt.StandardClaims
// 	NotificationType string              `json:"notificationType"`
// 	Subtype          string              `json:"subtype"`
// 	NotificationUUID string              `json:"notificationUUID"`
// 	Version          string              `json:"version"`
// 	Summary          NotificationSummary `json:"summary"`
// 	Data             NotificationData    `json:"data"`
// }

// type NotificationSummary struct {
// 	RequestIdentifier      string   `json:"requestIdentifier"`
// 	AppAppleId             string   `json:"appAppleId"`
// 	BundleId               string   `json:"bundleId"`
// 	ProductId              string   `json:"productId"`
// 	Environment            string   `json:"environment"`
// 	StoreFrontCountryCodes []string `json:"storefrontCountryCodes"`
// 	FailedCount            int64    `json:"failedCount"`
// 	SucceededCount         int64    `json:"succeededCount"`
// }

// type NotificationData struct {
// 	AppAppleId            int    `json:"appAppleId"`
// 	BundleId              string `json:"bundleId"`
// 	BundleVersion         string `json:"bundleVersion"`
// 	Environment           string `json:"environment"`
// 	SignedRenewalInfo     string `json:"signedRenewalInfo"`
// 	SignedTransactionInfo string `json:"signedTransactionInfo"`
// 	Status                int32  `json:"status"`
// }

type TransactionInfo struct {
	jwt.StandardClaims
	AppAccountToken             string `json:"appAccountToken"`
	BundleId                    string `json:"bundleId"`
	Environment                 string `json:"environment"`
	ExpiresDate                 int    `json:"expiresDate"`
	InAppOwnershipType          string `json:"inAppOwnershipType"`
	IsUpgraded                  bool   `json:"isUpgraded"`
	OfferIdentifier             string `json:"offerIdentifier"`
	OfferType                   int32  `json:"offerType"`
	OriginalPurchaseDate        int    `json:"originalPurchaseDate"`
	OriginalTransactionId       string `json:"originalTransactionId"`
	ProductId                   string `json:"productId"`
	PurchaseDate                int    `json:"purchaseDate"`
	Quantity                    int32  `json:"quantity"`
	RevocationDate              int    `json:"revocationDate"`
	RevocationReason            int32  `json:"revocationReason"`
	SignedDate                  int    `json:"signedDate"`
	StoreFront                  string `json:"storefront"`
	StoreFrontId                string `json:"storefrontId"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier"`
	TransactionId               string `json:"transactionId"`
	TransactionReason           string `json:"transactionReason"`
	Type                        string `json:"type"`
	WebOrderLineItemId          string `json:"webOrderLineItemId"`
}

type RenewalInfo struct {
	jwt.StandardClaims
	AutoRenewProductId          string `json:"autoRenewProductId"`
	AutoRenewStatus             int32  `json:"autoRenewStatus"`
	Environment                 string `json:"environment"`
	ExpirationIntent            int32  `json:"expirationIntent"`
	GracePeriodExpiresDate      int    `json:"gracePeriodExpiresDate"`
	IsInBillingRetryPeriod      bool   `json:"isInBillingRetryPeriod"`
	OfferIdentifier             string `json:"offerIdentifier"`
	OfferType                   int32  `json:"offerType"`
	OriginalTransactionId       string `json:"originalTransactionId"`
	PriceIncreaseStatus         int32  `json:"priceIncreaseStatus"`
	ProductId                   string `json:"productId"`
	RecentSubscriptionStartDate int    `json:"recentSubscriptionStartDate"`
	RenewalDate                 int    `json:"renewalDate"`
	SignedDate                  int    `json:"signedDate"`
}

// DecodeSignedPayload 解析SignedPayload数据
func DecodeSignedPayload(signedPayload string) (payload *NotificationV2Payload, err error) {
	if signedPayload == "" {
		return nil, fmt.Errorf("signedPayload is empty")
	}
	payload = &NotificationV2Payload{}
	if err = ExtractClaims(signedPayload, payload); err != nil {
		return nil, err
	}
	return
}

// ExtractClaims 解析jws格式数据
// signedPayload：jws格式数据
// tran：指针类型的结构体，用于接收解析后的数据
func ExtractClaims(signedPayload string, tran jwt.Claims) (err error) {
	valueOf := reflect.ValueOf(tran)
	if valueOf.Kind() != reflect.Ptr {
		return errors.New("tran must be ptr struct")
	}
	tokenStr := signedPayload
	rootCertStr, err := extractHeaderByIndex(tokenStr, 2)
	if err != nil {
		return err
	}
	intermediaCertStr, err := extractHeaderByIndex(tokenStr, 1)
	if err != nil {
		return err
	}
	if err = verifyCert(rootCertStr, intermediaCertStr); err != nil {
		return err
	}
	_, err = jwt.ParseWithClaims(tokenStr, tran, func(token *jwt.Token) (interface{}, error) {
		return extractPublicKeyFromToken(tokenStr)
	})
	if err != nil {
		return err
	}
	return nil
}

// Per doc: https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.6
func extractPublicKeyFromToken(tokenStr string) (*ecdsa.PublicKey, error) {
	certStr, err := extractHeaderByIndex(tokenStr, 0)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(certStr)
	if err != nil {
		return nil, err
	}
	switch pk := cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		return pk, nil
	default:
		return nil, errors.New("appstore public key must be of type ecdsa.PublicKey")
	}
}

func extractHeaderByIndex(tokenStr string, index int) ([]byte, error) {
	if index > 2 {
		return nil, errors.New("invalid index")
	}
	tokenArr := strings.Split(tokenStr, ".")
	headerByte, err := base64.RawStdEncoding.DecodeString(tokenArr[0])
	if err != nil {
		return nil, err
	}
	type Header struct {
		Alg string   `json:"alg"`
		X5c []string `json:"x5c"`
	}
	header := &Header{}
	err = json.Unmarshal(headerByte, header)
	if err != nil {
		return nil, err
	}
	if len(header.X5c) < index {
		return nil, fmt.Errorf("index[%d] > header.x5c slice len(%d)", index, len(header.X5c))
	}
	certByte, err := base64.StdEncoding.DecodeString(header.X5c[index])
	if err != nil {
		return nil, err
	}
	return certByte, nil
}

func verifyCert(certByte, intermediaCertStr []byte) error {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		return errors.New("failed to parse root certificate")
	}
	interCert, err := x509.ParseCertificate(intermediaCertStr)
	if err != nil {
		return errors.New("failed to parse intermedia certificate")
	}
	intermedia := x509.NewCertPool()
	intermedia.AddCert(interCert)
	cert, err := x509.ParseCertificate(certByte)
	if err != nil {
		return err
	}
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermedia,
	}
	_, err = cert.Verify(opts)
	return err
}

const (
	// 通知类型常量
	// https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
	NotificationTypeV2ConsumptionRequest     = "CONSUMPTION_REQUEST"
	NotificationTypeV2DidChangeRenewalPref   = "DID_CHANGE_RENEWAL_PREF"
	NotificationTypeV2DidChangeRenewalStatus = "DID_CHANGE_RENEWAL_STATUS"
	NotificationTypeV2DidFailToRenew         = "DID_FAIL_TO_RENEW"
	NotificationTypeV2DidRenew               = "DID_RENEW"
	NotificationTypeV2Expired                = "EXPIRED"
	NotificationTypeV2GracePeriodExpired     = "GRACE_PERIOD_EXPIRED"
	NotificationTypeV2OfferRedeemed          = "OFFER_REDEEMED"
	NotificationTypeV2PriceIncrease          = "PRICE_INCREASE"
	NotificationTypeV2Refund                 = "REFUND"
	NotificationTypeV2RefundDeclined         = "REFUND_DECLINED"
	NotificationTypeV2RenewalExtended        = "RENEWAL_EXTENDED"
	NotificationTypeV2Revoke                 = "REVOKE"
	NotificationTypeV2Subscribed             = "SUBSCRIBED"

	// 子类型常量
	// https://developer.apple.com/documentation/appstoreservernotifications/subtype
	SubTypeV2InitialBuy        = "INITIAL_BUY"
	SubTypeV2Resubscribe       = "RESUBSCRIBE"
	SubTypeV2Downgrade         = "DOWNGRADE"
	SubTypeV2Upgrade           = "UPGRADE"
	SubTypeV2AutoRenewEnabled  = "AUTO_RENEW_ENABLED"
	SubTypeV2AutoRenewDisabled = "AUTO_RENEW_DISABLED"
	SubTypeV2Voluntary         = "VOLUNTARY"
	SubTypeV2BillingRetry      = "BILLING_RETRY"
	SubTypeV2PriceIncrease     = "PRICE_INCREASE"
	SubTypeV2GracePeriod       = "GRACE_PERIOD"
	SubTypeV2BillingRecovery   = "BILLING_RECOVERY"
	SubTypeV2Pending           = "PENDING"
	SubTypeV2Accepted          = "ACCEPTED"
)

// https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2
type NotificationV2Req struct {
	SignedPayload string `json:"signedPayload"`
}

// https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2decodedpayload
type NotificationV2Payload struct {
	jwt.StandardClaims
	NotificationType    string `json:"notificationType"`
	Subtype             string `json:"subtype"`
	NotificationUUID    string `json:"notificationUUID"`
	NotificationVersion string `json:"notificationVersion"`
	Data                *Data  `json:"data"`
}

func (d *NotificationV2Payload) DecodeRenewalInfo() (ri *RenewalInfo, err error) {
	if d.Data == nil {
		return nil, fmt.Errorf("data is nil")
	}
	if d.Data.SignedRenewalInfo == "" {
		return nil, fmt.Errorf("data.signedRenewalInfo is empty")
	}
	ri = &RenewalInfo{}
	if err = ExtractClaims(d.Data.SignedRenewalInfo, ri); err != nil {
		return nil, err
	}
	return
}

func (d *NotificationV2Payload) DecodeTransactionInfo() (ti *TransactionInfo, err error) {
	if d.Data == nil {
		return nil, fmt.Errorf("data is nil")
	}
	if d.Data.SignedTransactionInfo == "" {
		return nil, fmt.Errorf("data.signedTransactionInfo is empty")
	}
	ti = &TransactionInfo{}
	if err = ExtractClaims(d.Data.SignedTransactionInfo, ti); err != nil {
		return nil, err
	}
	return
}

// https://developer.apple.com/documentation/appstoreservernotifications/data
type Data struct {
	AppAppleID            int    `json:"appAppleId"`
	BundleID              string `json:"bundleId"`
	BundleVersion         string `json:"bundleVersion"`
	Environment           string `json:"environment"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}
