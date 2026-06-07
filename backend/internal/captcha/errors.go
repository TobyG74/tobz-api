package captcha

import "errors"

var (
	ErrCaptchaRequired = errors.New("captcha wajib diisi")
	ErrCaptchaFailed   = errors.New("verifikasi captcha gagal")
)
