package ssl

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
)

func PollDomain(ctx context.Context, domain string) (*data.DomainTrackingInfo, error) {
	var (
		start      = time.Now()
		resultChan = make(chan data.DomainTrackingInfo)
		config     = &tls.Config{}
	)
	go func() {
		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", domain), config)
		if err != nil {
			info := data.DomainTrackingInfo{
				LastPollAt: time.Now(),
				Error:      err.Error(),
				Latency:    int(time.Since(start).Milliseconds()),
			}
			if IsVerificationError(err) {
				info.Status = data.StatusInvalid
				resultChan <- info
			}
			if IsConnectionRefused(err) {
				info.Status = data.StatusOffline
				resultChan <- info
			}
			return
		}
		defer conn.Close()
		var (
			state     = conn.ConnectionState()
			cert      = state.PeerCertificates[0]
			keyUsages = make([]string, len(cert.ExtKeyUsage))
			i         = 0
		)
		for _, usage := range cert.ExtKeyUsage {
			keyUsages[i] = extKeyUsageToString(usage)
			i++
		}
		resultChan <- data.DomainTrackingInfo{
			PublicKeyAlgo: cert.PublicKeyAlgorithm.String(),
			SignatureAlgo: cert.SignatureAlgorithm.String(),
			KeyUsage:      keyUsageToString(cert.KeyUsage),
			ExtKeyUsages:  keyUsages,
			PublicKey:     publicKeyFromCert(cert),
			EncodedPEM:    encodedPEMFromCert(cert),
			Signature:     sha1Hex(cert.Signature),
			Expires:       cert.NotAfter,
			DNSNames:      strings.Join(cert.DNSNames, ", "),
			Issuer:        cert.Issuer.Organization[0],
			LastPollAt:    time.Now(),
			Latency:       int(time.Since(start).Milliseconds()),
			Status:        getStatus(cert.NotAfter),
		}
	}()

	select {
	case <-ctx.Done():
		return &data.DomainTrackingInfo{
			Error:      ctx.Err().Error(),
			LastPollAt: time.Now(),
			Status:     data.StatusUnresponsive,
		}, nil
	case result := <-resultChan:
		return &result, nil
	}
}

func extKeyUsageToString(usage x509.ExtKeyUsage) string {
	switch usage {
	case x509.ExtKeyUsageAny:
		return "Any"
	case x509.ExtKeyUsageServerAuth:
		return "Server Auth"
	case x509.ExtKeyUsageClientAuth:
		return "Client Auth"
	case x509.ExtKeyUsageCodeSigning:
		return "Code Signing"
	case x509.ExtKeyUsageEmailProtection:
		return "Email Protection"
	case x509.ExtKeyUsageIPSECEndSystem:
		return "IPSEC End System"
	case x509.ExtKeyUsageIPSECTunnel:
		return "IPSEC Tunnel"
	case x509.ExtKeyUsageIPSECUser:
		return "IPSEC User"
	case x509.ExtKeyUsageTimeStamping:
		return "Time Stamping"
	case x509.ExtKeyUsageOCSPSigning:
		return "OCSP Signing"
	case x509.ExtKeyUsageMicrosoftServerGatedCrypto:
		return "Microsoft Server Gated Crypto"
	case x509.ExtKeyUsageNetscapeServerGatedCrypto:
		return "Netscape Server Gated Crypto"
	case x509.ExtKeyUsageMicrosoftCommercialCodeSigning:
		return "Microsoft Commercial Code Signing"
	case x509.ExtKeyUsageMicrosoftKernelCodeSigning:
		return "Microsoft Kernel Code Signing"
	default:
		return ""
	}
}

func keyUsageToString(usage x509.KeyUsage) string {
	switch usage {
	case x509.KeyUsageDigitalSignature:
		return "Digital Signature"
	case x509.KeyUsageContentCommitment:
		return "Content Commitment"
	case x509.KeyUsageKeyEncipherment:
		return "Key Encipherment"
	case x509.KeyUsageDataEncipherment:
		return "Data Encipherment"
	case x509.KeyUsageKeyAgreement:
		return "Key Agreement"
	case x509.KeyUsageCertSign:
		return "Cert Sign"
	case x509.KeyUsageCRLSign:
		return "CRL Sign"
	case x509.KeyUsageEncipherOnly:
		return "Encipher Only"
	case x509.KeyUsageDecipherOnly:
		return "Decipher Only"
	default:
		return ""
	}
}

func publicKeyFromCert(cert *x509.Certificate) string {
	publicKeyDer, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err == nil {
		return sha1Hex(publicKeyDer)
	}
	return ""
}

func encodedPEMFromCert(cert *x509.Certificate) string {
	certBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	return string(pem.EncodeToMemory(&certBlock))
}

func sha1Hex(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func getStatus(certificateNotAfter time.Time) string {
	if certificateNotAfter.After(time.Now()) {
		return data.StatusHealthy
	} else {
		return data.StatusExpires
	}
}

func IsVerificationError(err error) bool {
	return errors.Is(err, syscall.EINVAL)
}

func IsConnectionRefused(err error) bool {
	return errors.Is(err, syscall.ECONNREFUSED)
}
