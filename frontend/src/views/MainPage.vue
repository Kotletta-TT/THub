<template>
    <div class="main-page">
        <!-- Заголовок с информацией о пользователе и кнопкой выхода -->
        <header class="header">
            <div class="user-info">
                <div class="avatar"></div>
                <span>{{ user.name }}</span>
                <button @click="logout">Logout</button>
            </div>
        </header>

        <!-- Таблица с информацией о нодах -->
        <table>
            <tr>
                <th>Node ID</th>
                <th>Users Count</th>
                <th>Users Limit</th>
                <th>State</th>
                <th>Actions</th>
            </tr>
            <tr v-for="node in nodes" :key="node.nodeId">
                <td>{{ node.nodeId }}</td>
                <td>{{ node.users_count }}</td>
                <td>{{ node.users_limit }}</td>
                <td>{{ node.state }}</td>
                <td>
                    <div class="action-buttons">
                        <button @click="openTerminal(node.nodeId)" :disabled="node.state !== 'online'">Play</button>
                        <button @click="disconnectAll(node.nodeId)" :disabled="node.state !== 'online'">Disconnect
                            All</button>
                        <button @click="restart(node.nodeId)" :disabled="node.state !== 'online'">Restart</button>
                    </div>
                </td>
            </tr>
        </table>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    data() {
        return {
            user: { name: 'User Name' },
            nodes: []
        };
    },
    async created() {
        this.fetchNodes();
        setInterval(this.fetchNodes, 5000);
    },
    methods: {
        async fetchNodes() {
            const response = await axios.get('/api/user/list/nodes?limit=10&offset=0', {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
            });
            this.nodes = response.data;
        },
        openTerminal(nodeId) {
            this.$router.push({ name: 'Terminal', params: { nodeId } });
        },
        async disconnectAll(nodeId) {
            await axios.post('/api/user/node/disconnect-all', { nodeId }, {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
            });
            this.fetchNodes();
        },
        async restart(nodeId) {
            await axios.post('/api/user/node/restart', { nodeId }, {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
            });
            this.fetchNodes();
        },
        logout() {

            // Удаляем все данные по аутентификации
            localStorage.removeItem('token');
            localStorage.removeItem('user');

            // Перенаправляем пользователя на страницу входа
            this.$router.replace('/').catch(err => {
                console.error("Routing error:", err);
            });
        }

    }
}
</script>

<style>
.main-page {
    padding: 20px;
}

.header {
    position: absolute;
    top: 20px;
    right: 20px;
    display: flex;
    align-items: center;
    gap: 10px;
}

.user-info {
    display: flex;
    align-items: center;
}

.avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: #ccc;
    /* Серый цвет для аватара */
}

table {
    width: 100%;
    background-color: #333;
    border-collapse: collapse;
    margin-top: 80px;
}

th,
td {
    padding: 15px;
    text-align: left;
    border-bottom: 1px solid #444;
    color: #fff;
}

.action-buttons {
    display: flex;
    gap: 10px;
}

button {
    background-color: #3b82f6;
    color: #fff;
    border: none;
    padding: 10px 20px;
    border-radius: 4px;
    cursor: pointer;
}

button:disabled {
    background-color: #777;
    cursor: not-allowed;
}
</style>