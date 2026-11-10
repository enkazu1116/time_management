package attestation

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"testing"
	"time"

	"time_management/domain/time_tracking/integrity"
)

type fakeChain struct {
	payload Payload
	uid     [32]byte
	err     error
}

func (f *fakeChain) Attest(_ context.Context, payload Payload) ([32]byte, error) {
	f.payload = payload
	return f.uid, f.err
}

func TestWriterAttestSendsOnlyDigestCommitmentAndReference(t *testing.T) {
	chain := &fakeChain{uid: [32]byte{1}}
	t.Setenv(HMACSecretEnv, "test-hmac-secret-with-at-least-32-bytes")
	writer, err := NewWriterFromEnv(chain)
	if err != nil {
		t.Fatal(err)
	}
	records := emptyTokyoDay(t, "private-user-id")
	refUID := [32]byte{2}

	uid, err := writer.Attest(context.Background(), records, refUID)
	if err != nil {
		t.Fatal(err)
	}
	if uid != chain.uid {
		t.Fatalf("Attest() UID = %x, want %x", uid, chain.uid)
	}
	if chain.payload.RefUID != refUID {
		t.Fatalf("Attest() ref UID = %x, want %x", chain.payload.RefUID, refUID)
	}
	wantDigest, err := integrity.DigestDay(records)
	if err != nil {
		t.Fatal(err)
	}
	if chain.payload.DayDigest != wantDigest {
		t.Fatalf("Attest() digest = %x, want %x", chain.payload.DayDigest, wantDigest)
	}

	otherUser := emptyTokyoDay(t, "another-private-user-id")
	_, err = writer.Attest(context.Background(), otherUser, [32]byte{})
	if err != nil {
		t.Fatal(err)
	}
	if chain.payload.SubjectCommitment == ([32]byte{}) {
		t.Fatal("Attest() sent an empty subject commitment")
	}
}

func TestSubjectCommitmentGoldenVector(t *testing.T) {
	records := emptyTokyoDay(t, "private-user-id")
	got := subjectCommitment([]byte("0123456789abcdef0123456789abcdef"), records)
	const want = "ab0190b493544cd8f5f2a2b0ac191866bf12383ab79520950002aea876f6e18a"
	if formatted := hex.EncodeToString(got[:]); formatted != want {
		t.Fatalf("subjectCommitment() = %s, want %s", formatted, want)
	}
}

func TestWriterCommitmentChangesWithSecretDateUserAndLocation(t *testing.T) {
	base := emptyTokyoDay(t, "user-a")
	baseCommitment := subjectCommitment([]byte("secret-a"), base)

	cases := []struct {
		name   string
		secret []byte
		change func(*integrity.DayRecords)
	}{
		{"secret", []byte("secret-b"), func(*integrity.DayRecords) {}},
		{"date", []byte("secret-a"), func(records *integrity.DayRecords) {
			records.Day = records.Day.AddDate(0, 0, 1)
		}},
		{"user", []byte("secret-a"), func(records *integrity.DayRecords) {
			records.UserID = "user-b"
		}},
		{"location", []byte("secret-a"), func(records *integrity.DayRecords) {
			records.Location = time.UTC
		}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			records := base
			test.change(&records)
			if got := subjectCommitment(test.secret, records); got == baseCommitment {
				t.Fatalf("commitment did not change with %s", test.name)
			}
		})
	}
}

func TestWriterDoesNotLeakSecretsInLogsOrErrors(t *testing.T) {
	const (
		hmacSecret = "hmac-secret-that-must-not-leak"
		privateKey = "private-key-that-must-not-leak"
	)
	var logs bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(originalOutput) })

	chain := &fakeChain{err: errors.New(hmacSecret + " " + privateKey)}
	writer := newWriter(chain, []byte(hmacSecret))
	_, err := writer.Attest(context.Background(), emptyTokyoDay(t, "user-a"), [32]byte{})
	if !errors.Is(err, ErrAttestationFailed) {
		t.Fatalf("Attest() error = %v, want %v", err, ErrAttestationFailed)
	}
	for _, secret := range []string{hmacSecret, privateKey} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("returned error leaked %q", secret)
		}
		if strings.Contains(logs.String(), secret) {
			t.Fatalf("logs leaked %q", secret)
		}
	}
}

func TestNewWriterFromEnvRequiresSecret(t *testing.T) {
	t.Setenv(HMACSecretEnv, "")
	if _, err := NewWriterFromEnv(&fakeChain{}); !errors.Is(err, ErrHMACSecretRequired) {
		t.Fatalf("NewWriterFromEnv() error = %v, want %v", err, ErrHMACSecretRequired)
	}
}

func TestNewWriterFromEnvRejectsShortSecret(t *testing.T) {
	t.Setenv(HMACSecretEnv, "too-short")
	if _, err := NewWriterFromEnv(&fakeChain{}); !errors.Is(err, ErrHMACSecretTooShort) {
		t.Fatalf("NewWriterFromEnv() error = %v, want %v", err, ErrHMACSecretTooShort)
	}
}

func emptyTokyoDay(t *testing.T, userID string) integrity.DayRecords {
	t.Helper()
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	return integrity.DayRecords{
		UserID:   userID,
		Day:      time.Date(2026, time.November, 10, 0, 0, 0, 0, tokyo),
		Location: tokyo,
	}
}
