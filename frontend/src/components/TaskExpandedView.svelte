<script>
  import { onMount, onDestroy } from "svelte";
  import { API_BASE_URL } from "../config";

  export let todo;
  export let onClose;
  export let onUpdate;

  let editedTodo = { ...todo };
  let chatMessages = [];
  let newMessage = "";
  let ws;

  onMount(() => {
    connectWebSocket();
  });

  onDestroy(() => {
    if (ws) ws.close();
  });

  function connectWebSocket() {
    ws = new WebSocket(
      `ws://${API_BASE_URL.replace(/^https?:\/\//, "")}/ws/chat/${todo.id}`
    );

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      chatMessages = [...chatMessages, message];
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
    };
  }

  async function sendMessage() {
    if (newMessage.trim() && ws.readyState === WebSocket.OPEN) {
      ws.send(
        JSON.stringify({
          senderToken: localStorage.getItem("token"),
          message: newMessage,
        })
      );
      newMessage = "";
    }
  }

  async function saveEdit() {
    const response = await fetch(`${API_BASE_URL}/todos/${todo.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify(editedTodo),
    });
    if (response.ok) {
      onUpdate();
    } else {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }
</script>

<div class="expanded-view">
  <button class="close-btn" on:click={onClose}>
    <span class="material-icons">close</span>
  </button>
  <div class="task-details">
    <h2>Task Details</h2>
    <input type="text" bind:value={editedTodo.name} placeholder="Task name" />
    <textarea
      bind:value={editedTodo.description}
      placeholder="Description"
      rows="3"
    ></textarea>
    <div class="task-meta">
      <input type="date" bind:value={editedTodo.due_date} />
      <select bind:value={editedTodo.status}>
        <option value="pending">Pending</option>
        <option value="in_progress">In Progress</option>
        <option value="completed">Completed</option>
      </select>
    </div>
    <button class="save-btn" on:click={saveEdit}>Save Changes</button>
  </div>

  <div class="chat-room">
    <h2>Chat Room</h2>
    <div class="chat-messages">
      {#each chatMessages as message}
        <div class="message">
          <strong>{message.sender}:</strong>
          {message.content}
        </div>
      {/each}
    </div>
    <div class="message-input">
      <input
        type="text"
        bind:value={newMessage}
        placeholder="Type a message..."
      />
      <button on:click={sendMessage}>Send</button>
    </div>
  </div>
</div>

<style>
  .expanded-view {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    max-width: 600px;
    width: 90%;
    max-height: 90vh;
    overflow-y: auto;
  }

  .task-details {
    margin-bottom: 2rem;
  }

  h2 {
    margin-bottom: 1rem;
  }

  input,
  textarea,
  select {
    width: 100%;
    padding: 0.5rem;
    margin-bottom: 1rem;
    border: 1px solid #e8ecee;
    border-radius: 4px;
  }

  .task-meta {
    display: flex;
    gap: 1rem;
  }

  .save-btn,
  .close-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
  }

  .save-btn {
    background-color: #7b68ee;
    color: white;
  }

  .close-btn {
    background-color: #e8ecee;
    color: #2c3e50;
    position: absolute;
    top: 10px;
    right: 10px;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 24px;
  }

  .close-btn .material-icons {
    font-size: 24px;
  }

  .chat-room {
    margin-top: 2rem;
  }
  .chat-messages {
    max-height: 300px;
    overflow-y: auto;
    height: 200px;
    border: 1px solid #e8ecee;
    padding: 1rem;
    margin-bottom: 1rem;
  }

  .message {
    margin-bottom: 0.5rem;
  }

  .message-input {
    display: flex;
    gap: 0.5rem;
  }

  .message-input input {
    flex: 1;
  }
</style>
