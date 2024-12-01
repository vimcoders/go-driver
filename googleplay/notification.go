package googleplay

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/vimcoders/go-driver/log"

	"google.golang.org/api/pubsub/v1"
)

// // https://developer.android.com/google/play/billing/rtdn-reference#sub
type SubscriptionNotificationType int

const (
	SubscriptionNotificationTypeRecovered SubscriptionNotificationType = iota + 1
	SubscriptionNotificationTypeRenewed
	SubscriptionNotificationTypeCanceled
	SubscriptionNotificationTypePurchased
	SubscriptionNotificationTypeAccountHold
	SubscriptionNotificationTypeGracePeriod
	SubscriptionNotificationTypeRestarted
	SubscriptionNotificationTypePriceChangeConfirmed
	SubscriptionNotificationTypeDeferred
	SubscriptionNotificationTypePaused
	SubscriptionNotificationTypePauseScheduleChanged
	SubscriptionNotificationTypeRevoked
	SubscriptionNotificationTypeExpired
)

// // https://developer.android.com/google/play/billing/rtdn-reference#one-time
type OneTimeProductNotificationType int

const (
	OneTimeProductNotificationTypePurchased OneTimeProductNotificationType = iota + 1
	OneTimeProductNotificationTypeCanceled
)

// DeveloperNotification is sent by a Pub/Sub topic.
// Detailed description is following.
// https://developer.android.com/google/play/billing/rtdn-reference#json_specification
type DeveloperNotification struct {
	Version                    string                     `json:"version"`
	PackageName                string                     `json:"packageName"`
	EventTimeMillis            string                     `json:"eventTimeMillis"`
	SubscriptionNotification   SubscriptionNotification   `json:"subscriptionNotification,omitempty"`
	OneTimeProductNotification OneTimeProductNotification `json:"oneTimeProductNotification,omitempty"`
	TestNotification           TestNotification           `json:"testNotification,omitempty"`
}

// SubscriptionNotification https://developer.android.google.cn/google/play/billing/rtdn-reference#sub
type SubscriptionNotification struct {
	Version          string `json:"version,omitempty"`          // 版本
	NotificationType int    `json:"notificationType,omitempty"` // 变化类型
	PurchaseToken    string `json:"purchaseToken,omitempty"`    // 支付令牌
	SubscriptionId   string `json:"subscriptionId,omitempty"`   // 商品 id
}

type OneTimeProductNotification struct {
	Version          string `json:"version,omitempty"`          // 版本
	NotificationType int    `json:"notificationType,omitempty"` // 变化类型
	PurchaseToken    string `json:"purchaseToken,omitempty"`    // 支付令牌
	SKU              string `json:"sku,omitempty"`
}

type TestNotification struct {
	Version string `json:"version,omitempty"` // 版本
}

type SubscriptionData struct {
	Version                    string                      `json:"version,omitempty"`
	PackageName                string                      `json:"packageName,omitempty"`
	EventTimeMillis            string                      `json:"eventTimeMillis,omitempty"`
	SubscriptionNotification   *SubscriptionNotification   `json:"subscriptionNotification,omitempty"`
	OneTimeProductNotification *OneTimeProductNotification `json:"oneTimeProductNotification,omitempty"`
	TestNotification           *TestNotification           `json:"testNotification,omitempty"`
}

type PushRequest struct {
	Message      pubsub.PubsubMessage `json:"message"`
	Subscription string               `json:"subscription"`
}

func GooglePayNotify(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var pushRequest PushRequest
	if err := json.Unmarshal(body, &pushRequest); err != nil {
		log.Error(err.Error())
		return
	}
	dataByte, err := base64.StdEncoding.DecodeString(pushRequest.Message.Data)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var subscriptionData = &SubscriptionData{}
	err = json.Unmarshal(dataByte, subscriptionData)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if subscriptionData.SubscriptionNotification == nil {
		log.Error(err.Error())
		return
	}
	switch subscriptionData.OneTimeProductNotification.NotificationType {
	case int(OneTimeProductNotificationTypePurchased):
	case int(OneTimeProductNotificationTypeCanceled):
		// 自行根据状态来处理
	}
	// 返回 200 状态通知 google 这个通知已经接受
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// // receiveMessagesHandler validates authentication token and caches the Pub/Sub
// // message received.
// func receiveMessagesHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Verify that the request originates from the application.
// 	// a.pubsubVerificationToken = os.Getenv("PUBSUB_VERIFICATION_TOKEN")
// 	// if token, ok := r.URL.Query()["token"]; !ok || len(token) != 1 || token[0] != a.pubsubVerificationToken {
// 	// 	http.Error(w, "Bad token", http.StatusBadRequest)
// 	// 	return
// 	// }

// 	// Get the Cloud Pub/Sub-generated JWT in the "Authorization" header.
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" || len(strings.Split(authHeader, " ")) != 2 {
// 		http.Error(w, "Missing Authorization header", http.StatusBadRequest)
// 		return
// 	}
// 	token := strings.Split(authHeader, " ")[1]
// 	// Verify and decode the JWT.
// 	// If you don't need to control the HTTP client used you can use the
// 	// convenience method idtoken.Validate instead of creating a Validator.
// 	v, err := idtoken.NewValidator(r.Context(), option.WithHTTPClient(http.DefaultClient))
// 	if err != nil {
// 		http.Error(w, "Unable to create Validator", http.StatusBadRequest)
// 		return
// 	}
// 	// Please change http://example.com to match with the value you are
// 	// providing while creating the subscription.
// 	payload, err := v.Validate(r.Context(), token, "http://example.com")
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Invalid Token: %v", err), http.StatusBadRequest)
// 		return
// 	}
// 	if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
// 		http.Error(w, "Wrong Issuer", http.StatusBadRequest)
// 		return
// 	}

// 	// IMPORTANT: you should validate claim details not covered by signature
// 	// and audience verification above, including:
// 	//   - Ensure that `payload.Claims["email"]` is equal to the expected service
// 	//     account set up in the push subscription settings.
// 	//   - Ensure that `payload.Claims["email_verified"]` is set to true.
// 	if payload.Claims["email"] != "test-service-account-email@example.com" || payload.Claims["email_verified"] != true {
// 		http.Error(w, "Unexpected email identity", http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Fprint(w, "OK")
// }
