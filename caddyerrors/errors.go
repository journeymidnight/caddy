package caddyerrors

import "net/http"

type HandleError interface {
	error
	CaddyErrorCode() string
	Description() string
	HttpStatusCode() int
}

type HandleErrorStruct struct {
	CaddyErrorCode string
	Description    string
	HttpStatusCode int
}

type HandleErrorCode int

const (
	ErrAccessDenied HandleErrorCode = iota
	ErrInvalidDnsResolution
	ErrInvalidJwtParams
	ErrExpiredToken
	ErrInvalidRequestMethod
	ErrJwtParameterParsing
	ErrParameterParsing
	ErrInvalidBucketDomain
	ErrInvalidBindBucketDomain
	ErrInvalidBucketPermission
	ErrNoSuchKey
	ErrGetMarshal
	ErrInvalidHostDomain
	ErrInvalidSql
	ErrSqlTransaction
	ErrSqlInsert
	ErrSqlDelete
	ErrTooManyHostDomainWithBucket
)

var ErrorCodeResponse = map[HandleErrorCode]HandleErrorStruct{
	ErrAccessDenied: {
		CaddyErrorCode: "AccessDenied",
		Description:    "Access Denied.",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrInvalidDnsResolution: {
		CaddyErrorCode: "InvalidDnsResolution",
		Description:    "Query DNS domain name resolution failed.",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrInvalidJwtParams: {
		CaddyErrorCode: "InvalidJwtParams",
		Description:    "JWT parameter conversion error.",
		HttpStatusCode: http.StatusInternalServerError,
	},
	ErrExpiredToken: {
		CaddyErrorCode: "ExpiredToken",
		Description:    "JWT token has expired.",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrInvalidRequestMethod: {
		CaddyErrorCode: "InvalidRequestMethod",
		Description:    "The request was made using the wrong request method.",
		HttpStatusCode: http.StatusMethodNotAllowed,
	},
	ErrJwtParameterParsing: {
		CaddyErrorCode: "ErrJwtParameterParsing",
		Description:    "Parameter parsing carried by JWT failed.",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrParameterParsing: {
		CaddyErrorCode: "ErrJwtParameterParsing",
		Description:    "Parameter parsing carried by host failed.",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrInvalidBucketDomain: {
		CaddyErrorCode: "InvalidBucketDomain",
		Description:    "Bucket and bucket domain names do not match.",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrInvalidBindBucketDomain: {
		CaddyErrorCode: "InvalidBindBucketDomain",
		Description:    "The bound domain name does not match the request server domain name.",
		HttpStatusCode: http.StatusPreconditionFailed,
	},
	ErrInvalidBucketPermission: {
		CaddyErrorCode: "InvalidBucketPermission",
		Description:    "No bucket operation permission.",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrNoSuchKey: {
		CaddyErrorCode: "NoSuchKey",
		Description:    "The specified key does not exist.",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrGetMarshal: {
		CaddyErrorCode: "GetMarshalFailure",
		Description:    "Return parameter parsing failed.",
		HttpStatusCode: http.StatusInternalServerError,
	},
	ErrInvalidHostDomain: {
		CaddyErrorCode: "InvalidHostDomain",
		Description:    "The added custom domain name is already bound.",
		HttpStatusCode: http.StatusAlreadyReported,
	},
	ErrInvalidSql: {
		CaddyErrorCode: "InvalidSql",
		Description:    "Database query has some exceptions.",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrSqlTransaction: {
		CaddyErrorCode: "SqlTransactionErr",
		Description:    "Database query has some exceptions.",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrSqlInsert: {
		CaddyErrorCode: "SqlInsertErr",
		Description:    "Database query has some exceptions.",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrSqlDelete: {
		CaddyErrorCode: "SqlDeleteErr",
		Description:    "Database query has some exceptions.",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrTooManyHostDomainWithBucket: {
		CaddyErrorCode: "TooManyHostDomainWithBucket",
		Description:    "Bind up to 20 custom domain names per bucket.",
		HttpStatusCode: http.StatusForbidden,
	},
}

func (e HandleErrorCode) CaddyErrorCode() string {
	CaddyError, ok := ErrorCodeResponse[e]
	if !ok {
		return "InternalError"
	}
	return CaddyError.CaddyErrorCode
}

func (e HandleErrorCode) Description() string {
	caddyError, ok := ErrorCodeResponse[e]
	if !ok {
		return "We encountered an internal error, please try again."
	}
	return caddyError.Description
}

func (e HandleErrorCode) Error() string {
	return e.Description()
}

func (e HandleErrorCode) HttpStatusCode() int {
	caddyError, ok := ErrorCodeResponse[e]
	if !ok {
		return http.StatusInternalServerError
	}
	return caddyError.HttpStatusCode
}
