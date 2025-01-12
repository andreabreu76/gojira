# CLI Task Description Generator

Este é um aplicativo CLI escrito em Go para gerar descrições detalhadas de tarefas usando IA (OpenAI). Ele foi criado para melhorar a produtividade ao criar descrições para tarefas, bugs e épicos, mas o texto gerado deve sempre passar por uma revisão, pois é gerado por IA.

## Instalação de binarios

Execute o seguinte comando no terminal para baixar e instalar o aplicativo automaticamente:

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/andreabreu76/gojira/main/install.sh)"
```

Este script detecta o sistema operacional e a arquitetura automaticamente, baixa o binário correto da última release do repositório GitHub e o instala em `/usr/local/bin`.

## Requisitos para compilação

Certifique-se de que você possui as seguintes dependências instaladas no seu sistema:

1. **Go (Golang)**:
   - Instale o Go seguindo as instruções oficiais: [Download Go](https://go.dev/dl/).
   - Após instalar, verifique a instalação:
     ```bash
     go version
     ```

2. **Chave da OpenAI (OPENAI_API_KEY)**:
   - Para que o aplicativo funcione, é necessário ter uma chave válida da OpenAI.
   - Defina a chave no ambiente do sistema:
     - No Linux/MacOS:
       ```bash
       export OPENAI_API_KEY="sua-chave-aqui"
       ```
     - No Windows:
       ```cmd
       set OPENAI_API_KEY=sua-chave-aqui
       ```

   - Alternativamente, crie um arquivo `.env` no mesmo diretório do aplicativo com o seguinte conteúdo:
     ```
     OPENAI_API_KEY=sua-chave-aqui
     ```

## Como Baixar e Compilar

1. Clone este repositório em sua máquina:
   ```bash
   git clone git@github.com:andreabreu76/gojira.git
   cd gojira
   ```

2. Compile o binário:
   ```bash
   go build -o gojira
   ```

3. Crie um link simbólico para facilitar a execução:
   ```bash
   sudo ln -s $(pwd)/gojira /usr/local/bin/gojira
   ```

   Agora você pode executar o programa de qualquer lugar usando:
   ```bash
   gojira
   ```

## Como Usar

Execute o programa com as opções necessárias. Exemplo:

```bash
cli-task-gen --title "Corrigir erro no login" --type "BUG" --description "Usuários não conseguem acessar o sistema"
```

- **`--title`**: O título da tarefa (obrigatório).
- **`--type`**: O tipo da tarefa (EPICO, BUG ou TASK) (obrigatório).
- **`--description`**: Uma breve descrição da tarefa (opcional).

O programa gerará uma descrição detalhada e a copiará automaticamente para o clipboard.

### Transcrição de Commits Git

Se você estiver em um repositório Git e utilizar o tipo `Commit`, o programa irá:
1. Detectar automaticamente as alterações não comitadas (incluindo novos arquivos e modificações).
2. Gerar uma mensagem de commit detalhada no estilo **gitemoji**, contendo:
    - Um título sucinto.
    - Uma lista de alterações detectadas no `git diff`.
    - Menções a novos arquivos criados, com breves descrições baseadas no propósito da branch.

Exemplo de uso:

```bash
gojira --type "Commit"
```

#### Exemplo de saída:
Para uma branch chamada `feat/TOP-123-adicionar-funcionalidade`, o programa pode gerar uma mensagem como:

```plaintext
:sparkles: feat(TOP-123-adicionar-funcionalidade) Implementação da nova funcionalidade

- Adicionado `services/openai.go` para integrar com a API da OpenAI.
- Criado `utils/commons/env.go` para gerenciamento de variáveis de ambiente.
- Modularizada a estrutura do projeto para melhor organização.
```

Se a branch não seguir o padrão Git Flow (`tipo/nome-da-branch`), o programa utiliza um tipo genérico como `chore` e o nome da branch completo.

#### Requisitos para o tipo Commit:
- O comando deve ser executado dentro de um repositório Git válido.
- As alterações devem estar visíveis no `git status`.

## Importante

Este aplicativo foi criado para melhorar a produtividade, mas os textos gerados devem sempre ser revisados antes de serem usados, pois são produzidos por um modelo de IA e podem conter inconsistências.

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.
