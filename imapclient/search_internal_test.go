package imapclient

import (
	"bufio"
	"strings"
	"testing"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/internal/imapwire"
)

func encodeSearchKey(t *testing.T, criteria *imap.SearchCriteria) string {
	t.Helper()
	var sb strings.Builder
	enc := imapwire.NewEncoder(bufio.NewWriter(&sb), imapwire.ConnSideClient)
	writeSearchKey(enc, criteria)
	if err := enc.CRLF(); err != nil {
		t.Fatalf("CRLF() = %v", err)
	}
	return strings.TrimSuffix(sb.String(), "\r\n")
}

func TestWriteSearchKeyGmailRaw(t *testing.T) {
	for _, tc := range []struct {
		name     string
		criteria imap.SearchCriteria
		want     string
	}{
		{
			name:     "raw only",
			criteria: imap.SearchCriteria{GmailRaw: "from:foo has:attachment"},
			want:     `X-GM-RAW "from:foo has:attachment"`,
		},
		{
			name: "combined with flag",
			criteria: imap.SearchCriteria{
				Flag:     []imap.Flag{imap.FlagSeen},
				GmailRaw: "subject:hello",
			},
			want: `SEEN X-GM-RAW "subject:hello"`,
		},
		{
			name:     "empty raw omitted",
			criteria: imap.SearchCriteria{},
			want:     "ALL",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := encodeSearchKey(t, &tc.criteria); got != tc.want {
				t.Errorf("writeSearchKey = %q, want %q", got, tc.want)
			}
		})
	}
}
