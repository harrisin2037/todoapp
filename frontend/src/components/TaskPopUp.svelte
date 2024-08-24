<script>
  export let date;
  export let tasks = [];
  export let onClose;
  export let onToggleComplete;

  function handleToggle(task) {
    onToggleComplete(task.id);
  }

  function formatDate(dateString) {
    const options = {
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "numeric",
    };
    return new Date(dateString).toLocaleString("en-US", options);
  }

  console.log("TaskPopup rendered with date:", date, "and tasks:", tasks);
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="popup-overlay" on:click={onClose}>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="popup-content" on:click|stopPropagation>
    <div class="popup-header">
      <h2>
        Tasks for {date.toLocaleDateString("en-US", {
          month: "long",
          day: "numeric",
          year: "numeric",
        })}
      </h2>
      <button class="close-button" on:click={onClose}>
        <span class="material-icons">close</span>
      </button>
    </div>
    <div class="task-list">
      {#each tasks as task}
        <div class="task-item">
          <label class="checkbox-container">
            <input
              type="checkbox"
              checked={task.status === "completed"}
              on:change={() => handleToggle(task)}
            />
            <span class="checkmark"></span>
          </label>
          <div class="task-details">
            <span class="task-name">{task.name}</span>
            {#if task.due_date}
              <span class="task-date">{formatDate(task.due_date)}</span>
            {/if}
          </div>
        </div>
      {:else}
        <div class="no-tasks">No tasks for this day.</div>
      {/each}
    </div>
  </div>
</div>

<svelte:head>
  <link
    href="https://fonts.googleapis.com/icon?family=Material+Icons"
    rel="stylesheet"
  />
</svelte:head>

<style>
  .material-icons {
    font-family: "Material Icons";
    font-weight: normal;
    font-style: normal;
    font-size: 24px; /* Preferred icon size */
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

  .popup-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .popup-content {
    background-color: white;
    border-radius: 8px;
    max-width: 400px;
    width: 100%;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #e8ecee;
  }

  .popup-header h2 {
    font-size: 18px;
    font-weight: 600;
    color: #2c3e50;
    margin: 0;
  }

  .close-button {
    background: none;
    border: none;
    cursor: pointer;
    color: #7f8c8d;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 4px;
    border-radius: 50%;
    transition: background-color 0.3s;
  }

  .close-button:hover {
    background-color: #f0f0f0;
  }

  .task-list {
    padding: 16px 20px;
    overflow-y: auto;
    max-height: calc(80vh - 60px);
  }

  .task-item {
    display: flex;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid #e8ecee;
  }

  .task-item:last-child {
    border-bottom: none;
  }

  .checkbox-container {
    display: block;
    position: relative;
    padding-left: 30px;
    margin-right: 12px;
    cursor: pointer;
    font-size: 16px;
    user-select: none;
  }

  .checkbox-container input {
    position: absolute;
    opacity: 0;
    cursor: pointer;
    height: 0;
    width: 0;
  }

  .checkmark {
    position: absolute;
    top: 0;
    left: 0;
    height: 20px;
    width: 20px;
    background-color: #fff;
    border: 2px solid #7b68ee;
    border-radius: 4px;
  }

  .checkbox-container:hover input ~ .checkmark {
    background-color: #f0f0f0;
  }

  .checkbox-container input:checked ~ .checkmark {
    background-color: #7b68ee;
  }

  .checkmark:after {
    content: "";
    position: absolute;
    display: none;
  }

  .checkbox-container input:checked ~ .checkmark:after {
    display: block;
  }

  .checkbox-container .checkmark:after {
    left: 6px;
    top: 2px;
    width: 5px;
    height: 10px;
    border: solid white;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg);
  }

  .task-details {
    flex-grow: 1;
  }

  .task-name {
    font-size: 14px;
    color: #2c3e50;
    margin-right: 8px;
  }

  .task-date {
    font-size: 12px;
    color: #7f8c8d;
  }

  .no-tasks {
    color: #7f8c8d;
    font-size: 14px;
    text-align: center;
    padding: 16px 0;
  }
</style>
