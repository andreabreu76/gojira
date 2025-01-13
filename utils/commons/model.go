package commons

import "strings"

func GetModel(taskType string) string {
	switch strings.ToUpper(taskType) {
	case "EPICO", "TASK":
		return `
Objetivo:
Como:
Critérios de Aceite:
Testes:
Informações para o time de Infra:
Outras observações:
`
	case "BUG":
		return `
Resumo

Título do bug:
ID do bug:
Data de identificação:
Cliente:
Autor do ticket:

Descrição:
Passos para reproduzir:
Comportamento esperado:
Comportamento real:
Capturas de tela:

Ambiente:
Sistema operacional:
Navegador:
Versão do App/Dashboard/Emissor:
Dispositivo:

Informações para o time de Infra:

Outras observações:
`
	default:
		return ""
	}
}
