<template>
  <div class="terminal-page">
    <div id="terminal"></div>
    <div class="terminal-controls">
      <button @click="showDownloadDialog">Download File</button>
      <button @click="showUploadDialog">Upload File</button>
      <button @click="disconnect">Disconnect</button>
    </div>

    <div v-if="downloadDialog" class="dialog">
      <input v-model="downloadPath" placeholder="Path to file" class="centered-input"/>
      <div class="dialog-buttons">
        <button @click="downloadFile">Download</button>
        <button @click="closeDownloadDialog">Cancel</button>
      </div>
    </div>

    <div v-if="uploadDialog" class="dialog">
      <input type="file" @change="selectFile" />
      <div class="dialog-buttons">
        <button @click="uploadFile">Upload</button>
        <button @click="closeUploadDialog">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
import { Terminal } from 'xterm';
import 'xterm/css/xterm.css';
import { AttachAddon } from 'xterm-addon-attach';
import { FitAddon } from 'xterm-addon-fit';
import axios from 'axios';

export default {
  props: ['nodeId'],
  data() {
    return {
      terminal: null,
      socket: null,
      fitAddon: null,
      downloadDialog: false,
      uploadDialog: false,
      downloadPath: '',
      file: null
    };
  },
  mounted() {
    this.initializeTerminal();
  },
  beforeUnmount() {
    if (this.socket) {
      this.socket.close();
    }
    window.removeEventListener('resize', this.resizeScreen);
  },
  methods: {
    initializeTerminal() {
      this.terminal = new Terminal();
      this.fitAddon = new FitAddon();
      this.terminal.loadAddon(this.fitAddon);

      this.socket = new WebSocket(`/api/user/node/connect?nodeId=${this.nodeId}&t_width=800&t_height=400`);
      const attachAddon = new AttachAddon(this.socket);
      this.terminal.loadAddon(attachAddon);

      this.terminal.open(document.getElementById('terminal'));
      this.fitAddon.fit();

      this.sendTerminalSize();
      this.socket.onopen = this.sendTerminalSize;
      window.addEventListener('resize', this.resizeScreen, false);
    },
    sendTerminalSize() {
      const windowSize = { high: this.terminal.rows, width: this.terminal.cols };
      const blob = new Blob([JSON.stringify(windowSize)], { type: 'application/json' });
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(blob);
      }
    },
    resizeScreen() {
      this.fitAddon.fit();
      this.sendTerminalSize();
    },
    showDownloadDialog() { this.downloadDialog = true; },
    closeDownloadDialog() { this.downloadDialog = false; },
    showUploadDialog() { this.uploadDialog = true; },
    closeUploadDialog() { this.uploadDialog = false; },
    async downloadFile() {
      try {
        const response = await axios.post(
          '/api/user/node/download',
          { nodeId: this.nodeId, download_path: this.downloadPath },
          { headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }, responseType: 'blob' }
        );

        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', this.downloadPath.split('/').pop());
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        this.closeDownloadDialog();
      } catch (error) {
        console.error("Error downloading file:", error);
        alert("Error downloading file. Please check the path.");
      }
    },
    selectFile(event) {
      this.file = event.target.files[0];
    },
    async uploadFile() {
      if (!this.file) {
        alert("Please select a file to upload.");
        return;
      }

      const formData = new FormData();
      formData.append("file", this.file);
      formData.append("nodeId", this.nodeId);
      formData.append("name", this.file.name);

      try {
        await axios.post('/api/user/node/upload', formData, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'multipart/form-data'
          }
        });
        alert("File uploaded successfully.");
        this.closeUploadDialog();
      } catch (error) {
        console.error("Error uploading file:", error);
        alert("Error uploading file.");
      }
    },
    disconnect() {
      if (this.socket) {
        this.socket.close();
      }
      this.$router.push('/main');
    }
  }
}
</script>

<style>
.terminal-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  background-color: #2a2a2a;
  color: #fff;
  min-height: 100vh;
}

#terminal {
  width: 800px;
  height: 400px;
  background-color: #000;
  margin-bottom: 15px;
}

.terminal-controls {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.dialog {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #222;
  padding: 20px;
  border-radius: 8px;
  width: 300px; /* Увеличенная ширина окна */
}

.centered-input {
  width: calc(100% - 20px); /* Оставляем место для полей слева и справа */
  padding: 10px;
  margin-bottom: 10px;
  border-radius: 4px;
  border: none;
  text-align: center; /* Центрируем текст внутри поля */
}

.dialog-buttons {
  display: flex;
  gap: 10px;
  justify-content: center; /* Центрируем кнопки */
}

button {
  background-color: #3b82f6;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}
</style>