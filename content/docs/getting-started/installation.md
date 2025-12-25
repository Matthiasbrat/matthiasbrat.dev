---
title: "Installation"
description: "How to install and set up the project"
order: 2
---

This guide will walk you through installing and setting up the project.

## Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.21 or later
- Git

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/example/site.git
cd site
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Build the Project

```bash
go build ./cmd/site
```

## Verify Installation

Run the following command to verify everything is working:

```bash
./site --version
```

You should see output like:

```
site version 1.0.0
```

> [!TIP]
> If you encounter any issues, check our troubleshooting guide or open an issue on GitHub.

## Next Steps

Now that you have the project installed, continue to the Configuration guide to set up your site.
