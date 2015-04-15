package jujusvg

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gc "gopkg.in/check.v1"
	"gopkg.in/juju/charm.v5-unstable"

	"github.com/juju/jujusvg/assets"
)

func Test(t *testing.T) { gc.TestingT(t) }

type newSuite struct{}

var _ = gc.Suite(&newSuite{})

var bundle = `
services:
  mongodb:
    charm: "cs:precise/mongodb-21"
    num_units: 1
    annotations:
      "gui-x": "940.5"
      "gui-y": "388.7698359714502"
    constraints: "mem=2G cpu-cores=1"
  elasticsearch:
    charm: "cs:~charming-devs/precise/elasticsearch-2"
    num_units: 1
    annotations:
      "gui-x": "490.5"
      "gui-y": "369.7698359714502"
    constraints: "mem=2G cpu-cores=1"
  charmworld:
    charm: "cs:~juju-jitsu/precise/charmworld-58"
    num_units: 1
    expose: true
    annotations:
      "gui-x": "813.5"
      "gui-y": "112.23016402854975"
    options:
      charm_import_limit: -1
      source: "lp:~bac/charmworld/ingest-local-charms"
      revno: 511
relations:
  - - "charmworld:essearch"
    - "elasticsearch:essearch"
  - - "charmworld:database"
    - "mongodb:database"
series: precise
`

func iconURL(ref *charm.Reference) string {
	return "http://0.1.2.3/" + ref.Path() + ".svg"
}

type emptyFetcher struct{}

func (f *emptyFetcher) FetchIcons(*charm.BundleData) (map[string][]byte, error) {
	return nil, nil
}

type errFetcher string

func (f *errFetcher) FetchIcons(*charm.BundleData) (map[string][]byte, error) {
	return nil, fmt.Errorf("%s", *f)
}

func (s *newSuite) TestNewFromBundle(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	cvs, err := NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.IsNil)

	var buf bytes.Buffer
	cvs.Marshal(&buf)
	c.Logf("%s", buf.String())
	assertXMLEqual(c, buf.Bytes(), []byte(`
<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="639" height="465"
     style="font-family:Ubuntu, sans-serif;" viewBox="0 0 639 465"
     xmlns="http://www.w3.org/2000/svg" 
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
<g id="serviceBlock" transform="scale(0.8)" >`+assets.ServiceModule+`
</g>
<g id="healthCircle">
<circle cx="10" cy="10" r="10" style="stroke:#38B44A;fill:none;stroke-width:2px"/>
<circle cx="10" cy="10" r="5" style="fill:#38B44A"/>
</g>
<g id="icon-1" >
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/~juju-jitsu/precise/charmworld-58.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg></g>
<g id="icon-2" >
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/~charming-devs/precise/elasticsearch-2.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg></g>
<g id="icon-3" >
<svg:svg xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
&#x9;&#x9;&#x9;&#x9;&#x9;<svg:image width="96" height="96" xlink:href="http://0.1.2.3/precise/mongodb-21.svg"></svg:image>
&#x9;&#x9;&#x9;&#x9;</svg:svg></g>
</defs>
<g id="relations">
<line x1="417" y1="189" x2="189" y2="351" stroke="#38B44A" stroke-width="2px" stroke-dasharray="129.85, 20" />
<use x="293" y="260" xlink:href="#healthCircle" />
<line x1="417" y1="189" x2="544" y2="276" stroke="#38B44A" stroke-width="2px" stroke-dasharray="66.97, 20" />
<use x="470" y="222" xlink:href="#healthCircle" />
</g>
<g id="services">
<use x="323" y="0" xlink:href="#serviceBlock" id="charmworld" />
<use x="369" y="46" xlink:href="#icon-1" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="417" y="31" >charmworld</text>
</g>
<use x="0" y="257" xlink:href="#serviceBlock" id="elasticsearch" />
<use x="46" y="303" xlink:href="#icon-2" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="94" y="288" >elasticsearch</text>
</g>
<use x="450" y="276" xlink:href="#serviceBlock" id="mongodb" />
<use x="496" y="322" xlink:href="#icon-3" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="544" y="307" >mongodb</text>
</g>
</g>
</svg>
`))
}

