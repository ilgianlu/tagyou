package jwt

type VerificationResult struct {
	Valid  bool
	UserId int64
	Err    error
}
