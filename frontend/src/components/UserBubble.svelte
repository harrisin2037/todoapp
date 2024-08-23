<script>
  export let user;
</script>

<div class="user-bubble" title={user}>
  {user.charAt(0).toUpperCase()}
</div>

<style>
  .user-bubble {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    width: 30px;
    height: 30px;
    border-radius: 50%;
    background-color: #007bff;
    color: white;
    font-weight: bold;
    cursor: pointer;
    margin-right: 5px;
  }
</style>
