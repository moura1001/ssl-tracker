package ssl

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/notify"
	"github.com/robfig/cron/v3"
)

const pollEvery = "@every 10s"

func StartCron() {
	c := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	entry, err := c.AddFunc(pollEvery, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := PollAllDomains(ctx); err != nil {
			logger.Log("msg", "stopping cron", "error", err)
			c.Stop()
		}
	})

	if err != nil {
		logger.Log("msg", "cron addFunc", "error", err)
		return
	}
	logger.Log("msg", "starting domain polling cron", "entryId", entry)
	c.Start()
}

func PollAllDomains(ctx context.Context) error {
	start := time.Now()

	trackings, err := db.Store.Domain.GetAllTrackingsWithAccount()
	if err != nil {
		return err
	}

	var (
		updatedTrackings = []*data.DomainTracking{}
		wg               = sync.WaitGroup{}
	)

	for _, tracking := range trackings {
		wg.Add(1)
		go func(tracking data.TrackingAndAccount) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer func() {
				cancel()
				wg.Done()
			}()

			trackingInfo, err := PollDomain(ctx, tracking.DomainName)
			if err != nil || trackingInfo.Error != "" {
				logger.Log("error", "failed to poll domain", "err", trackingInfo.Error, "domain", tracking.DomainName)
				return
			}
			domainTracking := tracking.DomainTracking
			domainTracking.DomainTrackingInfo = *trackingInfo
			updatedTrackings = append(updatedTrackings, &domainTracking)

			expires := trackingInfo.Expires
			notifyUpfront := time.Hour * 24 * time.Duration(tracking.NotifyUpfront)
			if time.Until(expires) <= notifyUpfront {
				notifiers := []notify.Notifier{}
				for _, notifier := range notifiers {
					start := time.Now()
					if err := notifier.Notify(context.Background(), tracking); err != nil {
						logger.Log("error", "notifier error", "err", err, "kind", notifier.Kind())
						return
					}
					logger.Log("msg", "notification", "kind", notifier.Kind(), "took", time.Since(start))
				}
			}
		}(tracking)
	}
	wg.Wait()

	if err := db.Store.Domain.UpdateAllTrackings(updatedTrackings); err != nil {
		return err
	}

	logger.Log("msg", "finished polling all trackings", "count", len(trackings), "took", time.Since(start))

	return nil
}

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
				LastPollAt: start,
				Error:      err.Error(),
				Latency:    int(time.Since(start).Milliseconds()),
				Issuer:     "n/a",
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
			LastPollAt:    start,
			Latency:       int(time.Since(start).Milliseconds()),
			Status:        getStatus(cert.NotAfter),
		}
	}()

	select {
	case <-ctx.Done():
		//fmt.Println(ctx.Err())
		//return nil, ctx.Err()
		return &data.DomainTrackingInfo{
			Error:  ctx.Err().Error(),
			Status: data.StatusUnresponsive,
			Issuer: "n/a",
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
	return base64.RawStdEncoding.EncodeToString(hasher.Sum(nil))
}

var loomingThreshold = time.Hour * 24 * 7 * 2 // 2 weeks

func getStatus(expires time.Time) string {
	if expires.Before(time.Now()) {
		return data.StatusExpired
	} else if time.Now().Add(loomingThreshold).After(expires) {
		return data.StatusExpires
	} else {
		return data.StatusHealthy
	}
}

func IsVerificationError(err error) bool {
	return strings.Contains(err.Error(), "tls: failed to verify")
}

func IsConnectionRefused(err error) bool {
	return strings.Contains(err.Error(), "dial tcp")
}
