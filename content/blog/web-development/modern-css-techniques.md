---
title: "Modern CSS Techniques You Should Know"
description: "CSS Grid, Flexbox, custom properties, and more"
date: 2025-02-01
---

CSS has evolved dramatically. Here are techniques every modern web developer should master.

## CSS Grid

The most powerful layout system CSS has ever had:

```css
.container {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
}
```

## CSS Custom Properties (Variables)

```css
:root {
    --primary-color: #3498db;
    --spacing-unit: 8px;
}

.button {
    background: var(--primary-color);
    padding: calc(var(--spacing-unit) * 2);
}
```

> [!TIP]
> Custom properties cascade and can be overridden in different scopes.

## Container Queries

Style based on parent container size, not viewport:

```css
.card-container {
    container-type: inline-size;
}

@container (min-width: 400px) {
    .card {
        display: flex;
    }
}
```

## Modern Selectors

```css
/* Has selector */
.card:has(img) {
    padding: 0;
}

/* Where selector (zero specificity) */
:where(.btn, .link) {
    cursor: pointer;
}

/* Is selector (takes highest specificity) */
:is(h1, h2, h3) {
    line-height: 1.2;
}
```

> [!NOTE]
> Check browser support before using newer features in production.
