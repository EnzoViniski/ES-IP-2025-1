# Projeto Sobre Vidas: Gestão de Acompanhamento de Câncer Cervical

## 🎯 Sobre o Projeto
O "Sobre Vidas" é um sistema web completo desenhado para facilitar a gestão e o acompanhamento de pacientes na prevenção do câncer cervical. A plataforma é destinada a profissionais de saúde, permitindo um controle centralizado e eficiente de todo o ciclo de cuidado da paciente, desde o cadastro inicial até o lançamento de resultados de exames.

O principal objetivo do projeto é fornecer uma ferramenta robusta e intuitiva que auxilie na organização de dados clínicos, otimizando o tempo do profissional e garantindo um histórico detalhado e acessível para cada paciente, o que é fundamental para um diagnóstico e acompanhamento eficazes.

## ✨ Funcionalidades Principais

O sistema oferece um portal completo para o profissional de saúde, com as seguintes funcionalidades:

* **Autenticação Segura:** Sistema de login e registro para profissionais com hash de senhas para garantir a segurança dos dados.
* **Gestão de Profissionais:** Os profissionais podem se cadastrar e gerenciar suas informações de perfil, como nome e CRM.
* **Gestão Completa de Pacientes:**
    * **Cadastro Detalhado:** Formulário completo para registrar novas pacientes com dados demográficos, de contato e endereço.
    * **Busca Rápida:** Ferramenta de busca para localizar pacientes por nome ou CPF.
    * **Visualização de Histórico:** Tela de detalhes da paciente que centraliza seus dados cadastrais e o histórico completo de exames.
    * **Edição de Dados:** Possibilidade de atualizar as informações cadastrais da paciente a qualquer momento.
* **Agenda Inteligente de Exames:**
    * Uma visualização em formato de calendário que exibe todos os exames agendados por mês, permitindo ao profissional consultar sua agenda de forma rápida.
* **Agendamento de Exames com Anamnese:**
    * Formulário para agendar novos exames (ex: Papanicolau), especificando data, hora e local.
    * Registro de uma **anamnese clínica** completa no momento do agendamento, coletando informações vitais como histórico de preventivos, gravidez, uso de medicamentos e outros dados relevantes.
* **Lançamento de Resultados:**
    * Fluxo dedicado para que o profissional possa lançar o laudo descritivo e o diagnóstico sugestivo de exames realizados, atualizando o status do exame (de "Agendado" para "Concluído", por exemplo).

## 💻 Tecnologias Utilizadas

O projeto foi construído com uma arquitetura cliente-servidor, utilizando as seguintes tecnologias:

### **Backend**
* **Linguagem:** Go (Golang)
* **Servidor:** `net/http` (biblioteca padrão do Go)
* **Banco de Dados:** PostgreSQL
* **Driver do Banco de Dados:** `github.com/lib/pq`
* **Segurança:** `golang.org/x/crypto/bcrypt` para hashing de senhas

### **Frontend**
* **Estrutura:** HTML5, CSS3 e JavaScript puro (Vanilla JS)
* **Comunicação com API:** Fetch API para realizar requisições assíncronas ao backend
* **Sessão do Usuário:** `sessionStorage` para manter os dados do profissional logado no navegador.

## 🚀 Como Rodar o Projeto

Para executar o projeto em sua máquina local, siga os passos abaixo.

