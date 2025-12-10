<script>
  import { page } from '$app/stores';
  import { onMount } from 'svelte';

  export let navLinks = [];
  
  let user = null;
  
  onMount(() => {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      user = JSON.parse(userStr);
    }
  });
  
  function handleLogout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
  }
  
  function getInitials(name) {
    if (!name) return 'AD';
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
  }
  
  function getRoleName(type) {
    if (type === 'main_admin') return 'Main Admin';
    if (type === 'admin') return 'Admin';
    return 'User';
  }
</script>

<nav class="navbar">
  <div class="nav-brand">
    <div class="logo">F</div>
    <div class="brand-text">
      <span class="title">Furniture Admin</span>
      <span class="subtitle">Control Center</span>
    </div>
  </div>

  <div class="nav-links">
    {#each navLinks as link}
      {#if link.disabled}
        <span class="nav-link disabled" title="Access restricted">
          <span>{link.label}</span>
        </span>
      {:else}
        <a
          href={link.href}
          class:active={
            $page.url.pathname === link.href ||
            ($page.url.pathname.startsWith(link.href) && link.href !== '/')
          }
        >
          <span>{link.label}</span>
        </a>
      {/if}
    {/each}
  </div>

  <div class="nav-actions">
    <div class="profile-wrapper">
      <button class="profile">
        <span class="avatar" aria-hidden="true">{user ? getInitials(user.full_name || user.username) : 'AD'}</span>
        <span class="meta">
          <span class="name">{user?.full_name || user?.username || 'Admin'}</span>
          <span class="role">{getRoleName(user?.account_type)}</span>
        </span>
      </button>
      <button class="logout-btn" on:click={handleLogout} title="Logout">
        Logout
      </button>
    </div>
  </div>
</nav>

<style>
  .navbar {
    position: sticky;
    top: 0;
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 2rem;
    background: rgba(21, 24, 44, 0.92);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    color: #f7f9fc;
  }

  .nav-brand {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .logo {
    width: 42px;
    height: 42px;
    border-radius: 14px;
    display: grid;
    place-items: center;
    background: linear-gradient(135deg, #f8d417, #ff7e5f);
    font-size: 1.25rem;
    box-shadow: 0 12px 24px rgba(255, 126, 95, 0.25);
  }

  .brand-text {
    display: flex;
    flex-direction: column;
    line-height: 1.1;
  }

  .title {
    font-weight: 700;
    font-size: 1rem;
    letter-spacing: 0.5px;
  }

  .subtitle {
    font-size: 0.75rem;
    opacity: 0.65;
    letter-spacing: 0.4px;
  }

  .nav-links {
    display: flex;
    align-items: center;
    gap: 1rem;
    background: rgba(255, 255, 255, 0.04);
    padding: 0.35rem;
    border-radius: 999px;
  }

  .nav-links a,
  .nav-links .nav-link {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.55rem 1.1rem;
    border-radius: 999px;
    font-size: 0.9rem;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.7);
    text-decoration: none;
    transition: transform 0.2s ease, background 0.2s ease, color 0.2s ease;
  }

  .nav-links a:hover {
    color: #fff;
    transform: translateY(-1px);
    background: rgba(255, 255, 255, 0.08);
  }

  .nav-links a.active {
    background: linear-gradient(135deg, #50c9c3, #96deda);
    color: #14162b;
    box-shadow: 0 12px 24px rgba(80, 201, 195, 0.25);
  }
  
  .nav-links .nav-link.disabled {
    opacity: 0.4;
    cursor: not-allowed;
    pointer-events: none;
  }

  .icon {
    font-size: 1rem;
    line-height: 1;
  }

  .nav-actions {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .profile {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.4rem 0.9rem;
    border-radius: 999px;
    border: none;
    background: rgba(255, 255, 255, 0.08);
    color: inherit;
    cursor: pointer;
    transition: background 0.2s ease, transform 0.2s ease;
  }

  .profile:hover {
    background: rgba(255, 255, 255, 0.12);
    transform: translateY(-1px);
  }

  .avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea, #764ba2);
    display: grid;
    place-items: center;
    font-weight: 700;
    letter-spacing: 0.5px;
  }

  .meta {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    font-size: 0.75rem;
    line-height: 1.2;
  }

  .name {
    font-weight: 600;
    font-size: 0.8rem;
  }

  .role {
    opacity: 0.6;
  }
  
  .profile-wrapper {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  
  .logout-btn {
    padding: 0.5rem 1rem;
    background: rgba(239, 68, 68, 0.15);
    color: #f87171;
    border: 1px solid rgba(239, 68, 68, 0.3);
    border-radius: 8px;
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s ease;
  }
  
  .logout-btn:hover {
    background: rgba(239, 68, 68, 0.25);
  }

  @media (max-width: 900px) {
    .navbar {
      flex-wrap: wrap;
      gap: 1rem;
    }

    .nav-links {
      order: 3;
      width: 100%;
      justify-content: space-between;
      padding: 0.3rem;
    }

    .nav-links a {
      flex: 1;
      justify-content: center;
    }

    .nav-actions {
      margin-left: auto;
    }
  }

  @media (max-width: 620px) {
    .navbar {
      padding: 0.75rem 1rem;
    }

    .nav-brand {
      width: 100%;
      justify-content: space-between;
    }

    .nav-actions {
      width: 100%;
      justify-content: flex-end;
    }

    .nav-links {
      gap: 0.25rem;
    }

    .nav-links a {
      padding: 0.5rem 0.7rem;
      font-size: 0.8rem;
    }
  }
</style>
