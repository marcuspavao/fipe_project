document.addEventListener('DOMContentLoaded', () => {
    const tabelaSelect1 = document.getElementById('tabelaSelect');  
    const tabelaSelect2 = document.getElementById('tabelaSelect2');  
    const marcaSelect = document.getElementById('marcaSelect');
    const modeloSelect = document.getElementById('modeloSelect');
    const compareButton = document.getElementById('compareButton');
    const resultadoComparacao = document.getElementById('resultadoComparacao');
    let refPeriodo1 = '', refPeriodo2 = '';
    let tabelasData = [];

    fetch('/api/tabelas')
        .then(response => response.json())
        .then(tabelas => {
            tabelasData = tabelas; 

            tabelas.forEach(tabela => {
                const option1 = document.createElement('option');
                option1.value = tabela.codigo; // Campo "codigo" da TabelaReferencia
                option1.textContent = `Tabela ${formatarPeriodoParaExibicao(tabela.mes)}`
                tabelaSelect1.appendChild(option1);

                const option2 = document.createElement('option');
                option2.value = tabela.codigo;
                option2.textContent = `Tabela ${formatarPeriodoParaExibicao(tabela.mes)}`
                tabelaSelect2.appendChild(option2);
            });
        })
        .catch(err => console.error('Erro ao carregar tabelas:', err));

    tabelaSelect1.addEventListener('change', () => {
        const selectedTabela = tabelasData.find(t => t.codigo == tabelaSelect1.value);
        if (selectedTabela) refPeriodo1 = `${formatarPeriodoParaExibicao(selectedTabela.mes)}`;
        carregarMarcas();
    });

    tabelaSelect2.addEventListener('change', () => {
        const selectedTabela = tabelasData.find(t => t.codigo == tabelaSelect2.value);
        if (selectedTabela) refPeriodo2 = `${formatarPeriodoParaExibicao(selectedTabela.mes)}`;
        carregarMarcas();
    });

    function carregarMarcas() {
        const tabelaVal1 = tabelaSelect1.value;
        const tabelaVal2 = tabelaSelect2.value;

        if (tabelaVal1 && tabelaVal2) {
            marcaSelect.disabled = false;
            marcaSelect.innerHTML = '<option value="">Selecione uma marca</option>';
            modeloSelect.innerHTML = '<option value="">Selecione um modelo</option>';
            modeloSelect.disabled = true;
            resultadoComparacao.innerHTML = '';

            fetch(`/api/marcas?tabela=${tabelaVal1}`)
                .then(response => response.json())
                .then(marcas => {
                    marcas.forEach(marca => {
                        const option = document.createElement('option');
                        option.value = marca.brandCode;
                        option.textContent = marca.brandName;
                        marcaSelect.appendChild(option);
                    });
                })
                .catch(err => console.error('Erro ao carregar marcas:', err));
        } else {
            marcaSelect.disabled = true;
            modeloSelect.disabled = true;
            marcaSelect.innerHTML = '<option value="">Selecione uma marca</option>';
            modeloSelect.innerHTML = '<option value="">Selecione um modelo</option>';
            resultadoComparacao.innerHTML = '';
        }
    }

    marcaSelect.addEventListener('change', () => {
        const marcaVal = marcaSelect.value;
        const tabelaVal1 = tabelaSelect1.value;

        if (marcaVal && tabelaVal1) {
            modeloSelect.disabled = false;
            modeloSelect.innerHTML = '<option value="">Selecione um modelo</option>';
            resultadoComparacao.innerHTML = '';

            fetch(`/api/modelos/${marcaVal}?tabela=${tabelaVal1}`)
                .then(response => response.json())
                .then(modelos => {
                    modelos.forEach(modelo => {
                        const option = document.createElement('option');
                        option.value = modelo.modelCode;
                        option.textContent = modelo.modelName;
                        modeloSelect.appendChild(option);
                    });
                })
                .catch(err => console.error('Erro ao carregar modelos:', err));
        } else {
            modeloSelect.disabled = true;
            modeloSelect.innerHTML = '<option value="">Selecione um modelo</option>';
            resultadoComparacao.innerHTML = '';
        }
    });

    compareButton.addEventListener('click', () => {
        const modeloVal = modeloSelect.value;
        const tabelaVal1 = tabelaSelect1.value;
        const tabelaVal2 = tabelaSelect2.value;

        if (modeloVal && tabelaVal1 && tabelaVal2) {
            resultadoComparacao.innerHTML = '<p>Carregando comparação...</p>';

            Promise.all([
                fetch(`/api/veiculos?modelo=${modeloVal}&tabela=${tabelaVal1}`).then(res => res.json()),
                fetch(`/api/veiculos?modelo=${modeloVal}&tabela=${tabelaVal2}`).then(res => res.json()),
            ])
                .then(([dadosTabela1, dadosTabela2]) => {
                    if (dadosTabela1.length > 0 && dadosTabela2.length > 0) {
                        let html = '<h2>Comparação de Valores</h2>';
                        html += `<h3>${modeloSelect.options[modeloSelect.selectedIndex].text}</h3>`;

                        dadosTabela1.forEach((veiculo2025, index) => {
                            const veiculo2024 = dadosTabela2[index];
                            if (veiculo2024) {
                                html += `
                                    <div class="vehicle-card">
                                        <p><strong>Ano do Carro:</strong> ${veiculo2025.year === 32000 ? '0km' : veiculo2025.year}</p>
                                        <p><strong>${refPeriodo1}:</strong> ${veiculo2025.price.replace(/"/g, '')}</p>
                                        <p><strong>${refPeriodo2}:</strong> ${veiculo2024.price.replace(/"/g, '')}</p>
                                    </div>`;
                            }
                        });

                        resultadoComparacao.innerHTML = html;
                    } else {
                        resultadoComparacao.innerHTML =
                            '<p>Nenhuma comparação disponível para os períodos selecionados.</p>';
                    }
                })
                .catch(err => {
                    console.error('Erro ao comparar valores:', err);
                    resultadoComparacao.innerHTML =
                        '<p>Erro ao buscar os dados para comparação.</p>';
                });
        } else {
            resultadoComparacao.innerHTML =
                '<p>Por favor, selecione os períodos e o modelo para comparação.</p>';
        }
    });
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