### **Pré-requisitos**
* [**Go**](https://go.dev/doc/install) (versão 1.18 ou superior)
* [**PostgreSQL**](https://www.postgresql.org/download/)
* [**Git**](https://git-scm.com/downloads/)

### **1. Clone o Repositório**
```bash
git clone [https://github.com/seu-usuario/nome-do-repositorio.git](https://github.com//EnzoViniski/ES-IP-2025-1/edit/main/PROJETO/.git)
cd PROJETO

2. Configure o Banco de Dados
Crie um Banco de Dados: Abra o psql ou sua ferramenta de gerenciamento de PostgreSQL e crie um banco de dados chamado Projeto.

SQL

CREATE DATABASE "Projeto";
Ajuste a Conexão (se necessário): As credenciais do banco de dados estão no arquivo main.go. Se as suas forem diferentes, altere as constantes no início do arquivo.

Go

// main.go
const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "postgres"
    dbname   = "Projeto"
)
Crie as Tabelas: Execute o script SQL abaixo em seu banco de dados Projeto para criar as tabelas necessárias (usuario, paciente e exame). Este schema foi derivado das estruturas e queries presentes no código-fonte.

SQL

-- Tabela para os profissionais de saúde
CREATE TABLE usuario (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    crm VARCHAR(50),
    senha_hash VARCHAR(255) NOT NULL,
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para as pacientes
CREATE TABLE paciente (
    id SERIAL PRIMARY KEY,
    nome_completo VARCHAR(255) NOT NULL,
    apelido VARCHAR(100),
    data_nascimento DATE,
    escolaridade VARCHAR(100),
    raca VARCHAR(50),
    nacionalidade VARCHAR(100),
    cpf VARCHAR(20) UNIQUE,
    rg VARCHAR(30),
    cartao_sus VARCHAR(30),
    email VARCHAR(255),
    telefone VARCHAR(25),
    cep VARCHAR(10),
    logradouro VARCHAR(255),
    numero VARCHAR(20),
    complemento VARCHAR(100),
    uf VARCHAR(2),
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para os exames e anamneses
CREATE TABLE exame (
    id SERIAL PRIMARY KEY,
    id_paciente INTEGER NOT NULL REFERENCES paciente(id),
    id_usuario_marcou INTEGER NOT NULL REFERENCES usuario(id),
    tipo_exame VARCHAR(255) NOT NULL,
    local_exame VARCHAR(255),
    data_exame DATE NOT NULL,
    horario_exame TIME,
    status VARCHAR(50) DEFAULT 'Agendado',
    observacoes TEXT,
    resultado TEXT,
    diagnostico_sugestivo TEXT,
    -- Campos da Anamnese
    motivo_exame TEXT,
    fez_preventivo VARCHAR(20),
    usa_diu VARCHAR(20),
    esta_gravida VARCHAR(20),
    usa_pilula VARCHAR(20),
    trata_menopausa VARCHAR(20),
    fez_radioterapia VARCHAR(20),
    data_ultima_menstruacao VARCHAR(50),
    sangramento_pos_relacao VARCHAR(20),
    sangramento_pos_menopausa VARCHAR(20),
    criado_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
3. Crie um Usuário Profissional (Opcional)
Para testar o login, você pode inserir um usuário diretamente no banco de dados. O projeto inclui um utilitário (hasher.go) para gerar um hash de senha seguro.

Gere o Hash:

Bash

go run hasher.go
Copie o hash gerado no terminal.

Insira o Usuário: Execute o comando SQL abaixo, substituindo os valores e o hash gerado.

SQL

INSERT INTO usuario (nome, email, crm, senha_hash) VALUES 
('Dr. Exemplo', 'exemplo@email.com', '12345-GO', 'COLE_SEU_HASH_GERADO_AQUI');
4. Execute o Backend (Servidor Go)
No terminal, na raiz do projeto (onde está o main.go), execute:

Bash

go run main.go
Você verá a mensagem de confirmação:
✅ Ligado à base de dados com sucesso!
✅ Servidor Backend a ser executado na porta 8080

O backend estará rodando em http://localhost:8080.

5. Execute o Frontend
Os arquivos de frontend (HTML/CSS/JS) precisam ser servidos por um servidor web simples. Você pode usar a extensão Live Server no VS Code ou um servidor Python. O backend está configurado para aceitar requisições de outras origens (CORS).

Usando Python (se tiver instalado):

Navegue até a pasta que contém os arquivos HTML (ex: a raiz do projeto ou uma pasta frontend).

Inicie o servidor:

Bash

python -m http.server 8000
Abra seu navegador e acesse http://localhost:8000.

Agora você pode acessar a aplicação, se registrar como novo profissional ou usar a conta que criou para fazer login e explorar todas as funcionalidades!

👥 Contribuidores
Este projeto foi desenvolvido com dedicação por:

Enzo Viniski de Castro

Gabriel Guimarães

Guilherme Vieira

Vitor Queiroz
