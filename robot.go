package neato

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

type Robot struct {
	session *PasswordSession

	Serial                            string   `json:"serial"`
	Prefix                            *string  `json:"prefix"`
	Name                              string   `json:"name"`
	Model                             *string  `json:"model"`
	Firmware                          *string  `json:"firmware"`
	Timezone                          *string  `json:"timezone"`
	SecretKey                         string   `json:"secret_key"`
	PurchasedAt                       *string  `json:"purchased_at"`
	LinkedAt                          *string  `json:"linked_at"`
	NucleoURL                         string   `json:"nucleo_url"`
	Traits                            []string `json:"traits"`
	ProofOfPurchaseURL                string   `json:"proof_of_purchase_url"`
	ProofOfPurchaseURLValidForSeconds int      `json:"proof_of_purchase_url_valid_for_seconds"`
	ProofOfPurchaseGeneratedAt        *string  `json:"proof_of_purchase_generated_at"`
	MACAddress                        *string  `json:"mac_address"`
	CreatedAt                         *string  `json:"created_at"`
}

func (r *Robot) String() string {
	model := "<not set>"
	if r.Model != nil {
		model = *r.Model
	}
	return fmt.Sprintf("Name: '%s', Serial: %s, Model: %s", r.Name, r.Serial, model)
}

func (r *Robot) Header(body []byte) *url.Values {
	header := url.Values{}

	header.Set("Accept", "application/vnd.neato.nucleo.v1")

	// RFC2616-formatted date and time
	date := time.Now().In(time.UTC).Format("Tue, 2 Jan 2006 15:04:05 GMT")
	header.Set("Date", date)

	msg := fmt.Sprintf("%s\n%s\n%s", strings.ToLower(r.Serial), date, body)
	h := hmac.New(sha256.New, []byte(r.SecretKey))
	h.Write([]byte(msg))
	signature := hex.EncodeToString(h.Sum(nil))
	header.Set("Authorization", "NEATOAPP "+signature)

	log.Printf("Auth %s", signature)
	log.Printf("Date %s", date)

	return &header
}

type Result string

var (
	ResultOK              Result = "ok"
	ResultInvalidJSON     Result = "invalid_json"
	ResultBadRequest      Result = "bad_request"
	ResultCommandNotFound Result = "command_not_found"
	ResultCommandRejected Result = "command_rejected"
	ResultKO              Result = "ko"
	ResultNotOnChargeBase Result = "not_on_charge_base"
)

type RobotState struct {
	Version  int         `json:"version"`
	ReqID    string      `json:"reqId"`
	Result   Result      `json:"result"`
	Data     interface{} `json:"data"`
	State    int         `json:"state"`
	Action   int         `json:"action"`
	Error    *string     `json:"error"`
	Alert    *string     `json:"alert"`
	Cleaning struct {
		Category       Category       `json:"category"`
		Mode           CleaningMode   `json:"mode"`
		Modifier       int            `json:"modifier"`
		NavigationMode NavigationMode `json:"navigationMode"`
		SpotWidth      int            `json:"spotWidth"`
		SpotHeight     int            `json:"spotHeight"`
	} `json:"cleaning"`
	Details struct {
		IsCharging        bool `json:"isCharging"`
		IsDocked          bool `json:"isDocked"`
		DockHasBeenSeen   bool `json:"dockHasBeenSeen"`
		Charge            int  `json:"charge"`
		IsScheduleEnabled bool `json:"isScheduleEnabled"`
	} `json:"details"`
	AvailableCommands struct {
		Start    bool `json:"start"`
		Stop     bool `json:"stop"`
		Pause    bool `json:"pause"`
		Resume   bool `json:"resume"`
		GoToBase bool `json:"goToBase"`
	} `json:"availableCommands"`
	AvailableServices struct {
		FindMe         string `json:"findMe"`
		GeneralInfo    string `json:"generalInfo"`
		HouseCleaning  string `json:"houseCleaning"`
		LocalStats     string `json:"localStats"`
		ManualCleaning string `json:"manualCleaning"`
		Maps           string `json:"maps"`
		Preferences    string `json:"preferences"`
		Schedule       string `json:"schedule"`
		SpotCleaning   string `json:"spotCleaning"`
		IECTest        string `json:"IECTest"`
		LogCopy        string `json:"logCopy"`
		SoftwareUpdate string `json:"softwareUpdate"`
		Wifi           string `json:"wifi"`
	} `json:"availableServices"`
	Meta struct {
		ModelName string `json:"modelName"`
		Firmware  string `json:"firmware"`
	} `json:"meta"`
}

func (r *Robot) State() (*RobotState, error) {
	var resp RobotState
	dataMap := map[string]string{
		"reqId": "1",
		"cmd":   "getRobotState",
	}
	if err := r.post(dataMap, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type CleaningMode int

var (
	CleaningModeEco   CleaningMode = 1
	CleaningModeTurbo CleaningMode = 2
)

type NavigationMode int

var (
	NavigationModeNormal    NavigationMode = 1
	NavigationModeExtraCare NavigationMode = 2
	NavigationModeDeep      NavigationMode = 3
)

type Category int

var (
	CategoryNonPersistentMap Category = 2
	CategoryPersistentMap    Category = 4
)

type CleaningOptions struct {
	CleaningMode   CleaningMode
	NavigationMode NavigationMode
	Category       Category
	BoundaryID     int
	MapID          int
}

func NewCleaningOptions(hasPersistentMaps bool) *CleaningOptions {
	category := CategoryNonPersistentMap
	if hasPersistentMaps {
		category = CategoryPersistentMap
	}
	return &CleaningOptions{
		CleaningMode:   CleaningModeEco,
		NavigationMode: NavigationModeNormal,
		Category:       category,
	}
}

func (r *Robot) HasPersistentMaps() bool {
	/*
		serviceVersion := r.State.AvailableServices
		return r.HasPersistentMaps && (serviceVersion == "basic-3" || serviceVersion == "basic-4")
	*/
	return false
}

func (r *Robot) Start(opts *CleaningOptions) error {
	/*
		if opts == nil {
			opts = NewCleaningOptions(r.HasPersistentMaps())
		}
	*/
	return nil
}

func (r *Robot) post(dataMap map[string]string, response interface{}) error {
	// remove port from nucleo URL
	uri, err := url.Parse(r.NucleoURL)
	if err != nil {
		return fmt.Errorf("failed to parse NucleoURL '%s': %v", r.NucleoURL, err)
	}
	host, _, err := net.SplitHostPort(uri.Host)
	if err != nil {
		return fmt.Errorf("failed to split host:port for '%s': %v", uri.Host, err)
	}
	uri.Host = host
	uri.Path += "/vendors/neato/robots/" + r.Serial + "/messages"

	// skip TLS verification :( This will otherwise fail with the message
	//   "x509: “*.neatocloud.com” certificate is not trusted"
	skipVerification := true

	body, err := json.Marshal(dataMap)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	return httpPost(uri.String(), r.Header(body), dataMap, skipVerification, response)
}
