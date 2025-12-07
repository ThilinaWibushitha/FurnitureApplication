<script>
  import { page } from "$app/stores";
  import Navbar from "$lib/components/Navbar.svelte";
  import { onMount } from "svelte";

  let user = null;
  let navLinks = [];

  onMount(() => {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      user = JSON.parse(userStr);
    }

    const allLinks = [
      { href: "/", label: "Dashboard", icon: "📊", roles: ["main_admin"] },
      // Sub-admins probably shouldn't see the main sales dashboard if it shows financial data
      // But if they have their own dashboard, we can route differently. For now, hide Main Dashboard.
      // Actually, request says "Main admin should be able to view daily and monthly sales".
      // "These other admins... View client profiles, Manage items...". Doesn't list sales.

      {
        href: "/clients",
        label: "Clients",
        icon: "👥",
        roles: ["main_admin", "admin"],
      },
      {
        href: "/items",
        label: "Items",
        icon: "🪑",
        roles: ["main_admin", "admin"],
      },
      {
        href: "/orders",
        label: "Orders",
        icon: "📦",
        roles: ["main_admin", "admin"],
      },
      // { href: '/payments', label: 'Payments', icon: '💳', roles: ['main_admin'] }, // Maybe integrated in Orders
      { href: "/admins", label: "Admins", icon: "🛡️", roles: ["main_admin"] },
    ];

    if (user) {
      navLinks = allLinks.filter((link) =>
        link.roles.includes(user.account_type),
      );
    } else {
      navLinks = [];
    }
  });

  $: showNavbar = $page.url.pathname !== "/login";
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
    font-family:
      "Inter",
      "Segoe UI",
      system-ui,
      -apple-system,
      sans-serif;
    background: #0b1020;
    color: #111827;
  }

  .layout {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background: radial-gradient(
        circle at top left,
        rgba(80, 201, 195, 0.12),
        transparent 45%
      ),
      radial-gradient(
        circle at top right,
        rgba(150, 222, 218, 0.1),
        transparent 50%
      ),
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
