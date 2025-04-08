document.addEventListener('DOMContentLoaded', () => {
    const tabelaSelect1 = document.getElementById('tabelaSelect');  
    const tabelaSelect2 = document.getElementById('tabelaSelect2');  
    const marcaSelect = document.getElementById('marcaSelect');
    const modeloSelect = document.getElementById('modeloSelect');
    const compareButton = document.getElementById('compareButton');
    const resultadoComparacao = document.getElementById('resultadoComparacao');

    fetch('/api/tabelas')
        .then(response => response.json())
        .then(tabelas => {
            tabelas.forEach(tabela => {
                const option1 = document.createElement('option');
                option1.value = tabela.codigo; // Campo "codigo" da TabelaReferencia
                option1.textContent = `Tabela ${tabela.codigo} - ${tabela.mes}`;
                tabelaSelect1.appendChild(option1);

                const option2 = document.createElement('option');
                option2.value = tabela.codigo;
                option2.textContent = `Tabela ${tabela.codigo} - ${tabela.mes}`;
                tabelaSelect2.appendChild(option2);
            });
        })
        .catch(err => console.error('Erro ao carregar tabelas:', err));

    tabelaSelect1.addEventListener('change', () => {
        carregarMarcas();
    });

    tabelaSelect2.addEventListener('change', () => {
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
                                        <p><strong>Ano:</strong> ${veiculo2025.year}</p>
                                        <p><strong>Período 1:</strong> ${veiculo2025.price.replace(/"/g, '')}</p>
                                        <p><strong>Período 2:</strong> ${veiculo2024.price.replace(/"/g, '')}</p>
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
