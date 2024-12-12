
# CLI Task Description Generator

Este é um aplicativo CLI escrito em Go para gerar descrições detalhadas de tarefas usando IA (OpenAI). Ele foi criado para melhorar a produtividade ao criar descrições para tarefas, bugs e épicos, mas o texto gerado deve sempre passar por uma revisão, pois é gerado por IA.

## Requisitos

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

## Como Instalar

### Instalação Automática via Script

Execute o seguinte comando no terminal para baixar e instalar o aplicativo automaticamente:

```bash
sh -c "$(curl -fsSL https://github.com/andreabreu76/gojira/install.sh)"
```

Este script detecta o sistema operacional e a arquitetura automaticamente, baixa o binário correto da última release do repositório GitHub e o instala em `/usr/local/bin`.

### Instalação Manual

1. Clone este repositório em sua máquina:
   ```bash
   git clone <URL_DO_REPOSITORIO>
   cd <NOME_DO_DIRETORIO>
   ```

2. Compile o binário:
   ```bash
   go build -o cli-task-gen
   ```

3. Crie um link simbólico para facilitar a execução:
   ```bash
   sudo ln -s $(pwd)/cli-task-gen /usr/local/bin/cli-task-gen
   ```

   Agora você pode executar o programa de qualquer lugar usando:
   ```bash
   cli-task-gen
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

## Importante

Este aplicativo foi criado para melhorar a produtividade, mas os textos gerados devem sempre ser revisados antes de serem usados, pois são produzidos por um modelo de IA e podem conter inconsistências.

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

## Licença

Este projeto está sob a licença MIT. Consulte o arquivo `LICENSE` para mais detalhes.
