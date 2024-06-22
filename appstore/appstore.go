package appstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	SAND_BOX_URL   string = "https://sandbox.itunes.apple.com/verifyReceipt"
	PRODUCTION_URL string = "https://buy.itunes.apple.com/verifyReceipt"
)

var (
	UNABLE_JSON_OBJECT     = errors.New("The App Store could not read the JSON object you provided.")
	DATA_MISSING           = errors.New("The data in the receipt-data property was malformed or missing.")
	UNABLE_AUTHENTICATED   = errors.New("The receipt could not be authenticated.")
	KEY_IS_INCORRECT       = errors.New("The shared secret you provided does not match the shared secret on file for your account.")
	SERVER_IS_UNABLE       = errors.New("The receipt server is not currently available.")
	TEST_ENVIRONMENT       = errors.New("This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead.")
	PRODUCTION_ENVIRONMENT = errors.New("This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead.")
	ARE_YOU_KIDDING_ME     = errors.New("This receipt could not be authorized. Treat this the same as if a purchase was never made.")
	INTERNAL_ERROR         = errors.New("Internal data access error.")
	UNKNOWN_ERROR          = errors.New("An unknown error occurred")
)

// https://developer.apple.com/library/content/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
// The IAPRequest type has the request parameter
type IAPRequest struct {
	ReceiptData string `json:"receipt-data"`
	// Only used for receipts that contain auto-renewable subscriptions.
	Password string `json:"password,omitempty"`
	// Only used for iOS7 style app receipts that contain auto-renewable or non-renewing subscriptions.
	// If value is true, response includes only the latest renewal transaction for any subscriptions.
	ExcludeOldTransactions bool `json:"exclude-old-transactions"`
}

// The IAPResponse type has the response properties
// We defined each field by the current IAP response, but some fields are not mentioned
// in the following Apple's document;
// https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html
// If you get other types or fields from the IAP response, you should use the struct you defined.
type IAPResponse struct {
	Status      int     `json:"status"`
	Environment string  `json:"environment"`
	Receipt     Receipt `json:"receipt"`
	IsRetryable bool    `json:"is-retryable,omitempty"`
}

func (this *IAPResponse) GetOrder(transactionID string) (orderID, productID, packageName string, purchaseTime int64, err error) {
	for _, inApp := range this.Receipt.InApps {
		if inApp.TransactionID != transactionID {
			continue
		}
		unix, err := strconv.Atoi(inApp.PurchaseDateMS)
		if err != nil {
			return inApp.TransactionID, inApp.ProductID, this.Receipt.BundleID, time.Now().Unix(), nil
		}
		return inApp.TransactionID, inApp.ProductID, this.Receipt.BundleID, int64(unix) / int64(1000), nil
	}
	return orderID, productID, packageName, purchaseTime, errors.New("No TransactionID")
}

// The ReceiptCreationDate type indicates the date when the app receipt was created.
type ReceiptCreationDate struct {
	CreationDate    string `json:"receipt_creation_date"`
	CreationDateMS  string `json:"receipt_creation_date_ms"`
	CreationDatePST string `json:"receipt_creation_date_pst"`
}

// The RequestDate type indicates the date and time that the request was sent
type RequestDate struct {
	RequestDate    string `json:"request_date"`
	RequestDateMS  string `json:"request_date_ms"`
	RequestDatePST string `json:"request_date_pst"`
}

// The PurchaseDate type indicates the date and time that the item was purchased
type PurchaseDate struct {
	PurchaseDate    string `json:"purchase_date"`
	PurchaseDateMS  string `json:"purchase_date_ms"`
	PurchaseDatePST string `json:"purchase_date_pst"`
}

