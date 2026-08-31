# 🗺️ Your Guide to Understanding SPECTRE (After Documentation Improvements)

This guide helps you find exactly what you need and explains how the improved documentation works.

---

## 🚀 Quick Navigation

### **I'm brand new to SPECTRE**
Start here and read in this order:
1. [README.md](README.md) (main page) → "What is SPECTRE?" section
2. [Getting Started Guide](docs/GETTING_STARTED.md) → Do your first investigation
3. [Installation Guide](docs/INSTALLATION.md) → Set up SPECTRE properly
4. [Architecture Guide](docs/ARCHITECTURE.md) → Understand how it works

**Estimated time:** 30 minutes to understand everything

---

### **Something is broken / I got an error**
Go straight here:
1. [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)
2. Find your error message in the list
3. Follow the "How to fix" steps
4. If still stuck, check the related guide mentioned

**Estimated time:** 5-10 minutes to fix your issue

---

### **I want to use SPECTRE right now**
Quick path:
1. [Installation Guide](docs/INSTALLATION.md) → Install it
2. [Getting Started Guide](docs/GETTING_STARTED.md) → Step 1: Create a Case
3. Run these commands:
   ```bash
   spectre case new "My Investigation"
   spectre collect all example.com --case <CASE_ID>
   spectre chat --case <CASE_ID>
   ```
4. Ask questions in the chat!

**Estimated time:** 10 minutes

---

### **I want to understand the architecture**
Best docs:
1. [Architecture Guide](docs/ARCHITECTURE.md) → Simple explanation with analogies
2. [API Documentation](docs/API_DOCUMENTATION.md) → Technical details (if needed)

**Key insight:** Go (backend) and Python (AI/visualization) work together via JSON messages

---

### **I want to build a custom collector**
Head here:
1. [Plugin Development Guide](docs/PLUGIN_DEVELOPMENT.md) → Learn how
2. [API Documentation](docs/API_DOCUMENTATION.md) → Understand the bridge
3. Check `plugins/` folder for examples

---

### **I want to use SPECTRE on a server**
Read this:
1. [Deployment Guide](docs/DEPLOYMENT.md) → Server setup
2. [Troubleshooting](docs/TROUBLESHOOTING.md) → Platform-specific tips

---

## 📚 What Each Guide Covers

### 1️⃣ **README.md** (Main Project Page)
**What:** The overview and quick start
**Who:** Everyone
**Read when:** First time visiting, need inspiration, want quick commands

**Contains:**
- What is SPECTRE (plain English)
- 3 ways to use it (Quick scan, Terminal, Web)
- Key features table
- Common commands

**Quick start:**
```bash
# Option 1: Quick investigation
spectre investigate scanme.nmap.org

# Option 2: Create a case for detailed work
spectre case new "My Investigation"
```

---

### 2️⃣ **docs/GETTING_STARTED.md** (Your First Investigation)
**What:** Step-by-step walkthrough
**Who:** New users
**Read when:** You want to run your first real investigation

**Contains:**
- 6 simple steps from case creation to report export
- What each command does
- Where files are stored
- Quick reference of common tasks

**Main flow:**
1. Create case
2. Collect data (passive)
3. Collect data (active) — optional
4. Chat with AI
5. Visualize the graph
6. Export report

---

### 3️⃣ **docs/INSTALLATION.md** (How to Set Up)
**What:** Installation instructions for all operating systems
**Who:** Anyone setting up SPECTRE
**Read when:** First time installing

**Contains:**
- What you need (Python, Go, Git)
- Windows setup (1 command: `.\install.ps1`)
- Linux/macOS setup (manual commands)
- Verification steps
- Proxy/Tor setup for privacy

**Easiest path (Windows):**
```powershell
.\install.ps1
# Done! Ready to use SPECTRE
```

---

### 4️⃣ **docs/ARCHITECTURE.md** (How It Works)
**What:** How SPECTRE is built and how data flows
**Who:** Technical users, developers, curious minds
**Read when:** You want to understand the system design

**Contains:**
- Simple "Chef & Assistants" analogy
- What each part does (Go, Python)
- How data flows through the system
- Where everything is stored
- Safety features explained

**Simple version:**
- Go = Fast, organized coordinator
- Python = Smart specialists (visualization, AI)
- Data = Stays on your computer

---

### 5️⃣ **docs/README.md** (Documentation Index)
**What:** Guide to all other documentation
**Who:** Anyone looking for a specific guide
**Read when:** You need something but don't know which document to read

**Contains:**
- Organized sections (Getting Started, Understanding, Building, etc.)
- Quick reference table ("I want to... → Read this")
- One-line descriptions of each guide

**Most useful:**
The "Quick Reference" table that maps your needs to documents

---

### 6️⃣ **docs/TROUBLESHOOTING.md** (Fix Problems)
**What:** Common problems and solutions
**Who:** Anyone who got an error or needs help
**Read when:** Something isn't working

