document.addEventListener('DOMContentLoaded', () => {
    const tabelaSelect = document.getElementById('tabelaSelect');
    const marcaSelect = document.getElementById('marcaSelect');
    const modeloSelect = document.getElementById('modeloSelect');
    const resultado = document.getElementById('resultado');
  
    fetch('/api/tabelas')
      .then(response => response.json())
      .then(tabelas => {
        tabelas.forEach(tabela => {
          const option = document.createElement('option');
          option.value = tabela.codigo; 
          option.textContent = `Tabela ${tabela.codigo} - ${tabela.mes}`;
          tabelaSelect.appendChild(option);
        });
      })
      .catch(err => console.error('Erro ao carregar tabelas:', err));
  
    tabelaSelect.addEventListener('change', () => {
      const tabelaVal = tabelaSelect.value;
      if (tabelaVal) {
        marcaSelect.disabled = false;
        marcaSelect.innerHTML = '<option value="">Selecione uma marca</option>';
        modeloSelect.innerHTML = '<option value="">Selecione um modelo</option>';
        modeloSelect.disabled = true;
        resultado.innerHTML = '';
        fetch(`/api/marcas?tabela=${tabelaVal}`)
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
        resultado.innerHTML = '';
      }
    });
  
    marcaSelect.addEventListener('change', () => {
      const marcaVal = marcaSelect.value;
      const tabelaVal = tabelaSelect.value;
      if (marcaVal && tabelaVal) {
        modeloSelect.disabled = false;
        modeloSelect.innerHTML = '<option value="">Selecione um modelo</option>';
        resultado.innerHTML = '';
        fetch(`/api/modelos/${marcaVal}?tabela=${tabelaVal}`)
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
        resultado.innerHTML = '';
      }
    });
  
    modeloSelect.addEventListener('change', () => {
      const modeloVal = modeloSelect.value;
      const tabelaVal = tabelaSelect.value;
      if (modeloVal && tabelaVal) {
        resultado.innerHTML = '<p>Carregando informações do veículo...</p>';
        console.log(modeloVal)
        console.log(tabelaVal)
        fetch(`/api/veiculos?modelo=${modeloVal}&tabela=${tabelaVal}`)
        .then(response => response.json())
        .then(veiculos => {
          if (Array.isArray(veiculos) && veiculos.length > 0) {
            let html = '<h2>Resultados da Consulta</h2>';
            html += `<h3>${modeloSelect.options[modeloSelect.selectedIndex].text}</h3>`;
            
            veiculos.forEach(v => {
              html += `
                <div class="vehicle-card">
                  <span class="year">${v.year}</span>
                  <p class="price">${v.price.replace(/"/g, '')}</p>
                  <p>Referência: ${v.monthReference.replace(/"/g, '')}</p>
                </div>
              `;
            });
            
            resultado.innerHTML = html;
          } else {
            resultado.innerHTML = '<div class="vehicle-card"><p>Nenhum veículo encontrado.</p></div>';
          }
        })
        .catch(err => {
          console.error('Erro ao carregar veículos:', err);
          resultado.innerHTML = '<div class="vehicle-card"><p>Erro ao buscar os dados.</p></div>';
        });
      } else {
        resultado.innerHTML = '';
      }
    });

  });
  