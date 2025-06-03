<template>
  <div>
    <DashboardControls @fetchDashboard="handleFetchDashboard" @status="updateUserMessage" @error="handleErrorCondition" />
    <DashboardDisplay 
      :dashboardData="data" 
      :loading="loading" 
      :error="error"
      :statusMessage="userMessage" 
      />
      <!-- :tabelas="tabelasForDisplay" prop removed as it's not strictly needed by DashboardDisplay's current logic -->
  </div>
</template>

<script setup>
import { ref } from 'vue';
import DashboardControls from '../components/dashboard/DashboardControls.vue';
import DashboardDisplay from '../components/dashboard/DashboardDisplay.vue';
import { getDashboardData } from '../services/apiService';

const data = ref([]);
const loading = ref(false);
const error = ref('');
const userMessage = ref(''); // Unified message for DashboardDisplay (can be status, loading hint, or no data message)

// onMounted: No specific onMounted logic needed here for now. 
// DashboardControls handles its own initial data loading (tabelas, marcas).

async function handleFetchDashboard(payload) {
  const { tabela1, tabela2, marca } = payload;

  loading.value = true;
  error.value = '';
  userMessage.value = 'Buscando dados do dashboard...'; // This will be shown by DashboardDisplay via its loading prop + statusMessage
  data.value = [];

  try {
    const result = await getDashboardData(tabela1, tabela2, marca);
    data.value = result;
    if (result.length === 0) {
      userMessage.value = 'Nenhum dado encontrado para os critérios selecionados.';
    } else {
      userMessage.value = ''; // Clear message if data found; DashboardDisplay will just show data.
    }
  } catch (errCatch) {
    console.error(errCatch);
    error.value = errCatch.message || 'Erro ao buscar dados do dashboard.';
    userMessage.value = ''; // Error prop will be used by DashboardDisplay
    data.value = [];
  } finally {
    loading.value = false;
  }
}

// Called when DashboardControls emits a status (e.g., "select periods", "loading marcas")
function updateUserMessage(message) {
  userMessage.value = message;
  error.value = ''; // Clear any previous error
  data.value = []; // Clear data as controls state changed
  loading.value = false; // Ensure loading is off if it was a control-driven status
}

// Called when DashboardControls emits an error (e.g., "periods cannot be same")
function handleErrorCondition(errorMessage) {
    error.value = errorMessage;
    userMessage.value = ''; // Let DashboardDisplay show the error
    data.value = [];
    loading.value = false;
}

</script>

<style scoped>
/* No specific styles needed for DashboardMarcasView itself. */
/* Global/shared styles for .status, .error, .loading are in DashboardDisplay.vue's scoped CSS with :deep() or should be made global. */
</style>
