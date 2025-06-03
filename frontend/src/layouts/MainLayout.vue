<template>
  <div class="container">
    <header>
      <h1><i class="fas fa-car"></i> Consulta Tabela FIPE</h1>
      <p class="subtitle">Valores de referência para veículos no Brasil</p>
      <nav>
        <router-link to="/">Consulta FIPE</router-link> |
        <router-link to="/comparar">Comparar Valores</router-link> |
        <router-link to="/dashboard">Dashboard Marcas</router-link>
      </nav>
    </header>
    <slot></slot>
    <footer>
      Dados fornecidos pela Tabela FIPE &copy; 2025
    </footer>
  </div>
</template>

<script setup>
// No script needed for now
</script>

<style scoped>
/* Basic styling for the layout can be added here if needed */
</style>
