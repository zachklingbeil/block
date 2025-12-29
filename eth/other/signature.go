package other

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SigProvider struct {
	BaseURL string
	Client  *http.Client
}

func NewSigProvider(baseURL string) *SigProvider {
	return &SigProvider{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// DecodeFunctionInput calls /api/v1/abi/function?txInput=...
func (s *SigProvider) DecodeFunctionInput(txInput string) ([]map[string]any, error) {
	url := fmt.Sprintf("%s/api/v1/abi/function?txInput=%s", s.BaseURL, txInput)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeEvent calls /api/v1/abi/event?data=...&topics=...
func (s *SigProvider) DecodeEvent(data, topics string) ([]map[string]any, error) {
	url := fmt.Sprintf("%s/api/v1/abi/event?data=%s&topics=%s", s.BaseURL, data, topics)
	resp, err := s.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
