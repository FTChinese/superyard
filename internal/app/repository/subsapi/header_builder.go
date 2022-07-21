package subsapi

import (
	"github.com/FTChinese/superyard/pkg/ids"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"net/http"
)

func ReaderIDsHeader(id ids.UserIDs) http.Header {
	b := xhttp.NewHeaderBuilder()
	if id.FtcID.Valid {
		b.WithFtcID(id.FtcID.String)
	}

	if id.UnionID.Valid {
		b.WithUnionID(id.UnionID.String)
	}

	return b.Build()
}
