package shortlink

import (
	"github.com/jiaqi-yin/go-url-shortener/utils"
	"gopkg.in/validator.v2"
)

type EncodedID struct {
	Shortlink string `uri:"shortlink" binding:"required"`
}

type ShortenRequest struct {
	URL string `json:"url" validate:"nonzero"`
}

type ShortlinkResponse struct {
	Shortlink string `json:"shortlink"`
}

func (req *ShortenRequest) Validate() utils.RestErr {
	if err := validator.Validate(req); err != nil {
		return utils.NewBadRequestError("missing url in the request")
	}
	return nil
}
