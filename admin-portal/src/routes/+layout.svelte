<script>
  import { page } from '$app/stores';
  import Navbar from '$lib/components/Navbar.svelte';

  const navLinks = [
    { href: '/', label: 'Dashboard', icon: '📊' },
    { href: '/clients', label: 'Clients', icon: '👥' },
    { href: '/items', label: 'Items', icon: '🪑' }
  ];

  $: showNavbar = $page.url.pathname !== '/login';
</script>

<div class="layout">
  {#if showNavbar}
    <Navbar {navLinks} />
  {/if}
  <main class:with-navbar={showNavbar}>
    <slot />
  </main>
</div>

<style>
  :global(body) {
    margin: 0;
    font-family: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
    background: #0b1020;
    color: #111827;
  }

  .layout {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background: radial-gradient(circle at top left, rgba(80, 201, 195, 0.12), transparent 45%),
      radial-gradient(circle at top right, rgba(150, 222, 218, 0.1), transparent 50%),
      #0b1020;
  }

  main {
    flex: 1;
    width: 100%;
    padding: 0;
  }

  main.with-navbar {
    padding-top: 1rem;
  }

  main :global(.container) {
    backdrop-filter: none;
  }
</style>
