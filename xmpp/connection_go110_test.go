// +build go1.10

package xmpp

import (
	"io"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

func (s *ConnectionXMPPSuite) Test_Dial_failsWhenStartingAHandshake(c *C) {
	tlsC, v := tlsConfigForRecordedHandshake()
	tlsC.Rand = fixedRand([]string{"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F"})

	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>",
	)}
	d := &dialer{
		JID:      "user@domain",
		password: "pass",

		verifier: v,
		config: data.Config{
			TLSConfig: tlsC,
		},
	}

	conn := &fullMockedConn{rw: rw}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	// Go 1.10 supports 2 additional signature algorithms, and they
	// can't be disabled. We need to rewrite this expected data
	// because (in comparison to connection_go16_test.go) to include
	// these signature algorithms in the ClientHello message.
	// See: https://github.com/golang/go/commit/96cd66b266a24d8ed1f66dfa10ddb86d88b50fca
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?><stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x92\x01\x00\x00\x8e\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
		"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x1c\xc0/\xc0+\xc00\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
		"\x00\x9c\x00\x9d\x00/\x005\xc0\x12\x00\n"+
		"\x01\x00\x00I\x00\x00\x00\v\x00\t\x00\x00\x06domain\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
		"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x12\x00\x10\x04\x01\x04\x03\x05\x01\x05\x03\x06\x01\x06\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00",
	)
}

//TODO: record the handshake again and update the tests
/*
func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceeds(c *C) {
	tlsC, v := tlsConfigForRecordedHandshake()
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      v,
		config: data.Config{
			TLSConfig: tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?><stream:stream to='www.olabini.se' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x96\x01\x00\x00\x92\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
		"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x1c\xc0/\xc0+\xc00\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
		"\x00\x9c\x00\x9d\x00/\x005\xc0\x12\x00\n"+
		"\x01\x00\x00M\x00\x00\x00\x13\x00\x11\x00\x00\x0ewww.olabini.se\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
		"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00\x16\x03\x03\x00F\x10\x00\x00BA\x04\b\xda\xf8\xb0\xab\xae5\t\xf3\\\xe1\xd31\x04\xcb\x01\xb9Qb̹\x18\xba\x1f\x81o8\xd3\x13\x0f\xb8\u007f\x92\xa3\b7\xf8o\x9e\xef\x19\u007fCy\xa5\n"+
		"b\x06\x82fy]\xb9\xf83\xea6\x1d\x03\xafT[\xe7\x92\x14\x03\x03\x00\x01\x01\x16\x03\x03\x00(\x00\x00\x00\x00\x00\x00\x00\x00j\xc9\xd4\xef'\xc2ڲ6e\x88JD$\x9d\x15O\x80\x15\u0099-.(T\x99\xb3\xbf\x869~\x11\x17\x03\x03\x00\xa5\x00\x00\x00\x00\x00\x00\x00\x01W'*.\vz\x85%%E\xee\x119\x82\xe4\xfe\xf5o\x111m|\x8d\xbb ģ\x06\xf2a_M\xd3\xebޤ:\xe3\x00\x92A\xd6\\\xa7\xce<\"F?tr\xbc1gȮD\\ՎG\xe7\x9e\xc8\u007f\x9e\x1eD\xdbJ&\xe5T2b\xa5J\xc9\xc99\"Ek\xc0\x0fa\xf4\x9b*\xb1\x9a\x10^^(\x8f\x1b\x8b\xdd\x04\xe8\xe6\xf5\xb6+V\x1d\xd7\xeb512\xc0*\x8e7\x15\xdd\t\x04\x1d-#\xf8\x9d^\x16k\xa4\x1d\xe0s\x90E\x8cNG\x93\xb0\xd5]\x1d\x01B\xe9\x18T\xc7@\x11\x1a\x17\xe3\xb9b\x19\x1a")
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceedsButFailsOnInvalidCertHash(c *C) {
	tlsC, _ := tlsConfigForRecordedHandshake()
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      &basicTLSVerifier{[]byte("aaaaa")},

		config: data.Config{
			TLSConfig: tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "tls: server certificate does not match expected hash (got: 82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7, want: 6161616161)")
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceedsButSucceedsOnValidCertHash(c *C) {
	tlsC, _ := tlsConfigForRecordedHandshake()
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      &basicTLSVerifier{bytesFromHex("82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7")},

		config: data.Config{
			TLSConfig: tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}
*/
