package ai

// Provider é uma interface que define os métodos que um provedor de IA deve implementar
type Provider interface {
	// GetName retorna o nome do provedor
	GetName() string

	// GetAvailableModels retorna a lista de modelos disponíveis para este provedor
	GetAvailableModels() []string

	// GetDefaultModel retorna o modelo padrão a ser usado
	GetDefaultModel() string

	// GetCompletions envia um prompt para o provedor de IA e retorna a resposta
	GetCompletions(prompt string, modelID string) (string, error)
}

// ProviderFactory é um mapa de funções que criam instâncias de provedores de IA
var ProviderFactory = map[string]func() Provider{
	"openai":    NewOpenAIProvider,
	"anthropic": NewAnthropicProvider,
}

// GetProvider retorna uma instância do provedor especificado
func GetProvider(providerName string) (Provider, bool) {
	factory, exists := ProviderFactory[providerName]
	if !exists {
		return nil, false
	}

	return factory(), true
}

// GetDefaultProvider retorna o provedor padrão
func GetDefaultProvider() Provider {
	// Por padrão, usamos OpenAI
	return NewOpenAIProvider()
}
