package tokenhelper

import (
	"strings"

	"github.com/disgoorg/snowflake/v2"

	"github.com/fluxergo/fluxergo/fluxer"
)

// IDFromToken returns the applicationID from the token
func IDFromToken(token string) (*snowflake.ID, error) {
	strs := strings.Split(token, ".")
	if len(strs) == 0 {
		return nil, fluxer.ErrInvalidBotToken
	}

	strID, err := snowflake.Parse(strs[0])
	if err != nil {
		return nil, err
	}
	return &strID, nil
}
