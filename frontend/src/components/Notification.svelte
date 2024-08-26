<script>
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  import { API_BASE_URL } from "../config";

  let notifications = [];
  let notificationId = 0;
  let ws;

  function addNotification(data) {
    const id = notificationId++;
    notifications = [...notifications, { id, data }];
    setTimeout(() => {
      removeNotification(id);
    }, 10000);
  }

  function removeNotification(id) {
    notifications = notifications.filter((n) => n.id !== id);
  }

  function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
  }

  function handleWebSocketMessage(event) {
    console.log("Raw event data:", event.data);
    try {
      const data = JSON.parse(event.data);
      console.log("Parsed data:", data);
      if (data.message === "todo updated") {
        addNotification(data);
      }
    } catch (error) {
      console.error("Error parsing WebSocket message:", error);
    }
  }

  function initializeWebSocket() {
    ws = new WebSocket(`${API_BASE_URL.replace("http", "ws")}/ws`);

    ws.onopen = () => {
      console.log("WebSocket connection established");
    };

    ws.onmessage = handleWebSocketMessage;

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
      console.log("WebSocket connection closed");
      setTimeout(initializeWebSocket, 5000);
    };
  }

  onMount(() => {
    initializeWebSocket();

    return () => {
      if (ws) {
        ws.close();
      }
    };
  });
</script>

<div class="notifications">
  {#each notifications as notification (notification.id)}
    <div class="notification" transition:fade={{ duration: 300 }}>
      <h3>{notification.data.message}</h3>
      <div class="todo-details">
        <p><strong>Name:</strong> {notification.data.todo.name}</p>
        <p>
          <strong>Description:</strong>
          {notification.data.todo.description}
        </p>
        <p>
          <strong>Due Date:</strong>
          {formatDate(notification.data.todo.due_date)}
        </p>
        <p><strong>Status:</strong> {notification.data.todo.status}</p>
        <p><strong>Owner ID:</strong> {notification.data.todo.owner_id}</p>
        <p><strong>Todo ID:</strong> {notification.data.todo.id}</p>
      </div>
      <button
        class="close-btn"
        on:click={() => removeNotification(notification.id)}>×</button
      >
    </div>
  {/each}
</div>

<style>
  .notifications {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 1000;
  }

  .notification {
    background-color: #f8f8f8;
    border: 1px solid #ddd;
    color: #333;
    padding: 16px;
    margin-bottom: 10px;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    max-width: 400px;
    word-wrap: break-word;
    position: relative;
  }

  .notification h3 {
    margin-top: 0;
    margin-bottom: 10px;
    color: #4caf50;
  }

  .todo-details p {
    margin: 5px 0;
  }

  .close-btn {
    position: absolute;
    top: 5px;
    right: 5px;
    background: none;
    border: none;
    font-size: 20px;
    cursor: pointer;
    color: #999;
  }

  .close-btn:hover {
    color: #333;
  }
</style>
