package opensearch

import (
	"log/slog"

	"github.com/sfomuseum/go-edtf"
	"github.com/sfomuseum/go-edtf/common"
	"github.com/sfomuseum/go-edtf/parser"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-flags"
	"github.com/whosonfirst/go-whosonfirst-flags/existential"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type SpelunkerRecordSPR struct {
	wof_spr.StandardPlacesResult

	// The problem with this is that we can't use it with api/spr.go
	// endpoint if it is using the Spelunker interface GetSPRWithId
	// method.

	props []byte
}

type SpelunkerStandardPlacesResults struct {
	wof_spr.StandardPlacesResults
	results []wof_spr.StandardPlacesResult
}

func (r *SpelunkerStandardPlacesResults) Results() []wof_spr.StandardPlacesResult {
	return r.results
}

func NewSpelunkerStandardPlacesResults(results []wof_spr.StandardPlacesResult) wof_spr.StandardPlacesResults {

	r := &SpelunkerStandardPlacesResults{
		results: results,
	}

	return r
}

func NewSpelunkerRecordSPR(props []byte) (wof_spr.StandardPlacesResult, error) {

	s := &SpelunkerRecordSPR{
		props: props,
	}

	return s, nil
}

func (s *SpelunkerRecordSPR) Id() string {
	return gjson.GetBytes(s.props, "wof:id").String()
}

func (s *SpelunkerRecordSPR) ParentId() string {
	return gjson.GetBytes(s.props, "wof:parent_id").String()
}

func (s *SpelunkerRecordSPR) Name() string {
	return gjson.GetBytes(s.props, "wof:name").String()
}

func (s *SpelunkerRecordSPR) Placetype() string {
	return gjson.GetBytes(s.props, "wof:placetype").String()
}

func (s *SpelunkerRecordSPR) Country() string {
	return gjson.GetBytes(s.props, "wof:country").String()
}

func (s *SpelunkerRecordSPR) Repo() string {
	return gjson.GetBytes(s.props, "wof:repo").String()
}

func (s *SpelunkerRecordSPR) Path() string {

	id := gjson.GetBytes(s.props, "wof:id").Int()
	path, _ := uri.Id2RelPath(id)
	return path
}

func (s *SpelunkerRecordSPR) URI() string {
	return s.Path()
}

func (s *SpelunkerRecordSPR) Inception() *edtf.EDTFDate {
	return s.edtfDate("edtf:inception")
}

func (s *SpelunkerRecordSPR) Cessation() *edtf.EDTFDate {
	return s.edtfDate("edtf:cessation")
}

func (s *SpelunkerRecordSPR) Latitude() float64 {
	return gjson.GetBytes(s.props, "geom:latitude").Float()
}

func (s *SpelunkerRecordSPR) Longitude() float64 {
	return gjson.GetBytes(s.props, "geom:longitude").Float()
}

func (s *SpelunkerRecordSPR) MinLatitude() float64 {
	return gjson.GetBytes(s.props, "geom:bbox.1").Float()
}

func (s *SpelunkerRecordSPR) MinLongitude() float64 {
	return gjson.GetBytes(s.props, "geom:bbox.0").Float()
}

func (s *SpelunkerRecordSPR) MaxLatitude() float64 {
	return gjson.GetBytes(s.props, "geom:bbox.3").Float()
}

func (s *SpelunkerRecordSPR) MaxLongitude() float64 {
	return gjson.GetBytes(s.props, "geom:bbox.2").Float()
}

func (s *SpelunkerRecordSPR) IsCurrent() flags.ExistentialFlag {
	fl_i := gjson.GetBytes(s.props, "mz:is_current").Int()
	return s.existentialFlag(fl_i)
}

func (s *SpelunkerRecordSPR) IsCeased() flags.ExistentialFlag {

	fl_i := int64(0)

	r := gjson.GetBytes(s.props, "edtf:cessation")

	if r.Exists() {

		switch r.String() {
		case edtf.UNKNOWN, edtf.UNKNOWN_2012:
			fl_i = -1
		default:
			fl_i = 1
		}
	}

	return s.existentialFlag(fl_i)
}

func (s *SpelunkerRecordSPR) IsDeprecated() flags.ExistentialFlag {

	fl_i := int64(0)

	r := gjson.GetBytes(s.props, "edtf:deprecated")

	if r.Exists() && r.String() != "" {
		fl_i = 1
	}

	return s.existentialFlag(fl_i)
}

func (s *SpelunkerRecordSPR) IsSuperseded() flags.ExistentialFlag {

	fl_i := int64(0)

	if len(s.SupersededBy()) > 0 {
		fl_i = 1
	}

	return s.existentialFlag(fl_i)
}

func (s *SpelunkerRecordSPR) IsSuperseding() flags.ExistentialFlag {

	fl_i := int64(0)

	if len(s.Supersedes()) > 0 {
		fl_i = 1
	}

	return s.existentialFlag(fl_i)
}

func (s *SpelunkerRecordSPR) SupersededBy() []int64 {

	return s.gatherIds("wof:superseded_by")
}

func (s *SpelunkerRecordSPR) Supersedes() []int64 {

	return s.gatherIds("wof:supersedes")
}

func (s *SpelunkerRecordSPR) BelongsTo() []int64 {

	return s.gatherIds("wof:belongsto")
}

func (s *SpelunkerRecordSPR) LastModified() int64 {

	return gjson.GetBytes(s.props, "wof:lastmodified").Int()
}

func (s *SpelunkerRecordSPR) edtfDate(path string) *edtf.EDTFDate {

	str_dt := gjson.GetBytes(s.props, path).String()
	dt, err := parser.ParseString(str_dt)

	if err != nil {
		slog.Error("Failed to parse date", "id", s.Id(), "path", path, "date", str_dt, "error", err)
		return s.unknownEDTF()
	}

	return dt
}

func (s *SpelunkerRecordSPR) unknownEDTF() *edtf.EDTFDate {

	sp := common.UnknownDateSpan()

	d := &edtf.EDTFDate{
		Start:   sp.Start,
		End:     sp.End,
		EDTF:    edtf.UNKNOWN,
		Level:   -1,
		Feature: "Unknown",
	}

	return d
}

func (s *SpelunkerRecordSPR) existentialFlag(fl_i int64) flags.ExistentialFlag {

	fl, err := existential.NewKnownUnknownFlag(fl_i)

	if err != nil {
		fl, _ = existential.NewNullFlag()
	}

	return fl
}

func (s *SpelunkerRecordSPR) gatherIds(path string) []int64 {

	ids := make([]int64, 0)

	r := gjson.GetBytes(s.props, "wof:superseded_by")

	if !r.Exists() {
		return ids
	}

	ra := r.Array()
	count := len(ra)

	if count == 0 {
		return ids
	}

	ids = make([]int64, count)

	for idx, a := range ra {
		ids[idx] = a.Int()
	}

	return ids
}
