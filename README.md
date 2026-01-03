# LinkedIn Automation

This repository contains a **Go-based proof of concept (PoC)** for LinkedIn browser automation using a real Chrome session powered by the **Rod** library.

## 📌 Project Scope

This PoC focuses on:

- Persistent browser session handling
- Login flow with manual checkpoint / CAPTCHA support
- LinkedIn People Search navigation
- Parsing visible profile URLs from search results
- Pagination with controlled limits
- SQLite-based state tracking
- Read-only profile visits with human-like behavior
---

## 🕵️ Stealth Architecture (Summary)

This project includes a **lightweight, modular stealth layer** designed to normalize the browser environment before any LinkedIn interaction occurs.

**Key characteristics:**
- One **stable human-like identity per session** (User-Agent + viewport)
- **Minimal, high-impact techniques** instead of aggressive fingerprint spoofing
- All stealth logic is **centralized and extensible**
- JavaScript overrides are applied **before page scripts execute**

**Implemented measures:**
- Masks automation flags (`navigator.webdriver`)
- Uses a realistic Chrome-on-Windows User-Agent
- Applies realistic viewport sizing with session-level consistency
- Ensures pages are fully loaded before interaction

**Design goal:**  
Remove common automation fingerprints while keeping the system reliable, explainable, and easy to extend.

---

## ⚙️ How It Works

### 1. Browser & Login
- Launches a real Chrome instance with a persistent profile
- Reuses login sessions across runs
- Allows manual completion of LinkedIn checkpoints if triggered

### 2. Search & Parse
- Builds LinkedIn People Search URLs
- Navigates result pages
- Extracts visible profile URLs
- Stores unique profiles in SQLite

### 3. Profile Visit Loop
- Retrieves unvisited profiles from state store
- Opens profiles in new tabs
- Scrolls and pauses to simulate reading
- Marks profiles as visited

---

## 🗃 State Management

The application uses **SQLite (pure Go)** to track profile state across runs.

Stored information includes:
- Profile URL
- Visit status
- Creation timestamp
- Last action timestamp

This prevents duplicate processing and enables controlled iteration.

---

## 🔐 Environment Configuration

Create a `.env` file using the provided template.

### `.env.example`

```env
# LinkedIn login credentials
LINKEDIN_EMAIL=your_email@example.com
LINKEDIN_PASSWORD=your_password
```
## Running the Project
```
go mod tidy

go run cmd/bot/main.go
```
---


