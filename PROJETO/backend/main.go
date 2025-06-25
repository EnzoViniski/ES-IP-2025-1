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

type ExameDetalheCompleto struct {
	ID                   int     `json:"id"`
	IDPaciente           int     `json:"id_paciente"`
	NomePaciente         string  `json:"nome_paciente"`
	TipoExame            string  `json:"tipo_exame"`
	DataExame            string  `json:"data_exame"`
	Status               string  `json:"status"`
	Resultado            *string `json:"resultado"`
	DiagnosticoSugestivo *string `json:"diagnostico_sugestivo"`
}

type ResultadoRequest struct {
	ID                   int    `json:"id"`
	Status               string `json:"status"`
	Resultado            string `json:"resultado"`
	DiagnosticoSugestivo string `json:"diagnostico_sugestivo"`
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

type ExameCompletoRequest struct {
	IDPaciente              int    `json:"id_paciente"`
	IDUsuarioMarcou         int    `json:"id_usuario_marcou"` // Incluído para corresponder à sua tabela
	TipoExame               string `json:"tipo_exame"`
	LocalExame              string `json:"local_exame"`
	DataExame               string `json:"data_exame"`
	HorarioExame            string `json:"horario_exame"`
	Status                  string `json:"status"`
	Observacoes             string `json:"observacoes"`
	MotivoExame             string `json:"motivo_exame"`
	FezPreventivo           string `json:"fez_preventivo"`
	UsaDiu                  string `json:"usa_diu"`
	EstaGravida             string `json:"esta_gravida"`
	UsaPilula               string `json:"usa_pilula"`
	TrataMenopausa          string `json:"trata_menopausa"`
	FezRadioterapia         string `json:"fez_radioterapia"`
	DataUltimaMenstruacao   string `json:"data_ultima_menstruacao"`
	SangramentoPosRelacao   string `json:"sangramento_pos_relacao"`
	SangramentoPosMenopausa string `json:"sangramento_pos_menopausa"`
}

type PacienteDetalhesResponse struct {
	Paciente PacienteCompleto `json:"paciente"`
	Exames   []ExameDetalhe   `json:"exames"`
}

// --- MAIN ---
func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/login/profissional", enableCORS(loginHandler))
	http.HandleFunc("/profissionais/registrar", enableCORS(registrarProfissionalHandler))
	http.HandleFunc("/profissionais/atualizar/nome", enableCORS(atualizarNomeHandler))
	http.HandleFunc("/profissionais/atualizar/crm", enableCORS(atualizarCrmHandler))
	http.HandleFunc("/pacientes/cadastrar", enableCORS(cadastrarPacienteHandler))
	http.HandleFunc("/pacientes/buscar", enableCORS(buscarPacientesHandler))
	http.HandleFunc("/pacientes/detalhes", enableCORS(pacienteDetalhesHandler))
	http.HandleFunc("/pacientes/editar", enableCORS(editarPacienteHandler))
	http.HandleFunc("/exames/por-mes", enableCORS(examesPorMesHandler))
	http.HandleFunc("/exames/agendar", enableCORS(agendarExameHandler))
	http.HandleFunc("/exames/detalhes", enableCORS(exameDetalhesHandler))
	http.HandleFunc("/exames/lancar-resultado", enableCORS(lancarResultadoHandler))

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

func exameDetalhesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	exameID := r.URL.Query().Get("id")
	if exameID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID do exame não fornecido"})
		return
	}

	var exame ExameDetalheCompleto
	query := `
		SELECT e.id, e.id_paciente, p.nome_completo, e.tipo_exame, e.data_exame, e.status, e.resultado, e.diagnostico_sugestivo 
		FROM exame e
		JOIN paciente p ON e.id_paciente = p.id
		WHERE e.id = $1`

	var dataExame time.Time
	err := db.QueryRow(query, exameID).Scan(
		&exame.ID, &exame.IDPaciente, &exame.NomePaciente, &exame.TipoExame, &dataExame,
		&exame.Status, &exame.Resultado, &exame.DiagnosticoSugestivo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Exame não encontrado"})
		} else {
			log.Println("Erro ao buscar detalhes do exame:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno"})
		}
		return
	}
	exame.DataExame = dataExame.Format("2006-01-02")
	json.NewEncoder(w).Encode(exame)
}

// Handler para salvar o resultado de um exame
func lancarResultadoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ResultadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "JSON inválido"})
		return
	}

	query := `
		UPDATE exame 
		SET status = $1, resultado = $2, diagnostico_sugestivo = $3
		WHERE id = $4`

	_, err := db.Exec(query, req.Status, req.Resultado, req.DiagnosticoSugestivo, req.ID)

	if err != nil {
		log.Println("Erro ao atualizar resultado do exame:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro ao salvar resultado no banco de dados"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Resultado salvo com sucesso"})
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
func agendarExameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Método não permitido"})
		return
	}

	var req ExameCompletoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "JSON inválido ou mal formatado"})
		return
	}

	// Validação mínima
	if req.IDPaciente == 0 || req.IDUsuarioMarcou == 0 || req.DataExame == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID da paciente, ID do usuário e data do exame são obrigatórios."})
		return
	}

	// Permite que o horário seja nulo/vazio, mas a sua tabela exige (NOT NULL), então o frontend deve garantir o envio.
	var horarioExame sql.NullString
	if req.HorarioExame != "" {
		horarioExame.String = req.HorarioExame
		horarioExame.Valid = true
	} else {
		// Se a sua coluna horario_exame for NOT NULL, este braço do else pode ser um problema.
		// O ideal é garantir que o frontend sempre envie um valor.
		horarioExame.Valid = false
	}

	if req.Status == "" {
		req.Status = "Agendado"
	}

	// Query COMPLETA para inserir todos os campos da anamnese
	query := `
        INSERT INTO exame (
            id_paciente, id_usuario_marcou, tipo_exame, local_exame, data_exame, horario_exame, status, observacoes,
            motivo_exame, fez_preventivo, usa_diu, esta_gravida, usa_pilula, trata_menopausa,
            fez_radioterapia, data_ultima_menstruacao, sangramento_pos_relacao, sangramento_pos_menopausa
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`

	_, err := db.Exec(query,
		req.IDPaciente, req.IDUsuarioMarcou, req.TipoExame, req.LocalExame, req.DataExame, horarioExame, req.Status, req.Observacoes,
		req.MotivoExame, req.FezPreventivo, req.UsaDiu, req.EstaGravida, req.UsaPilula, req.TrataMenopausa,
		req.FezRadioterapia, req.DataUltimaMenstruacao, req.SangramentoPosRelacao, req.SangramentoPosMenopausa,
	)

	if err != nil {
		log.Println("Erro ao inserir novo exame completo:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro interno ao salvar o exame. Verifique o log do servidor."})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Exame agendado e anamnese salva com sucesso!"})
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
