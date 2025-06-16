// Arquivo: frontend/js/main.js
function carregarPainel(abaAtiva) {
    fetch('components/painel_template.html')
        .then(response => {
            if (!response.ok) {
                throw new Error('Não foi possível carregar o template do painel.');
            }
            return response.text();
        })
        .then(data => {
            document.getElementById('panel-placeholder').innerHTML = data;
            
            // Lógica para destacar a aba correta
            if (abaAtiva) {
                const aba = document.getElementById(abaAtiva);
                if (aba) {
                    aba.classList.add('active');
                    // Também destaca o link principal "Painel do profissional" no header
                    document.getElementById('nav-painel').classList.add('active');
                }
            }
            
            // Lógica do usuário e logout
            const usuarioJSON = sessionStorage.getItem('usuarioLogado');
            if (usuarioJSON) {
                const usuario = JSON.parse(usuarioJSON);
                document.getElementById('user-name-display').textContent = `Dr(a). ${usuario.nome}`;
            } else {
                // Se não houver sessão, redireciona para o login
                window.location.href = 'index.html';
            }
        })
        .catch(error => {
            console.error('Erro ao carregar o painel:', error);
            document.getElementById('panel-placeholder').innerHTML = "<p>Erro ao carregar a interface principal. Tente novamente.</p>";
        });
}

// O código desta função que você me passou já estava correto.
function carregarLayout(paginaAtiva) {
    fetch('components/painel_template.html')
        .then(response => response.text())
        .then(data => {
            document.getElementById('layout-placeholder').innerHTML = data;
            
            if (paginaAtiva.startsWith('tab-')) {
                document.getElementById('nav-painel').classList.add('active');
                const aba = document.getElementById(paginaAtiva);
                if (aba) aba.classList.add('active');
            } else {
                const navLink = document.getElementById(paginaAtiva);
                if(navLink) navLink.classList.add('active');
            }
            
            const usuarioJSON = sessionStorage.getItem('usuarioLogado');
            const userNameDisplay = document.getElementById('user-name-display');

            if (usuarioJSON && userNameDisplay) {
                const usuario = JSON.parse(usuarioJSON);
                userNameDisplay.textContent = `Dr(a). ${usuario.nome}`;
            } else if (!usuarioJSON) {
                window.location.href = 'index.html';
            }
        });
}
