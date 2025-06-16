
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
            
            if (abaAtiva) {
                const aba = document.getElementById(abaAtiva);
                if (aba) {
                    aba.classList.add('active');
                    document.getElementById('nav-painel').classList.add('active');
                }
            }
            
            const usuarioJSON = sessionStorage.getItem('usuarioLogado');
            if (usuarioJSON) {
                const usuario = JSON.parse(usuarioJSON);
                document.getElementById('user-name-display').textContent = `Dr(a). ${usuario.nome}`;
            } else {
                window.location.href = 'index.html';
            }
        })
        .catch(error => {
            console.error('Erro ao carregar o painel:', error);
            document.getElementById('panel-placeholder').innerHTML = "<p>Erro ao carregar a interface principal. Tente novamente.</p>";
        });
}

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
