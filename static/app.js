const actionSelect = document.querySelector("#action");
const secondDoc = document.querySelector("#second-doc");
const exampleSelect = document.querySelector("#example");

function syncSecondDoc() {
  if (!actionSelect || !secondDoc) return;
  const selected = actionSelect.options[actionSelect.selectedIndex];
  const needsSecond = selected && selected.dataset.second === "true";
  secondDoc.classList.toggle("hidden", !needsSecond);
}

if (actionSelect) {
  actionSelect.addEventListener("change", syncSecondDoc);
  syncSecondDoc();
}

if (exampleSelect) {
  exampleSelect.addEventListener("change", () => {
    if (exampleSelect.value) {
      window.location.href = `/?example=${encodeURIComponent(exampleSelect.value)}`;
    }
  });
}
