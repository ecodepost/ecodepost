// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: common/v1/task.proto

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

// Validate checks the field values on TaskTotal with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *TaskTotal) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TaskTotal with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TaskTotalMultiError, or nil
// if none found.
func (m *TaskTotal) ValidateAll() error {
	return m.validate(true)
}

func (m *TaskTotal) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Uid

	// no validation rules for TaskType

	// no validation rules for InvitationRegisterTotal

	if len(errors) > 0 {
		return TaskTotalMultiError(errors)
	}

	return nil
}

// TaskTotalMultiError is an error wrapping multiple validation errors returned
// by TaskTotal.ValidateAll() if the designated constraints aren't met.
type TaskTotalMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TaskTotalMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TaskTotalMultiError) AllErrors() []error { return m }

// TaskTotalValidationError is the validation error returned by
// TaskTotal.Validate if the designated constraints aren't met.
type TaskTotalValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TaskTotalValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TaskTotalValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TaskTotalValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TaskTotalValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TaskTotalValidationError) ErrorName() string { return "TaskTotalValidationError" }

// Error satisfies the builtin error interface
func (e TaskTotalValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTaskTotal.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TaskTotalValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TaskTotalValidationError{}
