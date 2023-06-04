// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: common/v1/enum_drive.proto

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

// Validate checks the field values on DRIVE with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *DRIVE) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DRIVE with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in DRIVEMultiError, or nil if none found.
func (m *DRIVE) ValidateAll() error {
	return m.validate(true)
}

func (m *DRIVE) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return DRIVEMultiError(errors)
	}

	return nil
}

// DRIVEMultiError is an error wrapping multiple validation errors returned by
// DRIVE.ValidateAll() if the designated constraints aren't met.
type DRIVEMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DRIVEMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DRIVEMultiError) AllErrors() []error { return m }

// DRIVEValidationError is the validation error returned by DRIVE.Validate if
// the designated constraints aren't met.
type DRIVEValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DRIVEValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DRIVEValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DRIVEValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DRIVEValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DRIVEValidationError) ErrorName() string { return "DRIVEValidationError" }

// Error satisfies the builtin error interface
func (e DRIVEValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDRIVE.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DRIVEValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DRIVEValidationError{}
