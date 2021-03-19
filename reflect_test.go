//
// Copyright (c) 2020-2021 Markku Rossi
//
// All rights reserved.
//

package tabulate

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
)

type Outer struct {
	Name    string
	Comment string `tabulate:"@detail"`
	Age     int
	NPS     float64
	Address *Address `tabulate:"omitempty"`
	Info    []*Info
	Meta    *Info `tabulate:"omitempty"`
	Mapping map[string]string
}

type Address struct {
	Lines []string
}

type Info struct {
	Email string
	Work  bool
}

func reflectTest(flags Flags, tags []string, v interface{}) error {
	tab := New(Unicode)
	tab.Header("Field").SetAlign(MR)
	tab.Header("Value").SetAlign(ML)

	err := Reflect(tab, flags, tags, v)
	if err != nil {
		return err
	}

	tab.Print(os.Stdout)
	return nil
}

func TestReflect(t *testing.T) {
	err := reflectTest(OmitEmpty, nil, &Outer{
		Name: "Alyssa P. Hacker",
		Age:  45,
		Address: &Address{
			Lines: []string{"42 Hacker way", "03139 Cambridge", "MA"},
		},
		Info: []*Info{
			{
				Email: "mtr@iki.fi",
			},
			{
				Email: "markku.rossi@gmail.com",
				Work:  true,
			},
		},
		Mapping: map[string]string{
			"First":  "1st",
			"Second": "2nd",
		},
	})
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}

	data := &Outer{
		Name:    "Alyssa P. Hacker",
		Comment: "Structure and Interpretation of Computer Programs",
		Age:     45,
		Info: []*Info{
			nil,
			{
				Email: "markku.rossi@gmail.com",
				Work:  true,
			},
		},
	}

	err = reflectTest(OmitEmpty, nil, data)
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}
	err = reflectTest(0, []string{"detail"}, data)
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}
}

type Outer2 struct {
	Name        string
	Inner       *Inner
	Certificate *Certificate
}

type Inner struct {
	A int
	B int
}

func (in Inner) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("A=%v, B=%v", in.A, in.B)), nil
}

type Certificate struct {
	X509 *x509.Certificate
}

func (c Certificate) MarshalText() (text []byte, err error) {
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.X509.Raw,
	}
	return pem.EncodeToMemory(block), nil
}

