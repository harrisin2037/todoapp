<script>
  import { createEventDispatcher } from "svelte";
  import { API_BASE_URL } from "../config";

  export let todo;
  let isEditing = false;
  let editedTodo = { ...todo };

  const dispatch = createEventDispatcher();

  function startEditing() {
    isEditing = true;
    editedTodo = { ...todo };
  }

  function cancelEdit() {
    isEditing = false;
  }

  async function saveEdit() {
    const response = await fetch(`${API_BASE_URL}/todos/${todo.ID}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify(editedTodo),
    });
    if (response.ok) {
      isEditing = false;
      dispatch("update");
    } else {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }

  async function deleteTodo() {
    const response = await fetch(`${API_BASE_URL}/todos/${todo.ID}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });
    if (response.ok) {
      dispatch("delete");
    } else {
      const errorData = await response.json();
      alert(`Error: ${errorData.error}`);
    }
  }

  async function toggleStatus() {
    editedTodo.status =
      editedTodo.status === "completed" ? "pending" : "completed";
    await saveEdit();
  }
</script>

<div class="todo-item">
  <div class="todo-checkbox">
    <input
      type="checkbox"
      checked={todo.status === "completed"}
      on:change={toggleStatus}
    />
  </div>
  <div class="todo-content">
    {#if isEditing}
      <input type="text" bind:value={editedTodo.name} placeholder="Task name" />
      <textarea
        bind:value={editedTodo.description}
        placeholder="Description"
        rows="2"
      ></textarea>
      <div class="todo-details">
        <input type="date" bind:value={editedTodo.due_date} />
        <select bind:value={editedTodo.status}>
          <option value="pending">Pending</option>
          <option value="completed">Completed</option>
        </select>
      </div>
      <div class="todo-actions">
        <button class="save-btn" on:click={saveEdit}>Save</button>
        <button class="cancel-btn" on:click={cancelEdit}>Cancel</button>
      </div>
    {:else}
      <div class="todo-header">
        <h3 class:completed={todo.status === "completed"}>{todo.name}</h3>
        <div class="todo-actions">
          <button class="edit-btn" on:click={startEditing}>
            <span class="material-icons">edit</span>
          </button>
          <button class="delete-btn" on:click={deleteTodo}>
            <span class="material-icons">delete</span>
          </button>
        </div>
      </div>
      <p class="todo-description">{todo.description}</p>
      <div class="todo-details">
        <span class="due-date">
          <span class="material-icons">event</span>
          {todo.due_date}
        </span>
        <span class="status {todo.status}">
          <span class="material-icons">
            {todo.status === "completed" ? "check_circle" : "pending_actions"}
          </span>
          {todo.status}
        </span>
      </div>
    {/if}
  </div>
</div>

<style>
  @font-face {
    font-family: "Material Icons";
    font-style: normal;
    font-weight: 400;
    src: url(https://fonts.gstatic.com/s/materialicons/v142/flUhRq6tzZclQEJ-Vdg-IuiaDsNcIhQ8tQ.woff2)
      format("woff2");
  }

  .material-icons {
    font-family: "Material Icons";
    font-weight: normal;
    font-style: normal;
    font-size: 24px;
    display: inline-block;
    line-height: 1;
    text-transform: none;
    letter-spacing: normal;
    word-wrap: normal;
    white-space: nowrap;
    direction: ltr;
    -webkit-font-smoothing: antialiased;
    text-rendering: optimizeLegibility;
    -moz-osx-font-smoothing: grayscale;
    font-feature-settings: "liga";
  }

  .nav-button {
    background: none;
    border: none;
    cursor: pointer;
    color: #7b68ee;
    font-size: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    transition: background-color 0.3s;
  }

  .nav-button:hover {
    background-color: #f0f0f0;
  }

  .todo-item {
    display: flex;
    align-items: flex-start;
    background-color: #ffffff;
    border: 1px solid #e8ecee;
    border-radius: 8px;
    padding: 1rem;
    margin-bottom: 1rem;
    transition: box-shadow 0.3s ease;
  }

  .todo-item:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .todo-checkbox {
    margin-right: 1rem;
  }

  .todo-checkbox input[type="checkbox"] {
    width: 20px;
    height: 20px;
    cursor: pointer;
  }

  .todo-content {
    flex-grow: 1;
  }

  .todo-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  h3 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 500;
    color: #2c3e50;
  }

  h3.completed {
    text-decoration: line-through;
    color: #7f8c8d;
  }

  .todo-description {
    margin: 0 0 0.5rem 0;
    color: #34495e;
    font-size: 0.9rem;
  }

  .todo-details {
    display: flex;
    justify-content: space-between;
    font-size: 0.8rem;
    color: #7f8c8d;
  }

  .todo-details span {
    display: flex;
    align-items: center;
  }

  .todo-details .material-icons {
    font-size: 1rem;
    margin-right: 0.25rem;
  }

  .status {
    text-transform: capitalize;
  }

  .status.completed {
    color: #27ae60;
  }

  .status.pending {
    color: #e67e22;
  }

  .todo-actions {
    display: flex;
    gap: 0.5rem;
  }

  button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.25rem;
    border-radius: 4px;
    transition: background-color 0.3s ease;
  }

  button:hover {
    background-color: #f1f3f5;
  }

  .edit-btn,
  .delete-btn {
    color: #7f8c8d;
  }

  .save-btn,
  .cancel-btn {
    padding: 0.5rem 1rem;
    border-radius: 4px;
    font-weight: 500;
  }

  .save-btn {
    background-color: #7b68ee;
    color: white;
  }

  .save-btn:hover {
    background-color: #6c5ce7;
  }

  .cancel-btn {
    background-color: #e8ecee;
    color: #2c3e50;
  }

  .cancel-btn:hover {
    background-color: #d1d8e0;
  }

  input[type="text"],
  input[type="date"],
  select,
  textarea {
    width: 100%;
    padding: 0.5rem;
    margin-bottom: 0.5rem;
    border: 1px solid #e8ecee;
    border-radius: 4px;
    font-size: 0.9rem;
  }

  textarea {
    resize: vertical;
  }
</style>
