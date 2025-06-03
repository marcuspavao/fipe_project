<template>
  <div class="results-grid">
    <div v-if="loading" class="status loading">Buscando dados do dashboard...</div>
    <div v-else-if="error" class="status error">Falha ao carregar dashboard: {{ error }}</div>
    <div v-else-if="statusMessage && !loading && !error" class="status">{{ statusMessage }}</div>
    
    <template v-if="dashboardData && dashboardData.length > 0 && !loading && !error">
      <div v-for="brandEntry in dashboardData" :key="brandEntry.brandCode" class="brand-card">
        <h2>{{ brandEntry.brandName }} ({{ brandEntry.brandCode }})</h2>
        
        <div v-if="brandEntry.periodo1" v-html="renderPeriodStats(brandEntry.periodo1, 1, getPeriodRefLabel(brandEntry.periodo1.ref))" class="period-stats-container"></div>
        <div v-else class="period-stats-container period-stats period-1"><p>Dados não disponíveis para o primeiro período.</p></div>

        <div v-if="brandEntry.periodo2" v-html="renderPeriodStats(brandEntry.periodo2, 2, getPeriodRefLabel(brandEntry.periodo2.ref))" class="period-stats-container"></div>
        <div v-else class="period-stats-container period-stats period-2"><p>Dados não disponíveis para o segundo período.</p></div>

        <div class="diff-stats">
          <h3>Diferenças Percentuais (Período 2 vs Período 1)</h3>
          <p><strong>Valor Médio 0km:</strong> <span v-html="formatPercentage(brandEntry.diferencasPercentuais?.valorMedio0km)"></span></p>
          <p><strong>Total de Modelos:</strong> <span v-html="formatPercentage(brandEntry.diferencasPercentuais?.totalModelos)"></span></p>
        </div>
      </div>
    </template>
    
    <div v-if="!dashboardData || dashboardData.length === 0 && !loading && !error && !statusMessage" class="status">
      Nenhum dado encontrado para os critérios selecionados.
    </div>
  </div>
</template>

<script setup>
import { formatarPeriodoParaExibicao } from '../../services/apiService'; // Assuming this can be used for refs

const props = defineProps({
  dashboardData: Array,
  loading: Boolean,
  error: String,
  statusMessage: String
  // tabelas: Array prop removed as it's not directly needed for getPeriodRefLabel's current implementation
});

function getPeriodRefLabel(refCode) {
  // The dashboard API currently returns "ref": "299/fevereiro/2024" or similar directly in periodo1/2.
  // This function will format that string.
  if (!refCode) return 'Período Indefinido';
  
  const parts = refCode.split('/');
  // Handles "mes/ano" or "codigo/mes/ano" by taking the last two parts.
  if (parts.length >= 2) { 
      const mesAnoString = parts.slice(parts.length - 2).join('/'); 
      return formatarPeriodoParaExibicao(mesAnoString);
  }
  // If only one part (e.g., just a code, or already formatted string not matching "xx/yy")
  return formatarPeriodoParaExibicao(refCode); // Attempt to format, or it will return "Tabela [refCode]"
}

function formatPercentage(value) {
  if (value === null || value === undefined || isNaN(value)) {
    return '<span class="neutral">N/A</span>';
  }
  const formatted = value.toFixed(2).replace('.', ',');
  const sign = value > 0 ? '+' : '';
  const className = value > 0 ? 'positive' : (value < 0 ? 'negative' : 'neutral');
  return `<span class="${className}">${sign}${formatted}%</span>`;
}

