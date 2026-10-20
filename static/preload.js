(() => {
  if (navigator.connection?.saveData) return;
  const seen = new Set();
  const preload = (link) => {
    if (!link || link.target || link.hasAttribute("download")) return;
    const url = new URL(link.href, location.href);
    if (url.origin !== location.origin) return;
    if (url.pathname === location.pathname && url.search === location.search)
      return;
    if (seen.has(url.href)) return;
    seen.add(url.href);
    const prefetch = document.createElement("link");
    prefetch.rel = "prefetch";
    prefetch.href = url.href;
    document.head.append(prefetch);
  };
  const find = (event) => event.target.closest?.("a[href]");
  for (const type of ["pointerover", "pointerdown", "focusin"]) {
    document.addEventListener(type, (event) => preload(find(event)), {
      passive: true,
    });
  }
})();
