// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
)

// ValidateFields checks the field values on CreateMACSettingsProfileRequest
// with the rules defined in the proto definition for this message. If any
// rules are violated, an error is returned.
func (m *CreateMACSettingsProfileRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = CreateMACSettingsProfileRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "mac_settings_profile":

			if m.GetMacSettingsProfile() == nil {
				return CreateMACSettingsProfileRequestValidationError{
					field:  "mac_settings_profile",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetMacSettingsProfile()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return CreateMACSettingsProfileRequestValidationError{
						field:  "mac_settings_profile",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return CreateMACSettingsProfileRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// CreateMACSettingsProfileRequestValidationError is the validation error
// returned by CreateMACSettingsProfileRequest.ValidateFields if the
// designated constraints aren't met.
type CreateMACSettingsProfileRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateMACSettingsProfileRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateMACSettingsProfileRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateMACSettingsProfileRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateMACSettingsProfileRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateMACSettingsProfileRequestValidationError) ErrorName() string {
	return "CreateMACSettingsProfileRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateMACSettingsProfileRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateMACSettingsProfileRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateMACSettingsProfileRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateMACSettingsProfileRequestValidationError{}

// ValidateFields checks the field values on CreateMACSettingsProfileResponse
// with the rules defined in the proto definition for this message. If any
// rules are violated, an error is returned.
func (m *CreateMACSettingsProfileResponse) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = CreateMACSettingsProfileResponseFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "mac_settings_profile":

			if m.GetMacSettingsProfile() == nil {
				return CreateMACSettingsProfileResponseValidationError{
					field:  "mac_settings_profile",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetMacSettingsProfile()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return CreateMACSettingsProfileResponseValidationError{
						field:  "mac_settings_profile",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return CreateMACSettingsProfileResponseValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// CreateMACSettingsProfileResponseValidationError is the validation error
// returned by CreateMACSettingsProfileResponse.ValidateFields if the
// designated constraints aren't met.
type CreateMACSettingsProfileResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateMACSettingsProfileResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateMACSettingsProfileResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateMACSettingsProfileResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateMACSettingsProfileResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateMACSettingsProfileResponseValidationError) ErrorName() string {
	return "CreateMACSettingsProfileResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateMACSettingsProfileResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateMACSettingsProfileResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateMACSettingsProfileResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateMACSettingsProfileResponseValidationError{}

// ValidateFields checks the field values on GetMACSettingsProfileRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *GetMACSettingsProfileRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = GetMACSettingsProfileRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "mac_settings_profile_ids":

			if m.GetMacSettingsProfileIds() == nil {
				return GetMACSettingsProfileRequestValidationError{
					field:  "mac_settings_profile_ids",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetMacSettingsProfileIds()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return GetMACSettingsProfileRequestValidationError{
						field:  "mac_settings_profile_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "field_mask":

			if v, ok := interface{}(m.GetFieldMask()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return GetMACSettingsProfileRequestValidationError{
						field:  "field_mask",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return GetMACSettingsProfileRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// GetMACSettingsProfileRequestValidationError is the validation error returned
// by GetMACSettingsProfileRequest.ValidateFields if the designated
// constraints aren't met.
type GetMACSettingsProfileRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetMACSettingsProfileRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetMACSettingsProfileRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetMACSettingsProfileRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetMACSettingsProfileRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetMACSettingsProfileRequestValidationError) ErrorName() string {
	return "GetMACSettingsProfileRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetMACSettingsProfileRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetMACSettingsProfileRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetMACSettingsProfileRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetMACSettingsProfileRequestValidationError{}

// ValidateFields checks the field values on GetMACSettingsProfileResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *GetMACSettingsProfileResponse) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = GetMACSettingsProfileResponseFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "mac_settings_profile":

			if m.GetMacSettingsProfile() == nil {
				return GetMACSettingsProfileResponseValidationError{
					field:  "mac_settings_profile",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetMacSettingsProfile()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return GetMACSettingsProfileResponseValidationError{
						field:  "mac_settings_profile",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return GetMACSettingsProfileResponseValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// GetMACSettingsProfileResponseValidationError is the validation error
// returned by GetMACSettingsProfileResponse.ValidateFields if the designated
// constraints aren't met.
type GetMACSettingsProfileResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetMACSettingsProfileResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetMACSettingsProfileResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetMACSettingsProfileResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetMACSettingsProfileResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetMACSettingsProfileResponseValidationError) ErrorName() string {
	return "GetMACSettingsProfileResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetMACSettingsProfileResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetMACSettingsProfileResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetMACSettingsProfileResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetMACSettingsProfileResponseValidationError{}

// ValidateFields checks the field values on UpdateMACSettingsProfileRequest
// with the rules defined in the proto definition for this message. If any
// rules are violated, an error is returned.
func (m *UpdateMACSettingsProfileRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = UpdateMACSettingsProfileRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "mac_settings_profile_ids":

			if m.GetMacSettingsProfileIds() == nil {
				return UpdateMACSettingsProfileRequestValidationError{
					field:  "mac_settings_profile_ids",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetMacSettingsProfileIds()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return UpdateMACSettingsProfileRequestValidationError{
						field:  "mac_settings_profile_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "mac_settings_profile":

			if m.GetMacSettingsProfile() == nil {
				return UpdateMACSettingsProfileRequestValidationError{
					field:  "mac_settings_profile",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetMacSettingsProfile()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return UpdateMACSettingsProfileRequestValidationError{
						field:  "mac_settings_profile",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "field_mask":

			if v, ok := interface{}(m.GetFieldMask()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return UpdateMACSettingsProfileRequestValidationError{
						field:  "field_mask",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return UpdateMACSettingsProfileRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// UpdateMACSettingsProfileRequestValidationError is the validation error
// returned by UpdateMACSettingsProfileRequest.ValidateFields if the
// designated constraints aren't met.
type UpdateMACSettingsProfileRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateMACSettingsProfileRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateMACSettingsProfileRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateMACSettingsProfileRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateMACSettingsProfileRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateMACSettingsProfileRequestValidationError) ErrorName() string {
	return "UpdateMACSettingsProfileRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateMACSettingsProfileRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateMACSettingsProfileRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateMACSettingsProfileRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateMACSettingsProfileRequestValidationError{}

// ValidateFields checks the field values on UpdateMACSettingsProfileResponse
// with the rules defined in the proto definition for this message. If any
// rules are violated, an error is returned.
func (m *UpdateMACSettingsProfileResponse) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = UpdateMACSettingsProfileResponseFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "mac_settings_profile":

			if m.GetMacSettingsProfile() == nil {
				return UpdateMACSettingsProfileResponseValidationError{
					field:  "mac_settings_profile",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetMacSettingsProfile()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return UpdateMACSettingsProfileResponseValidationError{
						field:  "mac_settings_profile",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return UpdateMACSettingsProfileResponseValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// UpdateMACSettingsProfileResponseValidationError is the validation error
// returned by UpdateMACSettingsProfileResponse.ValidateFields if the
// designated constraints aren't met.
type UpdateMACSettingsProfileResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateMACSettingsProfileResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateMACSettingsProfileResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateMACSettingsProfileResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateMACSettingsProfileResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateMACSettingsProfileResponseValidationError) ErrorName() string {
	return "UpdateMACSettingsProfileResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateMACSettingsProfileResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateMACSettingsProfileResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateMACSettingsProfileResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateMACSettingsProfileResponseValidationError{}

// ValidateFields checks the field values on DeleteMACSettingsProfileRequest
// with the rules defined in the proto definition for this message. If any
// rules are violated, an error is returned.
func (m *DeleteMACSettingsProfileRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = DeleteMACSettingsProfileRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "mac_settings_profile_ids":

			if m.GetMacSettingsProfileIds() == nil {
				return DeleteMACSettingsProfileRequestValidationError{
					field:  "mac_settings_profile_ids",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetMacSettingsProfileIds()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return DeleteMACSettingsProfileRequestValidationError{
						field:  "mac_settings_profile_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return DeleteMACSettingsProfileRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// DeleteMACSettingsProfileRequestValidationError is the validation error
// returned by DeleteMACSettingsProfileRequest.ValidateFields if the
// designated constraints aren't met.
type DeleteMACSettingsProfileRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteMACSettingsProfileRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteMACSettingsProfileRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteMACSettingsProfileRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteMACSettingsProfileRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteMACSettingsProfileRequestValidationError) ErrorName() string {
	return "DeleteMACSettingsProfileRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteMACSettingsProfileRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteMACSettingsProfileRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteMACSettingsProfileRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteMACSettingsProfileRequestValidationError{}

// ValidateFields checks the field values on DeleteMACSettingsProfileResponse
// with the rules defined in the proto definition for this message. If any
// rules are violated, an error is returned.
func (m *DeleteMACSettingsProfileResponse) ValidateFields(paths ...string) error {
	if len(paths) > 0 {
		return fmt.Errorf("message DeleteMACSettingsProfileResponse has no fields, but paths %s were specified", paths)
	}
	return nil
}

// DeleteMACSettingsProfileResponseValidationError is the validation error
// returned by DeleteMACSettingsProfileResponse.ValidateFields if the
// designated constraints aren't met.
type DeleteMACSettingsProfileResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteMACSettingsProfileResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteMACSettingsProfileResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteMACSettingsProfileResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteMACSettingsProfileResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteMACSettingsProfileResponseValidationError) ErrorName() string {
	return "DeleteMACSettingsProfileResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteMACSettingsProfileResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteMACSettingsProfileResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteMACSettingsProfileResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteMACSettingsProfileResponseValidationError{}

// ValidateFields checks the field values on ListMACSettingsProfilesRequest
// with the rules defined in the proto definition for this message. If any
// rules are violated, an error is returned.
func (m *ListMACSettingsProfilesRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ListMACSettingsProfilesRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "application_ids":

			if m.GetApplicationIds() == nil {
				return ListMACSettingsProfilesRequestValidationError{
					field:  "application_ids",
					reason: "value is required",
				}
			}

			if v, ok := interface{}(m.GetApplicationIds()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ListMACSettingsProfilesRequestValidationError{
						field:  "application_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "field_mask":

			if v, ok := interface{}(m.GetFieldMask()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ListMACSettingsProfilesRequestValidationError{
						field:  "field_mask",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "order":

			if _, ok := _ListMACSettingsProfilesRequest_Order_InLookup[m.GetOrder()]; !ok {
				return ListMACSettingsProfilesRequestValidationError{
					field:  "order",
					reason: "value must be in list [ ids.profile_id -ids.profile_id]",
				}
			}

		case "limit":

			if m.GetLimit() > 1000 {
				return ListMACSettingsProfilesRequestValidationError{
					field:  "limit",
					reason: "value must be less than or equal to 1000",
				}
			}

		case "page":
			// no validation rules for Page
		default:
			return ListMACSettingsProfilesRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ListMACSettingsProfilesRequestValidationError is the validation error
// returned by ListMACSettingsProfilesRequest.ValidateFields if the designated
// constraints aren't met.
type ListMACSettingsProfilesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListMACSettingsProfilesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListMACSettingsProfilesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListMACSettingsProfilesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListMACSettingsProfilesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListMACSettingsProfilesRequestValidationError) ErrorName() string {
	return "ListMACSettingsProfilesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListMACSettingsProfilesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListMACSettingsProfilesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListMACSettingsProfilesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListMACSettingsProfilesRequestValidationError{}

var _ListMACSettingsProfilesRequest_Order_InLookup = map[string]struct{}{
	"":                {},
	"ids.profile_id":  {},
	"-ids.profile_id": {},
}

// ValidateFields checks the field values on ListMACSettingsProfilesResponse
// with the rules defined in the proto definition for this message. If any
// rules are violated, an error is returned.
func (m *ListMACSettingsProfilesResponse) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ListMACSettingsProfilesResponseFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "mac_settings_profiles":

			for idx, item := range m.GetMacSettingsProfiles() {
				_, _ = idx, item

				if v, ok := interface{}(item).(interface{ ValidateFields(...string) error }); ok {
					if err := v.ValidateFields(subs...); err != nil {
						return ListMACSettingsProfilesResponseValidationError{
							field:  fmt.Sprintf("mac_settings_profiles[%v]", idx),
							reason: "embedded message failed validation",
							cause:  err,
						}
					}
				}

			}

		default:
			return ListMACSettingsProfilesResponseValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ListMACSettingsProfilesResponseValidationError is the validation error
// returned by ListMACSettingsProfilesResponse.ValidateFields if the
// designated constraints aren't met.
type ListMACSettingsProfilesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListMACSettingsProfilesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListMACSettingsProfilesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListMACSettingsProfilesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListMACSettingsProfilesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListMACSettingsProfilesResponseValidationError) ErrorName() string {
	return "ListMACSettingsProfilesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListMACSettingsProfilesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListMACSettingsProfilesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListMACSettingsProfilesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListMACSettingsProfilesResponseValidationError{}