// The OriginalPurchaseDate type indicates the beginning of the subscription period
type OriginalPurchaseDate struct {
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	OriginalPurchaseDateMS  string `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePST string `json:"original_purchase_date_pst"`
}

// The ExpiresDate type indicates the expiration date for the subscription
type ExpiresDate struct {
	ExpiresDate             string `json:"expires_date,omitempty"`
	ExpiresDateMS           string `json:"expires_date_ms,omitempty"`
	ExpiresDatePST          string `json:"expires_date_pst,omitempty"`
	ExpiresDateFormatted    string `json:"expires_date_formatted,omitempty"`
	ExpiresDateFormattedPST string `json:"expires_date_formatted_pst,omitempty"`
}

// The CancellationDate type indicates the time and date of the cancellation by Apple customer support
type CancellationDate struct {
	CancellationDate    string `json:"cancellation_date,omitempty"`
	CancellationDateMS  string `json:"cancellation_date_ms,omitempty"`
	CancellationDatePST string `json:"cancellation_date_pst,omitempty"`
}

// The GracePeriodDate type indicates the grace period date for the subscription
type GracePeriodDate struct {
	GracePeriodDate    string `json:"grace_period_expires_date,omitempty"`
	GracePeriodDateMS  string `json:"grace_period_expires_date_ms,omitempty"`
	GracePeriodDatePST string `json:"grace_period_expires_date_pst,omitempty"`
}

// The InApp type has the receipt attributes
type InApp struct {
	Quantity              string `json:"quantity"`
	ProductID             string `json:"product_id"`
	TransactionID         string `json:"transaction_id"`
	OriginalTransactionID string `json:"original_transaction_id"`
	WebOrderLineItemID    string `json:"web_order_line_item_id,omitempty"`

	IsTrialPeriod        string `json:"is_trial_period"`
	IsInIntroOfferPeriod string `json:"is_in_intro_offer_period,omitempty"`
	ExpiresDate

	PurchaseDate
	OriginalPurchaseDate

	CancellationDate
	CancellationReason string `json:"cancellation_reason,omitempty"`
}

// The Receipt type has whole data of receipt
type Receipt struct {
	ReceiptType                string  `json:"receipt_type"`
	AdamID                     int64   `json:"adam_id"`
	AppItemID                  int32   `json:"app_item_id"`
	BundleID                   string  `json:"bundle_id"`
	ApplicationVersion         string  `json:"application_version"`
	DownloadID                 int64   `json:"download_id"`
	VersionExternalIdentifier  int32   `json:"version_external_identifier"`
	OriginalApplicationVersion string  `json:"original_application_version"`
	InApps                     []InApp `json:"in_app"`
	ReceiptCreationDate
	RequestDate
	OriginalPurchaseDate
}

// A pending renewal may refer to a renewal that is scheduled in the future or a renewal that failed in the past for some reason.
type PendingRenewalInfo struct {
	SubscriptionExpirationIntent   string `json:"expiration_intent"`
	SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
	SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
	SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
	SubscriptionPriceConsentStatus string `json:"price_consent_status"`
	ProductID                      string `json:"product_id"`
	OriginalTransactionID          string `json:"original_transaction_id"`

	GracePeriodDate
}

type AppStore struct {
	URL string
}

func (this *AppStore) Verify(receipt string) (response *IAPResponse, err error) {
	iap := IAPRequest{
		ReceiptData: receipt,
	}
	b := new(bytes.Buffer)
	if err = json.NewEncoder(b).Encode(iap); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", this.URL, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return nil, err
	}
	switch {
	case response.Status == 0:
		return response, nil
	case response.Status == 21000:
		return response, UNABLE_JSON_OBJECT
	case response.Status == 21002:
		return response, DATA_MISSING
	case response.Status == 21003:
		return response, UNABLE_AUTHENTICATED
	case response.Status == 21004:
		return response, KEY_IS_INCORRECT
	case response.Status == 21005:
		return response, SERVER_IS_UNABLE
	case response.Status == 21007:
		return response, TEST_ENVIRONMENT
	case response.Status == 21008:
		return response, PRODUCTION_ENVIRONMENT
	case response.Status == 21010:
		return response, ARE_YOU_KIDDING_ME
	case response.Status >= 21100 && response.Status <= 21199:
		return response, INTERNAL_ERROR
	default:
		return response, UNKNOWN_ERROR
	}
}
