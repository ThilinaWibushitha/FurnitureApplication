<script>
  import { page } from '$app/stores';

  export let navLinks = [];
</script>

<nav class="navbar">
  <div class="nav-brand">
    <div class="logo">🛋️</div>
    <div class="brand-text">
      <span class="title">Furniture Admin</span>
      <span class="subtitle">Control Center</span>
    </div>
  </div>

  <div class="nav-links">
    {#each navLinks as link}
      <a
        href={link.href}
        class:active={
          $page.url.pathname === link.href ||
          ($page.url.pathname.startsWith(link.href) && link.href !== '/')
        }
      >
        {#if link.icon}
          <span class="icon" aria-hidden="true">{link.icon}</span>
        {/if}
        <span>{link.label}</span>
      </a>
    {/each}
  </div>

  <div class="nav-actions">
    <button class="profile">
      <span class="avatar" aria-hidden="true">AD</span>
      <span class="meta">
        <span class="name">Admin</span>
        <span class="role">Administrator</span>
      </span>
    </button>
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

  .nav-links a {
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
