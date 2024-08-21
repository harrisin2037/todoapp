<script>
  import { createEventDispatcher } from "svelte";

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
    editedTodo = { ...todo };
  }

  function saveEdit() {
    isEditing = false;
    dispatch("update", editedTodo);
  }

  function deleteTodo() {
    dispatch("delete");
  }
</script>

<li>
  {#if isEditing}
    <div>
      <input
        type="text"
        bind:value={editedTodo.name}
        placeholder="Enter todo name"
      />
      <input
        type="text"
        bind:value={editedTodo.description}
        placeholder="Enter todo description"
      />
      <input
        type="date"
        bind:value={editedTodo.due_date}
        placeholder="Enter due date"
      />
      <select bind:value={editedTodo.status}>
        <option value="pending">Pending</option>
        <option value="completed">Completed</option>
      </select>
      <button class="save-btn" on:click={saveEdit}>Save</button>
      <button class="cancel-btn" on:click={cancelEdit}>Cancel</button>
    </div>
  {:else}
    <div>
      <p>{todo.name}</p>
      <p>{todo.description}</p>
      <p>{todo.due_date}</p>
      <p>{todo.status}</p>
      <button class="edit-btn" on:click={startEditing}>Edit</button>
      <button class="delete-btn" on:click={() => dispatch("delete")}
        >Delete</button
      >
    </div>
  {/if}
</li>

<style>
  li {
    background-color: #f9f9f9;
    margin-bottom: 10px;
    padding: 15px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: flex-start;
    border-radius: 5px;
  }

  input {
    margin-bottom: 10px;
  }

  .edit-btn,
  .delete-btn {
    margin-right: 5px;
  }

  .edit-btn {
    background-color: #4caf50;
    color: white;
  }

  .delete-btn {
    background-color: #f44336;
    color: white;
  }
</style>
