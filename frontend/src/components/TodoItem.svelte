<script>
  import { createEventDispatcher } from "svelte";
  import { API_BASE_URL } from "../config";

  export let allUsers = [];
  export let todo;

  let isEditing = false;
  let editedTodo = { ...todo };
  let isExpanded = false;
  let newAssignee = "";

  const dispatch = createEventDispatcher();

  function toggleExpanded() {
    isExpanded = !isExpanded;
  }

  function handleUpdate() {
    dispatch("update");
  }

  function getStatusText(status) {
    const statusMap = {
      pending: "Pending",
      in_progress: "In Progress",
      completed: "Completed",
    };
    return statusMap[status] || "Unknown Status";
  }

  function getInitials(name) {
    return name ? name.charAt(0).toUpperCase() : "?";
  }

  function startEditing() {
    editedTodo = {
      ...todo,
      assignees: Array.isArray(todo.assignees) ? [...todo.assignees] : [],
    };
    isEditing = true;
  }

  function cancelEdit() {
    isEditing = false;
  }

  async function saveEdit() {
    const todoToSave = {
      ...editedTodo,
      assignee_ids: editedTodo.assignees
        ? editedTodo.assignees.map((user) => user.id)
        : [],
    };

    const response = await fetch(`${API_BASE_URL}/todos/${todo.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify(todoToSave),
    });

    if (response.ok) {
      todo = await response.json();
      isEditing = false;
      dispatch("update");
    } else {
      console.error("Failed to update todo");
    }
  }

  function getUserById(id) {
    return allUsers.find((user) => user.id === id);
  }

  async function deleteTodo() {
    const response = await fetch(`${API_BASE_URL}/todos/${todo.id}`, {
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

  function addAssignee() {
    if (newAssignee) {
      const userToAdd = allUsers.find(
        (user) => user.id === parseInt(newAssignee)
      );
      if (
        userToAdd &&
        !editedTodo.assignees.some((user) => user.id === userToAdd.id)
      ) {
        editedTodo.assignees = [...editedTodo.assignees, userToAdd];
        newAssignee = "";
      }
    }
  }

  function removeAssignee(assigneeId) {
    editedTodo.assignees = editedTodo.assignees.filter(
      (user) => user.id !== assigneeId
    );
  }
</script>

<div class="todo-item {todo.status}">
  <div class="todo-header">
    <div class="todo-title">
      <input
        type="checkbox"
        checked={todo.status === "completed"}
        on:change={toggleStatus}
      />
      {#if isEditing}
        <input
          type="text"
          bind:value={editedTodo.name}
          placeholder="Task name"
        />
      {:else}
        <h3 class:completed={todo.status === "completed"}>{todo.name}</h3>
      {/if}
    </div>
    <div class="todo-actions">
      {#if isEditing}
        <button class="btn save-btn" on:click={saveEdit}>Save</button>
        <button class="btn cancel-btn" on:click={cancelEdit}>Cancel</button>
      {:else}
        <button class="btn edit-btn" on:click={startEditing}>
          <span class="material-icons">edit</span>
        </button>
        <button class="btn delete-btn" on:click={deleteTodo}>
          <span class="material-icons">delete</span>
        </button>
        <button class="btn expand-btn" on:click={toggleExpanded}>
          <span class="material-icons">
            {isExpanded ? "expand_less" : "expand_more"}
          </span>
        </button>
      {/if}
    </div>
  </div>

  {#if isEditing || isExpanded}
    <div class="todo-details">
      <div class="todo-description">
        {#if isEditing}
          <textarea
            bind:value={editedTodo.description}
            placeholder="Description"
            rows="3"
          ></textarea>
        {:else}
          <p>{todo.description}</p>
        {/if}
      </div>

      <div class="todo-meta">
        <div class="meta-item">
          <span class="material-icons">event</span>
          {#if isEditing}
            <input type="date" bind:value={editedTodo.due_date} />
          {:else}
            <span>{todo.due_date}</span>
          {/if}
        </div>

        <div class="meta-item">
          <span class="material-icons">
            {todo.status === "completed" ? "check_circle" : "pending_actions"}
          </span>
          {#if isEditing}
            <select bind:value={editedTodo.status}>
              <option value="pending">Pending</option>
              <option value="in_progress">In Progress</option>
              <option value="completed">Completed</option>
            </select>
          {:else}
            <span class="status {todo.status}"
              >{getStatusText(todo.status)}</span
            >
          {/if}
        </div>
      </div>

      <!-- <div class="assignees-section">
        <h4>Owner</h4>
        <div class="assignee-item">
          <div
            class="user-bubble"
            style="background-color: {todo.owner.color || '#ccc'};"
            title={todo.owner.username}
          >
            {getInitials(todo.owner.username)}
          </div>
          <span>{todo.owner.username}</span>
        </div>
      </div> -->

      <div class="assignees-section">
        <h4>Assignees</h4>
        {#if isEditing}
          <div class="assignee-input">
            <select bind:value={newAssignee}>
              <option value="">Select user</option>
              {#each allUsers as user}
                {#if !editedTodo.assignees.some((assignee) => assignee.id === user.id)}
                  <option value={user.id.toString()}>{user.username}</option>
                {/if}
              {/each}
            </select>
            <button class="btn add-btn" on:click={addAssignee}>Add</button>
          </div>
        {/if}
        <div class="assignees-list">
          {#if editedTodo.assignees && editedTodo.assignees.length > 0}
            {#each editedTodo.assignees as assignee}
              <div class="assignee-item">
                <div
                  class="user-bubble"
                  style="background-color: {assignee.color || '#ccc'};"
                  title={assignee.username}
                >
                  {getInitials(assignee.username)}
                </div>
                <span>{assignee.username}</span>
                {#if isEditing}
                  <button
                    class="btn remove-btn"
                    on:click={() => removeAssignee(assignee.id)}
                  >
                    <span class="material-icons">close</span>
                  </button>
                {/if}
              </div>
            {/each}
          {:else}
            <p>No assignees</p>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  @import url("https://fonts.googleapis.com/icon?family=Material+Icons");

  .todo-item {
    background-color: #ffffff;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 1rem;
    transition: all 0.3s ease;
    overflow: hidden;
  }

  .todo-item:hover {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }

  .todo-item.completed {
    opacity: 0.7;
  }

  .todo-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background-color: #f8f9fa;
    border-bottom: 1px solid #e9ecef;
  }

  .todo-title {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .todo-title input[type="checkbox"] {
    width: 20px;
    height: 20px;
  }

  h3 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 500;
    color: #343a40;
  }

  h3.completed {
    text-decoration: line-through;
    color: #6c757d;
  }

  .todo-actions {
    display: flex;
    gap: 0.5rem;
  }

  .btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.25rem;
    border-radius: 4px;
    transition: background-color 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .btn:hover {
    background-color: #e9ecef;
  }

  .add-btn,
  .save-btn,
  .cancel-btn {
    padding: 0.5rem 1rem;
    font-weight: 500;
  }

  .add-btn {
    background-color: #007bff;
    color: white;
  }

  .save-btn {
    background-color: #007bff;
    color: white;
  }

  .save-btn:hover {
    background-color: #0056b3;
  }

  .cancel-btn {
    background-color: #6c757d;
    color: white;
  }

  .cancel-btn:hover {
    background-color: #545b62;
  }

  .todo-details {
    padding: 1rem;
  }

  .todo-description {
    margin-bottom: 1rem;
  }

  .todo-description p {
    margin: 0;
    color: #495057;
  }

  .todo-meta {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #6c757d;
  }

  .status {
    text-transform: capitalize;
  }

  .status.completed {
    color: #28a745;
  }
  .status.pending {
    color: #ffc107;
  }
  .status.in_progress {
    color: #17a2b8;
  }

  .assignees-section {
    margin-top: 1rem;
  }

  .assignees-section h4 {
    margin: 0 0 0.5rem 0;
    font-size: 1rem;
    color: #343a40;
  }

  .assignee-input {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }

  .assignee-input select {
    flex-grow: 1;
  }

  .assignees-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .assignee-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background-color: #e9ecef;
    padding: 0.25rem 0.5rem;
    border-radius: 16px;
  }

  .user-bubble {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.8rem;
    color: white;
    font-weight: bold;
  }

  input[type="text"],
  input[type="date"],
  select,
  textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ced4da;
    border-radius: 4px;
    font-size: 0.9rem;
  }

  textarea {
    resize: vertical;
  }
</style>
