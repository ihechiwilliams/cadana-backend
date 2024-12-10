package constants

// ErrorCode ENUM(
//
// api_error
// account_is_not_active
// account_id_missing
// file_upload_error
// file_download_error
// invalid_account_status
// invalid_provider
// invalid_provider_data
// invalid_transaction_date
// invalid_transaction_id
// marshaling_error
// transaction_creation_failed
// validation_error,
//
// )
//
//go:generate go run github.com/abice/go-enum@v0.5.5
type ErrorCode string
