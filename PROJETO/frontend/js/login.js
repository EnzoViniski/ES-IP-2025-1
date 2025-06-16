

document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const tipoUsuarioLogin = document.getElementById('tipoUsuarioLogin');
    const campoUsuarioLogin = document.getElementById('campoUsuarioLogin');
    const senhaLogin = document.getElementById('senhaLogin');
    const loginErrorMessage = document.getElementById('login-error-message');

    tipoUsuarioLogin.addEventListener('change', atualizarInterfaceLogin);
    document.getElementById('toggleSenhaLogin').addEventListener('click', () => togglePassword('senhaLogin'));
    loginForm.addEventListener('submit', handleLogin);

    async function handleLogin(event) {
        event.preventDefault();
        loginErrorMessage.textContent = '';
        const tipoUsuario = tipoUsuarioLogin.value;
        const identificador = campoUsuarioLogin.value;
        const senha = senhaLogin.value;

        if (tipoUsuario === 'profissional') {
            try {
                const response = await fetch('http://localhost:8080/login/profissional', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email: identificador, senha: senha })
                });
                const resultado = await response.json();
                if (response.ok) {
                    alert(`Login bem-sucedido! Bem-vindo(a), ${resultado.nome}!`);
                    sessionStorage.setItem('usuarioLogado', JSON.stringify(resultado));
                    window.location.href = 'homeprofissional.html';
                } else {
                    loginErrorMessage.textContent = resultado.erro || 'Credenciais inválidas.';
                }
            } catch (error) {
                loginErrorMessage.textContent = 'Erro de conexão. O servidor está rodando?';
            }
        } else {
            alert('Login de paciente ainda não implementado.');
        }
    }

    function togglePassword(id) {
        const senhaInput = document.getElementById(id);
        senhaInput.type = senhaInput.type === "password" ? "text" : "password";
    }

    function atualizarInterfaceLogin() {
        const tipo = tipoUsuarioLogin.value;
        if (tipo === "profissional") {
            loginForm.classList.remove("verde");
            loginForm.classList.add("azul");
            campoUsuarioLogin.placeholder = "Informe seu e-mail ou CRM:";
        } else {
            loginForm.classList.remove("azul");
            loginForm.classList.add("verde");
            campoUsuarioLogin.placeholder = "Informe seu e-mail ou CPF:";
        }
    }

    atualizarInterfaceLogin();
});
