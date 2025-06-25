package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// --- CONFIGURAÇÕES DO BANCO DE DADOS ---
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "Projeto"
)

var db *sql.DB

// --- ESTRUTURAS ---
type Usuario struct {
	ID        int     `json:"id"`
	Nome      string  `json:"nome"`
	Email     string  `json:"email"`
	CRM       *string `json:"crm"`
	SenhaHash string  `json:"-"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type PacienteCompleto struct {
	ID             int     `json:"id"`
	NomeCompleto   string  `json:"nome_completo"`
	Email          *string `json:"email"`
	Telefone       *string `json:"telefone"`
	CPF            *string `json:"cpf"`
	RG             *string `json:"rg"`
	DataNascimento *string `json:"data_nascimento"`
	CartaoSUS      *string `json:"cartao_sus"`
	CEP            *string `json:"cep"`
	Logradouro     *string `json:"logradouro"`
	Numero         *string `json:"numero"`
	Complemento    *string `json:"complemento"`
	UF             *string `json:"uf"`
	Apelido        *string `json:"apelido"`
	Raca           *string `json:"raca"`
	Nacionalidade  *string `json:"nacionalidade"`
	Escolaridade   *string `json:"escolaridade"`
}

type ExameDetalhe struct {
	ID        int     `json:"id"`
	TipoExame string  `json:"tipo_exame"`
	DataExame string  `json:"data_exame"`
	Status    string  `json:"status"`
	Local     *string `json:"local_exame"`
}

type PacienteDetalhesResponse struct {
	Paciente PacienteCompleto `json:"paciente"`
	Exames   []ExameDetalhe   `json:"exames"`
}
type Paciente struct { // Nova struct para dados básicos do paciente para login
	ID        int    `json:"id"`
	Nome      string `json:"nome"` // Nome Completo do Paciente
	Email     string `json:"email"`
	CPF       string `json:"cpf"`
	SenhaHash string `json:"-"`
}

// Nova struct para a requisição de login do paciente
type LoginPacienteRequest struct {
	Identificador string `json:"identificador"` // Pode ser e-mail ou CPF
	Senha         string `json:"senha"`
}
// --- MAIN ---
func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/login/profissional", enableCORS(loginHandler))
	http.HandleFunc("/login/paciente", enableCORS(loginPacienteHandler))
	http.HandleFunc("/profissionais/registrar", enableCORS(registrarProfissionalHandler))
	http.HandleFunc("/profissionais/atualizar/nome", enableCORS(atualizarNomeHandler))
	http.HandleFunc("/profissionais/atualizar/crm", enableCORS(atualizarCrmHandler))
	http.HandleFunc("/pacientes/cadastrar", enableCORS(cadastrarPacienteHandler))
	http.HandleFunc("/pacientes/buscar", enableCORS(buscarPacientesHandler))
	http.HandleFunc("/pacientes/detalhes", enableCORS(pacienteDetalhesHandler))
	http.HandleFunc("/pacientes/editar", enableCORS(editarPacienteHandler))
	http.HandleFunc("/exames/por-mes", enableCORS(examesPorMesHandler))

	fmt.Println("✅ Servidor Backend a ser executado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- HANDLERS (PROFISSIONAL) ---

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Pedido inválido"})
		return
	}
	var u Usuario
	query := "SELECT id, nome, email, crm, senha_hash FROM usuario WHERE email = $1"
	err := db.QueryRow(query, req.Email).Scan(&u.ID, &u.Nome, &u.Email, &u.CRM, &u.SenhaHash)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"erro": "E-mail ou palavra-passe inválidos"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.SenhaHash), []byte(req.Senha))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"erro": "E-mail ou palavra-passe inválidos"})
		return
	}
	json.NewEncoder(w).Encode(u)
}

func registrarProfissionalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}

	var req struct {
		Nome  string `json:"nome"`
		CRM   string `json:"crm"`
		Email string `json:"email"`
		Senha string `json:"senha"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "JSON inválido"})
		return
	}
	if req.Email == "" || req.Senha == "" || req.Nome == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Nome, e-mail e palavra-passe são obrigatórios"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro ao gerar hash de palavra-passe:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno de segurança"})
		return
	}

	query := "INSERT INTO usuario (nome, crm, email, senha_hash) VALUES ($1, $2, $3, $4)"
	_, err = db.Exec(query, req.Nome, req.CRM, req.Email, string(hash))

	if err != nil {
		log.Println("Erro ao inserir novo utilizador:", err)
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Este e-mail já está registado"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Erro ao criar utilizador"})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Utilizador criado com sucesso"})
}

func atualizarNomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}

	var req struct {
		ID   int    `json:"id"`
		Nome string `json:"nome"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "JSON inválido"})
		return
	}

	_, err := db.Exec("UPDATE usuario SET nome = $1 WHERE id = $2", req.Nome, req.ID)
	if err != nil {
		log.Println("Erro ao atualizar nome do profissional:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Não foi possível atualizar o nome"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Nome atualizado com sucesso"})
}

func atualizarCrmHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}

	var req struct {
		ID  int    `json:"id"`
		CRM string `json:"crm"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "JSON inválido"})
		return
	}

	_, err := db.Exec("UPDATE usuario SET crm = $1 WHERE id = $2", req.CRM, req.ID)
	if err != nil {
		log.Println("Erro ao atualizar CRM do profissional:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Não foi possível atualizar o CRM"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "CRM atualizado com sucesso"})
}
// NOVO HANDLER PARA LOGIN DE PACIENTE
func loginPacienteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}
	var req LoginPacienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Pedido inválido"})
		return
	}

	var p Paciente
	query := "SELECT id, nome_completo, email, cpf, senha_hash FROM paciente WHERE email = $1 OR cpf = $2"
	err := db.QueryRow(query, req.Identificador, req.Identificador).Scan(&p.ID, &p.Nome, &p.Email, &p.CPF, &p.SenhaHash)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"erro": "E-mail/CPF ou palavra-passe inválidos"})
		} else {
			log.Println("Erro ao buscar paciente para login:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno no servidor"})
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(p.SenhaHash), []byte(req.Senha))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"erro": "E-mail/CPF ou palavra-passe inválidos"})
		return
	}

	// Retorna dados do paciente sem o hash da senha
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    p.ID,
		"nome":  p.Nome,
		"email": p.Email,
		"cpf":   p.CPF,
	})
}

