document.addEventListener('DOMContentLoaded', () => {
    // Elementos da UI
    const tabelaSelect1 = document.getElementById('tabelaSelect1');
    const tabelaSelect2 = document.getElementById('tabelaSelect2');
    const marcaSelect = document.getElementById('marcaSelect');
    const fetchButton = document.getElementById('fetchDashboardButton');
    const dashboardResultDiv = document.getElementById('dashboardResult');
    const statusMessageDiv = document.getElementById('statusMessage');

    const API_BASE_URL = ''; // Deixe vazio se for a mesma origem, ou coloque http://localhost:8080

    let tabelasData = []; // Armazena dados das tabelas {codigo, mes}
    let marcasData = [];  // Armazena dados das marcas {brandCode, brandName}

    // Função para exibir mensagens de status/erro
    function showStatus(message, isError = false) {
        statusMessageDiv.textContent = message;
        statusMessageDiv.className = isError ? 'status error' : 'status loading'; // Usa 'loading' para não erro
        statusMessageDiv.style.display = message ? 'block' : 'none';
    }

    function formatPeriodDisplay(mesAnoString) {
        const parts = mesAnoString.split('/');
        if (parts.length === 2) {
            const mes = parts[0].charAt(0).toUpperCase() + parts[0].slice(1).toLowerCase();
            const ano = parts[1];
            return `${mes} de ${ano}`;
        }
        return mesAnoString; 
    }

    async function loadTabelas() {
        try {
            const response = await fetch(`${API_BASE_URL}/api/tabelas`);
            if (!response.ok) throw new Error(`Erro ao buscar tabelas: ${response.statusText}`);
            tabelasData = await response.json();

            // Limpa opções padrão e popula os selects
            tabelaSelect1.innerHTML = '<option value="">Selecione o 1º Período</option>';
            tabelaSelect2.innerHTML = '<option value="">Selecione o 2º Período</option>';

            tabelasData.forEach(tabela => {
                const displayMesAno = formatPeriodDisplay(tabela.mes);
                const option1 = document.createElement('option');
                option1.value = tabela.codigo;
                option1.textContent = `Tabela ${formatarPeriodoParaExibicao(displayMesAno)}`;
                option1.dataset.ref = tabela.mes; // Guarda a ref original
                tabelaSelect1.appendChild(option1);

                const option2 = document.createElement('option');
                option2.value = tabela.codigo;
                option2.textContent = `Tabela ${formatarPeriodoParaExibicao(displayMesAno)}`;
                 option2.dataset.ref = tabela.mes; // Guarda a ref original
                tabelaSelect2.appendChild(option2);
            });
             enableControls(); // Habilita controles após carregar tabelas
        } catch (error) {
            console.error('Erro ao carregar tabelas:', error);
            showStatus(`Falha ao carregar períodos. ${error.message}`, true);
            tabelaSelect1.innerHTML = '<option value="">Erro ao carregar</option>';
            tabelaSelect2.innerHTML = '<option value="">Erro ao carregar</option>';
        }
    }

    // Função para carregar marcas baseado na primeira tabela selecionada
    async function loadMarcas() {
        const tabela1Value = tabelaSelect1.value;
        if (!tabela1Value) {
            marcaSelect.innerHTML = '<option value="all">Todas as Marcas</option>';
            marcaSelect.disabled = true;
            return;
        }

        marcaSelect.disabled = true; // Desabilita enquanto carrega
        marcaSelect.innerHTML = '<option value="">Carregando marcas...</option>';

        try {
            const response = await fetch(`${API_BASE_URL}/api/marcas?tabela=${tabela1Value}`);
            if (!response.ok) throw new Error(`Erro ao buscar marcas: ${response.statusText}`);
            marcasData = await response.json();

            // Popula o select de marcas
            marcaSelect.innerHTML = '<option value="all">-- Todas as Marcas --</option>'; // Opção "Todas"
            marcasData.sort((a, b) => a.brandName.localeCompare(b.brandName)); // Ordena alfabeticamente
            marcasData.forEach(marca => {
                const option = document.createElement('option');
                option.value = marca.brandCode;
                option.textContent = marca.brandName;
                marcaSelect.appendChild(option);
            });
            marcaSelect.disabled = false; // Habilita após carregar
             enableControls(); // Verifica se o botão principal pode ser habilitado

        } catch (error) {
            console.error('Erro ao carregar marcas:', error);
            showStatus(`Falha ao carregar marcas. ${error.message}`, true);
            marcaSelect.innerHTML = '<option value="all">Erro ao carregar</option>';
             marcaSelect.disabled = true;
        }
    }

     // Habilita o botão de comparar se ambos os períodos são selecionados e diferentes
     function enableControls() {
        const tabela1Valido = tabelaSelect1.value !== ""; false
        const tabela2Valido = tabelaSelect2.value !== ""; false
        const periodosDiferentes = tabelaSelect1.value !== tabelaSelect2.value; false
        const marcaSelectValido =  marcaSelect.value !== "all"; // Considera carregado se não estiver desabilitado 

        fetchButton.disabled = !(tabela1Valido && tabela2Valido && periodosDiferentes && marcaSelectValido);
        // Habilita select de marca apenas se o período 1 estiver selecionado
        if (tabela1Valido && marcaSelect.innerHTML.includes('Erro') || marcaSelect.innerHTML.includes('Carregando')) {
             // Não desabilita se já estiver carregando ou com erro,
             // mas não dispara loadMarcas de novo se só o período 2 mudou
        } else {
            marcaSelect.disabled = !tabela1Valido;
        }
    }


    // Função para formatar diferença percentual
    function formatPercentage(value) {
        if (value === null || value === undefined || isNaN(value)) {
            return '<span class="neutral">N/A</span>';
        }
        const formatted = value.toFixed(2).replace('.', ','); // Formato PT-BR
        const sign = value > 0 ? '+' : '';
        const className = value > 0 ? 'positive' : (value < 0 ? 'negative' : 'neutral');
        return `<span class="${className}">${sign}${formatted}%</span>`;
    }

    // Função para renderizar os dados do dashboard
    function displayDashboardData(data) {
        dashboardResultDiv.innerHTML = ''; // Limpa resultados anteriores

        if (!data || data.length === 0) {
            showStatus('Nenhum dado encontrado para os critérios selecionados.', false); // Não é erro, apenas sem dados
            return;
        }

        data.forEach(brandEntry => {
            const card = document.createElement('div');
            card.className = 'brand-card';

            const diffAvg = formatPercentage(brandEntry.diferencasPercentuais?.valorMedio0km);
            const diffModels = formatPercentage(brandEntry.diferencasPercentuais?.totalModelos);

            // Helper para renderizar estatísticas de um período
            const renderPeriodStats = (periodStats, periodNum) => {
                const refDisplay = formatPeriodDisplay(periodStats.ref);
                return `
                    <div class="period-stats period-${periodNum}">
                        <h3>${refDisplay}</h3>
                        <p><strong>Menor Preço 0km:</strong> <span class="price-detail">${periodStats.menorPreco0km.valorFmt || 'N/A'}</span></p>
                        ${periodStats.menorPreco0km.modelo && periodStats.menorPreco0km.modelo !== 'N/A' ? `<span class="modelo-destaque">(${periodStats.menorPreco0km.modelo})</span>` : ''}
                        <p><strong>Maior Preço 0km:</strong> <span class="price-detail">${periodStats.maiorPreco0km.valorFmt || 'N/A'}</span></p>
                        ${periodStats.maiorPreco0km.modelo && periodStats.maiorPreco0km.modelo !== 'N/A' ? `<span class="modelo-destaque">(${periodStats.maiorPreco0km.modelo})</span>` : ''}
                        <p><strong>Valor Médio 0km:</strong> <span class="price-detail">${periodStats.valorMedio0kmFmt || 'N/A'}</span></p>
                        <p><strong>Total Modelos:</strong> <span class="price-detail">${periodStats.totalModelos || 0}</span></p>
                    </div>
                `;
            };

            card.innerHTML = `
                <h2>${brandEntry.brandName} (${brandEntry.brandCode})</h2>
                ${renderPeriodStats(brandEntry.periodo1, 1)}
                ${renderPeriodStats(brandEntry.periodo2, 2)}
                <div class="diff-stats">
                    <h3>Variação (Período 1 vs Período 2)</h3>
                    <p><strong>Valor Médio 0km:</strong> <span class="diff-detail">${diffAvg}</span></p>
                    <p><strong>Total Modelos:</strong> <span class="diff-detail">${diffModels}</span></p>
                </div>
            `;
            dashboardResultDiv.appendChild(card);
        });
    }

    // Função principal para buscar e exibir o dashboard
    async function fetchAndDisplayDashboard() {
        const tabela1 = tabelaSelect1.value;
        const tabela2 = tabelaSelect2.value;
        const marca = marcaSelect.value; // Pode ser 'all' ou um brandCode

        if (!tabela1 || !tabela2) {
            showStatus('Por favor, selecione os dois períodos.', true);
            return;
        }
        if (tabela1 === tabela2) {
             showStatus('Os períodos de comparação devem ser diferentes.', true);
             return;
        }


        showStatus('Buscando dados do dashboard...', false);
        dashboardResultDiv.innerHTML = ''; // Limpa antes de buscar
        fetchButton.disabled = true; // Desabilita botão durante a busca

        let url = `${API_BASE_URL}/api/dashboard?tabela1=${tabela1}&tabela2=${tabela2}`;
        if (marca !== 'all' && marca !== '') {
            url += `&marca=${marca}`;
        }

        try {
            const response = await fetch(url);
            if (!response.ok) {
                 let errorMsg = `Erro ao buscar dashboard: ${response.statusText}`;
                 try { // Tenta pegar mensagem de erro do backend se houver
                     const errorBody = await response.json();
                     errorMsg += ` - ${errorBody.message || errorBody.error || JSON.stringify(errorBody)}`;
                 } catch(e) { /* Ignora se não conseguir parsear erro */ }
                throw new Error(errorMsg);
            }
            const data = await response.json();
            showStatus(''); // Limpa mensagem de status em caso de sucesso
            displayDashboardData(data);
        } catch (error) {
            console.error('Erro ao buscar dashboard:', error);
            showStatus(`Falha ao carregar dashboard: ${error.message}`, true);
             dashboardResultDiv.innerHTML = '<p style="text-align: center; color: var(--negative-color);">Não foi possível carregar os dados.</p>';
        } finally {
             enableControls(); // Reabilita o botão (ou mantém desabilitado se as condições mudaram)
        }
    }

    // Event Listeners
    tabelaSelect1.addEventListener('change', () => {
        loadMarcas(); // Recarrega marcas se período 1 mudar
        enableControls();
    });
    tabelaSelect2.addEventListener('change', enableControls);
    marcaSelect.addEventListener('change', enableControls); // Apenas para habilitar/desabilitar botão
    fetchButton.addEventListener('click', fetchAndDisplayDashboard);

    // Carregamento Inicial
    loadTabelas();
});

function formatarPeriodoParaExibicao(mesAnoString) {
    if (!mesAnoString || typeof mesAnoString !== 'string') {
        return 'Período inválido';
    }
    const parts = mesAnoString.split('/');
    if (parts.length === 2) {
        const mes = capitalizeFirstLetter(parts[0].trim()); // Garante capitalização
        const ano = parts[1].trim();
        // Verifica se o ano é numérico (básico)
        if (mes && /^\d{4}$/.test(ano)) {
             // Retorna o formato desejado: "Mês de Ano"
            return `${mes} de ${ano}`;
        }
    }
    // Retorna o original como fallback se o formato for inesperado
    return `Tabela ${mesAnoString}`;
  }
  
  function capitalizeFirstLetter(string) {
    if (!string) return '';
    return string.charAt(0).toUpperCase() + string.slice(1).toLowerCase();
  }