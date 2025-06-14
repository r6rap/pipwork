package checker

import (
    "crypto/tls"
    "strings"
    "time"
)

type SSLInfo struct {
    Name             string
    Valid            bool       // validation certificate
    ExpiredTime      time.Time  //
    DaysUntilExpiry  int
    Error            string
    Issuer           string
    Subject          string
    NotBefore        time.Time
}

func CheckSSL(address string) SSLInfo {
    address = strings.TrimPrefix(address, "https://")
	address = strings.TrimPrefix(address, "http://")

    if !strings.Contains(address, ":") {
        address += ":443"
    }

    // create connection TLS
    conn, err := tls.Dial("tcp", address, &tls.Config{
        InsecureSkipVerify: false, // validasi trusted cert
    })
    if err != nil {
        return SSLInfo{
            Valid: false,
            Error: err.Error(),
        }
    }
    defer conn.Close()

    // get certificate from connection
    certs := conn.ConnectionState().PeerCertificates
    if len(certs) == 0 {
        return SSLInfo{
            Valid: false,
            Error: "No certificates found",
        }
    }

    // get first certificate
    cert := certs[0]

    now := time.Now()

    // check the validity period of the certificate
    if now.Before(cert.NotBefore) {
        return SSLInfo{
            Valid: false,
            Error: "Certificate not valid yet",
        }
    }

    // count the remaining days
    expired := cert.NotAfter
    daysLeft := int(expired.Sub(now).Hours() / 24)

    return SSLInfo{
        Valid:           now.Before(expired),
        ExpiredTime:     expired,
        DaysUntilExpiry: daysLeft,
        Issuer:          cert.Issuer.CommonName,
        Subject:         cert.Subject.CommonName,
        NotBefore:       cert.NotBefore,
    }
}
