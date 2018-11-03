package dns01

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xenolf/lego/le"
	"github.com/xenolf/lego/le/api"
	"github.com/xenolf/lego/platform/tester"
)

type DNSProviderMock struct{}

func (*DNSProviderMock) Present(domain, token, keyAuth string) error { return nil }
func (*DNSProviderMock) CleanUp(domain, token, keyAuth string) error { return nil }

func TestDNSChallenge(t *testing.T) {
	_, apiURL, tearDown := tester.SetupFakeAPI()
	defer tearDown()

	privKey, err := rsa.GenerateKey(rand.Reader, 512)
	require.NoError(t, err)

	core, err := api.New(http.DefaultClient, "lego-test", apiURL, "", privKey)
	require.NoError(t, err)

	validate := func(_ *api.Core, _, _ string, _ le.Challenge) error { return nil }
	preCheck := func(fqdn, value string) (bool, error) { return true, nil }

	chlg := NewChallenge(core, validate, &DNSProviderMock{}, AddPreCheck(preCheck))

	clientChallenge := le.Challenge{Type: "dns01", Status: "pending", URL: apiURL + "/chlg", Token: "http8"}

	err = chlg.Solve(clientChallenge, "example.com")
	require.NoError(t, err)
}
