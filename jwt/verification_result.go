package jwt

type VerificationResult struct {
	Valid  bool
	UserId uint
	Err    error
}
