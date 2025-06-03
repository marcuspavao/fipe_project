<template>
  <div class="result-container">
    <div v-if="loading" class="vehicle-card">
      <p>Carregando comparação...</p>
    </div>
    <div v-else-if="error" class="vehicle-card error-card">
      <p>{{ error }}</p>
    </div>
    <div v-else-if="comparisonData && comparisonData.length > 0">
      <h2>Comparação de Valores</h2>
      <h3>{{ modeloNome }}</h3>
      <div v-for="(item, index) in comparisonData" :key="index + '-' + item.year" class="vehicle-card">
        <p class="year"><strong>Ano do Carro:</strong> {{ item.year === "32000" ? '0km' : item.year }}</p>
        <p><strong>{{ refPeriodo1Display }}:</strong> {{ item.price1 ? item.price1.replace(/"/g, '') : 'N/A' }}</p>
        <p><strong>{{ refPeriodo2Display }}:</strong> {{ item.price2 ? item.price2.replace(/"/g, '') : 'N/A' }}</p>
      </div>
    </div>
    <div v-else-if="!loading && !error && initialPrompt" class="vehicle-card">
      <p>Por favor, selecione os períodos e o modelo para comparação.</p>
    </div>
    <div v-else-if="!loading && !error && !initialPrompt && (!comparisonData || comparisonData.length === 0)" class="vehicle-card">
      <p>Nenhuma comparação disponível para os períodos e modelo selecionados, ou dados insuficientes para comparar.</p>
    </div>
  </div>
</template>

<script setup>
defineProps({
  comparisonData: Array,
  modeloNome: String,
  refPeriodo1Display: String,
  refPeriodo2Display: String,
  loading: Boolean,
  error: String,
  initialPrompt: Boolean
});
</script>

<style scoped>
/* Styles from ResultsDisplay.vue will be largely reused. */
.result-container {
  background-color: #fff;
  border-radius: 8px;
  padding: 25px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  min-height: 100px;
  margin-top: 30px;
}

.result-container h2 {
  color: #2c3e50;
  margin-bottom: 5px;
}
.result-container h3 {
  color: #3498db;
  margin-bottom: 20px;
  font-weight: normal;
}

.vehicle-card {
  border-left: 4px solid #27ae60; /* Green to distinguish from single FIPE consulta */
  padding: 15px 20px;
  margin-bottom: 15px;
  background-color: #f5f5f5;
  border-radius: 0 4px 4px 0;
  transition: all 0.3s ease;
}

.vehicle-card.error-card {
  border-left-color: #e74c3c; 
  background-color: #fbecec; 
  color: #e74c3c;
}

.vehicle-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 12px rgba(0,0,0,0.1);
}

.vehicle-card p {
  margin-bottom: 8px; /* More compact spacing for comparison lines */
}
.vehicle-card p:last-child {
  margin-bottom: 0;
}

.year { /* Re-styled from FipeConsulta for this context */
  font-weight: bold;
  font-size: 1.1rem;
  color: #333;
  margin-bottom: 10px !important;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.result-container {
  animation: fadeIn 0.5s ease-out;
}
</style>
