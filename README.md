# Comparei - Serviço de Listas 📋

O **Comparei - Serviço de Listas** é um microsserviço essencial do ecossistema "Comparei", focado na criação, manutenção e gerenciamento de listas de compras dos usuários. 

Além de fornecer uma API REST para as operações de CRUD (Criar, Ler, Atualizar e Deletar) das listas, este serviço escuta eventos de atualização de preços de forma assíncrona, garantindo que as listas dos usuários reflitam os valores mais recentes do mercado.

## 🛠️ Tecnologias Utilizadas

* **Linguagem:** [Go 1.23](https://golang.org/)
* **Banco de Dados Relacional:** **MySQL 8** (para persistência segura das listas e itens dos usuários).
* **Mensageria em Cache:** **Redis** (para arquitetura orientada a eventos, especificamente escutando alterações de preços).
* **Infraestrutura:** Docker e Docker Compose.
* **Roteamento HTTP:** [Gorilla Mux](https://github.com/gorilla/mux)
* **Segurança:** Validação e autenticação de rotas.

## ⚙️ Arquitetura do Sistema

A aplicação roda em duas frentes simultâneas:
1. **Servidor HTTP:** Expõe *endpoints* para gerenciar os dados das listas de compras.
2. **Subscriber (Mensageria):** Uma *goroutine* dedicada a ouvir o Redis em busca de eventos de "prices" (preços), garantindo a reatividade do sistema às flutuações de mercado.

## 🚀 Como Executar o Projeto Localmente

### Pré-requisitos
* [Go 1.23+](https://golang.org/dl/) instalado no seu ambiente local.
* [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados para subir os serviços de banco de dados e mensageria.

### Passo a Passo

1. **Clonar o repositório:**
```bash
   git clone https://github.com/ViniciusPereira03/comparei-servico-proconf
   cd comparei-servico-listas

```

2. **Configuração das Variáveis de Ambiente:**
A partir do arquivo de exemplo, crie o seu próprio arquivo `.env` na raiz do projeto:
```bash
cp .env.example .env

```


Preencha o `.env` com os valores adequados. Para rodar a aplicação localmente (fora do Docker), certifique-se de apontar os hosts para o `localhost`:
```env
# MySQL
MYSQL_HOST=localhost:3306
MYSQL_USER=root
MYSQL_PASSWORD=root
MYSQL_DB=listasdb

# Servidor HTTP
PORT=8086

# Redis
REDIS_MESSAGING_HOST=localhost
REDIS_MESSAGING_PORT=6379

```

*(Lembre-se de garantir que o contêiner do Redis também esteja rodando na sua rede local `comparei_net`).*
3. **Executar a Aplicação com `run.sh`:**
Dê a permissão de execução ao script (caso necessário) e inicialize o projeto:
```bash
chmod +x run.sh
./run.sh

```


> **⚠️ Nota Importante sobre o `wait-for-it.sh`:** > A inicialização do ambiente está configurada para utilizar o script `wait-for-it.sh`. Esse script é uma ferramenta inteligente que impede que a aplicação Go tente se conectar ao banco de dados ou ao Redis antes que eles estejam totalmente inicializados e prontos para receber conexões. Isso evita erros de "connection refused" durante o *startup* do microsserviço.


4. **Acompanhar os Logs:**
Se tudo ocorrer bem, você verá mensagens no terminal confirmando que as migrações foram aplicadas (se configuradas), o *subscriber* de preços iniciou e o servidor está rodando na porta `8086`.

## 📂 Estrutura de Diretórios (Resumo)

* `/internal`: Coração da aplicação.
    * `/app`: Regras de negócio (`lista_service.go`).
    * `/domain`: Entidades e interfaces do domínio (`lista.go`, `lista_repository.go`).
    * `/infrastructure`:
        * `/http`: *Routers*, *handlers*, *middlewares* e *DTOs*.
        * `/messaging`: Conexão com eventos (`subscriber/prices.go`).
        * `/repository`: Operações diretas com o MySQL (`mysql_repo.go`).
* `/migrations`: Scripts de criação das tabelas no banco de dados (`init.sql`).