const certPEM = `-----BEGIN CERTIFICATE-----
MIIFezCCBGOgAwIBAgIQBXgQRGTO9sioGXkO2OhJWjANBgkqhkiG9w0BAQsFADBG
MQswCQYDVQQGEwJVUzEPMA0GA1UEChMGQW1hem9uMRUwEwYDVQQLEwxTZXJ2ZXIg
Q0EgMUIxDzANBgNVBAMTBkFtYXpvbjAeFw0yMDA1MTQwMDAwMDBaFw0yMTA2MTQx
MjAwMDBaMCQxIjAgBgNVBAMTGXByaXZ4ZGVtby5zc2guZW5naW5lZXJpbmcwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDj0lLHarwcvVqy33Z4Xb1lM1wM
dQpFrc6Sa7L4x4F1H6KGKg9YLZaQGFFu/m7fq01cGhkFXn2A2AEmlcUk0G3Ul/BH
irUV91QBkUu2TKqpcOSLeI9JuczOjvVGI5IwJShAWPC7fqHjquKlInX7zlr6KscB
r69DiYhoU+rd1sV3Jkimlj1XT/r5VQ6UK3FLXBbhKpDdiNQkRzWqbGLDgXr65VJd
UYiatvzxcTOUevwOCEgOaGip0N7jkcXxnG+y5BmXxkLGhumZ24gL6en4ZL2HeQlw
kAxnDZb86MnebJxlwdo1ck93qFmz012uqEntRE7x06E8HL0y7dOb8Mv6+RcpAgMB
AAGjggKFMIICgTAfBgNVHSMEGDAWgBRZpGYGUqB7lZI8o5QHJ5Z0W/k90DAdBgNV
HQ4EFgQUkymG/chDUXc6cmPFimvQX0ursnQwJAYDVR0RBB0wG4IZcHJpdnhkZW1v
LnNzaC5lbmdpbmVlcmluZzAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYB
BQUHAwEGCCsGAQUFBwMCMDsGA1UdHwQ0MDIwMKAuoCyGKmh0dHA6Ly9jcmwuc2Nh
MWIuYW1hem9udHJ1c3QuY29tL3NjYTFiLmNybDAgBgNVHSAEGTAXMAsGCWCGSAGG
/WwBAjAIBgZngQwBAgEwdQYIKwYBBQUHAQEEaTBnMC0GCCsGAQUFBzABhiFodHRw
Oi8vb2NzcC5zY2ExYi5hbWF6b250cnVzdC5jb20wNgYIKwYBBQUHMAKGKmh0dHA6
Ly9jcnQuc2NhMWIuYW1hem9udHJ1c3QuY29tL3NjYTFiLmNydDAMBgNVHRMBAf8E
AjAAMIIBBAYKKwYBBAHWeQIEAgSB9QSB8gDwAHcAfT7y+I//iFVoJMLAyp5SiXkr
xQ54CX8uapdomX4i8NcAAAFyELLhaAAABAMASDBGAiEAh9jeGd6oQxZqjlQfsAwm
j1kDPHlxbvfM2QgrEdFORfMCIQDrHoTNJA/mMqLqjUN1VicPKTCL+jXewGBCNZfV
f+melQB1AFzcQ5L+5qtFRLFemtRW5hA3+9X6R9yhc5SyXub2xw7KAAABchCy4TcA
AAQDAEYwRAIgLh9u1+jarU2ombS6PpGU3fb/UuSSmvcI5VpW4uRpPtUCIDnYvE5+
hnfnkgimPfCnPpfycgDmfRUiO/eG1Avh6xgkMA0GCSqGSIb3DQEBCwUAA4IBAQB+
ntJ6DNnhoiN/kYCDI7+emL9jZ3aZTmC+OPCATZtCAPgNOViv8+kFldYsc3FQZEs4
Nd9lxVj78Pg50Lv0gcRIUkxyIum7bMRG7gl9Cc8A3yPSGdfATIpXccHswi2Yf3JT
GlMuxx+KS/4ixewB77PucbMwqYKJKOEXIowhwL31fNYCxY2X5fPFZCIVI3f+tjJb
BsMnFsKUVVbjsrO7qRs/k7p1JwxKuxc1FjK3TPxrD8j6zs1tc3xDhkB/8TDXTdEm
n9hAr+Ljvi20fbmWITdiqD+1i2Ue6Gwj1+l1EtFUUjFiKYGvxy2sksfNdNDrY1bo
iEk07xyUePVXINXPEnOQ
-----END CERTIFICATE-----
`

func decodeCertificate() (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("PEM decode failed")
	}
	return x509.ParseCertificate(block.Bytes)
}

func TestReflectTextMarshaler(t *testing.T) {
	c, err := decodeCertificate()
	if err != nil {
		t.Fatalf("failed to decode certificate: %s", err)
	}

	err = reflectTest(0, nil, &Outer2{
		Name: "ACME Corp.",
		Inner: &Inner{
			A: 100,
			B: 42,
		},
		Certificate: &Certificate{
			X509: c,
		},
	})
	if err != nil {
		t.Fatalf("Reflect failed: %s", err)
	}
}

func TestReflectArray(t *testing.T) {
	fmt.Printf("TestReflectArray\n")
	tab, err := Array(New(ASCII), [][]interface{}{
		{"a", "b", "c"},
		{"1", "2", "3"},
	})
	if err != nil {
		t.Fatalf("Array failed: %s", err)
	}
	tab.Print(os.Stdout)
}
