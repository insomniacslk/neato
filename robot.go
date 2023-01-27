package neato

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
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

type State int

var (
	StateInvalid State = 0
	StateIdle    State = 1
	StateBusy    State = 2
	StatePaused  State = 3
	StateError   State = 4
)

type Action int

var (
	// TODO validate that this is actually invalid. It seems to be used when the
	// robot is docked
	ActionInvalid              Action = 0
	ActionHouseCleaning        Action = 1
	ActionSpotCleaning         Action = 2
	ActionManualCleaning       Action = 3
	ActionDocking              Action = 4
	ActionUserMenuActive       Action = 5
	ActionSuspendedCleaning    Action = 6
	ActionUpdating             Action = 7
	ActionCopyingLogs          Action = 8
	ActionRecoveringLocation   Action = 9
	ActionIECTest              Action = 10
	ActionMapCleaning          Action = 11
	ActionExploringMap         Action = 12
	ActionAcquiringMapIDs      Action = 13
	ActionCreatingMap          Action = 14
	ActionSuspendedExploration Action = 15
)

func (a Action) String() string {
	switch a {
	case ActionInvalid:
		return "invalid"
	case ActionHouseCleaning:
		return "house cleaning"
	case ActionSpotCleaning:
		return "spot cleaning"
	case ActionManualCleaning:
		return "manual cleaning"
	case ActionDocking:
		return "docking"
	case ActionUserMenuActive:
		return "user menu active"
	case ActionSuspendedCleaning:
		return "suspended cleaning"
	case ActionUpdating:
		return "updating"
	case ActionCopyingLogs:
		return "copying logs"
	case ActionRecoveringLocation:
		return "recovering location"
	case ActionIECTest:
		return "IEC test"
	case ActionMapCleaning:
		return "map cleaning"
	case ActionExploringMap:
		return "exploring map"
	case ActionAcquiringMapIDs:
		return "acquiring map IDs"
	case ActionCreatingMap:
		return "creating map"
	case ActionSuspendedExploration:
		return "suspended exploration"
	default:
		return "unknown"
	}
}

func (s State) String() string {
	switch s {
	case StateInvalid:
		return "invalid"
	case StateIdle:
		return "idle"
	case StateBusy:
		return "busy"
	case StatePaused:
		return "paused"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

type RobotState struct {
	Version  int         `json:"version"`
	ReqID    string      `json:"reqId"`
	Result   Result      `json:"result"`
	Data     interface{} `json:"data"`
	State    State       `json:"state"`
	Action   Action      `json:"action"`
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

func (s *RobotState) String() string {
	errStr := "<not set>"
	if s.Error != nil {
		errStr = *s.Error
	}
	alert := "<not set>"
	if s.Alert != nil {
		alert = *s.Alert
	}
	return fmt.Sprintf("State: %s, Action: %s, Error: %s, Alert: %s", s.State, s.Action, errStr, alert)
}

func (r *Robot) State() (*RobotState, error) {
	var resp RobotState
	dataMap := map[string]interface{}{
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
	Category       *Category
	BoundaryID     int
	MapID          int
}

func NewCleaningOptions() *CleaningOptions {
	return &CleaningOptions{
		CleaningMode:   CleaningModeEco,
		NavigationMode: NavigationModeNormal,
		Category:       &CategoryNonPersistentMap,
	}
}

func (r *Robot) Start(opts *CleaningOptions) error {
	if opts == nil {
		opts = NewCleaningOptions()
	}
	state, err := r.State()
	if err != nil {
		return fmt.Errorf("failed to get robot state: %w", err)
	}
	serviceVersion := state.AvailableServices.HouseCleaning
	if opts.Category == nil {
		if serviceVersion == "basic-3" || serviceVersion == "basic-4" {
			opts.Category = &CategoryPersistentMap
		}
	}
	dataMap := map[string]interface{}{
		"reqId": "1",
		"cmd":   "startCleaning",
		"params": map[string]interface{}{
			"category": strconv.FormatInt(int64(*opts.Category), 10),
		},
	}
	switch serviceVersion {
	case "basic-1":
		dataMap["params"].(map[string]interface{})["mode"] = int(opts.CleaningMode)
		dataMap["params"].(map[string]interface{})["modifier"] = 1
	case "basic-2":
		dataMap["params"].(map[string]interface{})["mode"] = int(opts.CleaningMode)
		dataMap["params"].(map[string]interface{})["modifier"] = 1
		dataMap["params"].(map[string]interface{})["navigationMode"] = opts.NavigationMode
	case "minimal-2":
		dataMap["params"].(map[string]interface{})["navigationMode"] = opts.NavigationMode
	case "basic-3", "basic-4":
		dataMap["params"].(map[string]interface{})["mode"] = int(opts.CleaningMode)
		dataMap["params"].(map[string]interface{})["modifier"] = 1
		dataMap["params"].(map[string]interface{})["navigationMode"] = opts.NavigationMode
	}
	var resp RobotState
	if err := r.post(dataMap, &resp); err != nil {
		return fmt.Errorf("start request failed: %w", err)
	}
	if resp.Result != "ok" {
		return fmt.Errorf("start request failed: %s", resp.Result)
	}
	return nil
}

func (r *Robot) Stop() error {
	dataMap := map[string]interface{}{
		"reqId": "1",
		"cmd":   "stopCleaning",
	}
	var resp RobotState
	if err := r.post(dataMap, &resp); err != nil {
		return fmt.Errorf("stop request failed: %w", err)
	}
	if resp.Result != "ok" {
		return fmt.Errorf("stop request failed: %s", resp.Result)
	}
	return nil
}

func (r *Robot) post(dataMap map[string]interface{}, response interface{}) error {
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