**Sections:**
- Installation Problems (Python not found, permissions, etc.)
- Running SPECTRE Problems (database locked, collector not found, etc.)
- Common Tasks (backup, proxy setup, reset)
- Platform-specific tips (Windows, Linux, macOS)
- Performance issues

**Format:** Each problem has:
- "What it means" (explanation)
- "How to fix" (step-by-step)

---

### 7️⃣ **docs/ARCHITECTURE.md** (Technical Deep Dive)
**What:** Internal architecture and design
**Who:** Developers, DevOps, curious engineers
**Read when:** You need to understand implementation details

**Contains:**
- Component breakdown
- Data flow diagrams
- Storage explanation (SQLite + evidence files)
- Configuration details

---

## 🎯 Common Scenarios & Where to Look

| Scenario | Read This | Time |
|----------|-----------|------|
| "I want to download and install SPECTRE" | [Installation](docs/INSTALLATION.md) | 5 min |
| "I installed it, now what?" | [Getting Started](docs/GETTING_STARTED.md) | 10 min |
| "I got an error: database is locked" | [Troubleshooting](docs/TROUBLESHOOTING.md) | 2 min |
| "I want to build a custom collector" | [Plugin Development](docs/PLUGIN_DEVELOPMENT.md) | 20 min |
| "I want to deploy SPECTRE on a server" | [Deployment](docs/DEPLOYMENT.md) | 30 min |
| "I want to understand how it works" | [Architecture](docs/ARCHITECTURE.md) | 10 min |
| "I want a quick one-liner investigation" | [README.md](README.md) main page | 1 min |
| "I forgot which command does what" | [Getting Started](docs/GETTING_STARTED.md) Quick Ref | 2 min |
| "I want to run AI analysis locally" | [LLM Integration](docs/llm.md) | 15 min |
| "I want to see all features" | [Features List](docs/features.md) | 10 min |

---

## 📖 Key Improvements You'll Notice

### 1. **Simpler Language**
- ❌ Before: "LLM Synthesis," "Vector DB Ingestion," "Subprocess Bridge"
- ✅ After: "AI Analysis," "Searching," "Python talking to Go"

### 2. **Clear Examples**
Every guide includes:
- Real command examples you can copy-paste
- What the output looks like
- What happens next

### 3. **Better Organization**
- Numbered steps (1, 2, 3...)
- Clear section headings
- Emoji markers for quick scanning

### 4. **Helpful Callouts**
- 💡 Tips — Additional helpful information
- ⚠️ Warnings — Important things to remember
- ✓ Success indicators — When you did it right
- Example outputs — See what success looks like

### 5. **Quick Reference Sections**
Every guide has a "Common Tasks" or "Quick Reference" section so you don't have to read everything

---

## 💡 Pro Tips for Using These Docs

1. **Use your browser's search (Ctrl+F)** to find commands or concepts
2. **Read in order** if you're new (README → Getting Started → Architecture)
3. **Keep Troubleshooting bookmarked** — You'll use it often
4. **Check "Quick Reference"** sections before reading whole guides
5. **Copy-paste commands** from the guides directly into your terminal

---

## 🎓 Learning Path

### Absolute Beginner
```
Day 1: README → GETTING_STARTED (Steps 1-3)
Day 2: Complete GETTING_STARTED
Day 3: Try more investigations
Day 4: Read ARCHITECTURE
```

### Technical User
```
Day 1: README → Installation
Day 2: GETTING_STARTED (fast read)
Day 3: ARCHITECTURE
Day 4: Plugin Development (if interested)
```

### Developer / Sysadmin
```
Day 1: README → Installation
Day 2: ARCHITECTURE → API_DOCUMENTATION
Day 3: Deployment or Plugin Development
```

---

## ✨ The Main Changes You'll See

### Before vs. After

**BEFORE:** Complex technical explanation
```
"SPECTRE utilizes a hybrid architecture where Go orchestrates 
CLI framework (cobra), concurrent collection, and SQLite storage, 
while the intelligence layer manages AI analysis via subprocess bridge..."
```

**AFTER:** Simple explanation
```
"Go handles the coordination (CLI, collectors, storage).
Python handles the smart analysis (graphs, AI).
They talk to each other via simple messages."
```

---

## 🎯 Remember

**The documentation is now written for YOU — whether you're:**
- ✅ A beginner who's never heard of OSINT
- ✅ A technical person looking for implementation details
- ✅ A developer wanting to extend SPECTRE
- ✅ Someone who just wants to investigate something quickly

**Everything is clearly organized, simplified, and easy to find.**

---

## 🚀 Ready to Start?

1. Read [README.md](README.md) (5 min)
2. Follow [Getting Started](docs/GETTING_STARTED.md) (10 min)
3. Try it out! 

**That's it. You're ready to investigate.**

---

For questions about specific topics, check the "📚 What Each Guide Covers" section above! 🎓
