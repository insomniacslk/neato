package neato

import "fmt"

type Robot struct {
	session                           *PasswordSession
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

func (r *Robot) post(dataMap map[string]string, response interface{}) error {
	uri := r.NucleoURL + "/vendors/neato/robots/" + r.Serial + "/messages"
	return httpPost(uri, nil, dataMap, response)
}
