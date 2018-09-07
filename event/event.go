package event

import (
	"encoding/json"
)

// State of a host or service in a event.
type State float64

// Event state values.
const (
	StateOK       State = 0.0
	StateWarning  State = 1.0
	StateCritical State = 2.0
	StateUnknown  State = 3.0
)

func (s State) String() string {
	switch s {
	case StateOK:
		return "OK"
	case StateWarning:
		return "WARNING"
	case StateCritical:
		return "CRITICAL"
	case StateUnknown:
		return "UNKOWN"
	}
	return ""
}

// StateType of a host or service state.
//
// Events with StateTypeSoft are before max_check_attempts are done (no notification is sent). After all re-checks have also failed, StateTypeHard will be set.
//
// See http://docs.icinga.org/icinga2/latest/doc/module/icinga2/chapter/monitoring-basics#hard-soft-states for more details.
type StateType float64

// Possible StateType values
const (
	StateTypeSoft StateType = 0.0
	StateTypeHard StateType = 1.0
)

func (s StateType) String() string {
	if s == StateTypeSoft {
		return "SOFT"
	}
	return "HARD"
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

//type Event interface{}

type event struct {
	Timestamp float64
	Type      StreamType
}

type CheckResultVars struct {
	Attempt   float64
	Reachable bool
	State     State
	StateType StateType `json:"state_type"`
}

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
