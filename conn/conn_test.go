package conn

import "testing"

var sampleBoard = "x"
var sampleThread = 21676943

func TestCon(t *testing.T) {
	/* thread */
	c := NewCon("/tmp/boof")
	th := &Posts{}
	path := MakePath(sampleBoard, sampleThread)
	url := Url(Static, path)
	if err := c.Gather(url, th); err != nil {
		t.Fatalf("Failed to gather %s: %s", path, err)
	}

	/* image */
	path = MakeImgPath(th.P[0].Ext, sampleBoard, MakeImgN(th.P[0].Tim))
	url = Url(Img, path)
	resp, err := c.Get(url.String())
	if err != nil {
		t.Fatalf("Failed to get %s: %s", url.String(), err)
	} else if resp.StatusCode != 200 { // http.StatusOK
		t.Fatalf("Status not ok on %s: %d", url.String(), resp.StatusCode)
	}
	resp.Body.Close()
}
