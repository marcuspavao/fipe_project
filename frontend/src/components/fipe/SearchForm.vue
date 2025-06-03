<template>
  <div class="search-container">
    <div class="search-card">
      <div class="form-group">
        <label for="tabela-select"><i class="fas fa-calendar-alt"></i> Ano de Referência:</label>
        <select id="tabela-select" class="form-control" v-model="selectedTabela" @change="onTabelaChange" :disabled="loadingTabelas">
          <option value="">Selecione a tabela</option>
          <option v-for="tabela in tabelas" :key="tabela.codigo" :value="tabela.codigo">
            {{ formatarPeriodoParaExibicao(tabela.mes) }}
          </option>
        </select>
      </div>
      <div class="form-group">
        <label for="marca-select"><i class="fas fa-trademark"></i> Marca:</label>
        <select id="marca-select" class="form-control" v-model="selectedMarca" @change="onMarcaChange" :disabled="!selectedTabela || loadingMarcas">
          <option value="">Selecione uma marca</option>
          <option v-for="marca in marcas" :key="marca.brandCode" :value="marca.brandCode">
            {{ marca.brandName }}
          </option>
        </select>
      </div>
      <div class="form-group">
        <label for="modelo-select"><i class="fas fa-car-side"></i> Modelo:</label>
        <select id="modelo-select" class="form-control" v-model="selectedModelo" @change="onModeloChange" :disabled="!selectedMarca || loadingModelos">
          <option value="">Selecione um modelo</option>
          <option v-for="modelo in modelos" :key="modelo.modelCode" :value="modelo.modelCode">
            {{ modelo.modelName }}
          </option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { getTabelas, getMarcas, getModelos, formatarPeriodoParaExibicao } from '../../services/apiService';

const tabelas = ref([]);
const marcas = ref([]);
const modelos = ref([]);

const selectedTabela = ref('');
const selectedMarca = ref('');
const selectedModelo = ref('');

const loadingTabelas = ref(false);
const loadingMarcas = ref(false);
const loadingModelos = ref(false);

const emit = defineEmits(['modeloSelected', 'loading', 'error']);

async function fetchTabelas() {
  loadingTabelas.value = true;
  emit('loading', true);
  emit('error', ''); 
  try {
    tabelas.value = await getTabelas();
  } catch (err) {
    console.error(err);
    emit('error', err.message || 'Erro ao carregar tabelas de referência.');
  } finally {
    loadingTabelas.value = false;
    emit('loading', false);
  }
}

onMounted(fetchTabelas);

async function onTabelaChange() {
  selectedMarca.value = '';
  selectedModelo.value = '';
  marcas.value = [];
  modelos.value = [];
  emit('modeloSelected', null); // Clear previous results

  if (selectedTabela.value) {
    loadingMarcas.value = true;
    emit('loading', true);
    emit('error', '');
    try {
      marcas.value = await getMarcas(selectedTabela.value);
    } catch (err) {
      console.error(err);
      emit('error', err.message || 'Erro ao carregar marcas.');
    } finally {
      loadingMarcas.value = false;
      emit('loading', false);
    }
  }
}

async function onMarcaChange() {
  selectedModelo.value = '';
  modelos.value = [];
  emit('modeloSelected', null); // Clear previous results

  if (selectedMarca.value) {
    loadingModelos.value = true;
    emit('loading', true);
    emit('error', '');
    try {
      modelos.value = await getModelos(selectedTabela.value, selectedMarca.value);
    } catch (err) {
      console.error(err);
      emit('error', err.message || 'Erro ao carregar modelos.');
    } finally {
      loadingModelos.value = false;
      emit('loading', false);
    }
  }
}

function onModeloChange() {
  if (selectedModelo.value) {
    const modeloObj = modelos.value.find(m => m.modelCode === selectedModelo.value);
    const modeloNome = modeloObj ? modeloObj.modelName : '';
    emit('modeloSelected', { 
      tabela: selectedTabela.value, 
      modelo: selectedModelo.value,
      modeloNome: modeloNome 
    });
  } else {
    emit('modeloSelected', null); // Clear previous results
  }
}
</script>

<style scoped>
.search-container {
  margin-bottom: 30px;
}

.search-card {
  background-color: #fff; /* Assuming var(--white) is #fff */
  border-radius: 8px;
  padding: 25px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); /* Assuming var(--shadow) */
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #2980b9; /* Assuming var(--secondary-color) */
}

.form-group i {
  margin-right: 8px;
  color: #3498db; /* Assuming var(--primary-color) */
}

.form-control {
  width: 100%; /* Changed from 50% to be more responsive by default */
  padding: 12px;
  border: 1px solid #e0e0e0; /* Assuming var(--medium-gray) */
  border-radius: 4px;
  font-size: 16px;
  transition: all 0.3s ease; /* Assuming var(--transition) */
}

.form-control:focus {
  outline: none;
  border-color: #3498db; /* Assuming var(--primary-color) */
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.25);
}

.form-control:disabled {
  background-color: #f5f5f5; /* Assuming var(--light-gray) */
  cursor: not-allowed;
}
</style>
