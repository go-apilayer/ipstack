package ipstack

import "fmt"

type ApiErr struct {
	// Most unfortunate this "success" is also not returned in the successful state.
	Success *bool `json:"success"`

	Err struct {
		Code int       `json:"code,omitempty"`
		Type ErrorType `json:"type,omitempty"`
		Info string    `json:"info,omitempty"`
	} `json:"error,omitempty"`
}

func (e *ApiErr) Error() string {
	return fmt.Sprintf("%d: %s", e.Err.Code, e.Err.Info)
}

// ErrorType represents a ipstack error type.
type ErrorType string

const (
	ErrNotFound                 ErrorType = "404_not_found"
	ErrMissingAccessKey         ErrorType = "missing_access_key"
	ErrInvalidAccessKey         ErrorType = "invalid_access_key"
	ErrInactiveUser             ErrorType = "inactive_user"
	ErrInvalidAPIFunction       ErrorType = "invalid_api_function"
	ErrUsageLimitReached        ErrorType = "usage_limit_reached"
	ErrFunctionAccessRestricted ErrorType = "function_access_restricted"
	ErrHTTPSAccessRestricted    ErrorType = "https_access_restricted"
	ErrInvalidFields            ErrorType = "invalid_fields"
	ErrTooManyIPs               ErrorType = "too_many_ips"
	ErrBatchNotSupportedOnPlan  ErrorType = "batch_not_supported_on_plan"
)

func (e ErrorType) Error() string {
	return string(e)
}

// MarshalText satisfies TextMarshaler
func (e ErrorType) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

// UnmarshalText satisfies TextUnmarshaler
func (e *ErrorType) UnmarshalText(text []byte) error {
	typ := ErrorType(text)
	switch typ {
	case ErrNotFound:
		*e = typ
	case ErrMissingAccessKey:
		*e = typ
	case ErrInvalidAccessKey:
		*e = typ
	case ErrInactiveUser:
		*e = typ
	case ErrInvalidAPIFunction:
		*e = typ
	case ErrUsageLimitReached:
		*e = typ
	case ErrFunctionAccessRestricted:
		*e = typ
	case ErrHTTPSAccessRestricted:
		*e = typ
	case ErrInvalidFields:
		*e = typ
	case ErrTooManyIPs:
		*e = typ
	case ErrBatchNotSupportedOnPlan:
		*e = typ
	default:
		return fmt.Errorf("unknown ipstack api error type: %s", typ)
	}

	return nil
}

func codeFromErrorType(typ ErrorType) int {
	switch typ {
	case ErrNotFound:
		return 404 // The requested resource does not exist.
	case ErrMissingAccessKey, ErrInvalidAccessKey:
		// No API Key was specified.
		// or
		// No API Key was specified or an invalid API Key was specified.
		return 101
	case ErrInactiveUser:
		return 102 // The current user account is not active. User will be prompted to get in touch with Customer Support.
	case ErrInvalidAPIFunction:
		return 103 // The requested API endpoint does not exist.
	case ErrUsageLimitReached:
		return 104 // The maximum allowed amount of monthly API requests has been reached.
	case ErrFunctionAccessRestricted, ErrHTTPSAccessRestricted:
		// The current subscription plan does not support this API endpoint.
		// or
		// The user's current subscription plan does not support HTTPS Encryption.
		return 105
	case ErrInvalidFields:
		return 301 // One or more invalid fields were specified using the fields parameter.
	case ErrTooManyIPs:
		return 302 // Too many IPs have been specified for the Bulk Lookup Endpoint. (max. 50)
	case ErrBatchNotSupportedOnPlan:
		return 303 // The Bulk Lookup Endpoint is not supported on the current subscription plan
	default:
		return 0 // Not valid!
	}
}
