# 📘 RAG Doc Converter: Markdown to HTML for Dialogflow CX

Questo strumento è un'applicazione **Go** con interfaccia testuale (**TUI**) progettata per automatizzare la preparazione dei documenti per i **Data Store di Dialogflow CX (Vertex AI Search)**. 

Trasforma file Markdown in HTML "puro", ottimizzando la capacità di indicizzazione dell'agente IA e riducendo le allucinazioni durante la fase di Retrieval (RAG).

---

## 🚀 Perché questo tool?

Dialogflow CX e i motori di ricerca vettoriali di Google Cloud elaborano la struttura semantica dei documenti in modo più preciso quando sono in **HTML**. 

* **Gerarchia Chiara:** I tag `<h1>`, `<h2>`, `<ul>` forniscono segnali forti rispetto ai simboli Markdown.
* **Purezza dei Dati:** Questo converter elimina ogni stile CSS o script, lasciando solo il contenuto strutturato che serve al LLM.
* **Workflow Agile:** Gestisce la pulizia delle vecchie versioni e apre automaticamente i risultati per un caricamento rapido su Google Cloud Storage.

---

## 🛠️ Requisiti

* **Go** (versione 1.20 o superiore)
* Un terminale (Bash, Zsh, PowerShell o CMD)

---

## 📦 Installazione

1.  Clona il repository:
    ```bash
    git clone [git@github.com:gimox/mdToHtml.git](git@github.com:gimox/mdToHtml.git)
    cd mdToHtml
    ```

2.  Scarica le dipendenze:
    ```bash
    go mod tidy
    ```

---

## 📂 Struttura del Progetto

| Cartella / File | Descrizione |
| :--- | :--- |
| `input_md/` | Inserisci qui i tuoi file `.md` sorgente. |
| `output_html/` | Qui verranno generati i file `.html` pronti per il RAG. |
| `main.go` | Il cuore dell'applicazione (Logica + TUI). |
| `.gitignore` | Configurato per non caricare i tuoi documenti privati su GitHub. |

---

## 🎮 Utilizzo

L'applicazione utilizza un'interfaccia interattiva nel terminale.

1.  **Avvio:** Esegui l'app con `go run main.go`.
2.  **Navigazione:** Usa le **frecce direzionali** per muoverti nel menu.
3.  **Selezione:** Premi **Invio** per scegliere un'opzione o un file specifico.
4.  **Indietro:** Premi **ESC** per tornare al menu precedente (es. dal picker dei file al menu principale).
5.  **Pulizia:** Prima di ogni conversione, l'app ti chiederà se vuoi svuotare la cartella di output (premere **Y** per pulire, **N** per mantenere i file esistenti).
6.  **Uscita:** Premi **Q** o `Ctrl+C`.

---

## ✨ Funzionalità Avanzate

* **Conversione Bulk:** Processa centinaia di file in pochi millisecondi.
* **Feedback Visivo:** Spinner animato e progresso file per file.
* **Auto-Open:** Al termine del processo, la cartella dei risultati si apre automaticamente (supporta Windows, macOS e Linux).
* **Gestione Filtri:** Supporta la ricerca rapida dei file nella lista tramite digitazione.

---

## 🛡️ Sicurezza e Privacy

Il file `.gitignore` incluso è configurato per **escludere tutti i file** contenuti nelle cartelle `input_md/` e `output_html/`. 

> **IMPORTANTE:** Non caricare mai documenti riservati in repository pubblici. Questo tool è pensato per gestire i dati localmente prima dell'upload sicuro su Google Cloud Storage.

---

## 📜 Licenza
Distribuito sotto licenza MIT.