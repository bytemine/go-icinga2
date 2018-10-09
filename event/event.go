package event

import (
	"encoding/json"
)

// State of a host or service in a event.
type State float64

// Event state values.
//
// StateNil isn't an official Icinga2 value, but sometimes a nil value is required (e.g. we don't care for the value).
// We have to live with -1 as "zero" value here, as the real zero is already in use for OK.
const (
	StateNil      State = -1
	StateOK       State = 0.0
	StateWarning  State = 1.0
	StateCritical State = 2.0
	StateUnknown  State = 3.0

	StateStringNil      = ""
	StateStringOK       = "OK"
	StateStringWarning  = "WARNING"
	StateStringCritical = "CRITICAL"
	StateStringUnknown  = "UNKNOWN"
)

func (s State) String() string {
	switch s {
	case StateOK:
		return StateStringOK
	case StateWarning:
		return StateStringWarning
	case StateCritical:
		return StateStringCritical
	case StateUnknown:
		return StateStringUnknown
	}
	return StateStringNil
}

// NewState returns a new State by its name (StateString... constants).
func NewState(name string) State {
	switch name {
	case StateStringOK:
		return StateOK
	case StateStringWarning:
		return StateWarning
	case StateStringCritical:
		return StateCritical
	case StateStringUnknown:
		return StateUnknown
	default:
		return StateNil
	}
}

// StateType of a host or service state.
//
// Events with StateTypeSoft are before max_check_attempts are done (no notification is sent). After all re-checks have also failed, StateTypeHard will be set.
//
// See http://docs.icinga.org/icinga2/latest/doc/module/icinga2/chapter/monitoring-basics#hard-soft-states for more details.
type StateType float64

// Possible StateType values.
//
// Similar to State, we define StateTypeNil here.
const (
	StateTypeNil  StateType = -1
	StateTypeSoft StateType = 0.0
	StateTypeHard StateType = 1.0

	StateTypeStringSoft = "SOFT"
	StateTypeStringHard = "HARD"
	StateTypeStringNil  = ""
)

func (s StateType) String() string {
	switch s {
	case StateTypeSoft:
		return StateTypeStringSoft
	case StateTypeHard:
		return StateTypeStringHard
	default:
		return StateTypeStringNil
	}
}

// NewStateType returns a StateType by its name (StateTypeString... constants).
func NewStateType(name string) StateType {
	switch name {
	case StateTypeStringSoft:
		return StateTypeSoft
	case StateTypeStringHard:
		return StateTypeHard
	default:
		return StateTypeNil
	}
}

// NotificationType of a sent notification.
//
// See http://docs.icinga.org/icinga2/latest/doc/module/icinga2/toc#!/icinga2/latest/doc/module/icinga2/chapter/monitoring-basics#notifications for more details.
type NotificationType string

// Possible NotificationType values
const (
	NotificationDowntimeStart   NotificationType = "DOWNTIMESTART"
	NotificationDowntimeEnd     NotificationType = "DOWNTIMEEND"
	NotificationDowntimeRemoved NotificationType = "DOWNTIMECANCELLED"
	NotificationCustom          NotificationType = "CUSTOM"
	NotificationAcknowledgement NotificationType = "ACKNOWLEDGEMENT"
	NotificationProblem         NotificationType = "PROBLEM"
	NotificationRecovery        NotificationType = "RECOVERY"
	NotificationFlappingStart   NotificationType = "FLAPPINGSTART"
	NotificationFlappingEnd     NotificationType = "FLAPPINGEND"
	NotificationUnknown         NotificationType = "UNKNOWN_NOTIFICATION"
)

// AcknowledgementType of an acknowledgement.
//
// Alas, nothing more is found in the Icinga2 documentation.
type AcknowledgementType int

// Possible AcknowledgementType values
const (
	AcknowledgementNone   AcknowledgementType = 0
	AcknowledgementNormal AcknowledgementType = 1
	AcknowledgementSticky AcknowledgementType = 2
)

type event struct {
	Timestamp float64
	Type      StreamType
}

// CheckResultVars is used in CheckResultData for the fields VarsAfter and VarsBefore.
// It is exported to permit the manual construction of CheckResult values.
type CheckResultVars struct {
	Attempt   float64
	Reachable bool
	State     State
	StateType StateType `json:"state_type"`
}

// CheckResultData is used in CheckResult for the field CheckResult.
// It is exported to permit the manual construction of CheckResult values.
type CheckResultData struct {
	Active      bool
	CheckSource string `json:"check_source"`
	// Command is strings mixed with numbers (and maybe other types, who knows..), so use a []interface{} here.
	Command         []interface{}
	ExecutionEnd    float64 `json:"execution_end"`
	ExecutionStart  float64 `json:"execution_start"`
	ExitStatus      float64 `json:"exit_status"`
	Output          string
	PerformanceData json.RawMessage `json:"performance_data"`
	ScheduleEnd     float64         `json:"schedule_end"`
	ScheduleStart   float64         `json:"schedule_start"`
	State           State
	Type            StreamType
	VarsAfter       CheckResultVars `json:"vars_after"`
	VarsBefore      CheckResultVars `json:"vars_before"`
}

// CheckResult events are the results of a check of a host/service.
type CheckResult struct {
	event

	// Host affected by this event.
	Host string

	// Service affected by this event.
	Service string

	// CheckResult has further data for the check that triggered this event.
	CheckResult CheckResultData `json:"check_result"`
}

// StateChange events are the result of a state change due to a check failing / going back to OK.
type StateChange struct {
	event

	// Host affected by this event.
	Host string

	// Service affected by this event.
	Service string

	// CheckResult has further data for the check that triggered this event.
	CheckResult CheckResultData `json:"check_result"`

	// State for the host/service due to this event. See State constants.
	State     State
	StateType StateType `json:"state_type"`
}

// Notification sent event.
type Notification struct {
	event
	// Host affected by this event.
	Host string

	// Service affected by this event.
	Service string

	// CheckResult has further data for the check that triggered this event.
	CheckResult CheckResultData `json:"check_result"`

	// Users receiving a notification
	Users []string

	// Author of a notification or an acknowledgement
	Author string

	// Text of a notification
	Text string

	// NotificationType, see constants.
	NotificationType NotificationType `json:"notification_type"`
}

// AcknowledgementSet for to a notification.
type AcknowledgementSet struct {
	event

	// Host affected by this event.
	Host string

	// Service affected by this event.
	Service string

	// State for the host/service due to this event. See State constants.
	State     State
	StateType StateType `json:"state_type"`

	// Author of a notification or an acknowledgement
	Author string

	Comment string

	// AcknowledgementType, see constants.
	AcknowledgementType AcknowledgementType `json:"acknowledgement_type"`

	// Notify of an acknowledgement
	Notify bool

	// Expiry of an acknowledgement
	Expiry float64
}

// AcknowledgementCleared for a notification.
type AcknowledgementCleared struct {
	event

	// Host affected by this event.
	Host string

	// Service affected by this event.
	Service string

	// State for the host/service due to this event. See State constants.
	State     State
	StateType StateType `json:"state_type"`
}

type CommentAdded struct {
	event

	// Comment data
	Comment json.RawMessage
}

type CommentRemoved struct {
	event

	// Comment data
	Comment json.RawMessage
}

type DowntimeAdded struct {
	event

	// Downtime data
	Downtime json.RawMessage
}

type DowntimeRemoved struct {
	event

	// Downtime data
	Downtime json.RawMessage
}

type DowntimeTriggered struct {
	event

	// Downtime data
	Downtime json.RawMessage
}
