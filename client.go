package musicbrainz

import (
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid"

	md "github.com/ytsiuryn/ds-audiomd"
	srv "github.com/ytsiuryn/ds-microservice"
)

// AudioOnlineRequest описывает структуру запроса к микросервису.
type AudioOnlineRequest struct {
	Cmd     string      `json:"cmd"`
	Release *md.Release `json:"release"`
	// Actor
	// *md.Publishing
}

// AudioOnlineResponse описывает структуру ответа микросервиса.
type AudioOnlineResponse struct {
	SuggestionSet *md.SuggestionSet  `json:"suggestion_set,omitempty"`
	Error         *srv.ErrorResponse `json:"error,omitempty"`
}

// type AudioOnlineDBClient struct {
// 	*srv.RPCClient
// 	req *AudioOnlineRequest
// }

// Unwrap контроллирует значение ответа микросервиса, и, в случае ошибки,
// печатает сведения об ошибке и останавливает процесс с запущенным клиентом.
func (resp *AudioOnlineResponse) Unwrap() *md.SuggestionSet {
	if resp.Error != nil {
		srv.FailOnError(errors.New(resp.Error.Error), resp.Error.Context)
	}
	return resp.SuggestionSet
}

// CreateReleaseRequest формирует данные запроса поиска релиза по указанным метаданным.
func CreateReleaseRequest(r *md.Release) (string, []byte, error) {
	correlationID, _ := uuid.NewV4()
	req := AudioOnlineRequest{
		Cmd:     "release",
		Release: r}
	data, err := json.Marshal(&req)
	if err != nil {
		return "", nil, err
	}
	return correlationID.String(), data, nil
}

// ParseReleaseAnswer разбирает ответ с предложением метаданных релиза.
func ParseReleaseAnswer(data []byte) (*AudioOnlineResponse, error) {
	resp := AudioOnlineResponse{SuggestionSet: md.NewSuggestionSet()}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
