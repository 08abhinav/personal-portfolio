package geo

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type Location struct {
	City, Region, Country, Timezone, Org string
	Lat, Lon                             float64
	HasCoords                            bool
	ASN                                  uint
}

type Resolver struct {
	city *geoip2.Reader
	asn  *geoip2.Reader
}

func Open(cityPath, asnPath string) (*Resolver, error) {
	city, err := geoip2.Open(cityPath)
	if err != nil {
		return nil, fmt.Errorf("open city db: %w", err)
	}
	r := &Resolver{city: city}
	if asnPath != "" {
		asn, err := geoip2.Open(asnPath)
		if err != nil {
			city.Close()
			return nil, fmt.Errorf("open asn db: %w", err)
		}
		r.asn = asn
	}
	return r, nil
}

func (r *Resolver) Close() {
	r.city.Close()
	if r.asn != nil {
		r.asn.Close()
	}
}

func (r *Resolver) Lookup(ipStr string) (*Location, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid ip %q", ipStr)
	}
	rec, err := r.city.City(ip)
	if err != nil {
		return nil, err
	}

	l := &Location{
		City:      rec.City.Names["en"],
		Country:   rec.Country.Names["en"],
		Timezone:  rec.Location.TimeZone,
		Lat:       rec.Location.Latitude,
		Lon:       rec.Location.Longitude,
		HasCoords: rec.Location.Latitude != 0 || rec.Location.Longitude != 0,
	}
	if len(rec.Subdivisions) > 0 {
		l.Region = rec.Subdivisions[0].Names["en"]
	}
	if r.asn != nil {
		if a, err := r.asn.ASN(ip); err == nil {
			l.ASN = a.AutonomousSystemNumber
			l.Org = a.AutonomousSystemOrganization
		}
	}
	return l, nil
}
