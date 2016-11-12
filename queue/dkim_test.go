package queue

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"strings"

	"github.com/Freeaqingme/dkim"
)

const (
	sampleMessage = `To: example@example.org
Subject: Example E-Mail
Date: Thu, 10 Nov 2016 19:42:46 +0330
Content-Type: multipart/mixed; boundary=03502b653d124d893d171a853d60488f7385f7d9f2ca7cefc8a4d6cf83d8
MIME-Version: 1.0
Message-Id: <d88de6cb-85e3-4430-9520-c0c745f3bd00@hectane>
From: Hectane Postman <hectane@example.org>

--03502b653d124d893d171a853d60488f7385f7d9f2ca7cefc8a4d6cf83d8
Content-Type: multipart/alternative; boundary=4aa16bab6b14378df91dfd097c676e67d956631721db4cecde5f6e6ea7cf

--4aa16bab6b14378df91dfd097c676e67d956631721db4cecde5f6e6ea7cf
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

Some stuff

=C2=A9 2016 Hectane.

--4aa16bab6b14378df91dfd097c676e67d956631721db4cecde5f6e6ea7cf
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable


<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org=
/TR/xhtml1/DTD/xhtml1-strict.dtd">=20
<html xmlns=3D"http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3Dutf-8" />
<meta name=3D"viewport" content=3D"width=3Ddevice-width, initial-scale=3D1.=
0"/>
<title>Example E-Mail</title>
</head>
<body style=3D"width:100%; margin:0; padding:0; -webkit-text-size-adjust:10=
0%; -ms-text-size-adjust:100%;">
<h1>Some stuff</h1>
</body>
</html>

--4aa16bab6b14378df91dfd097c676e67d956631721db4cecde5f6e6ea7cf--

--03502b653d124d893d171a853d60488f7385f7d9f2ca7cefc8a4d6cf83d8--
`
	sampleFrom = "Hectane Postman <hectane@example.org>"

	privKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXwIBAAKBgQDwIRP/UC3SBsEmGqZ9ZJW3/DkMoGeLnQg1fWn7/zYtIxN2SnFC
jxOCKG9v3b4jYfcTNh5ijSsq631uBItLa7od+v/RtdC2UzJ1lWT947qR+Rcac2gb
to/NMqJ0fzfVjH4OuKhitdY9tf6mcwGjaNBcWToIMmPSPDdQPNUYckcQ2QIDAQAB
AoGBALmn+XwWk7akvkUlqb+dOxyLB9i5VBVfje89Teolwc9YJT36BGN/l4e0l6QX
/1//6DWUTB3KI6wFcm7TWJcxbS0tcKZX7FsJvUz1SbQnkS54DJck1EZO/BLa5ckJ
gAYIaqlA9C0ZwM6i58lLlPadX/rtHb7pWzeNcZHjKrjM461ZAkEA+itss2nRlmyO
n1/5yDyCluST4dQfO8kAB3toSEVc7DeFeDhnC1mZdjASZNvdHS4gbLIA1hUGEF9m
3hKsGUMMPwJBAPW5v/U+AWTADFCS22t72NUurgzeAbzb1HWMqO4y4+9Hpjk5wvL/
eVYizyuce3/fGke7aRYw/ADKygMJdW8H/OcCQQDz5OQb4j2QDpPZc0Nc4QlbvMsj
7p7otWRO5xRa6SzXqqV3+F0VpqvDmshEBkoCydaYwc2o6WQ5EBmExeV8124XAkEA
qZzGsIxVP+sEVRWZmW6KNFSdVUpk3qzK0Tz/WjQMe5z0UunY9Ax9/4PVhp/j61bf
eAYXunajbBSOLlx4D+TunwJBANkPI5S9iylsbLs6NkaMHV6k5ioHBBmgCak95JGX
GMot/L2x0IYyMLAz6oLWh2hm7zwtb0CgOrPo1ke44hFYnfc=
-----END RSA PRIVATE KEY-----`

	expectedDKIMHeaderPrefix = `v=1; a=rsa-sha256; c=relaxed/simple; d=example.org; q=dns/txt; s=test;`
)

// borrowed from
// https://github.com/kalloc/dkim/blob/acfed5d65dd1e8cebcac4b4d429efa30e4cdedae/dkim.go#L235
func findDkimHeader(r *bufio.Reader) (string, error) {
	var value, line string
	var kv []string
	var err error
	var dkimState = 0

	for {
		if line, err = r.ReadString('\n'); err != nil {
			break
		}
		if len(line) <= 2 {
			break
		}

		if line[0] == ' ' || line[0] == '\t' {
			if dkimState == 1 {
				value += line
			}
		} else if dkimState == 0 {
			kv = strings.SplitN(line, ":", 2)
			// skip invalid header
			if len(kv) == 2 {
				if strings.Contains(strings.ToLower(kv[0]), "dkim-signature") {
					value = kv[1]
					dkimState = 1
				}
			}
		} else {
			dkimState = 2
		}
	}

	if dkimState == 0 {
		return "", errors.New("not found")
	}
	return value, nil
}

func TestDKIMSigning(t *testing.T) {
	dkimInstances = make(map[string]*dkim.DKIM)
	config := Config{
		DKIMConfigs: make(map[string]DKIMConfig),
	}
	config.DKIMConfigs["example.org"] = DKIMConfig{
		PrivateKey:       privKey,
		Selector:         "test",
		Canonicalization: "relaxed/simple",
	}

	r := ioutil.NopCloser(bytes.NewBufferString(sampleMessage))
	signedEmail, err := dkimSigned(sampleFrom, r, &config)
	if err != nil {
		t.Fatal(err)
	}
	signedEmailReader := bufio.NewReader(signedEmail)
	header, err := findDkimHeader(signedEmailReader)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(header, expectedDKIMHeaderPrefix) {
		t.Logf("DKIM header: %s", header)
		t.Fatalf("The value for DKIM header was not expected")
	}
}

func TestDKIMNotSigning(t *testing.T) {
	dkimInstances = make(map[string]*dkim.DKIM)
	config := Config{}
	r := ioutil.NopCloser(bytes.NewBufferString(sampleMessage))
	signedEmail, err := dkimSigned(sampleFrom, r, &config)
	if err != nil {
		dkim, err2 := dkimFor(sampleFrom, &config)
		t.Logf("dkim: %v / err2: %s", dkim, err2)
		t.Fatal(err)
	}
	signedEmailContent, err := ioutil.ReadAll(signedEmail)
	if err != nil {
		t.Fatal(err)
	}
	if string(signedEmailContent) != sampleMessage {
		t.Fatal("Expecting the message to be untouched")
	}
}
