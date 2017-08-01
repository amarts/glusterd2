package e2e

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddRemovePeer(t *testing.T) {
	r := require.New(t)

	g1, err := spawnGlusterd("./config/1.yaml", true)
	r.Nil(err)
	defer g1.Stop()
	defer g1.EraseWorkdir()
	r.True(g1.IsRunning())

	g2, err := spawnGlusterd("./config/2.yaml", true)
	r.Nil(err)
	defer g2.Stop()
	defer g2.EraseWorkdir()
	r.True(g2.IsRunning())

	// add peer: ask g1 to add g2 as peer
	reqBody := strings.NewReader(fmt.Sprintf(`{"addresses": ["%s"]}`, g2.PeerAddress))
	resp, err := http.Post("http://"+g1.ClientAddress+"/v1/peers", "application/json", reqBody)
	r.Nil(err)
	defer resp.Body.Close()
	r.Equal(resp.StatusCode, 201)

	// remove peer: ask g1 to remove g2 as peer
	delURL := fmt.Sprintf("http://%s/v1/peers/%s", g1.ClientAddress, g2.PeerID())
	req, err := http.NewRequest("DELETE", delURL, nil)
	r.Nil(err)
	resp, err = http.DefaultClient.Do(req)
	r.Nil(err)
	defer resp.Body.Close()
	r.Equal(resp.StatusCode, 204)
}