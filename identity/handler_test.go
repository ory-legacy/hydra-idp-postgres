package identity

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/parnurzeal/gorequest"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ory-am/ladon/policy"
	"github.com/ory-am/hydra/herodot"
	"github.com/golang/mock/gomock"
)

var (
	hd *Handler
)

type payload struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Data     string `json:"data,omitempty"`
}

var policies = map[string][]policy.Policy{
	"allow-create": {&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:identities"}, []string{"create"}, nil}},
	"allow-create-get": {
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:identities"}, []string{"create"}, nil},
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"<.*>"}, []string{"get"}, nil},
	},
	"allow-all": {
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:identities"}, []string{"create"}, nil},
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:identities:<.*>"}, []string{"<.*>"}, nil},
	},
	"empty": {},
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hd.he = herodot.Herodot{}
	ts := httptest.NewServer(hd)

	for k, c := range []struct{
	} {
	} {

	}
}
