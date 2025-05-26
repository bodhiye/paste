package proto

// Error messages related to invalid user input or system issues
const (
	ErrInvalidArgs     = "invalid request parameter"
	ErrTooManyContent  = "content length exceeds the maximum allowed size of %d characters"
	ErrTooManyCount    = "content count exceeds the allowed limit of %d"
	ErrOverMaxSize     = "the size exceeds the allowed limit of %d MB"
	ErrPasteFailed     = "failed to paste content"
	ErrGetPasteFailed  = "failed to retrieve pasted content"
	ErrWrongPassword   = "incorrect password"
	ErrContentExpired  = "the requested content has expired"
	ErrInvalidFileType = "invalid file type, only images are allowed"
	ErrUploadFailed    = "failed to upload file"
)
