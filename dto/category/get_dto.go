package categorydto

import (
	utilDto "github.com/susatyo441/go-ta-utils/dto"
)

type GetCategoryOptionsResponse struct {
	CategoryOptions []utilDto.InterfaceOptionDTO `json:"categoryOptions"`
}
