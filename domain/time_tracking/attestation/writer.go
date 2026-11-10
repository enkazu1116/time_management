package attestation

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"os"

	"time_management/domain/time_tracking/integrity"
)

const (
	HMACSecretEnv           = "EAS_SUBJECT_HMAC_SECRET"
	subjectCommitmentPrefix = "time-management/subject-day/v1"
	minimumHMACSecretBytes  = 32
)

var (
	ErrHMACSecretRequired = errors.New("attestation: EAS_SUBJECT_HMAC_SECRET is required")
	ErrHMACSecretTooShort = errors.New("attestation: EAS_SUBJECT_HMAC_SECRET must be at least 32 bytes")
	ErrChainRequired      = errors.New("attestation: chain writer is required")
	ErrAttestationFailed  = errors.New("attestation: chain write failed")
)

// Payload is the complete application data sent to the chain adapter. It
// deliberately contains no user ID, date, names, or work content.
type Payload struct {
	DayDigest         [sha256.Size]byte
	SubjectCommitment [sha256.Size]byte
	RefUID            [sha256.Size]byte
}

type OnChainRecord struct {
	DayDigest         [sha256.Size]byte
	SubjectCommitment [sha256.Size]byte
}

// Chain hides transaction signing and EAS contract access from the domain.
type Chain interface {
	Attest(context.Context, Payload) ([sha256.Size]byte, error)
	Read(context.Context, [sha256.Size]byte) (OnChainRecord, error)
}

type Writer struct {
	chain      Chain
	hmacSecret []byte
}

// NewWriterFromEnv constructs a writer using only the HMAC secret environment
// variable. The private key is owned and loaded separately by the chain adapter.
func NewWriterFromEnv(chain Chain) (*Writer, error) {
	if chain == nil {
		return nil, ErrChainRequired
	}
	secret := os.Getenv(HMACSecretEnv)
	if secret == "" {
		return nil, ErrHMACSecretRequired
	}
	if len(secret) < minimumHMACSecretBytes {
		return nil, ErrHMACSecretTooShort
	}
	return newWriter(chain, []byte(secret)), nil
}

func newWriter(chain Chain, secret []byte) *Writer {
	return &Writer{
		chain:      chain,
		hmacSecret: append([]byte(nil), secret...),
	}
}

// Attest computes the canonical day digest and an unlinkable subject/day HMAC,
// then asks the chain adapter to create the EAS attestation. Chain errors are
// intentionally replaced with a stable error so provider or signer internals
// cannot leak secrets through an API response.
func (w *Writer) Attest(ctx context.Context, records integrity.DayRecords, refUID [sha256.Size]byte) ([sha256.Size]byte, error) {
	digest, err := integrity.DigestDay(records)
	if err != nil {
		return [sha256.Size]byte{}, err
	}

	payload := Payload{
		DayDigest:         digest,
		SubjectCommitment: subjectCommitment(w.hmacSecret, records),
		RefUID:            refUID,
	}
	uid, err := w.chain.Attest(ctx, payload)
	if err != nil {
		return [sha256.Size]byte{}, ErrAttestationFailed
	}
	return uid, nil
}

func (w *Writer) Read(ctx context.Context, uid [sha256.Size]byte) (OnChainRecord, error) {
	record, err := w.chain.Read(ctx, uid)
	if err != nil {
		return OnChainRecord{}, ErrAttestationFailed
	}
	return record, nil
}

func (w *Writer) SubjectCommitment(records integrity.DayRecords) [sha256.Size]byte {
	return subjectCommitment(w.hmacSecret, records)
}

func subjectCommitment(secret []byte, records integrity.DayRecords) [sha256.Size]byte {
	hash := hmac.New(sha256.New, secret)
	writeCommitmentPart(hash, subjectCommitmentPrefix)
	writeCommitmentPart(hash, records.Location.String())
	writeCommitmentPart(hash, records.Day.In(records.Location).Format("2006-01-02"))
	writeCommitmentPart(hash, records.UserID)
	var commitment [sha256.Size]byte
	copy(commitment[:], hash.Sum(nil))
	return commitment
}

type byteWriter interface {
	Write([]byte) (int, error)
}

func writeCommitmentPart(writer byteWriter, value string) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(value)))
	_, _ = writer.Write(length[:])
	_, _ = writer.Write([]byte(value))
}
