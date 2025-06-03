<template>
  <div class="result-container">
    <div v-if="loading" class="vehicle-card">
      <p>Carregando informações do veículo...</p>
    </div>
    <div v-else-if="error" class="vehicle-card error-card">
      <p>Erro ao buscar os dados: {{ error }}</p>
    </div>
    <div v-else-if="veiculos && veiculos.length > 0">
      <h2>Resultados da Consulta</h2>
      <h3>{{ modeloNome }}</h3>
      <div v-for="veiculo in veiculos" :key="veiculo.year + '-' + veiculo.price" class="vehicle-card">
        <span class="year">{{ veiculo.year === "32000" ? '0km' : veiculo.year }}</span>
        <p class="price">{{ veiculo.price ? veiculo.price.replace(/"/g, '') : 'N/A' }}</p>
        <p>Referência: {{ veiculo.monthReference ? veiculo.monthReference.replace(/"/g, '') : 'N/A' }}</p>
      </div>
    </div>
    <div v-else-if="!loading && !error && (!veiculos || veiculos.length === 0) && modeloNome" class="vehicle-card">
      <p>Nenhum veículo encontrado para "{{ modeloNome }}" neste período.</p>
    </div>
     <div v-else-if="!loading && !error && !modeloNome" class="vehicle-card">
      <p>Selecione os filtros acima para consultar os veículos.</p>
    </div>
  </div>
</template>

<script setup>
defineProps({
  veiculos: Array,
  modeloNome: String,
  loading: Boolean,
  error: String
});
</script>

<style scoped>
.result-container {
  background-color: #fff; /* Assuming var(--white) */
  border-radius: 8px;
  padding: 25px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); /* Assuming var(--shadow) */
  min-height: 100px;
  margin-top: 30px; /* Added to separate from search form */
}

.result-container h2 {
  color: #2c3e50; /* Dark blue, good for titles */
  margin-bottom: 5px;
}
.result-container h3 {
  color: #3498db; /* Primary color */
  margin-bottom: 20px;
  font-weight: normal;
}

.result-container p {
  margin-bottom: 15px;
}

.result-container p:last-child {
  margin-bottom: 0;
}

.vehicle-card {
  border-left: 4px solid #3498db; /* Assuming var(--primary-color) */
  padding: 15px 20px; /* Increased padding */
  margin-bottom: 15px;
  background-color: #f5f5f5; /* Assuming var(--light-gray) */
  border-radius: 0 4px 4px 0;
  transition: all 0.3s ease;
}

.vehicle-card.error-card {
  border-left-color: #e74c3c; /* var(--accent-color) */
  background-color: #fbecec; /* Lighter red */
  color: #e74c3c;
}


.vehicle-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 12px rgba(0,0,0,0.1);
}

.price {
  color: #e74c3c; /* Assuming var(--accent-color) */
  font-weight: bold;
  font-size: 1.2rem; /* Slightly larger price */
  margin-bottom: 5px !important; /* Ensure spacing */
}

.year {
  font-weight: 600;
  color: #777; /* Assuming var(--dark-gray) */
  font-size: 1.1rem;
  display: block; /* Make it take its own line for better readability */
  margin-bottom: 8px;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.result-container {
  animation: fadeIn 0.5s ease-out;
}
</style>
