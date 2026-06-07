package auth

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	pw := "Sup3rSecret!"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if hash == pw {
		t.Fatal("hash tidak boleh sama dengan plaintext")
	}

	ok, err := VerifyPassword(pw, hash)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if !ok {
		t.Fatal("password yang benar harus lolos verifikasi")
	}

	ok, _ = VerifyPassword("salah", hash)
	if ok {
		t.Fatal("password salah tidak boleh lolos")
	}
}

func TestHashIsSaltedUnique(t *testing.T) {
	h1, _ := HashPassword("samepassword1")
	h2, _ := HashPassword("samepassword1")
	if h1 == h2 {
		t.Fatal("dua hash dari password sama harus berbeda (salt acak)")
	}
}

func TestVerifyRejectsMalformed(t *testing.T) {
	if _, err := VerifyPassword("x", "not-a-valid-hash"); err == nil {
		t.Fatal("hash tidak valid harus mengembalikan error")
	}
}
