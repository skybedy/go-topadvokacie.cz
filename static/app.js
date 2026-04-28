const actionSelect = document.querySelector("#action");
const secondDoc = document.querySelector("#second-doc");
const exampleSelect = document.querySelector("#example");
const analysisForm = document.querySelector("form");
const analysisLoading = document.querySelector("#analysis-loading");
const analysisPlaceholder = document.querySelector("#analysis-placeholder");
const analysisIdle = document.querySelector("[data-analysis-idle]");
const submitButton = document.querySelector('button[type="submit"]');
const submitLabel = document.querySelector("[data-submit-label]");
const analysisButtonSpinner = document.querySelector("[data-analysis-button-spinner]");
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
  const box = document.getElementById(`${status.id}_box`);
  const spinner = box?.querySelector("[data-upload-spinner]");
  status.textContent = message;
  if (box) {
    box.classList.remove(
      "hidden",
      "border-blue-200",
      "bg-blue-50",
      "text-blue-800",
      "border-emerald-200",
      "bg-emerald-50",
      "text-emerald-800",
      "border-red-200",
      "bg-red-50",
      "text-red-800",
      "border-slate-200",
      "bg-slate-50",
      "text-slate-700",
    );
    const toneClasses = {
      muted: ["border-slate-200", "bg-slate-50", "text-slate-700"],
      loading: ["border-blue-200", "bg-blue-50", "text-blue-800"],
      success: ["border-emerald-200", "bg-emerald-50", "text-emerald-800"],
      error: ["border-red-200", "bg-red-50", "text-red-800"],
    }[tone] || ["border-slate-200", "bg-slate-50", "text-slate-700"];
    box.classList.add(...toneClasses);
  }
  if (spinner) {
    spinner.classList.toggle("hidden", tone !== "loading");
  }
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

if (analysisForm && analysisLoading) {
  analysisForm.addEventListener("submit", (event) => {
    if (analysisForm.dataset.submitting === "true") return;

    if (activeUploads > 0) {
      event.preventDefault();
      return;
    }

    if (!analysisForm.checkValidity()) return;

    event.preventDefault();
    analysisIdle?.classList.add("hidden");
    analysisLoading.classList.remove("hidden");
    analysisPlaceholder?.classList.remove("border-slate-200", "bg-white");
    analysisPlaceholder?.classList.add("border-blue-200", "bg-blue-50");
    analysisButtonSpinner?.classList.remove("hidden");
    if (submitButton) {
      submitButton.disabled = true;
      submitButton.classList.add("opacity-60", "cursor-wait");
    }
    if (submitLabel) {
      submitLabel.innerText = "Čekám na odpověď modelu...";
    }
    analysisLoading.scrollIntoView({ behavior: "smooth", block: "center" });
    window.setTimeout(() => {
      analysisForm.dataset.submitting = "true";
      analysisForm.requestSubmit();
    }, 120);
  });
}

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