func (s *newSuite) TestWithFetcher(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	cvs, err := NewFromBundle(b, iconURL, new(emptyFetcher))
	c.Assert(err, gc.IsNil)

	var buf bytes.Buffer
	cvs.Marshal(&buf)
	c.Logf("%s", buf.String())
	assertXMLEqual(c, buf.Bytes(), []byte(`
<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="639" height="465"
     style="font-family:Ubuntu, sans-serif;" viewBox="0 0 639 465"
     xmlns="http://www.w3.org/2000/svg" 
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
<g id="serviceBlock" transform="scale(0.8)" >`+assets.ServiceModule+`
</g>
<g id="healthCircle">
<circle cx="10" cy="10" r="10" style="stroke:#38B44A;fill:none;stroke-width:2px"/>
<circle cx="10" cy="10" r="5" style="fill:#38B44A"/>
</g>
</defs>
<g id="relations">
<line x1="417" y1="189" x2="189" y2="351" stroke="#38B44A" stroke-width="2px" stroke-dasharray="129.85, 20" />
<use x="293" y="260" xlink:href="#healthCircle" />
<line x1="417" y1="189" x2="544" y2="276" stroke="#38B44A" stroke-width="2px" stroke-dasharray="66.97, 20" />
<use x="470" y="222" xlink:href="#healthCircle" />
</g>
<g id="services">
<use x="323" y="0" xlink:href="#serviceBlock" id="charmworld" />
<image x="369" y="46" width="96" height="96" xlink:href="http://0.1.2.3/~juju-jitsu/precise/charmworld-58.svg" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="417" y="31" >charmworld</text>
</g>
<use x="0" y="257" xlink:href="#serviceBlock" id="elasticsearch" />
<image x="46" y="303" width="96" height="96" xlink:href="http://0.1.2.3/~charming-devs/precise/elasticsearch-2.svg" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="94" y="288" >elasticsearch</text>
</g>
<use x="450" y="276" xlink:href="#serviceBlock" id="mongodb" />
<image x="496" y="322" width="96" height="96" xlink:href="http://0.1.2.3/precise/mongodb-21.svg" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="544" y="307" >mongodb</text>
</g>
</g>
</svg>
`))
}

func (s *newSuite) TestDefaultHTTPFetcher(c *gc.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<svg></svg>")
	}))
	defer ts.Close()

	tsIconUrl := func(ref *charm.Reference) string {
		return ts.URL + "/" + ref.Path() + ".svg"
	}

	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	cvs, err := NewFromBundle(b, tsIconUrl, &HTTPFetcher{IconURL: tsIconUrl})
	c.Assert(err, gc.IsNil)

	var buf bytes.Buffer
	cvs.Marshal(&buf)
	c.Logf("%s", buf.String())
	assertXMLEqual(c, buf.Bytes(), []byte(`
<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="639" height="465"
     style="font-family:Ubuntu, sans-serif;" viewBox="0 0 639 465"
     xmlns="http://www.w3.org/2000/svg" 
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
<g id="serviceBlock" transform="scale(0.8)" >`+assets.ServiceModule+`
</g>
<g id="healthCircle">
<circle cx="10" cy="10" r="10" style="stroke:#38B44A;fill:none;stroke-width:2px"/>
<circle cx="10" cy="10" r="5" style="fill:#38B44A"/>
</g>
<g id="icon-1" >
<svg:svg xmlns:svg="http://www.w3.org/2000/svg"></svg:svg></g>
<g id="icon-2" >
<svg:svg xmlns:svg="http://www.w3.org/2000/svg"></svg:svg></g>
<g id="icon-3" >
<svg:svg xmlns:svg="http://www.w3.org/2000/svg"></svg:svg></g>
</defs>
<g id="relations">
<line x1="417" y1="189" x2="189" y2="351" stroke="#38B44A" stroke-width="2px" stroke-dasharray="129.85, 20" />
<use x="293" y="260" xlink:href="#healthCircle" />
<line x1="417" y1="189" x2="544" y2="276" stroke="#38B44A" stroke-width="2px" stroke-dasharray="66.97, 20" />
<use x="470" y="222" xlink:href="#healthCircle" />
</g>
<g id="services">
<use x="323" y="0" xlink:href="#serviceBlock" id="charmworld" />
<use x="369" y="46" xlink:href="#icon-1" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="417" y="31" >charmworld</text>
</g>
<use x="0" y="257" xlink:href="#serviceBlock" id="elasticsearch" />
<use x="46" y="303" xlink:href="#icon-2" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="94" y="288" >elasticsearch</text>
</g>
<use x="450" y="276" xlink:href="#serviceBlock" id="mongodb" />
<use x="496" y="322" xlink:href="#icon-3" width="96" height="96" />
<g style="font-size:18px;fill:#505050;text-anchor:middle">
<text x="544" y="307" >mongodb</text>
</g>
</g>
</svg>
`))

}

func (s *newSuite) TestFetcherError(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	err = b.Verify(nil)
	c.Assert(err, gc.IsNil)

	ef := errFetcher("bad-wolf")
	_, err = NewFromBundle(b, iconURL, &ef)
	c.Assert(err, gc.ErrorMatches, "bad-wolf")
}

func (s *newSuite) TestWithBadBundle(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)
	b.Relations[0][0] = "evil-unknown-service"
	cvs, err := NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.ErrorMatches, "cannot verify bundle: .*")
	c.Assert(cvs, gc.IsNil)
}

func (s *newSuite) TestWithBadPosition(c *gc.C) {
	b, err := charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)

	b.Services["charmworld"].Annotations["gui-x"] = "bad"
	cvs, err := NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.ErrorMatches, `service "charmworld" does not have a valid position`)
	c.Assert(cvs, gc.IsNil)

	b, err = charm.ReadBundleData(strings.NewReader(bundle))
	c.Assert(err, gc.IsNil)

	b.Services["charmworld"].Annotations["gui-y"] = "bad"
	cvs, err = NewFromBundle(b, iconURL, nil)
	c.Assert(err, gc.ErrorMatches, `service "charmworld" does not have a valid position`)
	c.Assert(cvs, gc.IsNil)
}
