package box

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// setup is a setup function used by some tests to ensure valid configuration.
var setup func() *SDK

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)

	setup = func() *SDK {
		sdk := new(SDK)
		sdk.NewConfig(&Config{
			BoxAppSettings: AppSettings{
				ClientID:     os.Getenv("CLIENT_ID"),
				ClientSecret: os.Getenv("CLIENT_SECRET"),
				AppAuth: AppAuth{
					PublicKeyID: os.Getenv("PUBLIC_KEY_ID"),
					PrivateKey:  os.Getenv("PRIVATE_KEY"),
					Passphrase:  os.Getenv("PASSPHRASE"),
				},
			},
			EnterpriseID: os.Getenv("ENTERPRISE_ID"),
		})
		return sdk
	}

	os.Exit(m.Run())
}
