<template>
  <div>
    <SearchForm @modeloSelected="handleModeloSelected" @loading="handleLoading" @error="handleError" />
    <ResultsDisplay :veiculos="veiculos" :modeloNome="modeloNome" :loading="loading" :error="error" />
  </div>
</template>

<script setup>
import { ref } from 'vue';
import SearchForm from '../components/fipe/SearchForm.vue';
import ResultsDisplay from '../components/fipe/ResultsDisplay.vue';
import { getVeiculos } from '../services/apiService';

const veiculos = ref([]);
const modeloNome = ref('');
const loading = ref(false);
const error = ref('');

async function handleModeloSelected(payload) {
  if (!payload) {
    veiculos.value = [];
    modeloNome.value = '';
    error.value = '';
    loading.value = false; // Corrected from loadingResults.value
    return;
  }

  const { tabela, modelo, modeloNome: mNome } = payload;
  
  loading.value = true;
  error.value = '';
  veiculos.value = [];
  modeloNome.value = mNome;

  try {
    const result = await getVeiculos(tabela, modelo);
    veiculos.value = result;
  } catch (err) {
    console.error(err);
    error.value = err.message || 'Erro ao buscar veículos.';
    veiculos.value = []; // Clear any partial results
  } finally {
    loading.value = false;
  }
}

function handleLoading(isLoading) {
  // This event comes from SearchForm. If it's loading, it means selections are changing.
  if (isLoading) {
    veiculos.value = [];
    modeloNome.value = '';
    // error.value = ''; // Don't clear error if form is just re-loading dependencies
    loading.value = true; // Show loading indicator in results as well
  } else {
    // If the SearchForm is done loading its own dependencies (marcas/modelos)
    // and we are not already in a result loading state, set loading to false.
    if (!loading.value) { // Avoid overriding if handleModeloSelected is active
      loading.value = false;
    }
  }
}

function handleError(errorMessage) {
  error.value = errorMessage;
  veiculos.value = [];
  modeloNome.value = ''; // Clear model name on error too
  loading.value = false; // Ensure loading is stopped
}
</script>

<style scoped>
/* No specific styles for FipeConsultaView for now */
</style>
