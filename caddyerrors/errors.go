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
	ErrInvalidStyleCode
	ErrInvalidStyleName
	ErrExpiredToken
	ErrInvalidRequestMethod
	ErrJwtParameterParsing
	ErrFormDataParameterParsing
	ErrImageStyleParsing
	ErrInvalidTls
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
	ErrSqlUpdate
	ErrTooManyHostDomainWithBucket
	ErrTooManyImageStyle
	ErrInvalidTlsPem
	ErrInvalidTlsKey
	ErrNoRouter
	ErrNoRow
	ErrTimeout
	ErrInternalServer
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
		HttpStatusCode: http.StatusConflict,
	},
	ErrInvalidStyleCode: {
		CaddyErrorCode: "InvalidStyleCode",
		Description:    "Incorrect Style code parameters.",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrInvalidStyleName: {
		CaddyErrorCode: "InvalidStyleName",
		Description:    "Incorrect Style name.",
		HttpStatusCode: http.StatusBadRequest,
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
	ErrImageStyleParsing: {
		CaddyErrorCode: "ErrImageStyleParsing",
		Description:    "Style code parsing error.",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrFormDataParameterParsing: {
		CaddyErrorCode: "ErrFormDataParameterParsing",
		Description:    "Parameter parsing carried by form-data failed.",
		HttpStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTls: {
		CaddyErrorCode: "InvalidTlsFailure",
		Description:    "The public key and key format of the certificate you uploaded is incorrect.",
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
	ErrSqlUpdate: {
		CaddyErrorCode: "SqlUpdateErr",
		Description:    "Database query has some exceptions.",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrTooManyHostDomainWithBucket: {
		CaddyErrorCode: "TooManyHostDomainWithBucket",
		Description:    "Bind up to 20 custom domain names per bucket.",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrTooManyImageStyle: {
		CaddyErrorCode: "TooManyImageStyle",
		Description:    "Bind up to 50 image style per bucket.",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrInvalidTlsPem: {
		CaddyErrorCode: "InvalidTlsPem",
		Description:    "The upload certificate is malformed or does not match the bound domain name.",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrInvalidTlsKey: {
		CaddyErrorCode: "InvalidTlsKey",
		Description:    "Certificate private key error.",
		HttpStatusCode: http.StatusForbidden,
	},
	ErrNoRouter: {
		CaddyErrorCode: "NoRouter",
		Description:    "The specified parameter does not exist.",
		HttpStatusCode: http.StatusNotFound,
	},
	ErrNoRow: {
		CaddyErrorCode: "NoRow",
		Description:    "No related data found in the database.",
		HttpStatusCode: http.StatusNoContent,
	},
	ErrTimeout: {
		CaddyErrorCode: "Timeout",
		Description:    "Request timed out.",
		HttpStatusCode: http.StatusRequestTimeout,
	},
	ErrInternalServer: {
		CaddyErrorCode: "InternalErr",
		Description:    "Internal server error.",
		HttpStatusCode: http.StatusInternalServerError,
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
