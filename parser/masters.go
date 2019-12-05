package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseMasters(d *xml.Decoder, limit int) []model.Master {
	cnt := 0
	masters := make([]model.Master, 0, 0)
	for t, err := d.Token(); t != nil && err == nil && cnt+1 != limit; t, err = d.Token() {
		if IsStartElementName(t, "master") {
			masters = append(masters, ParseMaster(t.(xml.StartElement), d))
			cnt++
		}
	}

	return masters
}

func ParseMaster(se xml.StartElement, tr xml.TokenReader) model.Master {
	master := model.Master{}

	if se.Name.Local != "master" {
		return master
	}

	master.Id = se.Attr[0].Value
	for {
		t, _ := tr.Token()
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				master.Images = ParseImages(se, tr)
			case "main_release":
				master.MainRelease = parseValue(tr)
			case "artists":
				master.Artists = parseArtists("artists", tr)
			case "genres":
				master.Genres = parseChildValues("genres", "genre", tr)
			case "styles":
				master.Styles = parseChildValues("styles", "style", tr)
			case "year":
				master.Year = parseValue(tr)
			case "title":
				master.Title = parseValue(tr)
			case "data_quality":
				master.DataQuality = parseValue(tr)
			case "videos":
				master.Videos = ParseVideos(tr)
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "master" {
			break
		}
	}

	return master
}
