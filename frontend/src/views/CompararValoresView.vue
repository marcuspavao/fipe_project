<template>
  <div>
    <ComparisonForm @compare="handleCompare" @loading="handleFormLoading" @error="handleError" />
    <ComparisonResults 
      :comparisonData="results" 
      :modeloNome="modeloNome" 
      :refPeriodo1Display="refPeriodo1Display" 
      :refPeriodo2Display="refPeriodo2Display" 
      :loading="loading" 
      :error="error"
      :initialPrompt="initialLoad" />
  </div>
</template>

<script setup>
import { ref } from 'vue';
import ComparisonForm from '../components/comparar/ComparisonForm.vue';
import ComparisonResults from '../components/comparar/ComparisonResults.vue';
import { getVeiculos } from '../services/apiService';

const results = ref([]);
const modeloNome = ref('');
const refPeriodo1Display = ref('');
const refPeriodo2Display = ref('');
const loading = ref(false);
const error = ref('');
const initialLoad = ref(true); // To show initial prompt in Results

async function handleCompare(payload) {
  const { tabela1, tabela2, modelo, marcaNome, modeloNome: mNome, refPeriodo1Display: rP1, refPeriodo2Display: rP2 } = payload;

  loading.value = true;
  error.value = '';
  results.value = [];
  initialLoad.value = false; // Comparison has been attempted

  modeloNome.value = mNome;
  refPeriodo1Display.value = rP1;
  refPeriodo2Display.value = rP2;

  try {
    const [veiculosT1, veiculosT2] = await Promise.all([
      getVeiculos(tabela1, modelo),
      getVeiculos(tabela2, modelo)
    ]);

    // Process and merge results
    const processedResults = [];
    const veiculosT2Map = new Map(veiculosT2.map(v => [v.year, v]));

    for (const v1 of veiculosT1) {
      const v2 = veiculosT2Map.get(v1.year);
      processedResults.push({
        year: v1.year,
        price1: v1.price,
        price2: v2 ? v2.price : 'N/A', // Handle if vehicle year doesn't exist in second period
      });
      // Remove from map to handle vehicles in T2 not in T1 later if needed, though current logic prioritizes T1 years
      if (v2) veiculosT2Map.delete(v1.year); 
    }
    
    // Optional: Add vehicles only present in T2 (if business logic requires)
    // for (const v2 of veiculosT2Map.values()) {
    //   processedResults.push({
    //     year: v2.year,
    //     price1: 'N/A',
    //     price2: v2.price,
    //   });
    // }
    // Sort results by year, e.g., numerically. Year "32000" (0km) should probably come first or last.
    processedResults.sort((a, b) => {
        const yearA = a.year === "32000" ? 0 : parseInt(a.year);
        const yearB = b.year === "32000" ? 0 : parseInt(b.year);
        return yearA - yearB;
    });


    results.value = processedResults;
    if (processedResults.length === 0) {
        // error.value = "Nenhum veículo encontrado para os critérios selecionados em ambos os períodos."
        // This case is handled by ComparisonResults component's "Nenhuma comparação disponível..." message
    }

  } catch (err) {
    console.error(err);
    error.value = err.message || 'Erro ao buscar dados para comparação.';
    results.value = [];
  } finally {
    loading.value = false;
  }
}

// This loading event comes from ComparisonForm when it's fetching Marcas or Modelos.
function handleFormLoading(isLoading) {
  if (isLoading) {
    // If the form starts loading its dependencies, we are effectively resetting the state.
    results.value = [];
    error.value = '';
    initialLoad.value = true; 
    // We don't set loading.value = true here, as that's for the main comparison action.
    // The form itself will show its disabled state for selects.
  }
  // When form loading finishes, no specific action needed here unless
  // we want to explicitly stop a global page loader if one was activated.
}

function handleError(errorMessage) {
  error.value = errorMessage;
  results.value = [];
  loading.value = false;
  initialLoad.value = false; // Error occurred, so not initial state anymore
}

</script>

<style scoped>
/* No specific styles for CompararValoresView for now */
</style>
