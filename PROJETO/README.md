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
git clone [https://github.com/seu-usuario/nome-do-repositorio.git](https://github.com/seu-usuario/nome-do-repositorio.git)
cd nome-do-repositorio
