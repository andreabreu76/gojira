# Gojira 🦖

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/gojira)](https://goreportcard.com/report/github.com/yourusername/gojira)
[![GitHub release](https://img.shields.io/github/release/yourusername/gojira.svg)](https://github.com/yourusername/gojira/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yourusername/gojira)](https://github.com/yourusername/gojira)
[![Contributors](https://img.shields.io/github/contributors/yourusername/gojira)](https://github.com/yourusername/gojira/graphs/contributors)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/yourusername/gojira/pulls)
[![Stars](https://img.shields.io/github/stars/yourusername/gojira)](https://github.com/yourusername/gojira/stargazers)
[![Forks](https://img.shields.io/github/forks/yourusername/gojira)](https://github.com/yourusername/gojira/network/members)
[![Issues](https://img.shields.io/github/issues/yourusername/gojira)](https://github.com/yourusername/gojira/issues)
[![GitHub last commit](https://img.shields.io/github/last-commit/yourusername/gojira)](https://github.com/yourusername/gojira/commits/main)

<div align="center">
  <p>
    <img src="https://raw.githubusercontent.com/yourusername/gojira/main/assets/gojira-logo.png" alt="Gojira Logo" width="250">
  </p>
  <h3>Ferramenta CLI com IA para Potencializar seu Desenvolvimento</h3>
  <p>
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/Made%20with-Go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go"></a>
    <a href="https://github.com/yourusername/gojira/issues"><img src="https://img.shields.io/github/issues/yourusername/gojira?style=flat-square" alt="Issues"></a>
    <a href="https://github.com/yourusername/gojira/stargazers"><img src="https://img.shields.io/github/stars/yourusername/gojira?style=flat-square" alt="Stars"></a>
    <a href="https://github.com/yourusername/gojira/network/members"><img src="https://img.shields.io/github/forks/yourusername/gojira?style=flat-square" alt="Forks"></a>
    <a href="https://discord.gg/gojira"><img src="https://img.shields.io/badge/Discord-Join%20Us-7289DA?style=flat-square&logo=discord&logoColor=white" alt="Discord"></a>
  </p>
</div>

## 📖 Descrição

<div align="center">
  <img src="https://raw.githubusercontent.com/yourusername/gojira/main/assets/demo.gif" alt="Gojira Demo" width="700">
</div>

<br>
Gojira é uma ferramenta CLI poderosa projetada para agilizar o processo de desenvolvimento de software, integrando recursos de inteligência artificial para automação de tarefas comuns. O Gojira permite gerar mensagens de commit, analisar código, criar documentação, interagir com o Jira e facilitar o gerenciamento de tarefas de desenvolvimento.

## ✨ Funcionalidades Principais
- **Múltiplos Provedores de IA**: Suporte para OpenAI (GPT-4) e Anthropic (Claude), facilmente extensível para outros provedores
- **Integração com Jira**: Geração de descrições para tarefas do Jira e criação automática de issues
- **Geração de Documentação**: README, análise de código e checklists
- **Gerenciamento de Git**: Criação de branches, geração de mensagens de commit padronizadas
- **Workflow de Desenvolvimento**: Comandos para iniciar tarefas, criar checklists e agilizar workflows
- **Visualização Kanban**: Mostra as tarefas do Jira em formato kanban no terminal
- **Explicação de Código**: Analisa e explica trechos de código complexos
- **Geração de Testes**: Cria testes automatizados para código existente
- **Criação de PRs**: Gera Pull Requests com descrições detalhadas
- **Resumo de Alterações**: Cria resumos das mudanças no código
- **Standups Automáticos**: Gera relatórios para daily standups

## 🛠️ Tecnologias Utilizadas
- **Go**: Linguagem de programação principal
- **APIs de IA**: OpenAI API e Anthropic API para geração de conteúdo
- **Jira API**: Para integração com o Atlassian Jira
- **Git**: Para integração com repositórios Git

## 🔧 Instalação

1. **Clone o Repositório**:
   ```bash
   git clone https://github.com/yourusername/gojira.git
   cd gojira
   ```

2. **Instale Dependências**:
   ```bash
   bash install.sh
   ```

3. **Compile o Projeto**:
   ```bash
   go build -o gojira
   ```

4. **Configure o Gojira**:
   ```bash
   ./gojira config --provider openai --jira-url https://your-jira-instance.atlassian.net --jira-token your-jira-token
   ```

## 📋 Uso

### ⚙️ Configuração
```bash
# Mostrar configuração atual
./gojira config show

# Listar provedores de IA disponíveis
./gojira config providers

# Configurar o provedor de IA
./gojira config --provider anthropic --model claude-3-5-sonnet-20240620

# Configurar integração com Jira
./gojira config --jira-url https://your-jira-instance.atlassian.net --jira-token your-jira-token --jira-project PROJ
```

### 📝 Geração de Documentação
```bash
# Gerar README.md para o projeto
./gojira generate readme

# Gerar análise de código-fonte
./gojira generate analysis
```

### 🌱 Integração com Git
```bash
# Gerar mensagem de commit baseada nas alterações atuais
./gojira commit

# Criar um Pull Request (PR) com título e descrição gerados automaticamente
./gojira pr --base main

# Criar um PR com título personalizado e como rascunho
./gojira pr --title "Implementar autenticação OAuth" --draft

# Resumir alterações desde um commit específico
./gojira summary --base HEAD~5

# Resumir alterações dos últimos 10 commits e salvar em um arquivo
./gojira summary --base HEAD~10 --save --output resumo-alteracoes.md
```

### 🔄 Integração com Jira
```bash
# Gerar descrição para uma tarefa do Jira
./gojira jira --title "Implementar autenticação OAuth" --type TASK

# Criar uma issue no Jira
./gojira jira --title "Corrigir bug na página de login" --type BUG --project PROJ

# Visualizar quadro Kanban do Jira no terminal
./gojira kanban --project PROJ

# Filtrar tarefas do Kanban por usuário
./gojira kanban --project PROJ --user johndoe@example.com

# Exibir tarefas em formato texto simples
./gojira kanban --project PROJ --format plain
```

### 🔄 Workflow de Desenvolvimento
```bash
# Iniciar trabalho em uma issue
./gojira dev start --issue PROJ-123

# Criar uma branch para uma issue
./gojira dev branch --issue PROJ-123 --name "implementar-oauth"

# Gerar checklist para uma issue
./gojira dev checklist --issue PROJ-123

# Gerar relatório de standup
./gojira standup

# Relatório de standup dos últimos 3 dias
./gojira standup --days 3

# Relatório de standup e exportar para arquivo
./gojira standup --output standup.md
```

### 💻 Entendimento e Geração de Código
```bash
# Explicar um arquivo de código
./gojira explain --file /caminho/para/arquivo.go

# Explicar linhas específicas de um arquivo
./gojira explain --file /caminho/para/arquivo.py --start 10 --end 50

# Explicar código para um nível específico de desenvolvedor
./gojira explain --file /caminho/para/arquivo.js --level beginner

# Gerar testes automaticamente
./gojira test-gen --source /caminho/para/arquivo.go

# Gerar testes com alta cobertura
./gojira test-gen --source /caminho/para/arquivo.py --coverage alta

# Especificar framework de testes
./gojira test-gen --source /caminho/para/arquivo.js --framework jest
```

## 📚 Comandos Detalhados

### 🔀 PR - Pull Request
Cria um Pull Request (PR) utilizando a configuração do Git e integração com plataformas como GitHub ou GitLab. Gera automaticamente título e descrição detalhada com IA se não forem fornecidos.

```bash
./gojira pr [flags]

Flags:
  -b, --branch string    Nome da branch de origem (padrão: branch atual)
  -B, --base string      Branch base para o PR (default "main")
  -d, --description string  Descrição do PR
  -D, --draft            Criar o PR como rascunho
  -r, --remote string    Repositório remoto (formato: owner/repo)
  -t, --title string     Título do PR
```

### 🔍 Explain - Explicação de Código
Analisa e explica o funcionamento de trechos de código, classes, funções ou arquivos inteiros, tornando mais fácil entender código complexo ou legado.

```bash
./gojira explain [flags]

Flags:
  -f, --file string      Caminho para o arquivo a ser explicado
  -s, --start int        Linha inicial (opcional)
  -e, --end int          Linha final (opcional)
  -o, --output string    Arquivo para salvar a explicação (opcional)
  -l, --level string     Nível de experiência do desenvolvedor (beginner, intermediate, expert) (default "intermediate")
```

### 🧪 Test-Gen - Geração de Testes
Analisa o código-fonte e gera testes automaticamente utilizando inteligência artificial, cobrindo funções, classes e métodos.

```bash
./gojira test-gen [flags]

Flags:
  -s, --source string     Arquivo fonte para o qual gerar testes
  -o, --output string     Arquivo de saída para os testes (opcional)
  -f, --framework string  Framework de testes a ser usado (opcional)
  -c, --coverage string   Nível de cobertura desejado (básica, média, alta) (default "alta")
```

### 📊 Kanban - Visualização de Tarefas
Exibe as tarefas do Jira em um formato de quadro kanban diretamente no terminal, permitindo visualizar o progresso das tarefas sem sair da linha de comando.

```bash
./gojira kanban [flags]

Flags:
  -p, --project string    Chave do projeto Jira
  -u, --user string       Filtrar por usuário
  -s, --status string     Filtrar por status
  -l, --limit int         Número máximo de tarefas por status (default 10)
  -f, --format string     Formato de saída (color, plain) (default "color")
```

### 📈 Summary - Resumo de Alterações
Analisa as alterações no código desde um commit ou branch específica e gera um resumo detalhado do que foi alterado, adicionado ou removido.

```bash
./gojira summary [flags]

Flags:
  -b, --base string       Commit ou branch base para comparação (padrão: HEAD~10)
  -f, --format string     Formato do relatório (markdown, jira, text, html) (default "markdown")
  -s, --save              Salvar relatório em um arquivo
  -o, --output string     Arquivo para salvar o relatório (padrão: alteracoes-resumo.md)
  -m, --max int           Número máximo de arquivos a incluir (0 para todos) (default 20)
  -c, --code              Incluir código detalhado no prompt (aumenta precisão, mas consome mais tokens)
```

### 📢 Standup - Relatórios Diários
Analisa as atividades recentes (commits, issues, etc.) e gera um relatório formatado para reuniões de standup diárias, detalhando o que foi feito, o que está planejado e quaisquer bloqueios.

```bash
./gojira standup [flags]

Flags:
  -d, --days int          Número de dias para incluir no relatório (default 1)
  -e, --email string      Email do usuário para filtrar as atividades (padrão: email do git config)
  -t, --team              Incluir atividades de toda a equipe, não apenas do usuário
  -i, --issues            Focar apenas em issues, ignorando commits
  -o, --output string     Arquivo para salvar o relatório (opcional)
```

## ⚙️ Configuração de Ambiente

O Gojira necessita das seguintes variáveis de ambiente ou arquivos de configuração:

1. Um arquivo `.env` na raiz do projeto ou variáveis de ambiente do sistema:
   - `OPENAI_API_KEY`: Chave de API para o OpenAI
   - `ANTHROPIC_API_KEY`: Chave de API para o Anthropic

2. Arquivo de configuração `~/.gojira.json` (criado automaticamente):
   - Provedor de IA preferido
   - Modelo de IA preferido
   - Configurações do Jira

### 🔑 Como obter as chaves de API

#### OpenAI API Key
1. Crie uma conta ou faça login em [OpenAI Platform](https://platform.openai.com/)
2. Navegue até "API Keys" no painel
3. Clique em "Create new secret key" 
4. Copie a chave gerada e salve-a em seu arquivo `.env` como `OPENAI_API_KEY=sua-chave-aqui`

#### Anthropic API Key
1. Crie uma conta ou faça login em [Anthropic Console](https://console.anthropic.com/)
2. Navegue até "API Keys" no painel de controle
3. Clique em "Create Key"
4. Copie a chave gerada e salve-a em seu arquivo `.env` como `ANTHROPIC_API_KEY=sua-chave-aqui`

## 👥 Como Contribuir

1. Fork o repositório
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Faça commit das suas alterações (`git commit -m 'feat: adicionar nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## 📄 Licença
Este projeto está licenciado sob a licença MIT.

## ⭐ Showcase

<div align="center">
  <table>
    <tr>
      <td align="center">
        <img src="https://raw.githubusercontent.com/yourusername/gojira/main/assets/terminal-kanban.png" width="400px" alt="Terminal Kanban View"/>
        <br />
        <i>Visualização Kanban no Terminal</i>
      </td>
      <td align="center">
        <img src="https://raw.githubusercontent.com/yourusername/gojira/main/assets/code-explanation.png" width="400px" alt="Code Explanation"/>
        <br />
        <i>Explicação de Código Detalhada</i>
      </td>
    </tr>
    <tr>
      <td align="center">
        <img src="https://raw.githubusercontent.com/yourusername/gojira/main/assets/test-generation.png" width="400px" alt="Test Generation"/>
        <br />
        <i>Geração Automática de Testes</i>
      </td>
      <td align="center">
        <img src="https://raw.githubusercontent.com/yourusername/gojira/main/assets/standup-report.png" width="400px" alt="Standup Report"/>
        <br />
        <i>Relatório de Standup Diário</i>
      </td>
    </tr>
  </table>
</div>

## 🙏 Agradecimentos
- Todos os contribuidores que ajudaram a tornar o Gojira melhor
- [OpenAI](https://openai.com/) e [Anthropic](https://www.anthropic.com/) pelos poderosos modelos de IA
- A comunidade open source por todas as ferramentas e bibliotecas utilizadas

---

<div align="center">
  <sub>Construído com ❤️ por todos os colaboradores</sub>
  <br>
  <sub>© 2023-2025</sub>
</div>