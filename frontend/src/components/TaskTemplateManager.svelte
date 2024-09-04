<script>
  import { onMount } from "svelte";
  import { API_BASE_URL } from "../config.js";

  let showCreateModal = false;
  let templates = [];
  let newTemplate = {
    name: "",
    description: "",
    fields: [],
  };
  let editingTemplate = null;

  function showCreateTemplateModal() {
    showCreateModal = true;
    newTemplate = { name: "", description: "", fields: [] };
  }

  function hideCreateModal() {
    showCreateModal = false;
  }

  function addField(template) {
    template.fields = [
      ...template.fields,
      { name: "", type: "text", default: "" },
    ];
  }

  function removeField(template, index) {
    template.fields = template.fields.filter((_, i) => i !== index);
  }

  onMount(fetchTemplates);

  async function fetchTemplates() {
    try {
      const response = await fetch(`${API_BASE_URL}/task-templates`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      if (!response.ok) {
        throw new Error("Failed to fetch templates");
      }
      templates = await response.json();
    } catch (error) {
      console.error("Error fetching templates:", error);
    }
  }

  async function createTemplate() {
    try {
      const response = await fetch(`${API_BASE_URL}/task-templates`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify(newTemplate),
      });

      if (!response.ok) {
        throw new Error("Failed to create template");
      }

      await fetchTemplates();
      newTemplate = {
        name: "",
        description: "",
        fields: [],
      };
    } catch (error) {
      console.error("Error creating template:", error);
    }
  }

  function editTemplate(template) {
    editingTemplate = { ...template };
  }

  async function updateTemplate() {
    try {
      const response = await fetch(
        `${API_BASE_URL}/task-templates/${editingTemplate.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify(editingTemplate),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to update template");
      }

      await fetchTemplates();
      cancelEditing();
    } catch (error) {
      console.error("Error updating template:", error);
    }
  }

  async function deleteTemplate(id) {
    if (confirm("Are you sure you want to delete this template?")) {
      try {
        const response = await fetch(`${API_BASE_URL}/task-templates/${id}`, {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        });

        if (!response.ok) {
          throw new Error("Failed to delete template");
        }

        await fetchTemplates();
      } catch (error) {
        console.error("Error deleting template:", error);
      }
    }
  }

  function cancelEditing() {
    editingTemplate = null;
  }
</script>

<main>
  <div class="template-management-header">
    <h2>Template Management</h2>
    <button class="btn btn-primary" on:click={showCreateTemplateModal}>
      <span class="plus-icon">+</span> Create New Template
    </button>
  </div>

  <section class="template-list">
    <h3>Existing Templates</h3>
    {#if templates.length === 0}
      <p class="no-templates">No templates found.</p>
    {:else}
      <ul>
        {#each templates as template}
          <li>
            <h4>{template.name}</h4>
            <p>{template.description}</p>
            <div class="template-actions">
              <button
                class="btn btn-secondary"
                on:click={() => editTemplate(template)}>Edit</button
              >
              <button
                class="btn btn-secondary"
                on:click={() => deleteTemplate(template.id)}>Delete</button
              >
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </section>

  {#if showCreateModal}
    <div class="modal">
      <div class="modal-content">
        <h3>Create New Template</h3>
        <form on:submit|preventDefault={createTemplate}>
          <input
            bind:value={newTemplate.name}
            placeholder="Template Name"
            required
          />
          <textarea
            bind:value={newTemplate.description}
            placeholder="Template Description"
          ></textarea>

          <h4>Fields</h4>
          {#each newTemplate.fields as field, index}
            <div class="field">
              <input
                bind:value={field.name}
                placeholder="Field Name"
                required
              />
              <select bind:value={field.type}>
                <option value="text">Text</option>
                <option value="number">Number</option>
                <option value="checkbox">Checkbox</option>
                <option value="list">List</option>
              </select>
              <input bind:value={field.default} placeholder="Default Value" />
              <button
                type="button"
                class="btn btn-secondary"
                on:click={() => removeField(newTemplate, index)}>Remove</button
              >
            </div>
          {/each}
          <button
            type="button"
            class="btn btn-secondary"
            on:click={() => addField(newTemplate)}>Add Field</button
          >

          <div class="modal-actions">
            <button type="submit" class="btn btn-primary"
              >Create Template</button
            >
            <button
              type="button"
              class="btn btn-secondary"
              on:click={hideCreateModal}>Cancel</button
            >
          </div>
        </form>
      </div>
    </div>
  {/if}

  {#if editingTemplate}
    <div class="modal">
      <div class="modal-content">
        <h3>Edit Template: {editingTemplate.name}</h3>
        <form on:submit|preventDefault={updateTemplate}>
          <input
            bind:value={editingTemplate.name}
            placeholder="Template Name"
            required
          />
          <textarea
            bind:value={editingTemplate.description}
            placeholder="Template Description"
          ></textarea>

          <h4>Fields</h4>
          {#each editingTemplate.fields as field, index}
            <div class="field">
              <input
                bind:value={field.name}
                placeholder="Field Name"
                required
              />
              <select bind:value={field.type}>
                <option value="text">Text</option>
                <option value="number">Number</option>
                <option value="checkbox">Checkbox</option>
                <option value="list">List</option>
              </select>
              <input bind:value={field.default} placeholder="Default Value" />
              <button
                type="button"
                class="btn btn-secondary"
                on:click={() => removeField(editingTemplate, index)}
                >Remove</button
              >
            </div>
          {/each}
          <button
            type="button"
            class="btn btn-secondary"
            on:click={() => addField(editingTemplate)}>Add Field</button
          >

          <div class="modal-actions">
            <button type="submit" class="btn btn-primary">Save Changes</button>
            <button
              type="button"
              class="btn btn-secondary"
              on:click={cancelEditing}>Cancel</button
            >
          </div>
        </form>
      </div>
    </div>
  {/if}
</main>

<style>
  :root {
    --primary-color: #7b68ee;
    --primary-hover-color: #6c5ce7;
    --background-color: #f7f9fb;
    --white: #ffffff;
    --border-color: #e8ecee;
    --text-color: #2c3e50;
    --secondary-bg: #e8ecee;
    --secondary-hover: #d1d8e0;
  }

  main {
    font-family: Arial, sans-serif;
    background-color: var(--background-color);
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  h2 {
    font-size: 1.2rem;
    font-weight: 600;
    color: var(--text-color);
    margin: 0;
  }

  h3 {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--text-color);
    margin-top: 0;
    margin-bottom: 1rem;
  }

  .template-management-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    background-color: var(--white);
    border-bottom: 1px solid var(--border-color);
  }

  .btn {
    padding: 0.7rem 1.5rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    transition: background-color 0.3s ease;
  }

  .btn-primary {
    background-color: var(--primary-color);
    color: var(--white);
  }

  .btn-primary:hover {
    background-color: var(--primary-hover-color);
  }

  .btn-secondary {
    background-color: var(--secondary-bg);
    color: var(--text-color);
  }

  .btn-secondary:hover {
    background-color: var(--secondary-hover);
  }

  .template-list {
    flex-grow: 1;
    overflow-y: auto;
    padding: 1rem;
  }

  .template-list ul {
    list-style-type: none;
    padding: 0;
  }

  .template-list li {
    background-color: var(--white);
    padding: 1rem;
    margin-bottom: 1rem;
    border-radius: 4px;
    border: 1px solid var(--border-color);
  }

  .template-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 1rem;
  }

  form {
    display: flex;
    flex-direction: column;
  }

  input,
  textarea,
  select {
    margin-bottom: 1rem;
    padding: 0.7rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    font-size: 1rem;
  }

  .field {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .modal-content {
    background-color: var(--white);
    padding: 2rem;
    border-radius: 8px;
    width: 90%;
    max-width: 500px;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1rem;
  }

  .no-templates {
    text-align: center;
    color: var(--text-color);
    font-style: italic;
    padding: 1rem;
  }

  @media (max-width: 600px) {
    .field {
      flex-direction: column;
      align-items: stretch;
    }

    .btn {
      width: 100%;
    }
  }

  .field {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    align-items: center;
    margin-bottom: 1rem;
  }

  .field input,
  .field select {
    flex: 1;
    min-width: 120px;
    margin-bottom: 0;
  }

  .field button {
    flex-shrink: 0;
  }

  h4 {
    font-size: 1rem;
    font-weight: 600;
    color: var(--text-color);
    margin-top: 1rem;
    margin-bottom: 0.5rem;
  }

  .modal-content {
    max-height: 90vh;
    overflow-y: auto;
  }

  @media (max-width: 600px) {
    .field {
      flex-direction: column;
      align-items: stretch;
    }

    .field input,
    .field select,
    .field button {
      width: 100%;
    }
  }
</style>
