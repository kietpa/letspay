package xendit

import "letspay/repository/provider"

type (
	providerRepo struct {
		url    string
		apiKey string
	}

	NewProviderRepoInput struct {
		Url    string
		ApiKey string
	}
)

func NewProviderRepo(input NewProviderRepoInput) provider.ProviderRepo {
	return &providerRepo{
		url:    input.Url,
		apiKey: input.ApiKey,
	}
}

func (p *providerRepo) ExecuteDisbursement() {}

func (p *providerRepo) GetDisbursementStatus() {}