function renderPeriodStats(periodStats, periodNum, periodRefLabel) {
  if (!periodStats) {
    return `<div class="period-stats period-${periodNum}"><p>Dados não disponíveis.</p></div>`;
  }
  // Use periodRefLabel directly as it's now processed before calling this function
  return `
    <div class="period-stats period-${periodNum}">
      <h3>${periodRefLabel}</h3>
      <p><strong>Menor Preço 0km:</strong> <span class="price-detail">${periodStats.menorPreco0km?.valorFmt || 'N/A'}</span></p>
      ${periodStats.menorPreco0km?.modelo && periodStats.menorPreco0km?.modelo !== 'N/A' ? `<span class="modelo-destaque">(${periodStats.menorPreco0km.modelo})</span>` : ''}
      <p><strong>Maior Preço 0km:</strong> <span class="price-detail">${periodStats.maiorPreco0km?.valorFmt || 'N/A'}</span></p>
      ${periodStats.maiorPreco0km?.modelo && periodStats.maiorPreco0km?.modelo !== 'N/A' ? `<span class="modelo-destaque">(${periodStats.maiorPreco0km.modelo})</span>` : ''}
      <p><strong>Valor Médio 0km:</strong> <span class="price-detail">${periodStats.valorMedio0kmFmt || 'N/A'}</span></p>
      <p><strong>Total Modelos:</strong> <span class="price-detail">${periodStats.totalModelos || 0}</span></p>
    </div>
  `;
}
</script>

<style scoped>
.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.brand-card {
  background-color: #fff; /* var(--card-bg) */
  border: 1px solid #dee2e6; /* var(--border-color) */
  border-radius: 6px;
  padding: 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05); /* Softer shadow */
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.brand-card h2 {
  margin: 0 0 15px 0;
  padding-bottom: 10px;
  border-bottom: 2px solid #3498db; /* var(--primary-color) */
  text-align: center;
  font-size: 1.4em;
  color: #3498db; /* var(--primary-color) */
}

.period-stats-container { /* Wrapper for v-html content */
  border: 1px dashed #dee2e6; /* var(--border-color) */
  padding: 15px;
  border-radius: 4px;
}

/* Styles for content inside v-html (won't be scoped effectively by default) */
/* These need to be global or duplicated if not using ::v-deep or similar */
/* For now, assuming they might be picked up if not too complex, or use :deep() in Vue 3 */
/* It's generally better to avoid v-html if extensive styling of its content is needed from scoped CSS */

:deep(.period-stats h3) {
  margin: 0 0 10px 0;
  font-size: 1.1em;
  color: #555;
}

:deep(.period-stats p),
:deep(.diff-stats p) { /* diff-stats is not inside v-html, so this is fine */
  margin: 5px 0;
  font-size: 0.95em;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.period-stats p strong),
:deep(.diff-stats p strong) {
  color: #444;
  margin-right: 8px;
  flex-shrink: 0;
}

:deep(.period-stats .price-detail) {
  font-weight: bold;
}

:deep(.modelo-destaque) {
  font-style: italic;
  font-size: 0.85em;
  color: #555;
  margin-left: 5px; /* This was intended to be under the price line probably */
  display: block; /* Make it appear on its own line below the price */
  text-align: right; /* Align with the price value */
  margin-top: -4px; /* Adjust to pull it closer to the price line */
}


.diff-stats {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #dee2e6; /* var(--border-color) */
}

.diff-stats h3 {
  margin: 0 0 10px 0;
  font-size: 1.1em;
  color: #555;
  text-align: center;
}

/* For spans inside v-html for formatPercentage */
:deep(.positive) { color: #28a745; /* var(--positive-color) */ }
:deep(.negative) { color: #dc3545; /* var(--negative-color) */ }
:deep(.neutral)  { color: #6c757d; /* var(--neutral-color) */ }


/* Status messages styling */
.status {
  text-align: center;
  margin: 20px 0;
  padding: 15px;
  font-weight: bold;
  border-radius: 4px;
  background-color: #f8f9fa; /* Light background for general status */
  border: 1px solid #e9ecef;
}

.status.error {
  color: #dc3545; /* var(--negative-color) */
  background-color: #f8d7da; /* Light red */
  border-color: #f5c6cb;
}

.status.loading {
  color: #3498db; /* var(--primary-color) */
  background-color: #e7f3fe; /* Light blue */
  border-color: #d0eaff;
}
</style>
