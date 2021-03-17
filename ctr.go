package whatcountry

import (
	"io/ioutil"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

type Countries struct {
	fc *geojson.FeatureCollection
}

type Country struct {
	Name       string
	Iso_a2     string
	Region     string
	Iso_3166_2 string
}

func LoadCountries(file string) (Countries, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return Countries{}, err
	}
	featureCollection, err := geojson.UnmarshalFeatureCollection(b)
	if err != nil {
		return Countries{}, err
	}
	return Countries{fc: featureCollection}, nil
}

func (c *Countries) FindPoint(lon, lat float64) Country {
	countries := c.findCountries(orb.Point{lon, lat})
	var ctr Country
	if len(countries) == 0 {
		return ctr
	}
	t := countries[0]
	ctr.Name = t.Properties["admin"].(string)
	ctr.Iso_a2 = t.Properties["iso_a2"].(string)
	ctr.Region = t.Properties["name"].(string)
	ctr.Iso_3166_2 = t.Properties["iso_3166_2"].(string)
	return ctr
}

func (c *Countries) findCountries(point orb.Point) []geojson.Feature {
	countries := []geojson.Feature{}
	for _, feature := range c.fc.Features {
		// Try on a MultiPolygon to begin
		multiPoly, isMulti := feature.Geometry.(orb.MultiPolygon)
		if isMulti {
			if planar.MultiPolygonContains(multiPoly, point) {
				countries = append(countries, *feature)
			}
		} else {
			// Fallback to Polygon
			polygon, isPoly := feature.Geometry.(orb.Polygon)
			if isPoly {
				if planar.PolygonContains(polygon, point) {
					countries = append(countries, *feature)
				}
			}
		}
	}
	return countries
}
