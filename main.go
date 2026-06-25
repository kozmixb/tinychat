package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const indexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>WebChat</title>
  <style>
    :root {
      color-scheme: dark;
      --bg: #0b0b0b;
      --panel: #1f1f1f;
      --panel-2: #2b2b2b;
      --text: #f7f7f7;
      --muted: #b7b7b7;
      --line: #5c5c5c;
      --soft-line: #3a3a3a;
      --danger: #ff7b7b;
      --radius: 22px;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font: 14px/1.5 system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background: var(--bg);
      color: var(--text);
      height: 100vh;
      overflow: hidden;
    }
    button, textarea, input {
      font: inherit;
      color: inherit;
    }
    button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      border: 0;
      background: var(--panel-2);
      border-radius: 999px;
      width: 34px;
      height: 34px;
      padding: 0;
      cursor: pointer;
      transition: background .15s ease, opacity .15s ease;
    }
    button:hover { background: #383838; }
    button.primary { background: #a8a8a8; color: #111111; }
    button.primary:hover { background: #c7c7c7; }
    button.temporary {
      background: #d6a84d;
      color: #14110a;
    }
    button.temporary:hover { background: #edc267; }
    button.wide {
      width: auto;
      padding: 0 16px;
      font-weight: 650;
    }
    button svg {
      width: 18px;
      height: 18px;
      stroke: currentColor;
      stroke-width: 2.2;
      stroke-linecap: round;
      stroke-linejoin: round;
      fill: none;
    }
    #attach svg {
      width: 23px;
      height: 23px;
      stroke-width: 2.4;
    }
    button:disabled {
      cursor: not-allowed;
      opacity: .55;
    }
    .app {
      display: grid;
      grid-template-rows: minmax(0, 1fr) auto;
      height: 100vh;
      padding: 24px;
    }
    h1 {
      font-size: 30px;
      line-height: 1.15;
      margin: 0;
      font-weight: 700;
      letter-spacing: 0;
    }
    .subtitle {
      margin-top: 7px;
      color: var(--muted);
      font-size: 18px;
    }
    .stage {
      min-height: 0;
      overflow: auto;
      padding: 0 0 24px;
      scroll-behavior: smooth;
    }
    .messages {
      display: grid;
      gap: 16px;
      max-width: 760px;
      margin: 0 auto;
      padding: 28px 0;
    }
    .empty {
      display: grid;
      place-items: center;
      min-height: calc(100vh - 280px);
      text-align: center;
    }
    .msg {
      position: relative;
      display: grid;
      gap: 7px;
      align-items: start;
    }
    .msg-footer {
      display: flex;
      align-items: center;
      gap: 8px;
      min-height: 24px;
      color: var(--muted);
    }
    .msg-actions {
      display: inline-flex;
      gap: 4px;
      opacity: 0;
      transition: opacity .15s ease;
    }
    .msg:hover .msg-actions,
    .msg:focus-within .msg-actions {
      opacity: 1;
    }
    .msg-action {
      width: auto;
      height: 22px;
      border-radius: 6px;
      padding: 0 7px;
      background: transparent;
      color: var(--muted);
      font-size: 11px;
    }
    .msg-action svg {
      width: 13px;
      height: 13px;
    }
    .msg-action span {
      margin-left: 4px;
    }
    .msg-action:hover {
      background: #303030;
      color: var(--text);
    }
    .thinking {
      color: var(--muted);
      font-size: 12px;
      white-space: pre-wrap;
      overflow-wrap: anywhere;
      border-left: 2px solid var(--soft-line);
      padding-left: 10px;
    }
    .thinking-label {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      color: var(--muted);
      font-size: 12px;
    }
    .thinking-label svg {
      width: 14px;
      height: 14px;
      stroke: currentColor;
      stroke-width: 2;
      stroke-linecap: round;
      stroke-linejoin: round;
      fill: none;
    }
    .stats {
      color: var(--muted);
      font-size: 11px;
    }
    .bubble {
      white-space: pre-wrap;
      overflow-wrap: anywhere;
      border-radius: 18px;
      padding: 12px 14px;
      max-width: 100%;
    }
    .bubble pre {
      margin: 10px 0;
      padding: 12px;
      overflow: auto;
      background: #151515;
      border: 1px solid var(--soft-line);
      border-radius: 8px;
      white-space: pre;
    }
    .bubble code {
      font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", monospace;
      font-size: 12px;
    }
    .bubble :not(pre) > code {
      padding: 2px 4px;
      border-radius: 4px;
      background: #151515;
    }
    .bubble p {
      margin: 0 0 10px;
    }
    .bubble p:last-child {
      margin-bottom: 0;
    }
    .bubble ul,
    .bubble ol {
      margin: 8px 0 10px 20px;
      padding: 0;
    }
    .bubble li {
      margin: 3px 0;
    }
    .bubble h1,
    .bubble h2,
    .bubble h3,
    .bubble h4,
    .bubble h5,
    .bubble h6 {
      margin: 14px 0 8px;
      line-height: 1.25;
      letter-spacing: 0;
    }
    .bubble h1 { font-size: 22px; }
    .bubble h2 { font-size: 19px; }
    .bubble h3 { font-size: 16px; }
    .bubble h4,
    .bubble h5,
    .bubble h6 { font-size: 14px; }
    .bubble hr {
      border: 0;
      border-top: 1px solid var(--soft-line);
      margin: 14px 0;
    }
    .bubble blockquote {
      margin: 8px 0 10px;
      padding-left: 10px;
      border-left: 2px solid var(--soft-line);
      color: var(--muted);
    }
    .code-block {
      position: relative;
      margin: 10px 0;
    }
    .code-copy {
      position: absolute;
      top: 8px;
      right: 8px;
      width: auto;
      height: 24px;
      padding: 0 8px;
      border-radius: 6px;
      background: #2c2c2c;
      color: var(--muted);
      font-size: 11px;
    }
    .code-copy:hover {
      background: #3a3a3a;
      color: var(--text);
    }
    .msg.user .bubble {
      justify-self: end;
      background: var(--panel);
      max-width: min(620px, 85%);
    }
    .msg.assistant .bubble {
      padding-left: 0;
    }
    .composer {
      width: min(768px, calc(100vw - 48px));
      margin: 0 auto;
      border: 1px solid var(--line);
      background: var(--panel);
      border-radius: var(--radius);
      min-height: 112px;
      padding: 12px;
    }
    .setup {
      width: min(768px, calc(100vw - 48px));
      margin: 0 auto 12px;
      display: none;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 8px;
      align-items: center;
    }
    .setup.visible {
      display: grid;
    }
    .setup input {
      width: 100%;
      height: 38px;
      border: 1px solid var(--line);
      outline: 0;
      background: var(--panel);
      color: var(--text);
      border-radius: 999px;
      padding: 0 14px;
    }
    .setup button {
      width: auto;
      padding: 0 14px;
      font-weight: 650;
    }
    .modal-backdrop {
      position: fixed;
      inset: 0;
      z-index: 20;
      display: none;
      place-items: center;
      padding: 18px;
      background: rgba(0, 0, 0, .72);
    }
    .modal-backdrop.visible {
      display: grid;
    }
    .modal {
      width: min(460px, 100%);
      border: 1px solid var(--line);
      border-radius: 8px;
      background: #181818;
      box-shadow: 0 22px 70px rgba(0, 0, 0, .45);
      padding: 20px;
    }
    .modal-header {
      display: grid;
      grid-template-columns: 32px 1fr;
      gap: 12px;
      align-items: center;
      margin-bottom: 12px;
    }
    .warning-icon {
      display: grid;
      place-items: center;
      width: 32px;
      height: 32px;
      color: #ffd073;
    }
    .warning-icon svg {
      width: 28px;
      height: 28px;
      stroke: currentColor;
      stroke-width: 2;
      stroke-linecap: round;
      stroke-linejoin: round;
      fill: none;
    }
    .modal h2 {
      margin: 0;
      font-size: 18px;
      line-height: 1.2;
      letter-spacing: 0;
    }
    .modal p {
      margin: 0 0 16px;
      color: var(--muted);
      text-align: left;
    }
    .modal p span {
      display: block;
      margin-top: 6px;
    }
    .modal p span:first-child {
      margin-top: 0;
    }
    .modal code {
      color: var(--text);
      font-size: 12px;
    }
    .modal form {
      display: grid;
      gap: 10px;
    }
    .modal input {
      width: 100%;
      height: 42px;
      border: 1px solid var(--line);
      outline: 0;
      background: var(--panel);
      color: var(--text);
      border-radius: 8px;
      padding: 0 12px;
    }
    .modal-actions {
      display: flex;
      justify-content: flex-end;
    }
    .composer-inner {
      display: grid;
      grid-template-rows: minmax(42px, auto) auto;
      gap: 8px;
    }
    #prompt {
      width: 100%;
      min-height: 42px;
      max-height: 160px;
      resize: none;
      border: 0;
      outline: 0;
      background: transparent;
      padding: 4px 8px;
      color: var(--text);
      font-size: 14px;
    }
    #prompt::placeholder { color: #a7a7a7; }
    .composer-actions {
      display: grid;
      grid-template-columns: auto 1fr auto auto;
      gap: 8px;
      align-items: center;
    }
    .model-wrap {
      justify-self: end;
      display: grid;
      grid-template-columns: minmax(112px, 300px);
      align-items: center;
      height: 26px;
      border-radius: 6px;
      background: #343434;
      padding: 0 8px;
      min-width: 0;
    }
    .model-text {
      min-width: 0;
      height: 24px;
      color: var(--text);
      font-size: 12px;
      line-height: 24px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .endpoint, .status {
      color: #9d9d9d;
      font-size: 11px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .error {
      color: var(--danger);
      font-size: 12px;
      overflow-wrap: anywhere;
      min-height: 18px;
      width: min(768px, calc(100vw - 48px));
      margin: 8px auto 0;
      text-align: left;
    }
    @media (max-width: 760px) {
      .app { padding: 14px; }
      h1 { font-size: 26px; }
      .subtitle { font-size: 15px; }
      .messages { padding-top: 14px; }
      .composer {
        width: 100%;
      }
      .setup {
        width: 100%;
        grid-template-columns: 1fr;
      }
      .composer-actions {
        grid-template-columns: auto minmax(0, 1fr) auto;
      }
      .model-wrap {
        grid-column: 2;
        max-width: 100%;
      }
      .status { display: none; }
    }
  </style>
</head>
<body>
  <div class="modal-backdrop" id="hostModal" role="dialog" aria-modal="true" aria-labelledby="hostModalTitle">
    <div class="modal">
      <div class="modal-header">
        <div class="warning-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24"><path d="M12 3 2.5 20h19L12 3Z"></path><path d="M12 9v5"></path><path d="M12 17h.01"></path></svg>
        </div>
        <h2 id="hostModalTitle">OpenAI-compatible endpoint required</h2>
      </div>
      <p>
        <span>No OpenAI-compatible endpoint was provided.</span>
        <span>You can set one for this browser session.</span>
        <span>For production or shared deployments, configure <code>OPENAI_CHAT_HOST</code> on the host server.</span>
      </p>
      <form id="hostForm">
        <input id="runtimeEndpoint" placeholder="https://api.example.com/v1" autocomplete="url" inputmode="url">
        <div class="modal-actions">
          <button class="temporary wide" type="submit">Use temporarily</button>
        </div>
      </form>
    </div>
  </div>
  <div class="app">
    <section class="stage" id="chat">
      <div class="messages" id="messages"></div>
    </section>

    <div>
      <form class="composer" id="form">
        <div class="composer-inner">
          <textarea id="prompt" placeholder="Type a message..."></textarea>
          <div class="composer-actions">
            <button id="attach" type="button" title="Attach file" aria-label="Attach file">
              <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 5v14M5 12h14"></path></svg>
            </button>
            <div class="endpoint" id="endpoint"></div>
            <div class="model-wrap">
              <div class="model-text" id="model" title="Model"></div>
            </div>
            <button class="primary" id="send" type="submit" title="Send" aria-label="Send">
              <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M5 12h14M13 6l6 6-6 6"></path></svg>
            </button>
          </div>
        </div>
      </form>
      <input id="file" type="file" hidden multiple>
      <div class="error" id="error"></div>
    </div>
  </div>

  <script>
    const els = {
      endpoint: document.querySelector('#endpoint'),
      hostModal: document.querySelector('#hostModal'),
      hostForm: document.querySelector('#hostForm'),
      runtimeEndpoint: document.querySelector('#runtimeEndpoint'),
      model: document.querySelector('#model'),
      attach: document.querySelector('#attach'),
      file: document.querySelector('#file'),
      error: document.querySelector('#error'),
      chat: document.querySelector('#chat'),
      messages: document.querySelector('#messages'),
      form: document.querySelector('#form'),
      prompt: document.querySelector('#prompt'),
      send: document.querySelector('#send')
    };

    let messages = [];
    let controller = null;
    let configuredEndpoint = '';
    let runtimeEndpoint = '';
    let selectedModel = null;
    let renderPending = false;
    let forceNextScroll = false;
    let showThinkingByMessage = {};
    const numberFormat = new Intl.NumberFormat();
    const runtimeEndpointCookie = 'webchat_runtime_endpoint';
    const runtimeEndpointCookieMaxAge = 60 * 60;

    function setError(text) { els.error.textContent = text || ''; }
    function nearBottom() {
      return els.chat.scrollHeight - els.chat.scrollTop - els.chat.clientHeight < 80;
    }
    function scrollBottom() { els.chat.scrollTop = els.chat.scrollHeight; }
    function activeEndpoint() { return configuredEndpoint || runtimeEndpoint; }
    function displayEndpoint(endpoint) { return endpoint.replace(/\/$/, ''); }
    function readCookie(name) {
      const prefix = name + '=';
      return document.cookie
        .split(';')
        .map(cookie => cookie.trim())
        .find(cookie => cookie.startsWith(prefix))
        ?.slice(prefix.length) || '';
    }
    function saveRuntimeEndpoint(endpoint) {
      const value = encodeURIComponent(endpoint);
      document.cookie = runtimeEndpointCookie + '=' + value + '; Max-Age=' + runtimeEndpointCookieMaxAge + '; Path=/; SameSite=Lax';
    }
    function showHostModal() {
      els.hostModal.classList.add('visible');
      window.setTimeout(() => els.runtimeEndpoint.focus(), 0);
    }
    function hideHostModal() {
      els.hostModal.classList.remove('visible');
    }
    function endpointHeaders() {
      const endpoint = activeEndpoint();
      return endpoint ? { 'X-OpenAI-Chat-Host': endpoint } : {};
    }
    function setModelOptions(models, placeholder) {
      const selectedId = selectedModel?.id || '';
      if (!models.length) {
        selectedModel = null;
        els.model.textContent = placeholder || 'No models found';
        els.model.title = els.model.textContent;
        return;
      }
      selectedModel = models.find(model => model.id === selectedId) || models[0];
      els.model.textContent = modelLabel(selectedModel);
      els.model.title = els.model.textContent;
    }
    function modelLabel(model) {
      const tokens = model.maxTokens ? ' - ' + numberFormat.format(model.maxTokens) + ' tokens' : '';
      return model.id + tokens;
    }
    function modelList(data) {
      if (Array.isArray(data)) {
        return data.map(modelInfo).filter(model => model.id);
      }
      if (!data || typeof data !== 'object') return [];
      const candidates = [data.data, data.models, data.model, data.tags].filter(Boolean);
      for (const candidate of candidates) {
        const items = Array.isArray(candidate) ? candidate : [candidate];
        const models = items.map(modelInfo).filter(model => model.id);
        if (models.length) return models;
      }
      return [];
    }
    function modelInfo(item) {
      if (typeof item === 'string') return { id: item, maxTokens: 0 };
      if (!item || typeof item !== 'object') return { id: '', maxTokens: 0 };
      return {
        id: item.id || item.name || item.model || item.digest || '',
        maxTokens: modelMaxTokens(item)
      };
    }
    function modelMaxTokens(item) {
      const value = item.max_tokens || item.maxTokens || item.context_length || item.contextLength || item.n_ctx || item.n_ctx_train || item.meta?.n_ctx || item.meta?.n_ctx_train || item.details?.n_ctx;
      const tokens = Number(value);
      return Number.isFinite(tokens) && tokens > 0 ? tokens : 0;
    }

    function render() {
      const shouldScroll = forceNextScroll || nearBottom();
      forceNextScroll = false;
      els.messages.innerHTML = '';
      for (const [index, msg] of messages.entries()) {
        const row = document.createElement('div');
        row.className = 'msg ' + msg.role;
        const bubble = document.createElement('div');
        bubble.className = 'bubble';
        bubble.innerHTML = renderMarkdown(msg.content || msg.status || '');
        row.append(bubble);
        if (msg.reasoning && (showThinkingByMessage[index] || !msg.done)) {
          const thinkingLabel = document.createElement('div');
          thinkingLabel.className = 'thinking-label';
          thinkingLabel.innerHTML = '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6S2 12 2 12Z"></path><path d="M12 9a3 3 0 1 1 0 6 3 3 0 0 1 0-6Z"></path></svg><span>Thinking</span>';
          const thinking = document.createElement('div');
          thinking.className = 'thinking';
          thinking.textContent = msg.reasoning;
          row.append(thinkingLabel, thinking);
        }
        const footer = document.createElement('div');
        footer.className = 'msg-footer';
        if (msg.stats) {
          const stats = document.createElement('div');
          stats.className = 'stats';
          stats.textContent = msg.stats;
          footer.append(stats);
        }
        const actions = document.createElement('div');
        actions.className = 'msg-actions';
        const copy = document.createElement('button');
        copy.className = 'msg-action';
        copy.type = 'button';
        copy.textContent = 'Copy';
        copy.title = 'Copy message';
        copy.addEventListener('click', () => copyMessage(msg, copy));
        actions.append(copy);
        if (msg.reasoning && msg.done) {
          const thinking = document.createElement('button');
          thinking.className = 'msg-action';
          thinking.type = 'button';
          thinking.innerHTML = '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6S2 12 2 12Z"></path><path d="M12 9a3 3 0 1 1 0 6 3 3 0 0 1 0-6Z"></path></svg><span>' + (showThinkingByMessage[index] ? 'Hide thinking' : 'Thinking') + '</span>';
          thinking.title = 'Toggle thinking';
          thinking.addEventListener('click', () => {
            showThinkingByMessage[index] = !showThinkingByMessage[index];
            render();
          });
          actions.append(thinking);
        }
        footer.append(actions);
        row.append(footer);
        els.messages.append(row);
      }
      if (!messages.length) {
        const empty = document.createElement('div');
        empty.className = 'empty';
        const box = document.createElement('div');
        const title = document.createElement('h1');
        title.textContent = 'Hello there';
        const subtitle = document.createElement('div');
        subtitle.className = 'subtitle';
        subtitle.textContent = 'Type a message or upload files to get started';
        box.append(title, subtitle);
        empty.append(box);
        els.messages.append(empty);
      }
      if (shouldScroll) scrollBottom();
    }
    async function copyMessage(msg, button) {
      const parts = [];
      if (msg.content) parts.push(msg.content);
      if (msg.reasoning) parts.push('Thinking:\n' + msg.reasoning);
      const text = parts.join('\n\n') || msg.status || '';
      if (!text) return;
      try {
        await navigator.clipboard.writeText(text);
        button.textContent = 'Copied';
        window.setTimeout(() => { button.textContent = 'Copy'; }, 900);
      } catch {
        setError('Copy failed.');
      }
    }
    function scheduleRender() {
      if (renderPending) return;
      renderPending = true;
      requestAnimationFrame(() => {
        renderPending = false;
        render();
      });
    }

    async function loadConfig() {
      const res = await fetch('/api/config');
      if (!res.ok) throw new Error('config request failed');
      const cfg = await res.json();
      configuredEndpoint = cfg.openai_chat_host || '';
      if (configuredEndpoint) {
        els.endpoint.textContent = displayEndpoint(configuredEndpoint);
        hideHostModal();
      } else {
        runtimeEndpoint = decodeURIComponent(readCookie(runtimeEndpointCookie));
        if (runtimeEndpoint) {
          els.runtimeEndpoint.value = runtimeEndpoint;
          els.endpoint.textContent = displayEndpoint(runtimeEndpoint);
          hideHostModal();
        } else {
          els.endpoint.textContent = 'runtime endpoint required';
          showHostModal();
        }
      }
    }

    async function loadModels() {
      if (!activeEndpoint()) {
        setModelOptions([], 'No endpoint');
        showHostModal();
        return;
      }
      setModelOptions([], 'Loading...');
      await checkHealth();
      const res = await fetch('/api/models', { headers: endpointHeaders() });
      const body = await res.text();
      if (!res.ok) throw new Error(body || 'model lookup failed');
      let data;
      try {
        data = JSON.parse(body);
      } catch {
        throw new Error('model lookup returned invalid JSON');
      }
      setModelOptions(modelList(data), 'No models found');
    }

    async function checkHealth() {
      const res = await fetch('/api/health', { headers: endpointHeaders() });
      const body = await res.text();
      if (!res.ok) throw new Error(body || 'health check failed');
    }

    async function sendMessage(text) {
      setError('');
      const model = selectedModel?.id || '';
      if (!activeEndpoint()) {
        setError('Provide an OpenAI-compatible API host first.');
        showHostModal();
        return;
      }
      if (!model) {
        setError('Select a model first.');
        return;
      }

      messages.push({ role: 'user', content: text });
      messages.push({ role: 'assistant', content: '', reasoning: '', done: false });
      forceNextScroll = true;
      render();

      const requestMessages = messages.slice(0, -1).map(msg => ({
        role: msg.role,
        content: msg.content || ''
      }));

      controller = new AbortController();
      els.send.disabled = true;
      els.send.innerHTML = '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M6 6l12 12M18 6L6 18"></path></svg>';
      els.send.title = 'Stop';
      els.send.setAttribute('aria-label', 'Stop');

      try {
        const res = await fetch('/api/chat', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json', ...endpointHeaders() },
          body: JSON.stringify({
            model,
            messages: requestMessages,
            stream: true
          }),
          signal: controller.signal
        });
        if (!res.ok || !res.body) throw new Error(await res.text());

        const reader = res.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;
          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split('\n');
          buffer = lines.pop() || '';
          for (const line of lines) {
            if (!line.startsWith('data:')) continue;
            const data = line.slice(5).trim();
            if (!data) continue;
            if (data === '[DONE]') {
              messages[messages.length - 1].done = true;
              messages[messages.length - 1].status = '';
              scheduleRender();
              continue;
            }
            const chunk = JSON.parse(data);
            const delta = chunk.choices?.[0]?.delta?.content || chunk.choices?.[0]?.message?.content || '';
            const reasoning = chunk.choices?.[0]?.delta?.reasoning_content || '';
            const finishReason = chunk.choices?.[0]?.finish_reason || '';
            const timings = chunk.timings;
            if (reasoning) {
              messages[messages.length - 1].reasoning += reasoning;
              if (!messages[messages.length - 1].content) {
                messages[messages.length - 1].status = 'Thinking...';
              }
              scheduleRender();
            }
            if (delta) {
              messages[messages.length - 1].content += delta;
              messages[messages.length - 1].status = '';
              scheduleRender();
            }
            if (timings) {
              messages[messages.length - 1].stats = timingStats(timings);
              scheduleRender();
            }
            if (finishReason) {
              messages[messages.length - 1].done = true;
              messages[messages.length - 1].status = '';
              scheduleRender();
            }
          }
        }
      } catch (err) {
        if (err.name !== 'AbortError') {
          setError(String(err.message || err));
        }
      } finally {
        controller = null;
        els.send.disabled = false;
        els.send.innerHTML = '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M5 12h14M13 6l6 6-6 6"></path></svg>';
        els.send.title = 'Send';
        els.send.setAttribute('aria-label', 'Send');
        els.prompt.focus();
      }
    }
    function timingStats(timings) {
      const promptTokens = Number(timings.prompt_n || 0);
      const generatedTokens = Number(timings.predicted_n || 0);
      const promptSpeed = Number(timings.prompt_per_second || 0);
      const generatedSpeed = Number(timings.predicted_per_second || 0);
      const prompt = promptTokens ? 'Prompt ' + numberFormat.format(promptTokens) + ' tokens' + speedLabel(promptSpeed) : '';
      const generated = generatedTokens ? 'Generated ' + numberFormat.format(generatedTokens) + ' tokens' + speedLabel(generatedSpeed) : '';
      return [prompt, generated].filter(Boolean).join(' - ');
    }
    function speedLabel(value) {
      return value > 0 ? ' @ ' + value.toFixed(1) + ' tok/s' : '';
    }
    function renderMarkdown(text) {
      if (!text) return '';
      const fence = String.fromCharCode(96, 96, 96);
      const parts = String(text).split(fence);
      return parts.map((part, index) => {
        if (index % 2) {
          const code = part.replace(/^\w+\n/, '');
          return codeBlock(code);
        }
        return renderInlineMarkdown(escapeHtml(part));
      }).join('');
    }
    function codeBlock(code) {
      return '<div class="code-block"><button class="code-copy" type="button" data-copy-code="' + encodeURIComponent(code) + '">Copy code</button><pre><code>' + escapeHtml(code) + '</code></pre></div>';
    }
    function renderInlineMarkdown(text) {
      const blocks = text.split(/\n{2,}/).filter(Boolean);
      return blocks.map(renderMarkdownBlock).join('');
    }
    function renderMarkdownBlock(block) {
      const lines = block.split('\n');
      if (/^#{1,6} /.test(lines[0])) {
        const level = Math.min(lines[0].match(/^#+/)[0].length, 6);
        const heading = '<h' + level + '>' + inlineMarkdown(lines[0].replace(/^#{1,6} /, '')) + '</h' + level + '>';
        const rest = lines.slice(1).join('\n').trim();
        return heading + (rest ? renderMarkdownBlock(rest) : '');
      }
      if (lines.length === 1 && /^[-*_]{3,}$/.test(lines[0].trim())) return '<hr>';
      if (lines.every(line => /^> ?/.test(line))) {
        return '<blockquote>' + lines.map(line => inlineMarkdown(line.replace(/^> ?/, ''))).join('<br>') + '</blockquote>';
      }
      if (lines.every(line => /^[-*] /.test(line))) {
        return '<ul>' + lines.map(line => '<li>' + inlineMarkdown(line.slice(2)) + '</li>').join('') + '</ul>';
      }
      if (lines.every(line => /^\d+\. /.test(line))) {
        return '<ol>' + lines.map(line => '<li>' + inlineMarkdown(line.replace(/^\d+\. /, '')) + '</li>').join('') + '</ol>';
      }
      return '<p>' + inlineMarkdown(lines.join('<br>')) + '</p>';
    }
    function inlineMarkdown(text) {
      const tick = String.fromCharCode(96);
      return text
        .replace(new RegExp(tick + '([^' + tick + ']+)' + tick, 'g'), '<code>$1</code>')
        .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
        .replace(/\*([^*]+)\*/g, '<em>$1</em>');
    }
    async function copyCode(button) {
      try {
        await navigator.clipboard.writeText(decodeURIComponent(button.dataset.copyCode || ''));
        button.textContent = 'Copied';
        window.setTimeout(() => { button.textContent = 'Copy code'; }, 900);
      } catch {
        setError('Copy failed.');
      }
    }
    function escapeHtml(text) {
      return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
    }

    els.form.addEventListener('submit', event => {
      event.preventDefault();
      if (controller) {
        controller.abort();
        return;
      }
      const text = els.prompt.value.trim();
      if (!text) return;
      els.prompt.value = '';
      sendMessage(text);
    });
    els.prompt.addEventListener('keydown', event => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault();
        els.form.requestSubmit();
      }
    });
    els.attach.addEventListener('click', () => els.file.click());
    els.messages.addEventListener('click', event => {
      const button = event.target.closest('[data-copy-code]');
      if (button) copyCode(button);
    });
    els.hostForm.addEventListener('submit', async event => {
      event.preventDefault();
      runtimeEndpoint = els.runtimeEndpoint.value.trim();
      if (!runtimeEndpoint) {
        setError('Enter an OpenAI-compatible API host.');
        showHostModal();
        return;
      }
      els.endpoint.textContent = displayEndpoint(runtimeEndpoint);
      setError('');
      try {
        await loadModels();
        saveRuntimeEndpoint(runtimeEndpoint);
        hideHostModal();
        els.prompt.focus();
      } catch (err) {
        showHostModal();
        setError(String(err.message || err));
      }
    });
    els.file.addEventListener('change', async () => {
      const files = Array.from(els.file.files || []);
      if (!files.length) return;
      const chunks = [];
      for (const file of files) {
        if (file.type.startsWith('text/') || file.size < 512 * 1024) {
          chunks.push('File: ' + file.name + '\n' + await file.text());
        } else {
          chunks.push('File attached: ' + file.name);
        }
      }
      els.prompt.value = [els.prompt.value.trim(), ...chunks].filter(Boolean).join('\n\n');
      els.file.value = '';
      els.prompt.focus();
    });

    (async function init() {
      try {
        await loadConfig();
        await loadModels();
        render();
      } catch (err) {
        setError(String(err.message || err));
        render();
      }
    })();
  </script>
</body>
</html>`

type appConfig struct {
	Addr          string `json:"-"`
	OpenAIChatURL string `json:"openai_chat_host"`
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndex)
	mux.HandleFunc("/api/config", handleConfig(cfg))
	mux.HandleFunc("/api/health", proxyRequest(cfg, http.MethodGet, "/health"))
	mux.HandleFunc("/api/models", proxyRequest(cfg, http.MethodGet, "/models"))
	mux.HandleFunc("/api/chat", proxyRequest(cfg, http.MethodPost, "/chat/completions"))

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("listening on http://%s", cfg.Addr)
	log.Printf("proxying OpenAI-compatible chat API at %s", cfg.OpenAIChatURL)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func loadConfig() (appConfig, error) {
	host := env("APP_HOST", "127.0.0.1")
	port := env("APP_PORT", "8080")
	openAIHost := strings.TrimSpace(os.Getenv("OPENAI_CHAT_HOST"))

	if openAIHost != "" {
		openAIHost = normalizeOpenAIBaseURL(openAIHost)
		if _, err := url.ParseRequestURI(openAIHost); err != nil {
			return appConfig{}, fmt.Errorf("invalid OPENAI_CHAT_HOST: %w", err)
		}
	}

	return appConfig{
		Addr:          net.JoinHostPort(host, port),
		OpenAIChatURL: strings.TrimRight(openAIHost, "/"),
	}, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func normalizeBaseURL(raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	return "http://" + raw
}

func normalizeOpenAIBaseURL(raw string) string {
	baseURL := strings.TrimRight(normalizeBaseURL(strings.TrimSpace(raw)), "/")
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}
	if parsed.Path == "" || parsed.Path == "/" {
		parsed.Path = "/v1"
		return strings.TrimRight(parsed.String(), "/")
	}
	return baseURL
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, indexHTML)
}

func handleConfig(cfg appConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, cfg)
	}
}

func proxyRequest(cfg appConfig, method, path string) http.HandlerFunc {
	client := &http.Client{Timeout: 0}
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		baseURL, err := requestBaseURL(cfg, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		target := baseURL + path
		body, err := readBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), method, target, bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream, application/json")
		log.Printf("proxy request method=%s target=%s content_type=%q accept=%q bytes=%d", method, target, req.Header.Get("Content-Type"), req.Header.Get("Accept"), len(body))

		resp, err := client.Do(req)
		if err != nil {
			status := http.StatusBadGateway
			if errors.Is(err, context.Canceled) {
				status = http.StatusRequestTimeout
			}
			log.Printf("proxy response method=%s target=%s status=%d error=%q duration=%s", method, target, status, err.Error(), time.Since(start).Round(time.Millisecond))
			http.Error(w, err.Error(), status)
			return
		}
		defer resp.Body.Close()
		log.Printf("proxy response method=%s target=%s status=%d content_type=%q duration=%s", method, target, resp.StatusCode, resp.Header.Get("Content-Type"), time.Since(start).Round(time.Millisecond))

		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		_, _ = copyResponse(w, resp.Body)
	}
}

func copyResponse(w http.ResponseWriter, r io.Reader) (int64, error) {
	flusher, _ := w.(http.Flusher)
	buf := make([]byte, 32*1024)
	var written int64
	for {
		nr, er := r.Read(buf)
		if nr > 0 {
			nw, ew := w.Write(buf[:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if flusher != nil {
				flusher.Flush()
			}
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
		}
		if er != nil {
			if errors.Is(er, io.EOF) {
				return written, nil
			}
			return written, er
		}
	}
}

func requestBaseURL(cfg appConfig, r *http.Request) (string, error) {
	if cfg.OpenAIChatURL != "" {
		return cfg.OpenAIChatURL, nil
	}

	raw := strings.TrimSpace(r.Header.Get("X-OpenAI-Chat-Host"))
	if raw == "" {
		return "", errors.New("OPENAI_CHAT_HOST is not configured; provide a runtime host such as https://api.example.com/v1")
	}

	baseURL := normalizeOpenAIBaseURL(raw)
	parsed, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid OpenAI chat host: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("OpenAI chat host must use http or https")
	}
	if parsed.Host == "" {
		return "", errors.New("OpenAI chat host must include a host")
	}

	return baseURL, nil
}

func readBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()
	return io.ReadAll(io.LimitReader(r.Body, 8<<20))
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		lower := strings.ToLower(key)
		if lower == "content-length" || lower == "connection" {
			continue
		}
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
