// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: common/v1/space.proto

package commonv1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
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
	_ = sort.Sort
)

// Validate checks the field values on SpaceOption with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SpaceOption) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SpaceOption with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SpaceOptionMultiError, or
// nil if none found.
func (m *SpaceOption) ValidateAll() error {
	return m.validate(true)
}

func (m *SpaceOption) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	// no validation rules for SpaceOptionId

	// no validation rules for Value

	// no validation rules for SpaceOptionType

	if len(errors) > 0 {
		return SpaceOptionMultiError(errors)
	}

	return nil
}

// SpaceOptionMultiError is an error wrapping multiple validation errors
// returned by SpaceOption.ValidateAll() if the designated constraints aren't met.
type SpaceOptionMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SpaceOptionMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SpaceOptionMultiError) AllErrors() []error { return m }

// SpaceOptionValidationError is the validation error returned by
// SpaceOption.Validate if the designated constraints aren't met.
type SpaceOptionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SpaceOptionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SpaceOptionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SpaceOptionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SpaceOptionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SpaceOptionValidationError) ErrorName() string { return "SpaceOptionValidationError" }

// Error satisfies the builtin error interface
func (e SpaceOptionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSpaceOption.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SpaceOptionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SpaceOptionValidationError{}

// Validate checks the field values on SpaceInfo with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SpaceInfo) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SpaceInfo with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SpaceInfoMultiError, or nil
// if none found.
func (m *SpaceInfo) ValidateAll() error {
	return m.validate(true)
}

func (m *SpaceInfo) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Guid

	// no validation rules for Name

	// no validation rules for IconType

	// no validation rules for Icon

	// no validation rules for SpaceType

	// no validation rules for SpaceLayout

	// no validation rules for Visibility

	// no validation rules for MemberCnt

	// no validation rules for SpaceGroupGuid

	for idx, item := range m.GetSpaceOptions() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, SpaceInfoValidationError{
						field:  fmt.Sprintf("SpaceOptions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, SpaceInfoValidationError{
						field:  fmt.Sprintf("SpaceOptions[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SpaceInfoValidationError{
					field:  fmt.Sprintf("SpaceOptions[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for ChargeType

	// no validation rules for OriginPrice

	// no validation rules for Price

	// no validation rules for Desc

	// no validation rules for HeadImage

	// no validation rules for Cover

	// no validation rules for IsAllowSet

	// no validation rules for Access

	// no validation rules for Sort

	// no validation rules for Link

	if len(errors) > 0 {
		return SpaceInfoMultiError(errors)
	}

	return nil
}

// SpaceInfoMultiError is an error wrapping multiple validation errors returned
// by SpaceInfo.ValidateAll() if the designated constraints aren't met.
type SpaceInfoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SpaceInfoMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SpaceInfoMultiError) AllErrors() []error { return m }

// SpaceInfoValidationError is the validation error returned by
// SpaceInfo.Validate if the designated constraints aren't met.
type SpaceInfoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SpaceInfoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SpaceInfoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SpaceInfoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SpaceInfoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SpaceInfoValidationError) ErrorName() string { return "SpaceInfoValidationError" }

// Error satisfies the builtin error interface
func (e SpaceInfoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSpaceInfo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SpaceInfoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SpaceInfoValidationError{}
