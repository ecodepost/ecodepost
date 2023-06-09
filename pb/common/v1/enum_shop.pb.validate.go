// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: common/v1/enum_shop.proto

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

// Validate checks the field values on GOOD with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *GOOD) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GOOD with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in GOODMultiError, or nil if none found.
func (m *GOOD) ValidateAll() error {
	return m.validate(true)
}

func (m *GOOD) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GOODMultiError(errors)
	}

	return nil
}

// GOODMultiError is an error wrapping multiple validation errors returned by
// GOOD.ValidateAll() if the designated constraints aren't met.
type GOODMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GOODMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GOODMultiError) AllErrors() []error { return m }

// GOODValidationError is the validation error returned by GOOD.Validate if the
// designated constraints aren't met.
type GOODValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GOODValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GOODValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GOODValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GOODValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GOODValidationError) ErrorName() string { return "GOODValidationError" }

// Error satisfies the builtin error interface
func (e GOODValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGOOD.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GOODValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GOODValidationError{}

// Validate checks the field values on ORDER with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ORDER) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ORDER with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ORDERMultiError, or nil if none found.
func (m *ORDER) ValidateAll() error {
	return m.validate(true)
}

func (m *ORDER) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return ORDERMultiError(errors)
	}

	return nil
}

// ORDERMultiError is an error wrapping multiple validation errors returned by
// ORDER.ValidateAll() if the designated constraints aren't met.
type ORDERMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ORDERMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ORDERMultiError) AllErrors() []error { return m }

// ORDERValidationError is the validation error returned by ORDER.Validate if
// the designated constraints aren't met.
type ORDERValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ORDERValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ORDERValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ORDERValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ORDERValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ORDERValidationError) ErrorName() string { return "ORDERValidationError" }

// Error satisfies the builtin error interface
func (e ORDERValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sORDER.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ORDERValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ORDERValidationError{}

// Validate checks the field values on CG with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *CG) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CG with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in CGMultiError, or nil if none found.
func (m *CG) ValidateAll() error {
	return m.validate(true)
}

func (m *CG) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return CGMultiError(errors)
	}

	return nil
}

// CGMultiError is an error wrapping multiple validation errors returned by
// CG.ValidateAll() if the designated constraints aren't met.
type CGMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CGMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CGMultiError) AllErrors() []error { return m }

// CGValidationError is the validation error returned by CG.Validate if the
// designated constraints aren't met.
type CGValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CGValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CGValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CGValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CGValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CGValidationError) ErrorName() string { return "CGValidationError" }

// Error satisfies the builtin error interface
func (e CGValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCG.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CGValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CGValidationError{}
