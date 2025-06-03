const API_BASE_URL = '/api'; // Assuming the Go backend is served on the same host or proxied

export async function getTabelas() {
    const response = await fetch(`${API_BASE_URL}/tabelas`);
    if (!response.ok) throw new Error('Erro ao carregar tabelas');
    return response.json();
}

export async function getMarcas(tabela) {
    const response = await fetch(`${API_BASE_URL}/marcas?tabela=${tabela}`);
    if (!response.ok) throw new Error('Erro ao carregar marcas');
    return response.json();
}

export async function getModelos(tabela, marca) {
    const response = await fetch(`${API_BASE_URL}/modelos/${marca}?tabela=${tabela}`);
    if (!response.ok) throw new Error('Erro ao carregar modelos');
    return response.json();
}

export async function getVeiculos(tabela, modelo) {
    const response = await fetch(`${API_BASE_URL}/veiculos?modelo=${modelo}&tabela=${tabela}`);
    if (!response.ok) throw new Error('Erro ao carregar veículos');
    return response.json();
}

// Helper function (can be moved to a utils file later if needed)
export function formatarPeriodoParaExibicao(mesAnoString) {
    if (!mesAnoString || typeof mesAnoString !== 'string') {
        return 'Período inválido';
    }
    const parts = mesAnoString.split('/');
    if (parts.length === 2) {
        const mes = capitalizeFirstLetter(parts[0].trim());
        const ano = parts[1].trim();
        if (mes && /^\d{4}$/.test(ano)) {
            return `${mes} de ${ano}`;
        }
    }
    return `Tabela ${mesAnoString}`;
}

function capitalizeFirstLetter(string) {
    if (!string) return '';
    return string.charAt(0).toUpperCase() + string.slice(1).toLowerCase();
}

export async function getDashboardData(tabela1, tabela2, marca) {
    let url = `${API_BASE_URL}/dashboard?tabela1=${tabela1}&tabela2=${tabela2}`;
    if (marca && marca !== 'all') {
        url += `&marca=${marca}`;
    }
    const response = await fetch(url);
    if (!response.ok) {
        let errorMsg = `Erro ao buscar dashboard: ${response.statusText}`;
        try {
            const errorBody = await response.json();
            errorMsg += ` - ${errorBody.message || errorBody.error || JSON.stringify(errorBody)}`;
        } catch(e) { /* Ignora se não conseguir parsear erro */ }
        throw new Error(errorMsg);
    }
    return response.json();
}
