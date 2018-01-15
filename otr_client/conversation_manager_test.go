package otr_client

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type ConversationManagerSuite struct{}

var _ = Suite(&ConversationManagerSuite{})

type testSender struct {
	peer     data.JIDWithoutResource
	resource data.JIDResource
	msg      string
	err      error
}

type testConvBuilder struct {
	fake *otr3.Conversation
}

func (cb *testConvBuilder) NewConversation(peer data.JID) *otr3.Conversation {
	return cb.fake
}

func (ts *testSender) Send(peer data.JIDWithoutResource, resource data.JIDResource, msg string) error {
	ts.peer = peer
	ts.msg = msg
	ts.resource = resource
	return ts.err
}

func (s *ConversationManagerSuite) Test_TerminateAll_willTerminate(c *C) {
	cb := &testConvBuilder{&otr3.Conversation{}}
	ts := &testSender{err: nil}
	mgr := NewConversationManager(cb, ts)
	conv, created := mgr.EnsureConversationWith(data.JIDNR("someone@whitehouse.gov"), data.JIDResource(""))

	c.Assert(created, Equals, true)
	c.Assert(conv, Not(IsNil))

	mgr.TerminateAll()

	c.Assert(ts.msg, Equals, "")
}