// --- HANDLERS (PACIENTE E EXAMES) ---
func cadastrarPacienteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}

	var req PacienteCompleto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "JSON inválido ou mal formatado"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Erro ao iniciar a transação:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno no servidor"})
		return
	}
	defer tx.Rollback()

	query := `
        INSERT INTO paciente (
            nome_completo, email, telefone, cpf, rg, data_nascimento, cartao_sus, 
            cep, logradouro, numero, complemento, uf,
            apelido, raca, nacionalidade, escolaridade
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	var dataNascimentoSQL interface{}
	if req.DataNascimento != nil && *req.DataNascimento != "" {
		dataNascimentoSQL = *req.DataNascimento
	} else {
		dataNascimentoSQL = nil
	}

	_, err = tx.Exec(query,
		req.NomeCompleto, req.Email, req.Telefone, req.CPF, req.RG, dataNascimentoSQL,
		req.CartaoSUS, req.CEP, req.Logradouro, req.Numero, req.Complemento, req.UF,
		req.Apelido, req.Raca, req.Nacionalidade, req.Escolaridade)

	if err != nil {
		log.Println("Erro ao inserir paciente na base de dados:", err)
		if strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"erro": "CPF já registado no sistema."})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno ao guardar paciente."})
		}
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Erro ao confirmar a transação:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno ao finalizar o salvamento."})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Paciente registada com sucesso!"})
}

func buscarPacientesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}
	termo := r.URL.Query().Get("termo")
	if termo == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "O termo de pesquisa não pode estar vazio"})
		return
	}
	query := `SELECT id, nome_completo, cpf FROM paciente WHERE nome_completo ILIKE $1 OR cpf = $2`
	rows, err := db.Query(query, "%"+termo+"%", termo)
	if err != nil {
		log.Println("Erro na consulta de pesquisa de pacientes:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno no servidor"})
		return
	}
	defer rows.Close()
	type PacienteBusca struct {
		ID           int    `json:"id"`
		NomeCompleto string `json:"nome_completo"`
		CPF          string `json:"cpf"`
	}
	var pacientes []PacienteBusca
	for rows.Next() {
		var p PacienteBusca
		var cpf sql.NullString
		if err := rows.Scan(&p.ID, &p.NomeCompleto, &cpf); err != nil {
			log.Println("Erro ao analisar a pesquisa de pacientes:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Erro ao processar os resultados"})
			return
		}
		if cpf.Valid {
			p.CPF = cpf.String
		} else {
			p.CPF = "Não informado"
		}
		pacientes = append(pacientes, p)
	}
	json.NewEncoder(w).Encode(pacientes)
}

func pacienteDetalhesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}
	pacienteID := r.URL.Query().Get("id")
	if pacienteID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID da paciente não fornecido"})
		return
	}

	var response PacienteDetalhesResponse
	queryPaciente := `
        SELECT id, nome_completo, email, telefone, cpf, rg, data_nascimento, cartao_sus,
               cep, logradouro, numero, complemento, uf, apelido, raca, nacionalidade, escolaridade
        FROM paciente WHERE id = $1`

	var dataNascimento, cpf, rg, email, telefone, cartaoSus, cep, logradouro, numero, complemento, uf, apelido, raca, nacionalidade, escolaridade sql.NullString

	err := db.QueryRow(queryPaciente, pacienteID).Scan(
		&response.Paciente.ID, &response.Paciente.NomeCompleto, &email, &telefone, &cpf, &rg, &dataNascimento,
		&cartaoSus, &cep, &logradouro, &numero, &complemento, &uf, &apelido, &raca, &nacionalidade, &escolaridade,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Paciente não encontrada"})
		} else {
			log.Println("Erro ao pesquisar detalhes da paciente:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno no servidor"})
		}
		return
	}

	if email.Valid {
		response.Paciente.Email = &email.String
	}
	if telefone.Valid {
		response.Paciente.Telefone = &telefone.String
	}
	if cpf.Valid {
		response.Paciente.CPF = &cpf.String
	}
	if rg.Valid {
		response.Paciente.RG = &rg.String
	}
	if dataNascimento.Valid && dataNascimento.String != "" {
		t, err := time.Parse(time.RFC3339, dataNascimento.String)
		if err == nil {
			dnStr := t.Format("2006-01-02")
			response.Paciente.DataNascimento = &dnStr
		}
	}
	if cartaoSus.Valid {
		response.Paciente.CartaoSUS = &cartaoSus.String
	}
	if cep.Valid {
		response.Paciente.CEP = &cep.String
	}
	if logradouro.Valid {
		response.Paciente.Logradouro = &logradouro.String
	}
	if numero.Valid {
		response.Paciente.Numero = &numero.String
	}
	if complemento.Valid {
		response.Paciente.Complemento = &complemento.String
	}
	if uf.Valid {
		response.Paciente.UF = &uf.String
	}
	if apelido.Valid {
		response.Paciente.Apelido = &apelido.String
	}
	if raca.Valid {
		response.Paciente.Raca = &raca.String
	}
	if nacionalidade.Valid {
		response.Paciente.Nacionalidade = &nacionalidade.String
	}
	if escolaridade.Valid {
		response.Paciente.Escolaridade = &escolaridade.String
	}

	queryExames := `SELECT id, tipo_exame, data_exame, status, local_exame FROM exame WHERE id_paciente = $1 ORDER BY data_exame DESC`
	rows, err := db.Query(queryExames, pacienteID)
	if err != nil {
		log.Println("Erro ao pesquisar exames da paciente:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro ao pesquisar o histórico de exames"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var exame ExameDetalhe
		var dataExame time.Time
		var localExame sql.NullString
		if err := rows.Scan(&exame.ID, &exame.TipoExame, &dataExame, &exame.Status, &localExame); err != nil {
			log.Println("Erro ao analisar exames:", err)
			continue
		}
		if localExame.Valid {
			exame.Local = &localExame.String
		}
		exame.DataExame = dataExame.Format("02/01/2006")
		response.Exames = append(response.Exames, exame)
	}

	json.NewEncoder(w).Encode(response)
}

func examesPorMesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}

	anoStr := r.URL.Query().Get("ano")
	mesStr := r.URL.Query().Get("mes")

	if anoStr == "" || mesStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Ano e mês são obrigatórios"})
		return
	}

	ano, _ := strconv.Atoi(anoStr)
	mes, _ := strconv.Atoi(mesStr)

	dataInicio := time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)
	dataFim := dataInicio.AddDate(0, 1, 0)

	query := `
        SELECT e.data_exame, e.horario_exame, e.tipo_exame, e.local_exame, p.nome_completo
        FROM exame e
        JOIN paciente p ON e.id_paciente = p.id
        WHERE e.data_exame >= $1 AND e.data_exame < $2`

	rows, err := db.Query(query, dataInicio, dataFim)
	if err != nil {
		log.Println("Erro na consulta de pesquisa de exames por mês:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro ao pesquisar exames"})
		return
	}
	defer rows.Close()

	type ExameAgenda struct {
		Type     string `json:"type"`
		Patient  string `json:"patient"`
		Time     string `json:"time"`
		Location string `json:"location"`
	}

	examesPorDia := make(map[string][]ExameAgenda)

	for rows.Next() {
		var dataExame time.Time
		var horarioExame sql.NullString
		var tipoExame, nomePaciente string
		var localExame sql.NullString

		if err := rows.Scan(&dataExame, &horarioExame, &tipoExame, &localExame, &nomePaciente); err != nil {
			log.Println("Erro ao analisar exames por mês:", err)
			continue
		}

		diaStr := dataExame.Format("2006-01-02")
		exame := ExameAgenda{
			Type:    tipoExame,
			Patient: nomePaciente,
		}
		if horarioExame.Valid {
			t, err := time.Parse("15:04:05", horarioExame.String)
			if err == nil {
				exame.Time = t.Format("15:04")
			}
		} else {
			exame.Time = "N/A"
		}
		if localExame.Valid {
			exame.Location = localExame.String
		} else {
			exame.Location = "Não informado"
		}

		examesPorDia[diaStr] = append(examesPorDia[diaStr], exame)
	}
	json.NewEncoder(w).Encode(examesPorDia)
}

func editarPacienteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}

	var req PacienteCompleto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "JSON inválido"})
		return
	}

	query := `
        UPDATE paciente SET
            nome_completo = $1, email = $2, telefone = $3, cpf = $4, rg = $5, data_nascimento = $6,
            cartao_sus = $7, cep = $8, logradouro = $9, numero = $10, complemento = $11, uf = $12,
            apelido = $13, raca = $14, nacionalidade = $15, escolaridade = $16
        WHERE id = $17`

	var dataNascimentoSQL interface{}
	if req.DataNascimento != nil && *req.DataNascimento != "" {
		dataNascimentoSQL = *req.DataNascimento
	} else {
		dataNascimentoSQL = nil
	}

	_, err := db.Exec(query,
		req.NomeCompleto, req.Email, req.Telefone, req.CPF, req.RG, dataNascimentoSQL,
		req.CartaoSUS, req.CEP, req.Logradouro, req.Numero, req.Complemento, req.UF,
		req.Apelido, req.Raca, req.Nacionalidade, req.Escolaridade,
		req.ID,
	)

	if err != nil {
		log.Println("Erro ao atualizar paciente:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno ao atualizar os dados."})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Paciente atualizada com sucesso!"})
}

// --- UTILITÁRIOS ---
func initDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Erro ao abrir a ligação:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao ligar à base de dados:", err)
	}
	fmt.Println("✅ Ligado à base de dados com sucesso!")
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}
