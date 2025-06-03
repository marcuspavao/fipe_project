<template>
  <div class="controls search-card"> <!-- Reusing search-card for similar base styling -->
    <div class="form-group">
      <label for="tabela1-dashboard-select"><i class="fas fa-calendar-alt"></i> Primeiro Período:</label>
      <select id="tabela1-dashboard-select" class="form-control" v-model="selectedTabela1" @change="onTabela1Change" :disabled="loadingTabelas">
        <option value="">Selecione o 1º Período</option>
        <option v-for="tabela in tabelas" :key="tabela.codigo" :value="tabela.codigo">
          {{ formatarPeriodoParaExibicao(tabela.mes) }}
        </option>
      </select>
    </div>

    <div class="form-group">
      <label for="tabela2-dashboard-select"><i class="fas fa-calendar-alt"></i> Segundo Período:</label>
      <select id="tabela2-dashboard-select" class="form-control" v-model="selectedTabela2" @change="onControlChange" :disabled="loadingTabelas">
        <option value="">Selecione o 2º Período</option>
        <option v-for="tabela in tabelas" :key="tabela.codigo" :value="tabela.codigo">
          {{ formatarPeriodoParaExibicao(tabela.mes) }}
        </option>
      </select>
    </div>

    <div class="form-group">
      <label for="marca-dashboard-select"><i class="fas fa-trademark"></i> Marca:</label>
      <select id="marca-dashboard-select" class="form-control" v-model="selectedMarca" @change="onControlChange" :disabled="!marcaSelectEnabled || loadingMarcas">
        <option value="all">-- Todas as Marcas --</option>
        <option v-for="marca in marcas" :key="marca.brandCode" :value="marca.brandCode">
          {{ marca.brandName }}
        </option>
      </select>
    </div>

    <button @click="onFetchDashboard" :disabled="!fetchButtonEnabled" class="fetchDashboardButton">Comparar Marcas</button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { getTabelas, getMarcas, formatarPeriodoParaExibicao } from '../../services/apiService';

const emit = defineEmits(['fetchDashboard', 'loading', 'error', 'status']);

const tabelas = ref([]);
const marcas = ref([]);

const selectedTabela1 = ref('');
const selectedTabela2 = ref('');
const selectedMarca = ref('all');

const loadingTabelas = ref(false);
const loadingMarcas = ref(false);

async function fetchInitialTabelas() {
  loadingTabelas.value = true;
  emit('status', 'Carregando períodos de referência...');
  emit('error', '');
  try {
    tabelas.value = await getTabelas();
    emit('status', 'Selecione os períodos para consulta.');
  } catch (err) {
    console.error(err);
    emit('error', 'Erro ao carregar períodos.');
    emit('status', 'Falha ao carregar dados iniciais. Tente recarregar a página.');
  } finally {
    loadingTabelas.value = false;
  }
}
onMounted(fetchInitialTabelas);

async function loadMarcas() {
  if (!selectedTabela1.value) {
    marcas.value = [];
    selectedMarca.value = 'all';
    return;
  }

  loadingMarcas.value = true;
  emit('status', `Carregando marcas para o período ${formatarPeriodoParaExibicao(tabelas.value.find(t => t.codigo === selectedTabela1.value)?.mes)}...`);
  emit('error', '');
  try {
    const fetchedMarcas = await getMarcas(selectedTabela1.value);
    marcas.value = fetchedMarcas.sort((a, b) => a.brandName.localeCompare(b.brandName));
    selectedMarca.value = 'all'; // Reset to all
    emit('status', 'Selecione os períodos e opcionalmente uma marca.');
  } catch (err) {
    console.error(err);
    marcas.value = [];
    selectedMarca.value = 'all';
    emit('error', 'Erro ao carregar marcas.');
    emit('status', 'Falha ao carregar marcas.');
  } finally {
    loadingMarcas.value = false;
  }
}

function onTabela1Change() {
  loadMarcas();
  onControlChange(); // To update button state
}

function onControlChange() {
  // This function primarily serves to trigger computed property updates if needed,
  // or to emit status messages.
  if (fetchButtonEnabled.value) {
    emit('status', 'Pronto para buscar dados do dashboard.');
  } else if (selectedTabela1.value && selectedTabela2.value && selectedTabela1.value === selectedTabela2.value) {
    emit('status', 'O primeiro e o segundo período não podem ser iguais.');
     emit('error', 'Períodos de referência não podem ser iguais.'); // Also set error for immediate feedback
  } else if (selectedTabela1.value && selectedTabela2.value) {
    emit('status', 'Selecione os períodos para consulta.');
  }
   else if (!selectedTabela1.value || !selectedTabela2.value) {
     emit('status', 'Por favor, selecione ambos os períodos de referência.');
   }
}

const marcaSelectEnabled = computed(() => !!selectedTabela1.value);
const fetchButtonEnabled = computed(() => {
  const enabled = !!(selectedTabela1.value && selectedTabela2.value && selectedTabela1.value !== selectedTabela2.value);
  if (enabled) {
      emit('error', ''); // Clear error if button becomes enabled
  }
  return enabled;
});

function onFetchDashboard() {
  if (!fetchButtonEnabled.value) {
    if (selectedTabela1.value && selectedTabela2.value && selectedTabela1.value === selectedTabela2.value) {
         emit('error', 'Os períodos de referência não podem ser iguais.');
    } else {
         emit('error', 'Por favor, selecione ambos os períodos de referência.');
    }
    return;
  }
  emit('fetchDashboard', {
    tabela1: selectedTabela1.value,
    tabela2: selectedTabela2.value,
    marca: selectedMarca.value
  });
}
</script>

<style scoped>
/* .search-card from styles.css is used as a base for the main div */
/* .form-group and .form-control from styles.css are used */

.controls { /* This class is on the root div, which also has .search-card */
  /* display: flex; from dashboard.css - search-card is not flex by default */
  /* flex-wrap: wrap; from dashboard.css */
  /* gap: 20px; from dashboard.css */
  /* align-items: flex-end; from dashboard.css */
  /* These flex properties can be added if the search-card's block layout isn't sufficient */
  /* For now, relying on search-card and standard form-group block display */
  padding-bottom: 20px; /* from dashboard.css .controls */
  border-bottom: 1px solid #dee2e6; /* var(--border-color) from dashboard.css .controls */
}

/* Ensuring form elements inside take full width as per search-card's .form-control style */
.form-control {
  width: 100%; 
}

.fetchDashboardButton {
  background-color: #3498db; /* var(--primary-color) */
  color: white;
  border: none;
  padding: 12px 24px;
  font-size: 16px;
  font-weight: bold;
  text-transform: uppercase;
  border-radius: 8px;
  cursor: pointer;
  /* margin-bottom: 25px; from dashboard.css #fetchDashboardButton - remove if not needed or adjust */
  margin-top: 10px; /* Add some top margin for spacing */
  transition: all 0.3s ease; /* var(--transition) */
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1); /* var(--shadow) */
  width: 100%; /* Make button full width */
}

.fetchDashboardButton:hover {
  background-color: #27ae60; /* Green hover, consistent with other primary buttons */
  box-shadow: 0px 6px 10px rgba(0, 0, 0, 0.15);
}

.fetchDashboardButton:disabled {
  background-color: #bdc3c7; /* Muted color from dashboard.css */
  color: #7f8c8d;
  cursor: not-allowed;
  box-shadow: none;
}
</style>
