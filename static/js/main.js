// Modern Site JavaScript
(() => {
  "use strict";

  // ========================================
  // Login Modal
  // ========================================
  const LoginModal = {
    modal: null,
    providersContainer: null,
    closeBtn: null,
    providers: null,
    redirectUrl: null,

    init() {
      this.modal = document.querySelector(".login-modal");
      this.providersContainer = document.querySelector(".login-providers");
      this.closeBtn = document.querySelector(".login-close");

      if (!this.modal) return;

      // Header sign in button
      const loginToggle = document.querySelector(".login-toggle");
      if (loginToggle) {
        loginToggle.addEventListener("click", () => this.open());
      }

      // Close button
      if (this.closeBtn) {
        this.closeBtn.addEventListener("click", () => this.close());
      }

      // Click outside to close
      this.modal.addEventListener("click", (e) => {
        if (e.target === this.modal) {
          this.close();
        }
      });

      // Escape to close
      document.addEventListener("keydown", (e) => {
        if (e.key === "Escape" && this.modal.classList.contains("active")) {
          this.close();
        }
      });
    },

    async open(redirectUrl = null) {
      if (!this.modal) return;

      this.redirectUrl = redirectUrl || window.location.pathname;
      this.modal.classList.add("active");

      if (!this.providers) {
        await this.loadProviders();
      } else {
        this.renderProviders();
      }
    },

    close() {
      if (!this.modal) return;
      this.modal.classList.remove("active");
    },

    async loadProviders() {
      this.providersContainer.innerHTML =
        '<div class="login-loading">Loading...</div>';

      try {
        const response = await fetch("/api/auth/providers");
        if (!response.ok) throw new Error("Failed to load providers");

        const data = await response.json();
        this.providers = data.providers || [];
        this.renderProviders();
      } catch (err) {
        console.error("Failed to load auth providers:", err);
        this.providersContainer.innerHTML =
          '<div class="login-loading">Failed to load sign in options</div>';
      }
    },

    renderProviders() {
      if (!this.providers || this.providers.length === 0) {
        this.providersContainer.innerHTML =
          '<div class="login-loading">No sign in options available</div>';
        return;
      }

      const icons = {
        google: `<svg width="20" height="20" viewBox="0 0 24 24"><path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/><path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/><path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/><path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/></svg>`,
        github: `<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>`,
      };

      this.providersContainer.innerHTML = this.providers
        .map((provider) => {
          const icon = icons[provider.id] || "";
          const redirectParam = encodeURIComponent(this.redirectUrl);
          return `<a href="/auth/${provider.id}?redirect=${redirectParam}" class="login-provider-btn">
                    ${icon}
                    <span>Continue with ${provider.name}</span>
                </a>`;
        })
        .join("");
    },
  };

  // ========================================
  // Page Components
  // ========================================
  const PageComponents = {
    init() {
      this.initMobileMenu();
      this.initSearch();
      this.initCollectionSort();
      this.initReactions();
      this.initComments();
    },

    initMobileMenu() {
      const menuToggle = document.querySelector(".menu-toggle");
      const sidebarOverlay = document.querySelector(".sidebar-overlay");
      const sidebar = document.querySelector(".sidebar");

      if (menuToggle && !menuToggle.dataset.initialized) {
        menuToggle.dataset.initialized = "true";
        menuToggle.addEventListener("click", () => {
          // On pages with a sidebar (post pages), toggle sidebar
          // On pages without sidebar (home, docs landing), toggle nav
          if (sidebar) {
            document.body.classList.toggle("sidebar-open");
            const isOpen = document.body.classList.contains("sidebar-open");
            menuToggle.setAttribute("aria-expanded", isOpen);
          } else {
            document.body.classList.toggle("nav-open");
            const isOpen = document.body.classList.contains("nav-open");
            menuToggle.setAttribute("aria-expanded", isOpen);
          }
        });
      }

      if (sidebarOverlay && !sidebarOverlay.dataset.initialized) {
        sidebarOverlay.dataset.initialized = "true";
        sidebarOverlay.addEventListener("click", () => {
          document.body.classList.remove("sidebar-open", "nav-open");
          if (menuToggle) {
            menuToggle.setAttribute("aria-expanded", "false");
          }
        });
      }
    },

    initSearch() {
      const searchToggle = document.querySelector(".search-toggle");
      const searchModal = document.querySelector(".search-modal");
      const searchInput = document.querySelector(".search-input");
      const searchClose = document.querySelector(".search-close");
      const searchResults = document.querySelector(".search-results");

      if (!searchToggle || !searchModal) return;

      let searchTimeout;

      let scrollPos = 0;

      // Open search modal
      searchToggle.addEventListener("click", () => {
        scrollPos = window.scrollY;
        document.body.style.top = `-${scrollPos}px`;
        document.body.classList.add("modal-open");
        searchModal.classList.add("active");
        searchInput.focus();
      });

      // Close search modal
      const closeSearch = () => {
        document.body.classList.remove("modal-open");
        document.body.style.top = "";
        window.scrollTo(0, scrollPos);
        searchModal.classList.remove("active");
        searchInput.value = "";
        searchResults.innerHTML = "";
      };

      searchClose.addEventListener("click", closeSearch);
      searchModal.addEventListener("click", (e) => {
        if (e.target === searchModal) {
          closeSearch();
        }
      });

      // Escape to close
      document.addEventListener("keydown", (e) => {
        if (e.key === "Escape" && searchModal.classList.contains("active")) {
          closeSearch();
        }
        // Cmd/Ctrl+K to open search
        if ((e.metaKey || e.ctrlKey) && e.key === "k") {
          e.preventDefault();
          searchModal.classList.add("active");
          searchInput.focus();
        }
      });

      // Perform search
      searchInput.addEventListener("input", (e) => {
        const query = e.target.value.trim();

        clearTimeout(searchTimeout);

        if (!query) {
          searchResults.innerHTML = "";
          return;
        }

        searchResults.innerHTML =
          '<div class="search-loading">Searching...</div>';

        searchTimeout = setTimeout(async () => {
          try {
            const response = await fetch(
              `/api/search?q=${encodeURIComponent(query)}`
            );
            if (!response.ok) throw new Error("Search failed");

            const results = await response.json();

            if (results.length === 0) {
              searchResults.innerHTML = `<div class="search-empty">No results for '${escapeHtml(
                query
              )}'</div>`;
              return;
            }

            searchResults.innerHTML = results
              .map(
                (result) => `
                            <a href="${result.URL}" class="search-result">
                                <div class="search-result-title">${escapeHtml(
                                  result.Title
                                )}</div>
                                <div class="search-result-meta">
                                    <span class="search-result-type">${
                                      result.Type
                                    }</span>
                                    ${
                                      result.Date
                                        ? `<span>•</span><span>${formatDate(
                                            result.Date
                                          )}</span>`
                                        : ""
                                    }
                                </div>
                                <div class="search-result-snippet">${
                                  result.Snippet
                                }</div>
                            </a>
                        `
              )
              .join("");

            // Close search modal on click
            searchResults.querySelectorAll(".search-result").forEach((link) => {
              link.addEventListener("click", () => {
                closeSearch();
              });
            });
          } catch (err) {
            console.error("Search error:", err);
            searchResults.innerHTML = `<div class="search-empty">No results for '${escapeHtml(
              query
            )}'</div>`;
          }
        }, 300);
      });

      // Helper functions
      function escapeHtml(text) {
        const div = document.createElement("div");
        div.textContent = text;
        return div.innerHTML;
      }

      function formatDate(dateStr) {
        const date = new Date(dateStr);
        return date.toLocaleDateString("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
        });
      }
    },

    initCollectionSort() {
      const sortSelect = document.querySelector(".sort-select");
      const postsList = document.querySelector(".posts-list-blog");

      if (!sortSelect || !postsList) return;

      // Check if this is a blog page (uses URL query params)
      const isBlogPage = sortSelect.hasAttribute("data-blog-sort");

      if (isBlogPage) {
        // Blog page: use URL query params
        const urlParams = new URLSearchParams(window.location.search);
        const sortFromUrl = urlParams.get("sort");
        const currentSort = sortFromUrl === "oldest" ? "oldest" : "newest";

        sortSelect.value = currentSort;
        this.sortPosts(postsList, currentSort);

        sortSelect.addEventListener("change", (e) => {
          const sortType = e.target.value;
          const url = new URL(window.location.href);

          if (sortType === "newest") {
            url.searchParams.delete("sort");
          } else {
            url.searchParams.set("sort", sortType);
          }

          window.history.replaceState({}, "", url);
          this.sortPosts(postsList, sortType);
        });
      } else {
        // Collection page: use localStorage
        const savedSort = localStorage.getItem("collectionSort") || "newest";
        sortSelect.value = savedSort;
        this.sortPosts(postsList, savedSort);

        sortSelect.addEventListener("change", (e) => {
          const sortType = e.target.value;
          localStorage.setItem("collectionSort", sortType);
          this.sortPosts(postsList, sortType);
        });
      }
    },

    sortPosts(postsList, sortType) {
      const posts = Array.from(postsList.querySelectorAll(".post-card"));

      posts.sort((a, b) => {
        const dateA = parseInt(a.dataset.date);
        const dateB = parseInt(b.dataset.date);
        const updatedA = parseInt(a.dataset.updated);
        const updatedB = parseInt(b.dataset.updated);

        switch (sortType) {
          case "newest":
            return dateB - dateA;
          case "oldest":
            return dateA - dateB;
          case "updated":
            return updatedB - updatedA;
          default:
            return dateB - dateA;
        }
      });

      // Re-append posts in new order
      posts.forEach((post) => postsList.appendChild(post));
    },

    initReactions() {
      const container = document.querySelector(".reactions");
      if (!container) return;

      const postSlug = container.dataset.post;
      if (!postSlug) return;

      // Skip if already initialized
      if (container.dataset.initialized === postSlug) return;
      container.dataset.initialized = postSlug;

      new Reactions(container, postSlug);
    },

    initComments() {
      const container = document.querySelector(".comments-section");
      if (!container) return;

      const postSlug = container.dataset.post;
      if (!postSlug) return;

      // Skip if already initialized
      if (container.dataset.initialized === postSlug) return;
      container.dataset.initialized = postSlug;

      new Comments(container, postSlug);
    },
  };

  // ========================================
  // Reactions Handler
  // ========================================
  class Reactions {
    constructor(container, postSlug) {
      this.container = container;
      this.postSlug = postSlug;
      this.userReactions = [];
      this.isLoggedIn = false;

      this.init();
    }

    async init() {
      this.attachHandlers();
      await Promise.all([this.fetchReactions(), this.fetchUserReactions()]);
    }

    attachHandlers() {
      this.container.querySelectorAll(".reaction-btn").forEach((btn) => {
        btn.addEventListener("click", () =>
          this.toggleReaction(btn.dataset.emoji)
        );
      });
    }

    async fetchReactions() {
      try {
        const response = await fetch(
          `/api/reactions?post=${encodeURIComponent(this.postSlug)}`
        );
        if (response.ok) {
          const data = await response.json();
          this.updateCounts(data);
        }
      } catch (err) {
        console.error("Failed to fetch reactions:", err);
      }
    }

    async fetchUserReactions() {
      try {
        const response = await fetch(
          `/api/reactions/user?post=${encodeURIComponent(this.postSlug)}`
        );
        if (response.ok) {
          this.userReactions = await response.json();
          this.isLoggedIn = true;
          this.updateActiveStates();
        } else if (response.status === 401) {
          this.isLoggedIn = false;
        }
      } catch (err) {
        console.error("Failed to fetch user reactions:", err);
      }
    }

    updateCounts(data) {
      const reactionData = Object.fromEntries(
        data.map((item) => [item.emoji, { count: item.count, users: item.users || [] }])
      );

      this.container.querySelectorAll(".reaction-btn").forEach((btn) => {
        const emoji = btn.dataset.emoji;
        const info = reactionData[emoji] || { count: 0, users: [] };
        const countEl = btn.querySelector(".count");
        if (countEl) {
          countEl.textContent = info.count;
        }
        btn.title = this.buildTooltip(info.users, info.count);
      });
    }

    buildTooltip(users, count) {
      if (count === 0 || users.length === 0) return "";
      if (count === 1) return users[0];
      if (count === 2) return `${users[0]} and ${users[1]}`;
      if (count === 3 && users.length === 3) return `${users[0]}, ${users[1]} and ${users[2]}`;
      const others = count - users.length;
      if (others > 0) {
        return `${users.slice(0, 2).join(", ")} and ${others + (users.length > 2 ? 1 : 0)} others`;
      }
      return `${users.slice(0, -1).join(", ")} and ${users[users.length - 1]}`;
    }

    updateActiveStates() {
      this.container.querySelectorAll(".reaction-btn").forEach((btn) => {
        btn.classList.toggle(
          "active",
          this.userReactions.includes(btn.dataset.emoji)
        );
      });
    }

    async toggleReaction(emoji) {
      if (!this.isLoggedIn) {
        LoginModal.open(window.location.pathname);
        return;
      }

      try {
        const response = await fetch("/api/reactions", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ post: this.postSlug, emoji }),
        });

        if (response.ok) {
          const data = await response.json();

          if (data.added) {
            this.userReactions.push(emoji);
          } else {
            this.userReactions = this.userReactions.filter((e) => e !== emoji);
          }

          this.updateActiveStates();
          await this.fetchReactions();
        }
      } catch (err) {
        console.error("Failed to toggle reaction:", err);
      }
    }
  }

  // ========================================
  // Comments Handler
  // ========================================
  class Comments {
    constructor(container, postSlug) {
      this.container = container;
      this.postSlug = postSlug;
      this.isLoggedIn = false;
      this.currentUser = null;
      this.comments = [];
      this.isPreviewMode = false;
      this.editingCommentId = null;

      this.form = container.querySelector(".comment-form");
      this.textarea = container.querySelector(".comment-input");
      this.preview = container.querySelector(".comment-preview");
      this.previewToggle = container.querySelector(".preview-toggle");
      this.submitBtn = container.querySelector(".comment-submit");
      this.commentsList = container.querySelector(".comments-list");
      this.avatarPlaceholder = container.querySelector(".avatar-placeholder");

      this.init();
    }

    async init() {
      this.attachHandlers();
      await Promise.all([this.checkAuth(), this.fetchComments()]);
      this.updateFormState();
    }

    attachHandlers() {
      // Submit button
      this.submitBtn.addEventListener("click", () => this.submitComment());

      // Preview toggle
      this.previewToggle.addEventListener("click", () => this.togglePreview());

      // Enable/disable submit based on content
      this.textarea.addEventListener("input", () => {
        this.submitBtn.disabled = !this.textarea.value.trim();
        if (this.isPreviewMode) {
          this.renderPreview();
        }
      });

      // Submit on Ctrl/Cmd+Enter
      this.textarea.addEventListener("keydown", (e) => {
        if (
          (e.ctrlKey || e.metaKey) &&
          e.key === "Enter" &&
          this.textarea.value.trim()
        ) {
          this.submitComment();
        }
      });
    }

    async checkAuth() {
      try {
        const response = await fetch("/api/me");
        if (response.ok) {
          this.currentUser = await response.json();
          this.isLoggedIn = true;
        }
      } catch (err) {
        console.error("Failed to check auth:", err);
      }
    }

    updateFormState() {
      if (this.isLoggedIn && this.currentUser) {
        this.textarea.disabled = false;
        this.textarea.placeholder = "Write a comment... (Markdown supported)";

        // Update avatar if available
        if (this.currentUser.avatar) {
          this.avatarPlaceholder.innerHTML = `<img src="${this.currentUser.avatar}" alt="${this.currentUser.name}" class="comment-avatar-img">`;
        } else {
          this.avatarPlaceholder.innerHTML = `<span class="avatar-initial">${(this
            .currentUser.name || "U")[0].toUpperCase()}</span>`;
        }
      } else {
        this.textarea.disabled = true;
        this.textarea.placeholder = "Sign in to comment";
        this.submitBtn.disabled = true;
      }
    }

    async fetchComments() {
      try {
        const response = await fetch(
          `/api/comments?post=${encodeURIComponent(this.postSlug)}`
        );
        if (response.ok) {
          this.comments = await response.json();
          this.renderComments();
        }
      } catch (err) {
        console.error("Failed to fetch comments:", err);
      }
    }

    renderComments() {
      if (!this.comments || this.comments.length === 0) {
        this.commentsList.innerHTML =
          '<p class="no-comments">No comments yet. Be the first to comment!</p>';
        return;
      }

      this.commentsList.innerHTML = this.comments
        .map((comment) => this.renderComment(comment))
        .join("");

      // Attach event handlers for edit/delete buttons
      this.commentsList.querySelectorAll(".comment-edit-btn").forEach((btn) => {
        btn.addEventListener("click", () =>
          this.startEdit(parseInt(btn.dataset.id))
        );
      });

      this.commentsList
        .querySelectorAll(".comment-delete-btn")
        .forEach((btn) => {
          btn.addEventListener("click", () =>
            this.deleteComment(parseInt(btn.dataset.id))
          );
        });

      this.commentsList.querySelectorAll(".comment-save-btn").forEach((btn) => {
        btn.addEventListener("click", () =>
          this.saveEdit(parseInt(btn.dataset.id))
        );
      });

      this.commentsList
        .querySelectorAll(".comment-cancel-btn")
        .forEach((btn) => {
          btn.addEventListener("click", () =>
            this.cancelEdit(parseInt(btn.dataset.id))
          );
        });
    }

    renderComment(comment) {
      const isOwn = this.currentUser && comment.userId === this.currentUser.id;
      const avatar = comment.userAvatar
        ? `<img src="${comment.userAvatar}" alt="${comment.userName}" class="comment-avatar-img">`
        : `<span class="avatar-initial">${(comment.userName ||
            "U")[0].toUpperCase()}</span>`;

      const date = new Date(comment.createdAt);
      const timeAgo = this.formatTimeAgo(date);
      const isEdited = comment.updatedAt !== comment.createdAt;

      return `
                <div class="comment" data-id="${comment.id}">
                    <div class="comment-avatar">${avatar}</div>
                    <div class="comment-main">
                        <div class="comment-header">
                            <span class="comment-author">${this.escapeHtml(
                              comment.userName
                            )}</span>
                            <span class="comment-time" title="${date.toLocaleString()}">${timeAgo}${
        isEdited ? " (edited)" : ""
      }</span>
                            ${
                              isOwn
                                ? `
                                <div class="comment-actions">
                                    <button class="comment-edit-btn" data-id="${comment.id}">Edit</button>
                                    <button class="comment-delete-btn" data-id="${comment.id}">Delete</button>
                                </div>
                            `
                                : ""
                            }
                        </div>
                        <div class="comment-body">${comment.contentHtml}</div>
                        <div class="comment-edit-form" style="display: none;">
                            <textarea class="comment-edit-input">${this.escapeHtml(
                              comment.content
                            )}</textarea>
                            <div class="comment-edit-actions">
                                <button class="comment-cancel-btn" data-id="${
                                  comment.id
                                }">Cancel</button>
                                <button class="comment-save-btn" data-id="${
                                  comment.id
                                }">Save</button>
                            </div>
                        </div>
                    </div>
                </div>
            `;
    }

    formatTimeAgo(date) {
      const seconds = Math.floor((new Date() - date) / 1000);

      if (seconds < 60) return "just now";
      if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
      if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
      if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`;

      return date.toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
      });
    }

    escapeHtml(text) {
      const div = document.createElement("div");
      div.textContent = text;
      return div.innerHTML;
    }

    togglePreview() {
      this.isPreviewMode = !this.isPreviewMode;

      if (this.isPreviewMode) {
        this.renderPreview();
        this.preview.style.display = "block";
        this.textarea.style.display = "none";
        this.previewToggle.textContent = "Edit";
      } else {
        this.preview.style.display = "none";
        this.textarea.style.display = "block";
        this.previewToggle.textContent = "Preview";
      }
    }

    renderPreview() {
      const content = this.textarea.value;
      if (!content.trim()) {
        this.preview.innerHTML =
          '<p class="preview-empty">Nothing to preview</p>';
        return;
      }

      // Simple client-side markdown rendering for preview
      this.preview.innerHTML = this.simpleMarkdown(content);
    }

    // Simple markdown renderer for live preview
    simpleMarkdown(text) {
      return (
        text
          // Escape HTML first
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;")
          // Code blocks (``` ... ```)
          .replace(/```(\w*)\n([\s\S]*?)```/g, "<pre><code>$2</code></pre>")
          // Inline code
          .replace(/`([^`]+)`/g, "<code>$1</code>")
          // Bold
          .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
          .replace(/__([^_]+)__/g, "<strong>$1</strong>")
          // Italic
          .replace(/\*([^*]+)\*/g, "<em>$1</em>")
          .replace(/_([^_]+)_/g, "<em>$1</em>")
          // Links
          .replace(
            /\[([^\]]+)\]\(([^)]+)\)/g,
            '<a href="$2" target="_blank" rel="noopener">$1</a>'
          )
          // Line breaks
          .replace(/\n/g, "<br>")
      );
    }

    async submitComment() {
      if (!this.isLoggedIn) {
        LoginModal.open(window.location.pathname);
        return;
      }

      const content = this.textarea.value.trim();
      if (!content) return;

      this.submitBtn.disabled = true;
      this.submitBtn.textContent = "Posting...";

      try {
        const response = await fetch("/api/comments", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ post: this.postSlug, content }),
        });

        if (response.ok) {
          const newComment = await response.json();
          this.comments.unshift(newComment);
          this.renderComments();
          this.textarea.value = "";
          if (this.isPreviewMode) {
            this.togglePreview();
          }
        } else {
          throw new Error("Failed to post comment");
        }
      } catch (err) {
        console.error("Failed to post comment:", err);
        alert("Failed to post comment. Please try again.");
      } finally {
        this.submitBtn.disabled = false;
        this.submitBtn.textContent = "Comment";
      }
    }

    startEdit(commentId) {
      const commentEl = this.container.querySelector(
        `.comment[data-id="${commentId}"]`
      );
      if (!commentEl) return;

      commentEl.querySelector(".comment-body").style.display = "none";
      commentEl.querySelector(".comment-actions").style.display = "none";
      commentEl.querySelector(".comment-edit-form").style.display = "block";
    }

    cancelEdit(commentId) {
      const commentEl = this.container.querySelector(
        `.comment[data-id="${commentId}"]`
      );
      if (!commentEl) return;

      const comment = this.comments.find((c) => c.id === commentId);
      if (comment) {
        commentEl.querySelector(".comment-edit-input").value = comment.content;
      }

      commentEl.querySelector(".comment-body").style.display = "block";
      commentEl.querySelector(".comment-actions").style.display = "flex";
      commentEl.querySelector(".comment-edit-form").style.display = "none";
    }

    async saveEdit(commentId) {
      const commentEl = this.container.querySelector(
        `.comment[data-id="${commentId}"]`
      );
      if (!commentEl) return;

      const textarea = commentEl.querySelector(".comment-edit-input");
      const content = textarea.value.trim();
      if (!content) return;

      const saveBtn = commentEl.querySelector(".comment-save-btn");
      saveBtn.disabled = true;
      saveBtn.textContent = "Saving...";

      try {
        const response = await fetch(`/api/comments/${commentId}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ content }),
        });

        if (response.ok) {
          const updatedComment = await response.json();
          const index = this.comments.findIndex((c) => c.id === commentId);
          if (index !== -1) {
            this.comments[index] = updatedComment;
          }
          this.renderComments();
        } else {
          throw new Error("Failed to update comment");
        }
      } catch (err) {
        console.error("Failed to update comment:", err);
        alert("Failed to update comment. Please try again.");
        saveBtn.disabled = false;
        saveBtn.textContent = "Save";
      }
    }

    async deleteComment(commentId) {
      if (!confirm("Are you sure you want to delete this comment?")) return;

      try {
        const response = await fetch(`/api/comments/${commentId}`, {
          method: "DELETE",
        });

        if (response.ok) {
          this.comments = this.comments.filter((c) => c.id !== commentId);
          this.renderComments();
        } else {
          throw new Error("Failed to delete comment");
        }
      } catch (err) {
        console.error("Failed to delete comment:", err);
        alert("Failed to delete comment. Please try again.");
      }
    }
  }

  // ========================================
  // Keyboard Shortcuts
  // ========================================
  const KeyboardShortcuts = {
    init() {
      document.addEventListener("keydown", (e) => {
        // Escape - close sidebar and nav
        if (e.key === "Escape") {
          if (
            document.body.classList.contains("sidebar-open") ||
            document.body.classList.contains("nav-open")
          ) {
            document.body.classList.remove("sidebar-open", "nav-open");
            const menuToggle = document.querySelector(".menu-toggle");
            if (menuToggle) menuToggle.setAttribute("aria-expanded", "false");
          }
        }
      });
    },
  };

  // ========================================
  // Profile Dropdown
  // ========================================
  const ProfileDropdown = {
    dropdown: null,
    toggle: null,
    user: null,

    init() {
      this.dropdown = document.querySelector(".profile-dropdown");
      this.toggle = document.querySelector(".profile-toggle");

      if (this.dropdown && this.toggle) {
        this.attachHandlers();
      }

      // Hydrate header with user data from API
      this.hydrateUser();
    },

    attachHandlers() {
      if (this.toggle.dataset.initialized) return;
      this.toggle.dataset.initialized = "true";

      this.toggle.addEventListener("click", (e) => {
        e.stopPropagation();
        this.toggleMenu();
      });

      // Close on click outside
      document.addEventListener("click", (e) => {
        if (this.dropdown && !this.dropdown.contains(e.target)) {
          this.close();
        }
      });

      // Close on Escape
      document.addEventListener("keydown", (e) => {
        if (
          e.key === "Escape" &&
          this.dropdown &&
          this.dropdown.classList.contains("open")
        ) {
          this.close();
        }
      });
    },

    async hydrateUser() {
      try {
        const response = await fetch("/api/me");
        if (!response.ok) return;

        this.user = await response.json();
        this.updateHeader();
      } catch (err) {
        // Not logged in, keep sign-in button
      }
    },

    updateHeader() {
      if (!this.user) return;

      const headerActions = document.querySelector(".header-actions");
      if (!headerActions) return;

      // Check if already showing profile dropdown
      if (headerActions.querySelector(".profile-dropdown")) return;

      // Build avatar HTML
      const avatarHtml = this.user.avatar
        ? `<img src="${this.user.avatar}" alt="${this.user.name}" class="profile-avatar" />`
        : `<span class="profile-avatar-initial">${(this.user.name || "U")[0].toUpperCase()}</span>`;

      // Replace sign-in button with profile dropdown
      headerActions.innerHTML = `
        <div class="profile-dropdown">
          <button class="profile-toggle" aria-expanded="false" aria-haspopup="true">
            ${avatarHtml}
            <svg class="profile-chevron" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>
          <div class="profile-menu">
            <div class="profile-menu-header">
              <span class="profile-menu-name">${this.escapeHtml(this.user.name)}</span>
              ${this.user.email ? `<span class="profile-menu-email">${this.escapeHtml(this.user.email)}</span>` : ""}
            </div>
            <div class="profile-menu-divider"></div>
            <a href="/auth/logout?redirect=${encodeURIComponent(window.location.pathname)}" class="profile-menu-item" data-no-router>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                <polyline points="16 17 21 12 16 7"></polyline>
                <line x1="21" y1="12" x2="9" y2="12"></line>
              </svg>
              Sign out
            </a>
          </div>
        </div>
      `;

      // Re-initialize dropdown handlers
      this.dropdown = headerActions.querySelector(".profile-dropdown");
      this.toggle = headerActions.querySelector(".profile-toggle");
      this.attachHandlers();
    },

    escapeHtml(text) {
      const div = document.createElement("div");
      div.textContent = text;
      return div.innerHTML;
    },

    toggleMenu() {
      if (!this.dropdown) return;
      const isOpen = this.dropdown.classList.toggle("open");
      this.toggle.setAttribute("aria-expanded", isOpen);
    },

    close() {
      if (!this.dropdown) return;
      this.dropdown.classList.remove("open");
      this.toggle.setAttribute("aria-expanded", "false");
    },
  };

  // ========================================
  // Client-Side Router with Aggressive Prefetching
  // ========================================
  const Router = {
    cache: new Map(),
    currentUser: null,
    isNavigating: false,
    prefetchQueue: new Set(),
    prefetching: new Set(),
    maxConcurrentPrefetch: 3,
    shouldPrefetch: true,
    onPageLoad: null,

    init() {
      // Check prefetch preferences
      this.checkPrefetchPreferences();

      // Load user state for persistence
      this.loadUserState();

      // Setup link interception
      document.addEventListener("click", (e) => {
        const link = e.target.closest("a");
        if (!link) return;

        const href = link.getAttribute("href");
        if (!this.shouldInterceptLink(link, href)) return;

        e.preventDefault();
        this.navigate(href);
      });

      // Handle browser back/forward
      window.addEventListener("popstate", (e) => {
        if (e.state && e.state.path) {
          this.loadPage(e.state.path, false);
        }
      });

      // Store initial state
      history.replaceState({ path: location.pathname }, "", location.pathname);

      // Setup prefetching
      this.setupPrefetching();
    },

    checkPrefetchPreferences() {
      const connection =
        navigator.connection ||
        navigator.mozConnection ||
        navigator.webkitConnection;

      // Don't prefetch on save-data
      if (connection?.saveData) {
        this.shouldPrefetch = false;
        return;
      }

      // Don't prefetch on slow connections
      if (
        connection?.effectiveType === "slow-2g" ||
        connection?.effectiveType === "2g"
      ) {
        this.shouldPrefetch = false;
        return;
      }
    },

    setupPrefetching() {
      if (!this.shouldPrefetch) return;

      // Delay prefetching until after initial page load to not compete for bandwidth
      const startPrefetching = () => {
        // Viewport prefetching with Intersection Observer
        this.setupViewportPrefetching();

        // Hover prefetching
        this.setupHoverPrefetching();

        // Touch prefetching for mobile
        this.setupTouchPrefetching();
      };

      // Wait for page to be fully loaded before prefetching
      if (document.readyState === "complete") {
        // Use idle callback if available, otherwise small delay
        if ("requestIdleCallback" in window) {
          requestIdleCallback(startPrefetching, { timeout: 2000 });
        } else {
          setTimeout(startPrefetching, 1000);
        }
      } else {
        window.addEventListener("load", () => {
          if ("requestIdleCallback" in window) {
            requestIdleCallback(startPrefetching, { timeout: 2000 });
          } else {
            setTimeout(startPrefetching, 1000);
          }
        });
      }
    },

    setupViewportPrefetching() {
      if (!("IntersectionObserver" in window)) return;

      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              const url = entry.target.href;
              if (this.shouldPrefetchUrl(url)) {
                this.prefetchQueue.add(url);
              }
            }
          });

          // Process queue when browser is idle
          if ("requestIdleCallback" in window) {
            requestIdleCallback(() => this.processPrefetchQueue());
          } else {
            setTimeout(() => this.processPrefetchQueue(), 0);
          }
        },
        {
          rootMargin: "50px",
          threshold: 0.01,
        }
      );

      const observeLinks = () => {
        document.querySelectorAll("a[href]").forEach((link) => {
          const href = link.getAttribute("href");
          if (this.shouldInterceptLink(link, href)) {
            observer.observe(link);
          }
        });
      };

      observeLinks();
      this.onPageLoad = observeLinks;
    },

    setupHoverPrefetching() {
      let hoverTimer;

      document.addEventListener(
        "mouseover",
        (e) => {
          const link = e.target.closest("a");
          if (!link) return;

          const href = link.getAttribute("href");
          if (!this.shouldInterceptLink(link, href)) return;

          hoverTimer = setTimeout(() => {
            this.prefetch(href);
          }, 200);
        },
        { passive: true }
      );

      document.addEventListener(
        "mouseout",
        () => {
          if (hoverTimer) {
            clearTimeout(hoverTimer);
            hoverTimer = null;
          }
        },
        { passive: true }
      );
    },

    setupTouchPrefetching() {
      document.addEventListener(
        "touchstart",
        (e) => {
          const link = e.target.closest("a");
          if (!link) return;

          const href = link.getAttribute("href");
          if (this.shouldInterceptLink(link, href)) {
            this.prefetch(href);
          }
        },
        { passive: true }
      );
    },

    processPrefetchQueue() {
      const currentlyPrefetching = this.prefetching.size;
      const available = this.maxConcurrentPrefetch - currentlyPrefetching;

      if (available <= 0) return;

      const toPrefetch = Array.from(this.prefetchQueue)
        .filter((url) => !this.cache.has(url) && !this.prefetching.has(url))
        .slice(0, available);

      toPrefetch.forEach((url) => {
        this.prefetchQueue.delete(url);
        this.prefetch(url);
      });
    },

    async prefetch(url) {
      if (!url || this.cache.has(url) || this.prefetching.has(url)) {
        return;
      }

      this.prefetching.add(url);

      try {
        const response = await fetch(url);
        if (response.ok) {
          const html = await response.text();
          this.cache.set(url, html);
        }
      } catch (err) {
        // Silent fail for prefetch
      } finally {
        this.prefetching.delete(url);

        if (this.prefetchQueue.size > 0) {
          this.processPrefetchQueue();
        }
      }
    },

    shouldInterceptLink(link, href) {
      if (!href) return false;
      if (href.startsWith("#")) return false;
      if (href.startsWith("http") && !href.startsWith(location.origin))
        return false;
      if (href.startsWith("mailto:") || href.startsWith("tel:")) return false;
      if (link.hasAttribute("download")) return false;
      if (link.target === "_blank") return false;
      if (link.closest("[data-no-router]")) return false;
      if (href.includes("/auth/") || href.includes("/api/")) return false;
      return true;
    },

    shouldPrefetchUrl(url) {
      if (!url || url === location.pathname) return false;
      if (this.cache.has(url)) return false;
      if (this.prefetching.has(url)) return false;
      return true;
    },

    async navigate(path) {
      if (this.isNavigating || path === location.pathname) return;

      this.isNavigating = true;
      document.body.classList.add("page-loading");

      try {
        history.pushState({ path }, "", path);
        await this.loadPage(path, true);
      } catch (err) {
        console.error("Navigation error:", err);
        location.href = path;
      } finally {
        this.isNavigating = false;
        document.body.classList.remove("page-loading");
      }
    },

    async loadPage(path, pushState = true) {
      // Check cache first (may have been prefetched)
      let html = this.cache.get(path);

      if (!html) {
        const response = await fetch(path);
        if (!response.ok) throw new Error(`Failed to load ${path}`);
        html = await response.text();
        this.cache.set(path, html);
      }

      this.updateContent(html, path);

      // Scroll to top
      if (!path.includes("#")) {
        window.scrollTo({ top: 0, behavior: "instant" });
      }

      // Close mobile menu
      document.body.classList.remove("sidebar-open", "nav-open");
      const menuToggle = document.querySelector(".menu-toggle");
      if (menuToggle) menuToggle.setAttribute("aria-expanded", "false");

      // Re-initialize page components
      PageComponents.init();
      ProfileDropdown.init();

      // Trigger viewport prefetching for new page
      if (this.onPageLoad) this.onPageLoad();
    },

    updateContent(html, path) {
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, "text/html");

      const newMain = doc.querySelector("#main-content");
      const currentMain = document.querySelector("#main-content");
      if (newMain && currentMain) {
        currentMain.innerHTML = newMain.innerHTML;
      }

      const newTitle = doc.querySelector("title");
      if (newTitle) {
        document.title = newTitle.textContent;
      }

      this.updateActiveNav(path);

      const canonical = document.querySelector('link[rel="canonical"]');
      const newCanonical = doc.querySelector('link[rel="canonical"]');
      if (canonical && newCanonical) {
        canonical.href = newCanonical.href;
      }

      const metaDesc = document.querySelector('meta[name="description"]');
      const newMetaDesc = doc.querySelector('meta[name="description"]');
      if (metaDesc && newMetaDesc) {
        metaDesc.content = newMetaDesc.content;
      }
    },

    updateActiveNav(path) {
      document.querySelectorAll(".site-nav a").forEach((link) => {
        link.classList.remove("active");
        const href = link.getAttribute("href");
        if (path === href || (href !== "/" && path.startsWith(href))) {
          link.classList.add("active");
        }
      });
    },

    loadUserState() {
      const userNameEl = document.querySelector(".profile-menu-name");
      if (userNameEl) {
        this.currentUser = {
          name: userNameEl.textContent,
        };
      }
    },
  };

  // ========================================
  // Initialize Everything
  // ========================================
  LoginModal.init();
  PageComponents.init();
  KeyboardShortcuts.init();
  ProfileDropdown.init();
  Router.init();
})();
