<script>
  import Splash from './components/Splash.svelte'
  import Home from './components/Home.svelte'
  import Backup from './components/Backup.svelte'
  import Restore from './components/Restore.svelte'
  import Cleaner from './components/Cleaner.svelte'
  import Wiper from './components/Wiper.svelte'

  let showSplash = $state(true)
  let currentView = $state('home') // 'home', 'backup', 'restore', 'cleaner', 'wiper'

  function handleSplashComplete() {
    showSplash = false
  }

  function navigateTo(view) {
    currentView = view
  }

  function goHome() {
    currentView = 'home'
  }
</script>

{#if showSplash}
  <Splash onComplete={handleSplashComplete} />
{/if}

<main class:hidden={showSplash}>
  {#if currentView === 'home'}
    <Home onNavigate={navigateTo} />
  {:else if currentView === 'backup'}
    <Backup onBack={goHome} />
  {:else if currentView === 'restore'}
    <Restore onBack={goHome} />
  {:else if currentView === 'cleaner'}
    <Cleaner onBack={goHome} />
  {:else if currentView === 'wiper'}
    <Wiper onBack={goHome} />
  {/if}
</main>

<style>
  main {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: column;
    opacity: 1;
    transition: opacity 0.3s ease-in;
  }

  main.hidden {
    opacity: 0;
  }
</style>
