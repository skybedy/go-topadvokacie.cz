const actionSelect = document.querySelector("#action");
const secondDoc = document.querySelector("#second-doc");
const exampleSelect = document.querySelector("#example");
const submitButton = document.querySelector('button[type="submit"]');
let activeUploads = 0;

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

function setUploadStatus(status, message, tone = "muted") {
  if (!status) return;
  status.textContent = message;
  status.classList.remove("text-slate-500", "text-blue-700", "text-emerald-700", "text-red-700");
  const toneClass = {
    muted: "text-slate-500",
    loading: "text-blue-700",
    success: "text-emerald-700",
    error: "text-red-700",
  }[tone] || "text-slate-500";
  status.classList.add(toneClass);
}

function setUploadBusy(isBusy) {
  activeUploads += isBusy ? 1 : -1;
  activeUploads = Math.max(activeUploads, 0);
  if (submitButton) {
    submitButton.disabled = activeUploads > 0;
    submitButton.classList.toggle("opacity-60", activeUploads > 0);
    submitButton.classList.toggle("cursor-wait", activeUploads > 0);
  }
}

document.querySelectorAll("[data-upload-input]").forEach((input) => {
  input.addEventListener("change", async () => {
    const file = input.files && input.files[0];
    if (!file) return;

    const textarea = document.getElementById(input.dataset.targetTextarea);
    const status = document.getElementById(input.dataset.statusTarget);
    if (!textarea) return;

    const formData = new FormData();
    formData.append("file", file);

    input.disabled = true;
    textarea.disabled = true;
    setUploadBusy(true);
    setUploadStatus(status, `Načítám ${file.name} do textového pole...`, "loading");

    try {
      const response = await fetch("/upload-text", {
        method: "POST",
        body: formData,
      });
      const payload = await response.json();
      if (!response.ok || payload.error) {
        throw new Error(payload.error || "Soubor se nepodařilo načíst.");
      }
      textarea.value = payload.text || "";
      input.value = "";
      textarea.focus();
      setUploadStatus(status, `Soubor ${file.name} je načtený a viditelný v poli. Teď můžete spustit analýzu.`, "success");
    } catch (error) {
      setUploadStatus(status, error.message || "Soubor se nepodařilo načíst.", "error");
    } finally {
      input.disabled = false;
      textarea.disabled = false;
      setUploadBusy(false);
    }
  });
});

document.querySelectorAll("[data-copy-button]").forEach((button) => {
  button.addEventListener("click", async () => {
    const card = button.closest("[data-copy-card]");
    const title = card?.querySelector("h3")?.innerText.trim() || "";
    const content = card?.querySelector("[data-copy-content]")?.innerText.trim() || "";
    const text = [title, content].filter(Boolean).join("\n\n");
    if (!text) return;

    try {
      await navigator.clipboard.writeText(text);
      const original = button.innerText;
      button.innerText = "Zkopírováno";
      window.setTimeout(() => {
        button.innerText = original;
      }, 1400);
    } catch {
      button.innerText = "Nelze kopírovat";
      window.setTimeout(() => {
        button.innerText = "Kopírovat";
      }, 1400);
    }
  });
});
