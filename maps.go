package neato

import (
	"fmt"
	"strconv"
)

type Map struct {
	Version                        *int     `json:"version"`
	ID                             string   `json:"id"`
	URL                            string   `json:"url"`
	URLValidForSeconds             *int     `json:"url_valid_for_seconds"`
	RunID                          *string  `json:"run_id"`
	Status                         *string  `json:"status"`
	LaunchedFrom                   *string  `json:"launched_from"`
	Error                          *string  `json:"error"`
	Category                       *int     `json:"category"`
	Mode                           *int     `json:"mode"`
	Modifier                       *int     `json:"modifier"`
	StartAt                        *string  `json:"start_at"`
	EndAt                          *string  `json:"end_at"`
	EndOrientationRelativeDegrees  int      `json:"end_orientation_relative_degrees"`
	RunChargeAtStart               int      `json:"run_charge_at_start"`
	RunChargeAtEnd                 int      `json:"run_charge_at_end"`
	SuspendedCleaningChargingCount *int     `json:"suspended_cleaning_charging_count"`
	TimeInSuspendedCleaning        *int     `json:"time_in_suspended_cleaning"`
	TimeInError                    *int     `json:"time_in_error"`
	TimeInPause                    *int     `json:"time_in_pause"`
	CleanedArea                    *float64 `json:"cleaned_area"`
	BaseCount                      *int     `json:"base_count"`
	IsDocked                       *bool    `json:"is_docked"`
	Delocalized                    *bool    `json:"delocalized"`
	GeneratedAt                    *string  `json:"generated_at"`
	PersistentMapID                *string  `json:"persistent_map_id"`
	ValidAsPersistentMap           *bool    `json:"calid_as_persistent_map"`
	NavigationMode                 *int     `json:"navigation_mode"`
}

func (m *Map) String() string {
	cleanedArea := "<not set>"
	if m.CleanedArea != nil {
		cleanedArea = strconv.FormatFloat(*m.CleanedArea, 'g', 2, 64)
	}
	errStr := "<not set>"
	if m.Error != nil {
		errStr = *m.Error
	}
	return fmt.Sprintf("ID: '%s', URL: %s, Error: %s, Cleaned area: %s sqm", m.ID, m.URL, errStr, cleanedArea)
}
